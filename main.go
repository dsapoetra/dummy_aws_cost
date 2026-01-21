package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"dummy_aws_cost/internal/database"
	"dummy_aws_cost/internal/handlers"
	"dummy_aws_cost/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed all:frontend/dist
var frontendFS embed.FS

func main() {
	if err := database.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.Close()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins != "" {
		corsConfig.AllowOrigins = []string{allowedOrigins}
	} else {
		corsConfig.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	api := r.Group("/api")
	{
		api.POST("/auth/login", handlers.Login)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/auth/me", handlers.GetCurrentUser)

			protected.GET("/articles", handlers.GetArticles)
			protected.GET("/articles/:id", handlers.GetArticle)
			protected.POST("/articles", handlers.CreateArticle)
			protected.PUT("/articles/:id", handlers.UpdateArticle)
			protected.DELETE("/articles/:id", handlers.DeleteArticle)

			protected.GET("/pages", handlers.GetPages)
			protected.GET("/pages/:id", handlers.GetPage)
			protected.POST("/pages", handlers.CreatePage)
			protected.PUT("/pages/:id", handlers.UpdatePage)
			protected.DELETE("/pages/:id", handlers.DeletePage)

			protected.GET("/media", handlers.GetMedia)
			protected.POST("/media", handlers.UploadMedia)
			protected.DELETE("/media/:id", handlers.DeleteMedia)
		}
	}

	r.Static("/uploads", handlers.GetUploadDir())

	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err == nil {
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if _, err := fs.Stat(distFS, path[1:]); err == nil {
				c.FileFromFS(path, http.FS(distFS))
				return
			}
			indexFile, _ := fs.ReadFile(distFS, "index.html")
			c.Data(http.StatusOK, "text/html", indexFile)
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
