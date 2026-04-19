package main

import (
	"fmt"
	"log"

	"github.com/fatWill/agent-team-blog/backend/config"
	"github.com/fatWill/agent-team-blog/backend/internal/album"
	"github.com/fatWill/agent-team-blog/backend/internal/article"
	"github.com/fatWill/agent-team-blog/backend/internal/auth"
	"github.com/fatWill/agent-team-blog/backend/internal/changelog"
	"github.com/fatWill/agent-team-blog/backend/internal/download"
	"github.com/fatWill/agent-team-blog/backend/internal/guestbook"
	"github.com/fatWill/agent-team-blog/backend/internal/material"
	"github.com/fatWill/agent-team-blog/backend/internal/perf"
	"github.com/fatWill/agent-team-blog/backend/internal/photo"
	"github.com/fatWill/agent-team-blog/backend/internal/profile"
	"github.com/fatWill/agent-team-blog/backend/internal/pv"
	"github.com/fatWill/agent-team-blog/backend/internal/renovation"
	"github.com/fatWill/agent-team-blog/backend/internal/theme"
	"github.com/fatWill/agent-team-blog/backend/internal/upload"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/fatWill/agent-team-blog/backend/pkg/ipgeo"
	"github.com/fatWill/agent-team-blog/backend/pkg/middleware"
	"github.com/fatWill/agent-team-blog/backend/pkg/rds"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化 SQLite
	if err := db.Init(&cfg.DB); err != nil {
		log.Fatalf("初始化 SQLite 失败: %v", err)
	}
	log.Printf("✅ SQLite 连接成功 (%s)", cfg.DB.Path)

	// 初始化 Redis
	if err := rds.Init(&cfg.Redis); err != nil {
		log.Fatalf("初始化 Redis 失败: %v", err)
	}
	log.Println("✅ Redis 连接成功")

	// 设置上传配置
	upload.SetUploadConfig(&cfg.Upload)
	upload.SetCOSConfig(&cfg.COS)

	// 初始化 IP 地理位置解析（非关键服务，失败不阻塞启动）
	ip2regionPath := cfg.Server.IP2RegionPath
	if ip2regionPath == "" {
		ip2regionPath = "data/ip2region.xdb"
	}
	if err := ipgeo.Init(ip2regionPath); err != nil {
		log.Printf("⚠️ ip2region 初始化失败（地理位置功能不可用）: %v", err)
	}

	// 创建 Gin 引擎
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.Server.CORSOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Set-Cookie"},
		AllowCredentials: true,
	}))

	// 路由注册
	api := r.Group("/api")
	registerRoutes(api)

	// 启动服务
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("🚀 Go 后端服务启动在 %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

func registerRoutes(api *gin.RouterGroup) {
	// 需要鉴权的中间件
	authMW := middleware.AuthRequired()

	// ========== 认证 ==========
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", auth.Login)
		authGroup.POST("/logout", auth.Logout)
		authGroup.GET("/check", authMW, auth.AuthCheck)
	}

	// ========== 文章 ==========
	articlesGroup := api.Group("/articles")
	{
		// 公开接口（注意：/random 和 /like-status-batch 必须在 /:id 之前注册）
		articlesGroup.GET("", article.GetArticles)
		articlesGroup.GET("/random", article.GetRandomArticle)
		articlesGroup.GET("/like-status-batch", article.GetArticleLikeStatusBatch)
		articlesGroup.GET("/:id", article.GetArticle)
		articlesGroup.GET("/:id/like-status", article.GetArticleLikeStatus)
		articlesGroup.POST("/:id/like", article.LikeArticle)
		articlesGroup.POST("/:id/view", article.ViewArticle)

		// 需鉴权接口
		articlesGroup.POST("", authMW, article.CreateArticle)
		articlesGroup.PUT("/:id", authMW, article.UpdateArticle)
		articlesGroup.DELETE("/:id", authMW, article.DeleteArticle)
	}

	// ========== 相册 ==========
	albumsGroup := api.Group("/albums")
	{
		// 公开接口
		albumsGroup.GET("", album.GetAlbums)
		albumsGroup.GET("/:id/photos", album.GetAlbumPhotos)
		albumsGroup.POST("/:id/verify-password", album.VerifyAlbumPassword)

		// 需鉴权接口
		albumsGroup.POST("", authMW, album.CreateAlbum)
		albumsGroup.PUT("/:id", authMW, album.UpdateAlbum)
		albumsGroup.DELETE("/:id", authMW, album.DeleteAlbum)
		albumsGroup.POST("/:id/photos", authMW, album.AddAlbumPhoto)
	}

	// ========== 照片 ==========
	photosGroup := api.Group("/photos")
	{
		// 公开接口
		photosGroup.GET("/:id/likes", photo.GetPhotoLikes)
		photosGroup.POST("/:id/likes", photo.PostPhotoLike)
		photosGroup.GET("/:id/dislikes", photo.GetPhotoDislikes)
		photosGroup.POST("/:id/dislikes", photo.PostPhotoDislike)
		photosGroup.POST("/:id/verify-password", photo.VerifyPhotoPassword)

		// 需鉴权接口
		photosGroup.PUT("/:id", authMW, photo.UpdatePhoto)
		photosGroup.DELETE("/:id", authMW, photo.DeletePhoto)
	}

	// ========== 上传 ==========
	uploadGroup := api.Group("/upload")
	uploadGroup.Use(authMW)
	{
		uploadGroup.POST("", upload.Upload)
		uploadGroup.POST("/chunk", upload.UploadChunk)
		uploadGroup.POST("/merge", upload.MergeChunks)
		uploadGroup.DELETE("/chunk", upload.CancelChunkUpload)
	}

	// ========== 个人资料 ==========
	api.GET("/profile", profile.GetProfile)
	api.PUT("/profile", authMW, profile.UpdateProfile)

	// ========== 留言板 ==========
	api.GET("/messages", guestbook.GetMessages)
	api.POST("/messages", middleware.RateLimit(3), guestbook.CreateMessage)
	api.DELETE("/messages/:id", authMW, guestbook.DeleteMessage)

	// ========== 主题 ==========
	api.GET("/theme", theme.GetTheme)
	api.POST("/theme", theme.SaveTheme)

	// ========== PV/UV 统计 ==========
	pvGroup := api.Group("/pv")
	{
		pvGroup.POST("/record", pv.RecordPV)
		pvGroup.GET("/trend", authMW, pv.GetTrend)
		pvGroup.GET("/top-pages", authMW, pv.GetTopPages)
		pvGroup.GET("/logs", authMW, pv.GetLogs)
		pvGroup.GET("/overview", authMW, pv.GetOverview)
		pvGroup.GET("/geo", authMW, pv.GetGeoDistribution)
	}

	// ========== 前端性能监控 ==========
	perfGroup := api.Group("/perf")
	{
		perfGroup.POST("/report", perf.Report)
		perfGroup.GET("/overview", authMW, perf.GetOverview)
		perfGroup.GET("/pages", authMW, perf.GetPages)
		perfGroup.GET("/trend", authMW, perf.GetTrend)
		perfGroup.GET("/logs", authMW, perf.GetLogs)
	}

	// ========== 装修文章 ==========
	renoGroup := api.Group("/renovation/articles")
	{
		// 公开接口
		renoGroup.GET("", renovation.GetArticles)
		renoGroup.GET("/:id", renovation.GetArticle)

		// 需鉴权接口
		renoGroup.POST("", authMW, renovation.CreateArticle)
		renoGroup.PUT("/:id", authMW, renovation.UpdateArticle)
		renoGroup.DELETE("/:id", authMW, renovation.DeleteArticle)
	}

	// ========== 材料清单 ==========
	materialGroup := api.Group("/materials")
	{
		// 需鉴权接口（查看也需要登录）
		materialGroup.GET("", authMW, material.GetMaterials)
		materialGroup.GET("/:id", authMW, material.GetMaterial)
		materialGroup.POST("", authMW, material.CreateMaterial)
		materialGroup.PUT("/:id", authMW, material.UpdateMaterial)
		materialGroup.DELETE("/:id", authMW, material.DeleteMaterial)
		materialGroup.PUT("/:id/sort", authMW, material.UpdateMaterialSort)
	}

	// ========== 更新日志 ==========
	api.GET("/changelog", changelog.GetChangelog)

	// ========== 文件下载代理 ==========
	api.GET("/download", download.Proxy)
}