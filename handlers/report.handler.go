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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	username, ok := c.Locals("username").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	d, err := services.FindDistrictUsername(username)
	if err != nil {
		return fiber.NewError(500, "cannot find user district")
	}
	if d.RoleId == 1 {
	} else if d.RoleId == 2 {
	} else if d.RoleId == 3 {
	}
	// var u *models.User
	// u, err := services.FindUserByUsername(user)
	// userDistrict := u.DistrictId
	println("printuser", d.DistrictId, d.RoleId)
	query := `
	SELECT ir.id,
	       ir.tracking_code,
	       ir.details,
	       ir.status,
		   ir.village,
		   CONCAT(ir.village,' ต.',sd.name_th,' อ.',dt.name_th,' จ.',p.name_th) AS fullarea,
		   ir.sub_district_id,
	       rs.name_status,
	       ir.created_at,
	       ir.updated_at
	FROM incident_reports ir
 LEFT JOIN sub_districts sd ON ir.sub_district_id = sd.id
 LEFT JOIN districts dt ON dt.id=sd.district_id
 INNER JOIN provinces p ON p.id=dt.province_id
 INNER JOIN report_status rs ON rs.id_status = ir.status 
 WHERE ir.status=1
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
			&r.Details,
			&r.Status,
			&r.Village,
			&r.Fullarea,
			&r.SubDistrictId,
			&r.NameStatus,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return fiber.NewError(500, "scan error")
		}
		reports = append(reports, r)
	}
	if reports == nil {
		reports = []models.NacorticsReport{}
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
