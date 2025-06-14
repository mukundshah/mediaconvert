package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	// Initialize S3 Client
	s3Client := s3.NewFromConfig(aws.Config{
		Region: cfg.S3Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.S3AccessKey,
			cfg.S3SecretKey,
			"",
		),
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               cfg.S3Endpoint,
					HostnameImmutable: true,
					Source:            aws.EndpointSourceCustom,
				}, nil
			},
		),
	}, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	// Connect to Redis
	redisClient, err := worker.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Create job processor
	processor := worker.NewJobProcessor(database, s3Client, cfg, redisClient)

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
