package handlers

import (
	"context"
	"time"

	"api-naco/db"
	"api-naco/models"
	"api-naco/services"

	"github.com/gofiber/fiber/v2"
)

func ReceiveReport(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	UPDATE incident_reports
	SET status = 'RECEIVED',
	    received_at = now()
	WHERE id = $1 AND status = 'PENDING'
	`

	cmd, err := db.DB.Exec(ctx, query, id)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	if cmd.RowsAffected() == 0 {
		return fiber.NewError(400, "report already received or not found")
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func ListReports(c *fiber.Ctx) error {
	println("sdfsdf")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT id,tracking_code,details,status,created_at
	FROM incident_reports

	`

	rows, err := db.DB.Query(ctx, query)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	defer rows.Close()

	var reports []models.NacorticsReport

	for rows.Next() {
		var r models.NacorticsReport
		if err := rows.Scan(
			&r.ID,
			&r.TrackingCode,
			// &r.Category,
			&r.Details,
			&r.Status,
			&r.CreatedAt,
		); err != nil {
			return fiber.NewError(500, "scan error")
		}
		reports = append(reports, r)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(reports),
		"data":    reports,
	})
}

type TrackRequest struct {
	TrackingCode string `json:"tracking_code"`
}

func TrackReport(c *fiber.Ctx) error {
	var req TrackRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var report models.NacorticsReport

	query := `
	SELECT ir.id,
	       ir.tracking_code,
	       ir.details,
	       ir.status,
	       rs.name_status,
	       ir.created_at,
	       ir.updated_at
	FROM incident_reports ir
	JOIN report_status rs ON rs.id_status = ir.status
	WHERE ir.tracking_code = $1
	LIMIT 1
	`

	err := db.DB.QueryRow(ctx, query, req.TrackingCode).Scan(
		&report.ID,
		&report.TrackingCode,
		&report.Details,
		&report.Status,
		&report.NameStatus,
		&report.CreatedAt,
		&report.UpdatedAt,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "report not found")
	}

	files, err := services.GetFilesByReportID(ctx, report.ID)
	if err != nil {
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"cannot load report files",
		)
	}

	report.Files = files

	return c.JSON(fiber.Map{
		"success": true,
		"data":    report,
	})
}
