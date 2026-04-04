package main

import (
	"fmt"
	"log"

	"github.com/fatWill/agent-team-blog/backend/config"
	"github.com/fatWill/agent-team-blog/backend/handlers"
	"github.com/fatWill/agent-team-blog/backend/middleware"
	"github.com/fatWill/agent-team-blog/backend/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化 MySQL
	if err := utils.InitDB(&cfg.DB); err != nil {
		log.Fatalf("初始化 MySQL 失败: %v", err)
	}
	log.Println("✅ MySQL 连接成功")

	// 初始化 Redis
	if err := utils.InitRedis(&cfg.Redis); err != nil {
		log.Fatalf("初始化 Redis 失败: %v", err)
	}
	log.Println("✅ Redis 连接成功")

	// 设置上传配置
	handlers.SetUploadConfig(&cfg.Upload)

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
	auth := middleware.AuthRequired()

	// ========== 认证 ==========
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", handlers.Login)
		authGroup.POST("/logout", handlers.Logout)
		authGroup.GET("/check", auth, handlers.AuthCheck)
	}

	// ========== 文章 ==========
	articlesGroup := api.Group("/articles")
	{
		// 公开接口
		articlesGroup.GET("", handlers.GetArticles)
		articlesGroup.GET("/like-status-batch", handlers.GetArticleLikeStatusBatch)
		articlesGroup.GET("/:id", handlers.GetArticle)
		articlesGroup.GET("/:id/like-status", handlers.GetArticleLikeStatus)
		articlesGroup.POST("/:id/like", handlers.LikeArticle)

		// 需鉴权接口
		articlesGroup.POST("", auth, handlers.CreateArticle)
		articlesGroup.PUT("/:id", auth, handlers.UpdateArticle)
		articlesGroup.DELETE("/:id", auth, handlers.DeleteArticle)
	}

	// ========== 相册 ==========
	albumsGroup := api.Group("/albums")
	{
		// 公开接口
		albumsGroup.GET("", handlers.GetAlbums)
		albumsGroup.GET("/:id/photos", handlers.GetAlbumPhotos)
		albumsGroup.POST("/:id/verify-password", handlers.VerifyAlbumPassword)

		// 需鉴权接口
		albumsGroup.POST("", auth, handlers.CreateAlbum)
		albumsGroup.PUT("/:id", auth, handlers.UpdateAlbum)
		albumsGroup.DELETE("/:id", auth, handlers.DeleteAlbum)
		albumsGroup.POST("/:id/photos", auth, handlers.AddAlbumPhoto)
	}

	// ========== 照片 ==========
	photosGroup := api.Group("/photos")
	{
		// 公开接口
		photosGroup.GET("/:id/likes", handlers.GetPhotoLikes)
		photosGroup.POST("/:id/likes", handlers.PostPhotoLike)
		photosGroup.GET("/:id/dislikes", handlers.GetPhotoDislikes)
		photosGroup.POST("/:id/dislikes", handlers.PostPhotoDislike)
		photosGroup.POST("/:id/verify-password", handlers.VerifyPhotoPassword)

		// 需鉴权接口
		photosGroup.PUT("/:id", auth, handlers.UpdatePhoto)
		photosGroup.DELETE("/:id", auth, handlers.DeletePhoto)
	}

	// ========== 上传 ==========
	uploadGroup := api.Group("/upload")
	uploadGroup.Use(auth)
	{
		uploadGroup.POST("", handlers.Upload)
		uploadGroup.POST("/chunk", handlers.UploadChunk)
		uploadGroup.POST("/merge", handlers.MergeChunks)
		uploadGroup.DELETE("/chunk", handlers.CancelChunkUpload)
	}

	// ========== 个人资料 ==========
	api.GET("/profile", handlers.GetProfile)
	api.PUT("/profile", auth, handlers.UpdateProfile)

	// ========== 留言板 ==========
	api.GET("/messages", handlers.GetMessages)
	api.POST("/messages", middleware.RateLimit(3), handlers.CreateMessage)
	api.PUT("/messages/:id", handlers.UpdateMessage)

	// ========== 主题 ==========
	api.GET("/theme", handlers.GetTheme)
	api.POST("/theme", handlers.SaveTheme)

	// ========== 更新日志 ==========
	api.GET("/changelog", handlers.GetChangelog)
}
