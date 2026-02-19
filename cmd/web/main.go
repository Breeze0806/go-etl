package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Breeze0806/go-etl/cmd/web/config"
	"github.com/Breeze0806/go-etl/cmd/web/database"
	"github.com/Breeze0806/go-etl/cmd/web/handlers"
	"github.com/Breeze0806/go-etl/cmd/web/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.New(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	authHandler := handlers.NewAuthHandler(db.DB, jwtMiddleware)
	userHandler := handlers.NewUserHandler(db.DB, jwtMiddleware)
	datasourceHandler := handlers.NewDataSourceHandler(db.DB, jwtMiddleware)
	syncTaskHandler := handlers.NewSyncTaskHandler(db.DB, jwtMiddleware)
	syncJobHandler := handlers.NewSyncJobHandler(db.DB, jwtMiddleware)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	protected := router.Group("/api")
	protected.Use(jwtMiddleware.AuthMiddleware())
	{
		protected.GET("/me", authHandler.Me)

		admin := protected.Group("")
		admin.Use(jwtMiddleware.RequireRole("admin"))
		{
			admin.GET("/users", userHandler.List)
			admin.POST("/users", userHandler.Create)
			admin.DELETE("/users/:id", userHandler.Delete)
		}

		protected.GET("/users/:id", userHandler.GetByID)
		protected.PUT("/users/:id", userHandler.Update)

		protected.GET("/datasources", datasourceHandler.List)
		protected.POST("/datasources", datasourceHandler.Create)
		protected.GET("/datasources/:id", datasourceHandler.GetByID)
		protected.PUT("/datasources/:id", datasourceHandler.Update)
		protected.DELETE("/datasources/:id", datasourceHandler.Delete)
		protected.POST("/datasources/:id/test", datasourceHandler.TestConnection)

		protected.GET("/synctasks", syncTaskHandler.List)
		protected.POST("/synctasks", syncTaskHandler.Create)
		protected.GET("/synctasks/:id", syncTaskHandler.GetByID)
		protected.PUT("/synctasks/:id", syncTaskHandler.Update)
		protected.DELETE("/synctasks/:id", syncTaskHandler.Delete)
		protected.POST("/synctasks/:id/start", syncJobHandler.Start)

		protected.GET("/syncjobs", syncJobHandler.List)
		protected.GET("/syncjobs/:id", syncJobHandler.GetByID)
		protected.DELETE("/syncjobs/:id", syncJobHandler.Stop)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
