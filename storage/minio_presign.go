package storage

import (
	"context"
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

	// 🔥 override ให้เป็น public
	u.Scheme = os.Getenv("MINIO_PUBLIC_SCHEME") // https
	u.Host = os.Getenv("MINIO_PUBLIC_HOST")     // vm11-skko.moph.go.th
	u.Path = os.Getenv("MINIO_PUBLIC_PATH") + u.Path
	// /drugnarco/minio + /bucket/...

	return u.String(), nil
}
