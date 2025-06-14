package worker

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const JobNotificationChannel = "job:notifications"

// RedisClient wraps redis client for job notifications
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(addr string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{client: client}, nil
}

// PublishJobNotification publishes a job notification to Redis
func (r *RedisClient) PublishJobNotification(jobID uint) error {
	ctx := context.Background()
	return r.client.Publish(ctx, JobNotificationChannel, fmt.Sprintf("%d", jobID)).Err()
}

// SubscribeToJobNotifications subscribes to job notifications
func (r *RedisClient) SubscribeToJobNotifications(ctx context.Context) *redis.PubSub {
	return r.client.Subscribe(ctx, JobNotificationChannel)
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}
