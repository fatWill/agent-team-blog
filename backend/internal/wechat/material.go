package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

// UploadImgResp 上传图文消息内图片的响应
type UploadImgResp struct {
	URL     string `json:"url"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// AddMaterialResp 上传永久素材的响应
type AddMaterialResp struct {
	MediaID string `json:"media_id"`
	URL     string `json:"url"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// DraftAddResp 创建草稿的响应
type DraftAddResp struct {
	MediaID string `json:"media_id"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// DraftArticle 草稿文章结构
type DraftArticle struct {
	Title            string `json:"title"`
	Author           string `json:"author"`
	Digest           string `json:"digest"`
	Content          string `json:"content"`
	ContentSourceURL string `json:"content_source_url"`
	ThumbMediaID     string `json:"thumb_media_id"`
	NeedOpenComment  int    `json:"need_open_comment"`
	OnlyFansComment  int    `json:"only_fans_can_comment"`
}

// UploadImg 上传图文消息内的图片（返回永久 URL，不占素材库配额）
// 用于文章正文中的图片
func UploadImg(imgData []byte, filename string) (string, error) {
	return uploadImgWithRetry(imgData, filename, true)
}

func uploadImgWithRetry(imgData []byte, filename string, canRetry bool) (string, error) {
	token, err := GetAccessToken()
	if err != nil {
		return "", fmt.Errorf("获取 access_token 失败: %w", err)
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=%s", token)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("media", filename)
	if err != nil {
		return "", fmt.Errorf("创建 multipart 失败: %w", err)
	}
	if _, err := part.Write(imgData); err != nil {
		return "", fmt.Errorf("写入图片数据失败: %w", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传图片请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result UploadImgResp
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// token 过期，强制刷新后重试一次
	if result.ErrCode == 40001 || result.ErrCode == 42001 {
		if canRetry {
			log.Printf("[wechat] uploadimg token 过期，强制刷新重试")
			ForceRefreshToken()
			return uploadImgWithRetry(imgData, filename, false)
		}
		return "", fmt.Errorf("微信接口错误（重试后仍失败）: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信接口错误: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	if result.URL == "" {
		return "", fmt.Errorf("微信返回空 URL")
	}

	return result.URL, nil
}

// AddMaterial 上传永久图片素材（用于封面图，返回 media_id）
func AddMaterial(imgData []byte, filename string) (string, error) {
	return addMaterialWithRetry(imgData, filename, true)
}

func addMaterialWithRetry(imgData []byte, filename string, canRetry bool) (string, error) {
	token, err := GetAccessToken()
	if err != nil {
		return "", fmt.Errorf("获取 access_token 失败: %w", err)
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/material/add_material?access_token=%s&type=image", token)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("media", filename)
	if err != nil {
		return "", fmt.Errorf("创建 multipart 失败: %w", err)
	}
	if _, err := part.Write(imgData); err != nil {
		return "", fmt.Errorf("写入图片数据失败: %w", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传素材请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result AddMaterialResp
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.ErrCode == 40001 || result.ErrCode == 42001 {
		if canRetry {
			log.Printf("[wechat] add_material token 过期，强制刷新重试")
			ForceRefreshToken()
			return addMaterialWithRetry(imgData, filename, false)
		}
		return "", fmt.Errorf("微信接口错误（重试后仍失败）: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信接口错误: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	return result.MediaID, nil
}

// AddDraft 创建草稿
func AddDraft(article DraftArticle) (string, error) {
	return addDraftWithRetry(article, true)
}

func addDraftWithRetry(article DraftArticle, canRetry bool) (string, error) {
	token, err := GetAccessToken()
	if err != nil {
		return "", fmt.Errorf("获取 access_token 失败: %w", err)
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/draft/add?access_token=%s", token)

	payload := map[string]interface{}{
		"articles": []DraftArticle{article},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建草稿请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result DraftAddResp
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.ErrCode == 40001 || result.ErrCode == 42001 {
		if canRetry {
			log.Printf("[wechat] draft/add token 过期，强制刷新重试")
			ForceRefreshToken()
			return addDraftWithRetry(article, false)
		}
		return "", fmt.Errorf("微信接口错误（重试后仍失败）: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信接口错误: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	return result.MediaID, nil
}

// UpdateDraft 更新草稿
func UpdateDraft(mediaID string, article DraftArticle) error {
	return updateDraftWithRetry(mediaID, article, true)
}

func updateDraftWithRetry(mediaID string, article DraftArticle, canRetry bool) error {
	token, err := GetAccessToken()
	if err != nil {
		return fmt.Errorf("获取 access_token 失败: %w", err)
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/draft/update?access_token=%s", token)

	payload := map[string]interface{}{
		"media_id": mediaID,
		"index":    0,
		"articles": article,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("更新草稿请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if result.ErrCode == 40001 || result.ErrCode == 42001 {
		if canRetry {
			log.Printf("[wechat] draft/update token 过期，强制刷新重试")
			ForceRefreshToken()
			return updateDraftWithRetry(mediaID, article, false)
		}
		return fmt.Errorf("微信接口错误（重试后仍失败）: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("微信接口错误: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	return nil
}

// DownloadImage 下载远程图片（用于上传到微信）
func DownloadImage(imgURL string) ([]byte, string, error) {
	resp, err := http.Get(imgURL)
	if err != nil {
		return nil, "", fmt.Errorf("下载图片失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("下载图片返回状态码 %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("读取图片数据失败: %w", err)
	}

	// 限制 2MB
	if len(data) > 2*1024*1024 {
		return nil, "", fmt.Errorf("图片超过 2MB 限制 (%d bytes)", len(data))
	}

	filename := filepath.Base(imgURL)
	if filename == "" || filename == "." {
		filename = "image.jpg"
	}

	return data, filename, nil
}
