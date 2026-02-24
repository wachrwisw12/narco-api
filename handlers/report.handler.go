package handlers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"api-naco/db"
	"api-naco/models"
	"api-naco/services"

	"github.com/gofiber/fiber/v2"
)

func ReceiveReport(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

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

	var conditions []string
	var args []interface{}
	argPos := 1

	switch d.RoleId {
	case 5:
		// super admin see all

	case 4:
		conditions = append(conditions, "sd.district_id = $"+strconv.Itoa(argPos))
		args = append(args, d.DistrictId)
		argPos++

		conditions = append(conditions, "ir.status <> $"+strconv.Itoa(argPos))
		args = append(args, 1)
		argPos++

	case 3:
		conditions = append(conditions, "dt.id = $"+strconv.Itoa(argPos))
		args = append(args, d.DistrictId)
		argPos++
	}

	// ✅ filter ตาม status จาก tab
	statusParam := c.Query("status")
	println("statusParam:", statusParam)
	if statusParam != "" {
		statusInt, err := strconv.Atoi(statusParam)
		if err != nil {
			return fiber.NewError(400, "invalid status")
		}

		conditions = append(conditions, "ir.status = $"+strconv.Itoa(argPos))
		args = append(args, statusInt)
		argPos++
	}
	// ✅ สำคัญ: ต้องสร้าง whereClause
	var whereClause string
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := `
	SELECT 
       ir.id,
       ir.tracking_code,
       ir.details,
       ir.status,
       COUNT(rf.id) AS file_count,
       ir.village,
       COALESCE(ir.village,'') || ' ต.' ||
       COALESCE(sd.name_th,'') || ' อ.' ||
       COALESCE(dt.name_th,'') || ' จ.' ||
       COALESCE(p.name_th,'') AS fullarea,
       ir.sub_district_id,
       rs.name_status,
       ir.created_at,
       ir.updated_at
FROM incident_reports ir
LEFT JOIN report_files rf 
       ON ir.id = rf.incident_report_id
LEFT JOIN sub_districts sd 
       ON ir.sub_district_id = sd.id
LEFT JOIN districts dt 
       ON dt.id = sd.district_id
LEFT JOIN provinces p 
       ON p.id = dt.province_id
INNER JOIN report_status rs 
       ON rs.id_status = ir.status
` + whereClause + `
GROUP BY 
       ir.id,
       ir.tracking_code,
       ir.details,
       ir.status,
       ir.village,
       sd.name_th,
       dt.name_th,
       p.name_th,
       ir.sub_district_id,
       rs.name_status,
       ir.created_at,
       ir.updated_at
ORDER BY ir.created_at DESC;
	`

	rows, err := db.DB.Query(ctx, query, args...)
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
			&r.FileCount,
			&r.Village,
			&r.Fullarea,
			&r.SubDistrictId,
			&r.NameStatus,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return fiber.NewError(500, err.Error()) // ส่ง error จริง
		}
		reports = append(reports, r)
	}

	// ✅ สำคัญมาก
	if err := rows.Err(); err != nil {
		return fiber.NewError(500, err.Error())
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

func GetReportById(c *fiber.Ctx) error {
	id := c.Params("id")

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
	WHERE ir.id = $1
	LIMIT 1
	`

	err := db.DB.QueryRow(ctx, query, id).Scan(
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
