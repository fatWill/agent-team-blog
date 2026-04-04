package album

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/fatWill/agent-team-blog/backend/internal/upload"
	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetAlbums GET /api/albums
func GetAlbums(c *gin.Context) {
	type albumRow struct {
		ID           uint64  `gorm:"column:id"`
		Name         string  `gorm:"column:name"`
		Description  *string `gorm:"column:description"`
		CoverURL     *string `gorm:"column:cover_url"`
		PasswordHash *string `gorm:"column:password_hash"`
		PhotoCount   int     `gorm:"column:photo_count"`
		CreatedAt    string  `gorm:"column:created_at"`
		UpdatedAt    string  `gorm:"column:updated_at"`
	}

	var rows []albumRow
	err := db.DB.Raw(`
		SELECT
			a.id, a.name, a.description, a.cover_url, a.password_hash,
			a.created_at, a.updated_at,
			IFNULL(pc.cnt, 0) AS photo_count
		FROM albums a
		LEFT JOIN (
			SELECT album_id, COUNT(*) AS cnt
			FROM photos
			GROUP BY album_id
		) pc ON pc.album_id = a.id
		ORDER BY a.created_at DESC
	`).Scan(&rows).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "查询相册列表失败"})
		return
	}

	list := make([]models.AlbumListItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, models.AlbumListItem{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			CoverURL:    r.CoverURL,
			PhotoCount:  r.PhotoCount,
			HasPassword: r.PasswordHash != nil && *r.PasswordHash != "",
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"list": list})
}

// CreateAlbum POST /api/albums
func CreateAlbum(c *gin.Context) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Password    string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	if strings.TrimSpace(body.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "相册名称不能为空"})
		return
	}

	a := models.Album{
		Name: strings.TrimSpace(body.Name),
	}

	if body.Description != "" {
		desc := strings.TrimSpace(body.Description)
		a.Description = &desc
	}

	if strings.TrimSpace(body.Password) != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(body.Password)), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "密码加密失败"})
			return
		}
		hashStr := string(hash)
		a.PasswordHash = &hashStr
	}

	if err := db.DB.Create(&a).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "创建相册失败"})
		return
	}

	c.JSON(http.StatusOK, models.AlbumListItem{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		CoverURL:    a.CoverURL,
		PhotoCount:  0,
		HasPassword: a.PasswordHash != nil,
		CreatedAt:   a.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt:   a.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	})
}

// UpdateAlbum PUT /api/albums/:id
func UpdateAlbum(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的相册 ID"})
		return
	}

	var body struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Password    *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	if body.Name != nil && strings.TrimSpace(*body.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "相册名称不能为空"})
		return
	}

	updates := map[string]interface{}{}
	if body.Name != nil {
		updates["name"] = strings.TrimSpace(*body.Name)
	}
	if body.Description != nil {
		updates["description"] = strings.TrimSpace(*body.Description)
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

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "没有需要更新的字段"})
		return
	}

	result := db.DB.Model(&models.Album{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "更新相册失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "相册不存在"})
		return
	}

	var alb models.Album
	db.DB.Where("id = ?", id).First(&alb)

	var photoCount int64
	db.DB.Model(&models.Photo{}).Where("album_id = ?", id).Count(&photoCount)

	c.JSON(http.StatusOK, models.AlbumListItem{
		ID:          alb.ID,
		Name:        alb.Name,
		Description: alb.Description,
		CoverURL:    alb.CoverURL,
		PhotoCount:  int(photoCount),
		HasPassword: alb.PasswordHash != nil && *alb.PasswordHash != "",
		CreatedAt:   alb.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt:   alb.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	})
}

// DeleteAlbum DELETE /api/albums/:id
func DeleteAlbum(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的相册 ID"})
		return
	}

	// 先收集所有照片 URL，用于后续 COS 清理
	var photoURLs []string
	db.DB.Model(&models.Photo{}).Where("album_id = ?", id).Pluck("url", &photoURLs)

	// 删除关联照片
	db.DB.Where("album_id = ?", id).Delete(&models.Photo{})

	result := db.DB.Where("id = ?", id).Delete(&models.Album{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "删除相册失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "相册不存在"})
		return
	}

	// 异步批量删除 COS 对象（不阻塞响应）
	if len(photoURLs) > 0 {
		go func() {
			upload.BatchDeleteFromCOS(photoURLs)
			log.Printf("✅ 相册 %d 的 %d 张照片 COS 对象已清理", id, len(photoURLs))
		}()
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GetAlbumPhotos GET /api/albums/:id/photos
func GetAlbumPhotos(c *gin.Context) {
	albumID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || albumID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的相册 ID"})
		return
	}

	var albumCount int64
	db.DB.Model(&models.Album{}).Where("id = ?", albumID).Count(&albumCount)
	if albumCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "相册不存在"})
		return
	}

	var photos []models.Photo
	db.DB.Where("album_id = ?", albumID).Order("created_at DESC").Find(&photos)

	list := make([]models.PhotoListItem, 0, len(photos))
	for _, p := range photos {
		list = append(list, models.PhotoListItem{
			ID:          p.ID,
			AlbumID:     p.AlbumID,
			URL:         p.URL,
			Caption:     p.Caption,
			HasPassword: p.PasswordHash != nil && *p.PasswordHash != "",
			CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
			UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"list": list})
}

// AddAlbumPhoto POST /api/albums/:id/photos
func AddAlbumPhoto(c *gin.Context) {
	albumID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || albumID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的相册 ID"})
		return
	}

	var albumCount int64
	db.DB.Model(&models.Album{}).Where("id = ?", albumID).Count(&albumCount)
	if albumCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "相册不存在"})
		return
	}

	var body struct {
		URL      string `json:"url"`
		Caption  string `json:"caption"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	if strings.TrimSpace(body.URL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "图片 URL 不能为空"})
		return
	}

	photo := models.Photo{
		AlbumID: albumID,
		URL:     strings.TrimSpace(body.URL),
	}

	if body.Caption != "" {
		caption := strings.TrimSpace(body.Caption)
		photo.Caption = &caption
	}

	if strings.TrimSpace(body.Password) != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(body.Password)), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "密码加密失败"})
			return
		}
		hashStr := string(hash)
		photo.PasswordHash = &hashStr
	}

	if err := db.DB.Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "添加照片失败"})
		return
	}

	// 更新相册封面
	UpdateAlbumCover(albumID)

	c.JSON(http.StatusOK, models.PhotoListItem{
		ID:          photo.ID,
		AlbumID:     photo.AlbumID,
		URL:         photo.URL,
		Caption:     photo.Caption,
		HasPassword: photo.PasswordHash != nil,
		CreatedAt:   photo.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt:   photo.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	})
}

// VerifyAlbumPassword POST /api/albums/:id/verify-password
func VerifyAlbumPassword(c *gin.Context) {
	albumID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || albumID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的相册 ID"})
		return
	}

	var body struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "请输入密码"})
		return
	}

	var alb models.Album
	if err := db.DB.Select("password_hash").Where("id = ?", albumID).First(&alb).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "相册不存在或未设置密码"})
		return
	}

	if alb.PasswordHash == nil || *alb.PasswordHash == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "相册不存在或未设置密码"})
		return
	}

	success := bcrypt.CompareHashAndPassword([]byte(*alb.PasswordHash), []byte(body.Password)) == nil
	c.JSON(http.StatusOK, gin.H{"success": success})
}

// UpdateAlbumCover 更新相册封面（取最新照片），导出供 photo handler 调用
func UpdateAlbumCover(albumID uint64) {
	var photo models.Photo
	err := db.DB.Select("url").Where("album_id = ?", albumID).Order("created_at DESC").First(&photo).Error

	if err != nil {
		db.DB.Model(&models.Album{}).Where("id = ?", albumID).Update("cover_url", nil)
	} else {
		db.DB.Model(&models.Album{}).Where("id = ?", albumID).Update("cover_url", photo.URL)
	}
}
