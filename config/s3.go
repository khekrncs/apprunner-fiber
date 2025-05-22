// config/s3.go
package config

import (
	"os"
)

type S3Config struct {
	Region              string
	Bucket              string
	MaxRetries          int
	EnablePresignedURLs bool
}

func GetS3Config() *S3Config {
	return &S3Config{
		Region:              getEnvOrDefault("AWS_REGION", "ap-south-1"),
		Bucket:              getEnvOrDefault("AWS_S3_BUCKET", "unorganizedbucket-dev-971709774307-ap-south-1"),
		MaxRetries:          getEnvAsInt("AWS_MAX_RETRIES", 3),
		EnablePresignedURLs: getEnvAsBool("AWS_S3_ENABLE_PRESIGNED_URLS", true),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue := parseInt(value); intValue != 0 {
			return intValue
		}
	}
	return defaultValue
}

func parseInt(s string) int {
	result := 0
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			return 0
		}
	}
	return result
}
