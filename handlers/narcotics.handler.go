package handlers

import (
	"context"
	"fmt"
	"strconv"
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

func SendReport(c *fiber.Ctx) error {
	// 1️⃣ รับ field ธรรมดา

	subDistrictStr := c.FormValue("sub_district_id")

	subDistrictID, err := strconv.Atoi(subDistrictStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid sub_district_id",
		})
	}

	req := models.SendReportRequest{
		Details:       c.FormValue("details"),
		SubDistrictId: subDistrictID,
		Village:       c.FormValue("village"),
	}
	if req.Details == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing details")
	}
	// println("ssdfsfdsdf", req.Details)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 2️⃣ insert incident_reports
	var report models.NacorticsReport
	errq := db.DB.QueryRow(ctx, `
		INSERT INTO incident_reports (details,sub_district_id,village)
		VALUES ($1,$2,$3)
		RETURNING id, details, tracking_code
	`, req.Details, req.SubDistrictId, req.Village).Scan(
		&report.ID,
		&report.Details,
		&report.TrackingCode,
	)
	if errq != nil {
		return fiber.NewError(500, errq.Error())
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

func UpdateStatusReport(c *fiber.Ctx) error {
	var status models.StatusRequest
	username, ok := c.Locals("username").(string)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.BodyParser(&status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "body invalid",
			"error":   err.Error(),
		})
	}
	query := `
	UPDATE incident_reports
	SET status = $1,
	    operate_name=$2,
	    received_at = now()
	WHERE id = $3 AND status = 1
	`
	cmd, err := db.DB.Exec(ctx, query, status.Action, username, status.Id)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	if cmd.RowsAffected() == 0 {
		return fiber.NewError(400, "report already received or not found")
	}

	return c.JSON(fiber.Map{
		"username": username,
		"status":   status,
		"ok":       ok,
	})
}

func Test(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
	})
}
