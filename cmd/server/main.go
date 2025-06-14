package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mukund/mediaconvert/internal/auth"
	"github.com/mukund/mediaconvert/internal/config"
	"github.com/mukund/mediaconvert/internal/db"
	"github.com/mukund/mediaconvert/internal/handlers"
	"github.com/mukund/mediaconvert/internal/system"
)

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Check System Dependencies
	if err := system.CheckDependencies(); err != nil {
		log.Fatalf("System dependency check failed: %v", err)
	}

	// Connect to DB
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run Migrations
	if err := db.Migrate(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Auth
	auth.InitAuth(cfg.JWTSecret)

	// Setup Handlers
	authHandler := handlers.NewAuthHandler(database)
	jobHandler := handlers.NewJobHandler(database)
	pipelineHandler := handlers.NewPipelineHandler(database)
	s3CredentialHandler := handlers.NewS3CredentialHandler(database)

	// Setup Router
	r := gin.Default()

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Auth routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/me", func(c *gin.Context) {
			userID, _ := auth.GetUserID(c)
			email, _ := c.Get("user_email")
			c.JSON(200, gin.H{
				"user_id": userID,
				"email":   email,
			})
		})

		// Job routes
		protected.GET("/jobs", jobHandler.ListJobs)
		protected.GET("/jobs/:id", jobHandler.GetJob)
		protected.POST("/jobs/:id/cancel", jobHandler.CancelJob)
		protected.POST("/jobs/:id/rerun", jobHandler.RerunJob)

		// Pipeline routes
		protected.POST("/pipelines", pipelineHandler.CreatePipeline)
		protected.GET("/pipelines", pipelineHandler.ListPipelines)
		protected.GET("/pipelines/:id", pipelineHandler.GetPipeline)
		protected.PUT("/pipelines/:id", pipelineHandler.UpdatePipeline)
		protected.DELETE("/pipelines/:id", pipelineHandler.DeletePipeline)

		// S3 Credential routes
		protected.POST("/s3-credentials", s3CredentialHandler.CreateCredentials)
		protected.GET("/s3-credentials", s3CredentialHandler.ListCredentials)
		protected.DELETE("/s3-credentials/:id", s3CredentialHandler.RevokeCredentials)
	}

	log.Printf("Starting server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
