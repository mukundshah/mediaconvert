package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mukund/mediaconvert/internal/analytics"
	"github.com/mukund/mediaconvert/internal/config"
	"github.com/mukund/mediaconvert/internal/db"
	"github.com/mukund/mediaconvert/internal/worker"
)

func main() {
	log.Println("Starting media processing worker...")

	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to DB
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

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
			if err := analyticsClient.InitSchema(context.Background()); err != nil {
				log.Printf("Warning: Failed to initialize ClickHouse schema: %v", err)
			}
		}
	}

	// Create job processor
	processor := worker.NewJobProcessor(database, minioClient, cfg, redisClient, analyticsClient)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutdown signal received, stopping worker...")
		cancel()
	}()

	// Subscribe to job notifications
	pubsub := redisClient.SubscribeToJobNotifications(ctx)
	defer pubsub.Close()

	log.Println("Worker ready, listening for job notifications...")

	// Listen for notifications
	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			// Parse job ID from message
			jobID, err := strconv.ParseUint(msg.Payload, 10, 32)
			if err != nil {
				log.Printf("Invalid job ID in notification: %s", msg.Payload)
				continue
			}

			log.Printf("Received job notification: %d", jobID)

			// Process job
			if err := processor.ProcessJob(uint(jobID)); err != nil {
				log.Printf("Failed to process job %d: %v", jobID, err)
			}

		case <-ctx.Done():
			log.Println("Worker stopped")
			return
		}
	}
}
