package config

import (
	"os"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL   string `mapstructure:"DATABASE_URL"`
	RedisURL      string `mapstructure:"REDIS_URL"`
	ClickHouseDSN string `mapstructure:"CLICKHOUSE_DSN"`
	Port          string `mapstructure:"PORT"`
	S3Endpoint    string `mapstructure:"S3_ENDPOINT"`
	S3AccessKey   string `mapstructure:"S3_ACCESS_KEY"`
	S3SecretKey   string `mapstructure:"S3_SECRET_KEY"`
	S3Bucket      string `mapstructure:"S3_BUCKET"`
	S3Region      string `mapstructure:"S3_REGION"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/mediaconvert?sslmode=disable")
	viper.SetDefault("REDIS_URL", "redis://localhost:6379")
	viper.SetDefault("CLICKHOUSE_DSN", "clickhouse://default@localhost:9000/default")
	viper.SetDefault("S3_ENDPOINT", "http://localhost:9000")
	viper.SetDefault("S3_ACCESS_KEY", "minioadmin")
	viper.SetDefault("S3_SECRET_KEY", "minioadmin")
	viper.SetDefault("S3_BUCKET", "media")
	viper.SetDefault("S3_REGION", "us-east-1")
	viper.SetDefault("JWT_SECRET", "change-this-secret-in-production")

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
