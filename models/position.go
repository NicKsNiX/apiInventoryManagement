package models

import (
	"context"
	"database/sql"
)

// GetMpuIDByName ค้นหา mpu_id โดยใช้ mpu_name (SQL Server ใช้ @p1)
func GetMpuIDByName(db *sql.DB, mpuName string) (int, error) {
	const query = `
		SELECT TOP (1) mpu_id
		FROM dbo.mst_position_user
		WHERE mpu_name = @p1 AND mpu_status_flag = 1;`

	var mpuID int
	err := db.QueryRowContext(context.Background(), query, mpuName).Scan(&mpuID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return mpuID, nil
}
