package profile

import (
	"net/http"
	"strings"

	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/gin-gonic/gin"
)

// GetProfile GET /api/profile
func GetProfile(c *gin.Context) {
	var p models.Profile
	if err := db.DB.Where("id = 1").First(&p).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"avatar": "", "bio": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatar": p.Avatar, "bio": p.Bio})
}

// UpdateProfile PUT /api/profile
func UpdateProfile(c *gin.Context) {
	var body struct {
		Avatar *string `json:"avatar"`
		Bio    *string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	if body.Avatar == nil && body.Bio == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "至少需要提供 avatar 或 bio 中的一个字段"})
		return
	}

	updates := map[string]interface{}{}
	if body.Avatar != nil {
		updates["avatar"] = strings.TrimSpace(*body.Avatar)
	}
	if body.Bio != nil {
		updates["bio"] = strings.TrimSpace(*body.Bio)
	}

	db.DB.Model(&models.Profile{}).Where("id = 1").Updates(updates)

	var p models.Profile
	db.DB.Where("id = 1").First(&p)

	c.JSON(http.StatusOK, gin.H{"avatar": p.Avatar, "bio": p.Bio})
}
