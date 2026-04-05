package photo

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/fatWill/agent-team-blog/backend/internal/album"
	"github.com/fatWill/agent-team-blog/backend/internal/upload"
	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UpdatePhoto PUT /api/photos/:id
func UpdatePhoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的照片 ID"})
		return
	}

	var body struct {
		Caption  *string `json:"caption"`
		Password *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	if body.Caption == nil && body.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "没有需要更新的字段"})
		return
	}

	updates := map[string]interface{}{}
	if body.Caption != nil {
		updates["caption"] = strings.TrimSpace(*body.Caption)
	}
	if body.Password != nil {
		if *body.Password == "" {
			updates["password_hash"] = nil
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(*body.Password)), 10)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "密码加密失败"})
				return
			}
			updates["password_hash"] = string(hash)
		}
	}

	result := db.DB.Model(&models.Photo{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "更新照片失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "照片不存在"})
		return
	}

	var p models.Photo
	db.DB.Where("id = ?", id).First(&p)

	c.JSON(http.StatusOK, models.PhotoListItem{
		ID:          p.ID,
		AlbumID:     p.AlbumID,
		URL:         p.URL,
		Caption:     p.Caption,
		HasPassword: p.PasswordHash != nil && *p.PasswordHash != "",
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	})
}

// DeletePhoto DELETE /api/photos/:id
func DeletePhoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的照片 ID"})
		return
	}

	var p models.Photo
	if err := db.DB.Select("id, album_id, url").Where("id = ?", id).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "照片不存在"})
		return
	}

	albumID := p.AlbumID
	photoURL := p.URL

	result := db.DB.Where("id = ?", id).Delete(&models.Photo{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "删除照片失败"})
		return
	}

	// 异步删除 COS 对象（不阻塞响应）
	go func() {
		if err := upload.DeleteFromCOS(photoURL); err != nil {
			log.Printf("⚠️  删除照片 COS 对象失败 [id=%d, url=%s]: %v", id, photoURL, err)
		}
	}()

	// 更新相册封面
	album.UpdateAlbumCover(albumID)

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GetPhotoLikes GET /api/photos/:id/likes
func GetPhotoLikes(c *gin.Context) {
	photoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || photoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的照片 ID"})
		return
	}

	deviceID := c.Query("deviceId")

	var count int64
	db.DB.Model(&models.PhotoLike{}).Where("photo_id = ?", photoID).Count(&count)

	liked := false
	if deviceID != "" {
		var likeCount int64
		db.DB.Model(&models.PhotoLike{}).Where("photo_id = ? AND device_id = ?", photoID, deviceID).Count(&likeCount)
		liked = likeCount > 0
	}

	c.JSON(http.StatusOK, gin.H{"count": count, "liked": liked})
}

// PostPhotoLike POST /api/photos/:id/likes
func PostPhotoLike(c *gin.Context) {
	photoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || photoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的照片 ID"})
		return
	}

	var body struct {
		DeviceID string `json:"deviceId"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	deviceID := strings.TrimSpace(body.DeviceID)
	if len(deviceID) < 8 || len(deviceID) > 64 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的设备 ID"})
		return
	}

	db.DB.Exec("INSERT OR IGNORE INTO photo_likes (photo_id, device_id) VALUES (?, ?)", photoID, deviceID)

	var count int64
	db.DB.Model(&models.PhotoLike{}).Where("photo_id = ?", photoID).Count(&count)

	c.JSON(http.StatusOK, gin.H{"count": count, "liked": true})
}

// GetPhotoDislikes GET /api/photos/:id/dislikes
func GetPhotoDislikes(c *gin.Context) {
	photoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || photoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的照片 ID"})
		return
	}

	deviceID := c.Query("deviceId")

	var count int64
	db.DB.Model(&models.PhotoDislike{}).Where("photo_id = ?", photoID).Count(&count)

	disliked := false
	if deviceID != "" {
		var dislikeCount int64
		db.DB.Model(&models.PhotoDislike{}).Where("photo_id = ? AND device_id = ?", photoID, deviceID).Count(&dislikeCount)
		disliked = dislikeCount > 0
	}

	c.JSON(http.StatusOK, gin.H{"count": count, "disliked": disliked})
}

// PostPhotoDislike POST /api/photos/:id/dislikes
func PostPhotoDislike(c *gin.Context) {
	photoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || photoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的照片 ID"})
		return
	}

	var body struct {
		DeviceID string `json:"deviceId"`
		Action   string `json:"action"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	deviceID := strings.TrimSpace(body.DeviceID)
	if len(deviceID) < 8 || len(deviceID) > 64 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的设备 ID"})
		return
	}

	action := body.Action
	if action == "" {
		action = "dislike"
	}

	if action == "dislike" {
		db.DB.Exec("INSERT OR IGNORE INTO photo_dislikes (photo_id, device_id) VALUES (?, ?)", photoID, deviceID)
	} else {
		db.DB.Where("photo_id = ? AND device_id = ?", photoID, deviceID).Delete(&models.PhotoDislike{})
	}

	var count int64
	db.DB.Model(&models.PhotoDislike{}).Where("photo_id = ?", photoID).Count(&count)

	var dislikeCount int64
	db.DB.Model(&models.PhotoDislike{}).Where("photo_id = ? AND device_id = ?", photoID, deviceID).Count(&dislikeCount)

	c.JSON(http.StatusOK, gin.H{"count": count, "disliked": dislikeCount > 0})
}

// VerifyPhotoPassword POST /api/photos/:id/verify-password
func VerifyPhotoPassword(c *gin.Context) {
	photoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || photoID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的照片 ID"})
		return
	}

	var body struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "请输入密码"})
		return
	}

	var p models.Photo
	if err := db.DB.Select("password_hash").Where("id = ?", photoID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "照片不存在或未设置密码"})
		return
	}

	if p.PasswordHash == nil || *p.PasswordHash == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "照片不存在或未设置密码"})
		return
	}

	success := bcrypt.CompareHashAndPassword([]byte(*p.PasswordHash), []byte(body.Password)) == nil
	c.JSON(http.StatusOK, gin.H{"success": success})
}