package services

import (
	"context"
	"time"

	"api-naco/db"
	"api-naco/models"
)

func GetMenusByRole(role string) ([]models.Menu, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT m.menu_id ,mu.code,mu.label,mu."path",mu.icon ,mu.sort_order FROM roles r 
LEFT JOIN role_menus m ON m.role_id=r."id"
LEFT JOIN menus mu ON mu."id" =m.menu_id
WHERE r.role_name=$1
	`

	rows, err := db.DB.Query(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	menus := []models.Menu{}

	for rows.Next() {
		var m models.Menu
		if err := rows.Scan(
			&m.ID,
			&m.Code,
			&m.Label,
			&m.Path,
			&m.Icon,
			&m.Sort,
		); err != nil {
			return nil, err
		}
		menus = append(menus, m)
	}

	return menus, nil
}
