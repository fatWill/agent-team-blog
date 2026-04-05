package upload

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatWill/agent-team-blog/backend/config"
	"github.com/gin-gonic/gin"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

var (
	allowedExtensions = map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	allowedMIMEs = map[string]bool{
		"image/jpeg": true, "image/png": true, "image/gif": true, "image/webp": true,
	}
	uploadIDRegex = regexp.MustCompile(`^[\w-]+$`)
)

var (
	uploadCfg *config.UploadConfig
	cosCfg    *config.COSConfig
	cosClient *cos.Client
)

// SetUploadConfig 设置上传配置
func SetUploadConfig(cfg *config.UploadConfig) {
	uploadCfg = cfg
}

// SetCOSConfig 设置 COS 配置并初始化客户端
func SetCOSConfig(cfg *config.COSConfig) {
	cosCfg = cfg
	if cfg.SecretID == "" || cfg.SecretKey == "" {
		log.Println("⚠️  COS_ID 或 COS_KEY 未设置，上传功能将不可用")
		return
	}
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cfg.Bucket, cfg.Region))
	cosClient = cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})
	log.Println("✅ COS 客户端初始化成功")
}

// RandomString 生成随机字符串（导出供其他模块使用）
func RandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)[:n]
}

// cosKey 生成 COS 存储路径：upload/时间戳-随机串.ext
func cosKey(ext string) string {
	now := time.Now()
	filename := fmt.Sprintf("%d-%s%s", now.UnixMilli(), RandomString(8), ext)
	return fmt.Sprintf("upload/%s", filename)
}

// uploadToCOS 将数据通过 PutObject 上传到 COS，返回公开访问 URL（适用于 ≤2MB 小文件）
func uploadToCOS(data []byte, key string) (string, error) {
	if cosClient == nil {
		return "", fmt.Errorf("COS 客户端未初始化，请检查 COS_SECRET_ID 和 COS_SECRET_KEY 环境变量")
	}
	_, err := cosClient.Object.Put(context.Background(), key, bytes.NewReader(data), nil)
	if err != nil {
		return "", fmt.Errorf("上传 COS 失败: %w", err)
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(cosCfg.BaseURL, "/"), key), nil
}

// cosMultipartChunkSize COS 分片上传每片大小（1MB）
const cosMultipartChunkSize = 1 * 1024 * 1024

// multipartUploadToCOS 使用 COS 原生分片上传（适用于 >2MB 大文件）
// 流程：InitiateMultipartUpload → UploadPart → CompleteMultipartUpload
func multipartUploadToCOS(data []byte, key string) (string, error) {
	if cosClient == nil {
		return "", fmt.Errorf("COS 客户端未初始化，请检查 COS_SECRET_ID 和 COS_SECRET_KEY 环境变量")
	}
	ctx := context.Background()

	// 1. 初始化分片上传
	initResult, _, err := cosClient.Object.InitiateMultipartUpload(ctx, key, nil)
	if err != nil {
		return "", fmt.Errorf("COS 初始化分片上传失败: %w", err)
	}
	uploadID := initResult.UploadID

	// 2. 逐片上传
	totalSize := len(data)
	var parts []cos.Object
	partNumber := 1
	for offset := 0; offset < totalSize; offset += cosMultipartChunkSize {
		end := offset + cosMultipartChunkSize
		if end > totalSize {
			end = totalSize
		}
		chunk := data[offset:end]

		resp, err := cosClient.Object.UploadPart(ctx, key, uploadID, partNumber, bytes.NewReader(chunk), nil)
		if err != nil {
			// 上传失败，尝试中止分片上传
			cosClient.Object.AbortMultipartUpload(ctx, key, uploadID)
			return "", fmt.Errorf("COS 上传分片 %d 失败: %w", partNumber, err)
		}
		etag := resp.Header.Get("ETag")
		parts = append(parts, cos.Object{
			PartNumber: partNumber,
			ETag:       etag,
		})
		partNumber++
	}

	// 3. 完成分片上传
	completeOpt := &cos.CompleteMultipartUploadOptions{
		Parts: parts,
	}
	_, _, err = cosClient.Object.CompleteMultipartUpload(ctx, key, uploadID, completeOpt)
	if err != nil {
		return "", fmt.Errorf("COS 合并分片失败: %w", err)
	}

	return fmt.Sprintf("%s/%s", strings.TrimRight(cosCfg.BaseURL, "/"), key), nil
}

// smartUploadSize 小文件/大文件分界线（2MB）
const smartUploadSize = 2 * 1024 * 1024

// smartUploadToCOS 根据文件大小智能选择上传方式
// ≤2MB: PutObject 直传
// >2MB: COS 原生分片上传
func smartUploadToCOS(data []byte, key string) (string, error) {
	if len(data) <= smartUploadSize {
		return uploadToCOS(data, key)
	}
	return multipartUploadToCOS(data, key)
}

// DeleteFromCOS 根据完整 URL 删除 COS 对象（导出供其他模块调用）
// url 格式：https://assets.fatwill.cloud/upload/20260405/xxx.jpg
// 提取 key：upload/20260405/xxx.jpg
func DeleteFromCOS(fileURL string) error {
	if cosClient == nil {
		return fmt.Errorf("COS 客户端未初始化")
	}
	if fileURL == "" {
		return nil
	}
	// 从完整 URL 中提取 COS key
	baseURL := strings.TrimRight(cosCfg.BaseURL, "/") + "/"
	if !strings.HasPrefix(fileURL, baseURL) {
		// 非 COS 管理的 URL，跳过删除
		return nil
	}
	key := strings.TrimPrefix(fileURL, baseURL)
	if key == "" {
		return nil
	}
	_, err := cosClient.Object.Delete(context.Background(), key)
	if err != nil {
		return fmt.Errorf("删除 COS 对象失败: %w", err)
	}
	return nil
}

// BatchDeleteFromCOS 批量删除 COS 对象
func BatchDeleteFromCOS(fileURLs []string) {
	for _, u := range fileURLs {
		if u != "" {
			if err := DeleteFromCOS(u); err != nil {
				log.Printf("⚠️  COS 删除失败 [%s]: %v", u, err)
			}
		}
	}
}

// Upload POST /api/upload
func Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "请选择要上传的文件"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType != "" && !allowedMIMEs[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "不支持的文件格式，仅支持 jpg/jpeg/png/gif/webp"})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "不支持的文件格式，仅支持 jpg/jpeg/png/gif/webp"})
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "读取文件失败"})
		return
	}

	key := cosKey(ext)
	fileURL, err := smartUploadToCOS(data, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": fileURL})
}

// UploadChunk POST /api/upload/chunk
// 分片仍然临时保存到本地，不上传 COS
func UploadChunk(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "缺少 file 字段"})
		return
	}
	defer file.Close()

	uploadID := c.PostForm("uploadId")
	chunkIndexStr := c.PostForm("chunkIndex")
	totalChunksStr := c.PostForm("totalChunks")
	_ = c.PostForm("filename")

	if uploadID == "" || chunkIndexStr == "" || totalChunksStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "缺少必要参数：uploadId, chunkIndex, totalChunks, filename"})
		return
	}

	chunkIndex, err1 := strconv.Atoi(chunkIndexStr)
	totalChunks, err2 := strconv.Atoi(totalChunksStr)
	if err1 != nil || err2 != nil || chunkIndex < 0 || totalChunks <= 0 || chunkIndex >= totalChunks {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "chunkIndex 或 totalChunks 参数无效"})
		return
	}

	if !uploadIDRegex.MatchString(uploadID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "uploadId 包含非法字符"})
		return
	}

	tmpDir := filepath.Join(uploadCfg.TmpDir, uploadID)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "创建临时目录失败"})
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "读取分片数据失败"})
		return
	}

	chunkPath := filepath.Join(tmpDir, strconv.Itoa(chunkIndex))
	if err := os.WriteFile(chunkPath, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "保存分片失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": chunkIndex, "total": totalChunks})
}

// MergeChunks POST /api/upload/merge
// 合并本地分片后上传到 COS，然后清理本地临时文件
func MergeChunks(c *gin.Context) {
	var body struct {
		UploadID    string `json:"uploadId"`
		TotalChunks int    `json:"totalChunks"`
		Filename    string `json:"filename"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.UploadID == "" || body.TotalChunks <= 0 || body.Filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "缺少必要参数：uploadId, totalChunks, filename"})
		return
	}

	if !uploadIDRegex.MatchString(body.UploadID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "uploadId 包含非法字符"})
		return
	}

	ext := strings.ToLower(filepath.Ext(body.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "不支持的文件格式，仅支持 jpg/jpeg/png/gif/webp"})
		return
	}

	tmpDir := filepath.Join(uploadCfg.TmpDir, body.UploadID)
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "上传会话不存在或已过期"})
		return
	}

	// 合并分片
	var merged []byte
	for i := 0; i < body.TotalChunks; i++ {
		chunkPath := filepath.Join(tmpDir, strconv.Itoa(i))
		data, err := os.ReadFile(chunkPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("分片 %d 缺失，请重新上传", i)})
			return
		}
		merged = append(merged, data...)
	}

	// 上传到 COS（根据大小智能选择直传或分片上传）
	key := cosKey(ext)
	fileURL, err := smartUploadToCOS(merged, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": err.Error()})
		return
	}

	// 清理本地临时文件
	os.RemoveAll(tmpDir)

	c.JSON(http.StatusOK, gin.H{"url": fileURL})
}

// CancelChunkUpload DELETE /api/upload/chunk
func CancelChunkUpload(c *gin.Context) {
	var body struct {
		UploadID string `json:"uploadId"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.UploadID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "缺少必要参数：uploadId"})
		return
	}

	if !uploadIDRegex.MatchString(body.UploadID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "uploadId 包含非法字符"})
		return
	}

	tmpDir := filepath.Join(uploadCfg.TmpDir, body.UploadID)
	os.RemoveAll(tmpDir)

	c.JSON(http.StatusOK, gin.H{"ok": true})
}