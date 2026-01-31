package services

import (
	"context"
	"mime/multipart"
	"path"
	"path/filepath"
	"time"

	"api-naco/db"
	"api-naco/models"
	"api-naco/storage"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type UploadedFile struct {
	ObjectKey string
	FileName  string
	MimeType  string
	FileSize  int64
}

func UploadReportImages(
	bucket string,
	basePath string,
	files []*multipart.FileHeader,
) ([]UploadedFile, error) {
	var uploaded []UploadedFile

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		ext := filepath.Ext(file.Filename)
		filename := uuid.New().String() + ext

		// ✅ object key ที่ถูกต้อง
		objectKey := path.Join(basePath, "images", filename)

		_, err = storage.Minio.PutObject(
			ctx,
			bucket,
			objectKey,
			src,
			file.Size,
			minio.PutObjectOptions{
				ContentType: file.Header.Get("Content-Type"),
			},
		)
		if err != nil {
			return nil, err
		}

		uploaded = append(uploaded, UploadedFile{
			ObjectKey: objectKey,
			FileName:  file.Filename,
			MimeType:  file.Header.Get("Content-Type"),
			FileSize:  file.Size,
		})
	}

	return uploaded, nil
}

func GetFilesByReportID(
	ctx context.Context,
	reportID int,
) ([]models.ReportFile, error) {
	query := `
	SELECT id,
	       object_key,
	       file_name,
	       file_type,
	       file_size,
	       created_at
	FROM report_files
	WHERE incident_report_id = $1
	ORDER BY created_at ASC
	`

	rows, err := db.DB.Query(ctx, query, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]models.ReportFile, 0)

	for rows.Next() {
		var f models.ReportFile
		if err := rows.Scan(
			&f.ID,
			&f.ObjectKey,
			&f.FileName,
			&f.FileType,
			&f.FileSize,
			&f.CreatedAt,
		); err != nil {
			return nil, err
		}

		streamURL, err := storage.PresignedGetURL(
			ctx,
			f.ObjectKey,
			10*time.Minute,
		)
		if err != nil {
			return nil, err
		}

		f.StreamURL = streamURL
		files = append(files, f)
	}

	return files, nil // ถ้าไม่พบ = []
}
