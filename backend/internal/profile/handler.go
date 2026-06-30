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
		c.JSON(http.StatusOK, gin.H{"name": "fatwill", "avatar": "", "bio": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": p.Name, "avatar": p.Avatar, "bio": p.Bio})
}

// UpdateProfile PUT /api/profile
func UpdateProfile(c *gin.Context) {
	var body struct {
		Name   *string `json:"name"`
		Avatar *string `json:"avatar"`
		Bio    *string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	if body.Name == nil && body.Avatar == nil && body.Bio == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "至少需要提供 name、avatar 或 bio 中的一个字段"})
		return
	}

	updates := map[string]interface{}{}
	if body.Name != nil {
		name := strings.TrimSpace(*body.Name)
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "昵称不能为空"})
			return
		}
		if len([]rune(name)) > 50 {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "昵称长度不能超过 50 个字符"})
			return
		}
		updates["name"] = name
	}
	if body.Avatar != nil {
		updates["avatar"] = strings.TrimSpace(*body.Avatar)
	}
	if body.Bio != nil {
		updates["bio"] = strings.TrimSpace(*body.Bio)
	}

	db.DB.Model(&models.Profile{}).Where("id = 1").Updates(updates)

	var p models.Profile
	db.DB.Where("id = 1").First(&p)

	c.JSON(http.StatusOK, gin.H{"name": p.Name, "avatar": p.Avatar, "bio": p.Bio})
}
