package models

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

func GetEmpPermission(db *sql.DB, empCode string) (bool, int, error) {
	const q = `
		SELECT sa_permission_app
		FROM dbo.sys_account
		WHERE sa_emp_code = @p1`
	var perm sql.NullInt64
	err := db.QueryRowContext(context.Background(), q, empCode).Scan(&perm)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		return false, 0, err
	}
	if perm.Valid {
		return true, int(perm.Int64), nil
	}
	return true, 0, nil
}

// models/department.go
func GetDepartment(db *sql.DB, division string) (int, error) {
	const q = `
        SELECT sd_id
        FROM dbo.sys_department
        WHERE sd_dept_code = @p1;`

	var id int
	err := db.QueryRowContext(context.Background(), q, division).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // ไม่พบ division → คืน 0
		}
		return 0, err
	}
	return id, nil
}

func UpdatePermissionApp(db *sql.DB, empCode string, newPerm int) error {
	const q = `
		UPDATE dbo.sys_account
		SET sa_permission_app = @p1
		WHERE sa_emp_code = @p2`
	_, err := db.ExecContext(context.Background(), q, newPerm, empCode)
	return err
}

type SysAccount struct {
	SpgID           int
	SaPermissionWeb int
	SaPermissionApp int
	MpuID           int
	SaUsernameAD    string
	SaEmpCode       string
	SaPassword      string
	SaFirstName     string
	SaLastName      string
	SaEmail         string
	SaStatusFlag    int
	SaCreatedDate   any
	SaCreatedBy     string
	SaUpdatedDate   any
	SaUpdatedBy     string
	SdID            int
}

func InsertAccount(db *sql.DB, a SysAccount) error {
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	var spgAppID sql.NullInt64
	const findSPG = `
        SELECT TOP 1 spg_id
        FROM sys_permission_group WITH (HOLDLOCK)
        WHERE spg_name = @p1 AND spg_permission_flag = @p2
        ORDER BY spg_id;
    `
	if err := tx.QueryRowContext(ctx, findSPG, "Member", 2).Scan(&spgAppID); err != nil {
		_ = tx.Rollback()
		if err == sql.ErrNoRows {
			return fmt.Errorf("no sys_permission_group matched spg_name=Member & spg_permission_flag=2")
		}
		return fmt.Errorf("query sys_permission_group: %w", err)
	}
	if !spgAppID.Valid {
		_ = tx.Rollback()
		return fmt.Errorf("spg_id is NULL for Member/flag=2")
	}

	// 2) Insert sys_account (ใช้ spgAppID เป็น sa_permission_app)
	const tsql = `
        INSERT INTO dbo.sys_account(
            sa_permission_web, sa_permission_app, mpu_id,
            sa_username_ad, sa_emp_code, sa_password,
            sa_firstname, sa_lastname, sa_email,
            sa_status_flag, sa_created_date, sa_created_by,
            sa_updated_date, sa_updated_by, sd_id
        )
        SELECT
            @p1,@p2, @p3, @p4,
            @p5, @p6, @p7,
            @p8, @p9, @p10,
            @p11, @p12, @p13,
            @p14, @p15
        WHERE NOT EXISTS (
            SELECT 1
            FROM dbo.sys_account WITH (UPDLOCK, HOLDLOCK)
            WHERE sa_emp_code = @p6
        );
    `
	if _, err := tx.ExecContext(
		ctx,
		tsql,
		a.SaPermissionWeb,   // @p2
		int(spgAppID.Int64), // @p3  <-- ใช้ spg_id ที่หาได้แทนค่าเดิม
		a.MpuID,             // @p4
		a.SaUsernameAD,      // @p5
		a.SaEmpCode,         // @p6
		a.SaPassword,        // @p7
		a.SaFirstName,       // @p8
		a.SaLastName,        // @p9
		a.SaEmail,           // @p10
		a.SaStatusFlag,      // @p11
		a.SaCreatedDate,     // @p12
		a.SaCreatedBy,       // @p13
		a.SaUpdatedDate,     // @p14
		a.SaUpdatedBy,       // @p15
		a.SdID,              // @p16
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("insert sys_account: %w", err)
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}
