package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	Minio  *minio.Client
	Bucket string
)

func InitMinio() error {
	endpoint := os.Getenv("MINIO_INTERNAL_ENDPOINT") // minio:9000
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	Bucket = os.Getenv("MINIO_BUCKET")

	if endpoint == "" || accessKey == "" || secretKey == "" || Bucket == "" {
		return fmt.Errorf("MinIO env not set")
	}

	var err error
	const maxRetries = 10

	for i := 1; i <= maxRetries; i++ {
		log.Printf("⏳ Connecting to MinIO (%d/%d)...", i, maxRetries)

		Minio, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: false, // 🔥 internal = HTTP เท่านั้น
		})
		if err != nil {
			log.Println("❌ MinIO client error:", err)
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		exists, err := Minio.BucketExists(ctx, Bucket)
		if err != nil {
			log.Println("❌ MinIO not ready:", err)
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}

		if !exists {
			if err := Minio.MakeBucket(ctx, Bucket, minio.MakeBucketOptions{}); err != nil {
				return fmt.Errorf("create bucket failed: %w", err)
			}
			log.Println("🪣 Bucket created:", Bucket)
		}

		log.Println("✅ MinIO connected:", endpoint)
		return nil
	}

	return fmt.Errorf("MinIO unavailable after %d retries", maxRetries)
}
