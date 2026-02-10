package storage

import (
	"context"
	"fmt"
	"os"
	"time"
)

func PresignedGetURL(
	ctx context.Context,
	objectPath string,
	expire time.Duration,
) (string, error) {
	u, err := Minio.PresignedGetObject(
		ctx,
		Bucket,
		objectPath,
		expire,
		nil,
	)
	if err != nil {
		return "", err
	}

	fmt.Println("RAW URL :", u.String())

	// ===== override public =====
	scheme := os.Getenv("MINIO_PUBLIC_SCHEME")
	host := os.Getenv("MINIO_PUBLIC_HOST")
	publicPath := os.Getenv("MINIO_PUBLIC_PATH")

	if scheme != "" {
		u.Scheme = scheme
	}
	if host != "" {
		u.Host = host
	}
	if publicPath != "" {
		u.Path = publicPath + u.Path
	}

	fmt.Println("PUBLIC URL:", u.String())

	return u.String(), nil
}
