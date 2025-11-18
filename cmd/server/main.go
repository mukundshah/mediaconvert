package main

import (
	"context"
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mukund/mediaconvert/internal/analytics"
	"github.com/mukund/mediaconvert/internal/auth"
	"github.com/mukund/mediaconvert/internal/config"
	"github.com/mukund/mediaconvert/internal/db"
	"github.com/mukund/mediaconvert/internal/handlers"
	"github.com/mukund/mediaconvert/internal/s3compat"
	"github.com/mukund/mediaconvert/internal/system"
	"github.com/mukund/mediaconvert/internal/worker"
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

	// Initialize MinIO Client
	s3Url, err := url.Parse(cfg.S3Endpoint)
	if err != nil {
		log.Fatalf("Failed to parse S3 endpoint: %v", err)
	}

	useSSL := s3Url.Scheme == "https"
	endpoint := s3Url.Host

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: useSSL,
		Region: cfg.S3Region,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, cfg.S3Bucket)
	if err != nil {
		log.Fatalf("Failed to check if bucket exists: %v", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.S3Bucket, minio.MakeBucketOptions{Region: cfg.S3Region})
		if err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
		log.Printf("Created bucket: %s", cfg.S3Bucket)
	}

	// Connect to Redis
	redisClient, err := worker.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Connect to ClickHouse (optional - continue if it fails)
	var analyticsClient *analytics.Client
	if cfg.ClickHouseDSN != "" {
		analyticsClient, err = analytics.NewClient(cfg.ClickHouseDSN)
		if err != nil {
			log.Printf("Warning: Failed to connect to ClickHouse: %v (analytics disabled)", err)
		} else {
			defer analyticsClient.Close()
			// Initialize schema
			if err := analyticsClient.InitSchema(ctx); err != nil {
				log.Printf("Warning: Failed to initialize ClickHouse schema: %v", err)
			}
		}
	}

	// Setup Handlers
	authHandler := handlers.NewAuthHandler(database)
	jobHandler := handlers.NewJobHandler(database)
	pipelineHandler := handlers.NewPipelineHandler(database)
	s3CredentialHandler := handlers.NewS3CredentialHandler(database)
	s3Handler := s3compat.NewS3Handler(database, minioClient, cfg, redisClient)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsClient)

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
		protected.GET("/s3-credentials/check-availability", s3CredentialHandler.CheckBucketAvailability)
		protected.DELETE("/s3-credentials/:id", s3CredentialHandler.RevokeCredentials)

		// Analytics routes
		protected.GET("/analytics/jobs/stats", analyticsHandler.GetJobStats)
		protected.GET("/analytics/jobs/timeline", analyticsHandler.GetJobTimeline)
		protected.GET("/analytics/pipelines/stats", analyticsHandler.GetPipelineStats)
	}

	// S3-Compatible API routes (separate from /api)
	s3Routes := r.Group("")
	s3Routes.Use(s3compat.S3AuthMiddleware(database))
	{
		// Object operations
		s3Routes.PUT("/:bucket/*key", s3Handler.PutObject)
		s3Routes.POST("/:bucket/*key", s3Handler.PutObject) // AWS CLI uses POST for multipart
		s3Routes.GET("/:bucket/*key", s3Handler.GetObject)
		s3Routes.HEAD("/:bucket/*key", s3Handler.HeadObject)
		s3Routes.DELETE("/:bucket/*key", s3Handler.DeleteObject)

		// List objects
		s3Routes.GET("/:bucket", s3Handler.ListObjects)
	}

	log.Printf("Starting server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
