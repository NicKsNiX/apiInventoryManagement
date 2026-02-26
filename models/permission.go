package models

import (
	"context"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

type SysPermissionGroup struct {
	SpgID             int
	SpgName           string
	SpgPermissionFlag int
}

func GetPermissionGroupIDsByName(db *sql.DB, name string) (webID int, appID int, err error) {
	const q = `
        SELECT spg_id, spg_permission_flag
        FROM dbo.sys_permission_group
        WHERE spg_name = @p1`

	rows, err := db.QueryContext(context.Background(), q, name)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var flag sql.NullInt64
		if err := rows.Scan(&id, &flag); err != nil {
			return 0, 0, err
		}
		if flag.Valid {
			if flag.Int64&1 != 0 {
				webID = id
			}
			if flag.Int64&2 != 0 {
				appID = id
			}
		}
	}
	if err := rows.Err(); err != nil {
		return 0, 0, err
	}
	return webID, appID, nil
}
