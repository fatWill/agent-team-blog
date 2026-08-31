package wechat

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
)

// syncQueue 同步任务队列（缓冲 channel）
var syncQueue = make(chan string, 100)

// StartSyncWorker 启动同步 worker（在 main 中调用）
func StartSyncWorker() {
	go func() {
		for articleID := range syncQueue {
			syncArticleToWechat(articleID)
		}
	}()
	log.Println("✅ 微信同步 worker 已启动")
}

// EnqueueSync 将文章加入同步队列
func EnqueueSync(articleID string) {
	select {
	case syncQueue <- articleID:
		log.Printf("[wechat-sync] 文章 %s 已加入同步队列", articleID)
	default:
		log.Printf("[wechat-sync] 同步队列已满，文章 %s 丢弃", articleID)
	}
}

// SyncArticle 手动触发同步（管理 API 使用）
func SyncArticle(articleID string) error {
	return syncArticleToWechat(articleID)
}

// syncArticleToWechat 执行文章同步到微信草稿箱
func syncArticleToWechat(articleID string) error {
	// 更新状态为 syncing
	db.DB.Model(&models.Article{}).Where("id = ?", articleID).
		Update("wechat_sync_status", "syncing")

	// 查询文章
	var article models.Article
	if err := db.DB.Where("id = ?", articleID).First(&article).Error; err != nil {
		markSyncFailed(articleID, "create_draft", "文章不存在: "+err.Error())
		return fmt.Errorf("文章不存在: %w", err)
	}

	// 检查是否启用自动同步
	if article.WechatAutoSync == 0 {
		log.Printf("[wechat-sync] 文章 %s 已关闭自动同步，跳过", articleID)
		return nil
	}

	// 转换 Tiptap JSON 为公众号 HTML
	convertResult, err := TiptapToWechatHTML([]byte(article.Content), articleID)
	if err != nil {
		markSyncFailed(articleID, "create_draft", "内容转换失败: "+err.Error())
		return fmt.Errorf("内容转换失败: %w", err)
	}

	// 上传文章中的图片到微信
	htmlContent := convertResult.HTML
	for _, imgURL := range convertResult.ImageURLs {
		wxURL, err := uploadArticleImage(imgURL)
		if err != nil {
			log.Printf("[wechat-sync] 图片上传失败 %s: %v（跳过该图）", imgURL, err)
			logSyncAction(articleID, "upload_image", "failed", fmt.Sprintf("图片 %s 上传失败: %v", imgURL, err))
			// 图片上传失败，保留原 URL 并标记
			htmlContent = strings.Replace(htmlContent, imgURL, imgURL, 1)
			continue
		}
		htmlContent = strings.Replace(htmlContent, imgURL, wxURL, 1)
		logSyncAction(articleID, "upload_image", "success", "")
	}

	// 清理 HTML
	htmlContent = CleanHTMLForWechat(htmlContent)

	// 上传封面图（如果有）
	thumbMediaID := ""
	if article.CoverImage != "" {
		mediaID, err := uploadCoverImage(article.CoverImage)
		if err != nil {
			log.Printf("[wechat-sync] 封面图上传失败: %v（使用空封面）", err)
			logSyncAction(articleID, "upload_image", "failed", fmt.Sprintf("封面图上传失败: %v", err))
		} else {
			thumbMediaID = mediaID
			logSyncAction(articleID, "upload_image", "success", "封面图上传成功")
		}
	}

	// 构建摘要（≤120字）
	digest := article.Summary
	if len([]rune(digest)) > 120 {
		digest = string([]rune(digest)[:117]) + "..."
	}

	// 构建草稿文章
	draftArticle := DraftArticle{
		Title:            article.Title,
		Author:           "fatwill",
		Digest:           digest,
		Content:          htmlContent,
		ContentSourceURL: articleURL(articleID),
		ThumbMediaID:     thumbMediaID,
		NeedOpenComment:  0,
		OnlyFansComment:  0,
	}

	// 判断是创建还是更新草稿
	if article.WechatDraftMediaID != nil && *article.WechatDraftMediaID != "" {
		// 更新已有草稿
		if err := UpdateDraft(*article.WechatDraftMediaID, draftArticle); err != nil {
			markSyncFailed(articleID, "update_draft", err.Error())
			return fmt.Errorf("更新草稿失败: %w", err)
		}
		logSyncAction(articleID, "update_draft", "success", "")
	} else {
		// 创建新草稿
		mediaID, err := AddDraft(draftArticle)
		if err != nil {
			markSyncFailed(articleID, "create_draft", err.Error())
			return fmt.Errorf("创建草稿失败: %w", err)
		}
		// 保存 media_id
		db.DB.Model(&models.Article{}).Where("id = ?", articleID).
			Update("wechat_draft_media_id", mediaID)
		logSyncAction(articleID, "create_draft", "success", "")
	}

	// 标记同步成功
	now := time.Now()
	db.DB.Model(&models.Article{}).Where("id = ?", articleID).Updates(map[string]interface{}{
		"wechat_sync_status": "success",
		"wechat_synced_at":   now,
		"wechat_sync_error":  nil,
	})

	log.Printf("[wechat-sync] 文章 %s 同步成功", articleID)
	return nil
}

// uploadArticleImage 下载并上传文章正文中的图片
func uploadArticleImage(imgURL string) (string, error) {
	data, filename, err := DownloadImage(imgURL)
	if err != nil {
		return "", err
	}
	return UploadImg(data, filename)
}

// uploadCoverImage 下载并上传封面图（永久素材）
func uploadCoverImage(imgURL string) (string, error) {
	data, filename, err := DownloadImage(imgURL)
	if err != nil {
		return "", err
	}
	return AddMaterial(data, filename)
}

// markSyncFailed 标记同步失败
func markSyncFailed(articleID, action, errMsg string) {
	db.DB.Model(&models.Article{}).Where("id = ?", articleID).Updates(map[string]interface{}{
		"wechat_sync_status": "failed",
		"wechat_sync_error":  errMsg,
	})
	logSyncAction(articleID, action, "failed", errMsg)
	log.Printf("[wechat-sync] 文章 %s 同步失败: %s", articleID, errMsg)
}

// logSyncAction 记录同步日志
func logSyncAction(articleID, action, status, errMsg string) {
	logEntry := models.WechatSyncLog{
		ArticleID: articleID,
		Action:    action,
		Status:    status,
	}
	if errMsg != "" {
		logEntry.Error = &errMsg
	}
	db.DB.Create(&logEntry)
}
