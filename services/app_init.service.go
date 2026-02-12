package services

import (
	"context"
	"time"

	"api-naco/db"
	"api-naco/models"
)

func DistrictsService(q string, provinceID string) ([]models.District, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
	sd.id,
	sd.district_id,
	p.id,
	sd.name_th,
	d.name_th AS district,
	p.name_th AS province
FROM sub_districts sd
JOIN districts d ON d.id = sd.district_id
JOIN provinces p ON p.id = d.province_id
WHERE sd.name_th ILIKE '%' || $1 || '%' AND d.id::text LIKE $2 || '%'

LIMIT 20

	`

	rows, err := db.DB.Query(ctx, query, q, provinceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.District

	for rows.Next() {
		var d models.District

		if err := rows.Scan(
			&d.ID,
			&d.DistrictId,
			&d.ProvinceId,
			&d.SubDistricts,
			&d.District,
			&d.Province,
		); err != nil {
			return nil, err
		}

		result = append(result, d)
	}

	return result, nil
}
