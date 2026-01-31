package handlers

import (
	"context"
	"fmt"
	"time"

	"api-naco/db"
	"api-naco/models"
	"api-naco/services"

	"github.com/gofiber/fiber/v2"
)

func ReportInit(c *fiber.Ctx) error {
	println("test innit")
	return nil
}

type SendReportRequest struct {
	Details string `json:"details"`
}

func SendReport(c *fiber.Ctx) error {
	// 1️⃣ รับ field ธรรมดา
	req := SendReportRequest{
		Details: c.FormValue("details"),
	}
	if req.Details == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing details")
	}
	println("ssdfsfdsdf", req.Details)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 2️⃣ insert incident_reports
	var report models.NacorticsReport
	err := db.DB.QueryRow(ctx, `
		INSERT INTO incident_reports (details)
		VALUES ($1)
		RETURNING id, details, tracking_code
	`, req.Details).Scan(
		&report.ID,
		&report.Details,
		&report.TrackingCode,
	)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	// base path สำหรับไฟล์
	// basePath := fmt.Sprintf("%s/", report.TrackingCode)
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	basePath := fmt.Sprintf("%s/%s/%s", year, month, report.TrackingCode)

	// 3️⃣ รับไฟล์ (ถ้ามี)
	var uploadedFiles []services.UploadedFile

	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["images"]
		if len(files) > 0 {
			uploadedFiles, err = services.UploadReportImages(
				"narcotics-report",
				basePath,
				files,
			)
			if err != nil {
				return fiber.NewError(500, err.Error())
			}
		}
	}

	// 4️⃣ insert report_files
	for _, f := range uploadedFiles {
		_, err := db.DB.Exec(ctx, `
	INSERT INTO report_files (
		incident_report_id,
		object_key,
		file_name,
		file_type,
		file_size,
		storage_bucket,
		storage_base_path,
		storage_version
	) VALUES ($1, $2, $3, $4, $5, 'narcotics-report', $6, 1)
`,
			report.ID,   // $1
			f.ObjectKey, // $2
			f.FileName,  // $3
			f.MimeType,  // $4
			f.FileSize,  // $5
			basePath,    // $6
		)
		if err != nil {
			return fiber.NewError(500, err.Error())
		}
	}

	// 5️⃣ response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":            report.ID,
			"tracking_code": report.TrackingCode,
			"files":         uploadedFiles,
		},
	})
}

func Test(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
	})
}
