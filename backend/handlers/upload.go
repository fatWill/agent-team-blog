package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatWill/agent-team-blog/backend/config"
	"github.com/gin-gonic/gin"
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

// uploadCfg 上传配置
var uploadCfg *config.UploadConfig

// SetUploadConfig 设置上传配置
func SetUploadConfig(cfg *config.UploadConfig) {
	uploadCfg = cfg
}

// randomString 生成随机字符串
func randomString(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:n]
}

// Upload POST /api/upload
func Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "请选择要上传的文件"})
		return
	}
	defer file.Close()

	// 校验 MIME 类型
	contentType := header.Header.Get("Content-Type")
	if contentType != "" && !allowedMIMEs[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "不支持的文件格式，仅支持 jpg/jpeg/png/gif/webp"})
		return
	}

	// 校验扩展名
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "不支持的文件格式，仅支持 jpg/jpeg/png/gif/webp"})
		return
	}

	// 确保上传目录存在
	if err := os.MkdirAll(uploadCfg.Dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "创建上传目录失败"})
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixMilli(), randomString(8), ext)
	filePath := filepath.Join(uploadCfg.Dir, filename)

	// 读取文件内容
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "读取文件失败"})
		return
	}

	// 写入磁盘
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "保存文件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": "/uploads/" + filename})
}

// UploadChunk POST /api/upload/chunk
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

	// 按序读取所有分片并合并
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

	// 确保上传目录存在
	os.MkdirAll(uploadCfg.Dir, 0755)

	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixMilli(), randomString(8), ext)
	filePath := filepath.Join(uploadCfg.Dir, filename)

	if err := os.WriteFile(filePath, merged, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "保存合并文件失败"})
		return
	}

	// 清理临时目录
	os.RemoveAll(tmpDir)

	c.JSON(http.StatusOK, gin.H{"url": "/uploads/" + filename})
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
