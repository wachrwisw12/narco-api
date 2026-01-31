package handlers

import (
	"context"
	"fmt"
	"time"

	"api-naco/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func UploadHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file required"})
	}

	// 🔐 limit size (เช่น 50MB)
	if file.Size > 50*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"error": "file too large"})
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	objectName := fmt.Sprintf(
		"uploads/%s/%s",
		time.Now().Format("2006/01/02"),
		uuid.NewString()+"_"+file.Filename,
	)

	_, err = storage.Minio.PutObject(
		context.Background(),
		storage.Bucket,
		objectName,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"path": objectName,
		"size": file.Size,
	})
}
