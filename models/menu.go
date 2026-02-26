package models

import (
	"context"
	"fmt"
	"inventory-management/database"

	_ "github.com/denisenkom/go-mssqldb"
)

// ฟังก์ชันในการดึงข้อมูลจากฐานข้อมูล
type MenuDetail struct {
	MenuDetailName string `json:"smd_name"`
	LinkController string `json:"link_controller"`
	IconName       string `json:"smg_icon_name"`
}

func GetMenuDetails(ctx context.Context, spgID int, employee string) ([]MenuDetail, error) {
	const query = `
        SELECT
            smd.smd_name AS menu_name,
            smd.smd_link_controller AS link_controller,
            smg.smg_icon_name
        FROM sys_permission_detail spd
        LEFT JOIN sys_permission_group spg ON spd.spg_id = spg.spg_id
        LEFT JOIN sys_account sa ON sa.sa_permission_app = spg.spg_id
        LEFT JOIN sys_menu_detail smd ON smd.smd_id = spd.smd_id
        LEFT JOIN sys_menu_group smg ON smg.smg_id = smd.smg_id
        WHERE
            spg.spg_id = @p1
            AND sa.sa_emp_code = @p2
            AND smd.smd_status_flag = '1'
            AND spd.spd_status_flag = '1'
        ORDER BY smd.smd_order_no ASC;
    	`
	rows, err := database.DB.QueryContext(ctx, query, spgID, employee)
	if err != nil {
		return nil, fmt.Errorf("query menu details: %w", err)
	}
	defer rows.Close()

	var menuDetails []MenuDetail
	for rows.Next() {
		var md MenuDetail
		if err := rows.Scan(&md.MenuDetailName, &md.LinkController, &md.IconName); err != nil {
			return nil, fmt.Errorf("scan menu detail: %w", err)
		}
		menuDetails = append(menuDetails, md)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return menuDetails, nil
}
