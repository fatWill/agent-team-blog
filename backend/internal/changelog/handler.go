package changelog

import (
	"net/http"

	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/gin-gonic/gin"
)

// GetChangelog GET /api/changelog
func GetChangelog(c *gin.Context) {
	var changelogs []models.Changelog
	db.DB.Order("id DESC").Find(&changelogs)

	c.JSON(http.StatusOK, gin.H{"changelog": changelogs})
}
