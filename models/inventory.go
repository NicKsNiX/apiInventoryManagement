package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type ItemDetails struct {
	IidID            int        `json:"iid_id"`
	MitID            int        `json:"mit_id"`
	IidQty           *float64   `json:"iid_qty,omitempty"`
	IidStatusFlag    *int       `json:"iid_status_flag,omitempty"`
	IidCreatedDate   *time.Time `json:"iid_created_date,omitempty"`
	IidCreatedBy     *string    `json:"iid_created_by,omitempty"`
	IidUpdatedDate   *time.Time `json:"iid_updated_date,omitempty"`
	IidUpdatedBy     *string    `json:"iid_updated_by,omitempty"`
	LidStatusFlag    *int       `json:"lid_status_flag,omitempty"`
	LidBeforeQty     *float64   `json:"lid_before_qty"`
	LidAfterQty      *float64   `json:"lid_after_qty"`
	MitPlantCode     *int       `json:"mit_plant_code,omitempty"`
	MitWarehouseCode *string    `json:"mit_warehouse_code,omitempty"`
	MitLocation      *string    `json:"mit_location,omitempty"`
	MitItemCode      *string    `json:"mit_item_code,omitempty"`
	MitItemModel     *string    `json:"mit_item_model,omitempty"`
	MitItemStatus    *int       `json:"mit_item_status,omitempty"`
	MitPic           *string    `json:"mit_pic,omitempty"`
	MitSourceCode    *string    `json:"mit_source_code,omitempty"`
	MitPrintCount    *int       `json:"mit_print_count,omitempty"`
	MitCreatedDate   *time.Time `json:"mit_created_date,omitempty"`
	MitCreatedBy     *string    `json:"mit_created_by,omitempty"`
	MitUpdatedDate   *time.Time `json:"mit_updated_date,omitempty"`
	MitUpdatedBy     *string    `json:"mit_updated_by,omitempty"`
	MitCheckBy       *string    `json:"mit_check_by,omitempty"`
	MitTagNo         *string    `json:"mit_tag_no,omitempty"`
	MitStockUnit     *string    `json:"mit_stock_unit,omitempty"`
}

func GetItemDetailsByItem(db *sql.DB, item_cd string) (*ItemDetails, error) {
	const query = `
			SELECT
				i.iid_id, 
				i.mit_id, 
				i.iid_qty, 
				i.iid_status_flag, 
				i.iid_created_date, 
				i.iid_created_by, 
				i.iid_updated_date, 
				i.iid_updated_by,
				m.mit_plant_code,
				m.mit_warehouse_code,
				m.mit_location,
				m.mit_item_code,
				m.mit_item_model,
				m.mit_item_status,
				m.mit_pic,
				m.mit_source_code,
				m.mit_print_count,
				m.mit_created_date AS mit_created_date,
				m.mit_created_by AS mit_created_by,
				m.mit_updated_date AS mit_updated_date,
				m.mit_updated_by AS mit_updated_by,
				m.mit_check_by,
				m.mit_stock_unit
			FROM info_inventory_detail i
			INNER JOIN mst_inventory_tag m ON i.mit_id = m.mit_id
			WHERE 
			m.mit_id = @p1

			;`

	var item ItemDetails
	err := db.QueryRow(query, item_cd).Scan(
		&item.IidID,
		&item.MitID,
		&item.IidQty,
		&item.IidStatusFlag,
		&item.IidCreatedDate,
		&item.IidCreatedBy,
		&item.IidUpdatedDate,
		&item.IidUpdatedBy,
		&item.MitPlantCode,
		&item.MitWarehouseCode,
		&item.MitLocation,
		&item.MitItemCode,
		&item.MitItemModel,
		&item.MitItemStatus,
		&item.MitPic,
		&item.MitSourceCode,
		&item.MitPrintCount,
		&item.MitCreatedDate,
		&item.MitCreatedBy,
		&item.MitUpdatedDate,
		&item.MitUpdatedBy,
		&item.MitCheckBy,
		&item.MitStockUnit,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to fetch item details: %v", err)
	}

	return &item, nil
}

func GetInventoryCheckList(db *sql.DB, employee string) ([]ItemDetails, error) {
	const query = `
        SELECT
            i.iid_id, 
            i.mit_id, 
            i.iid_qty, 
            i.iid_status_flag, 
            i.iid_created_date, 
            i.iid_created_by, 
            i.iid_updated_date, 
            i.iid_updated_by,
            m.mit_plant_code,
            m.mit_warehouse_code,
            m.mit_location,
            m.mit_item_code,
            m.mit_item_model,
            m.mit_item_status,
            m.mit_pic,
            COALESCE(NULLIF(TRIM(m.mit_source_name), ''), m.mit_source_code) AS mit_source,
            m.mit_print_count,
            m.mit_created_date AS mit_created_date,
            m.mit_created_by AS mit_created_by,
            m.mit_updated_date AS mit_updated_date,
            m.mit_updated_by AS mit_updated_by,
			m.mit_tag_no
        FROM info_inventory_detail i
        INNER JOIN mst_inventory_tag m ON i.mit_id = m.mit_id
        WHERE 
        m.mit_check_by = @p1
		ORDER BY m.mit_tag_no ASC
        ;`

	rows, err := db.QueryContext(context.Background(), query, employee)
	if err != nil {
		return nil, fmt.Errorf("query inventory checklist: %w", err)
	}
	defer rows.Close()

	var items []ItemDetails
	for rows.Next() {
		var item ItemDetails
		if err := rows.Scan(
			&item.IidID,
			&item.MitID,
			&item.IidQty,
			&item.IidStatusFlag,
			&item.IidCreatedDate,
			&item.IidCreatedBy,
			&item.IidUpdatedDate,
			&item.IidUpdatedBy,
			&item.MitPlantCode,
			&item.MitWarehouseCode,
			&item.MitLocation,
			&item.MitItemCode,
			&item.MitItemModel,
			&item.MitItemStatus,
			&item.MitPic,
			&item.MitSourceCode,
			&item.MitPrintCount,
			&item.MitCreatedDate,
			&item.MitCreatedBy,
			&item.MitUpdatedDate,
			&item.MitUpdatedBy,
			&item.MitTagNo,
		); err != nil {
			return nil, fmt.Errorf("scan inventory checklist: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err inventory checklist: %w", err)
	}
	// return empty slice (not error) when no rows matched
	return items, nil
}
func GetShowHistoryAuditor(db *sql.DB, employee string) ([]ItemDetails, error) {
	const query = `
        SELECT
            l.iid_id,
            m.mit_warehouse_code,
            m.mit_location,
            m.mit_item_code,
            m.mit_item_model,
            COALESCE(NULLIF(TRIM(m.mit_source_name), ''), m.mit_source_code) AS mit_source,
			m.mit_tag_no,
			l.lid_before_qty,
			l.lid_after_qty
			
        FROM log_inventory_detail l
		LEFT JOIN info_inventory_detail i ON l.iid_id = i.iid_id
        LEFT JOIN mst_inventory_tag m ON i.mit_id = m.mit_id
        WHERE 
            l.lid_created_by = @p1 AND
            l.lid_status_flag IN ('5','9')
			GROUP BY l.iid_id, m.mit_warehouse_code,m.mit_location, m.mit_item_code,mit_item_model,m.mit_source_name,m.mit_source_code,m.mit_print_count,m.mit_tag_no,lid_before_qty,lid_after_qty
			ORDER BY m.mit_tag_no ASC
		
        ;`

	rows, err := db.QueryContext(context.Background(), query, employee)
	if err != nil {
		return nil, fmt.Errorf("query inventory checklist: %w", err)
	}
	defer rows.Close()

	var items []ItemDetails
	for rows.Next() {
		var item ItemDetails
		if err := rows.Scan(
			&item.IidID,
			&item.MitWarehouseCode,
			&item.MitLocation,
			&item.MitItemCode,
			&item.MitItemModel,
			&item.MitSourceCode,
			&item.MitTagNo,
			&item.LidBeforeQty,
			&item.LidAfterQty,
		); err != nil {
			return nil, fmt.Errorf("scan inventory checklist: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err inventory checklist: %w", err)
	}
	// return empty slice (not error) when no rows matched
	return items, nil
}

func UpdateItemQuantity(db *sql.DB, MitID int, qty float64, employee string, IidID int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// อ่านค่าเดิม + สถานะเดิม
	var (
		beforeQty  sql.NullFloat64
		beforeFlag sql.NullInt64
		querySel   string
		argsSel    []interface{}
	)

	if IidID > 0 {
		querySel = `
            SELECT i.iid_qty, i.iid_status_flag
            FROM info_inventory_detail i WITH (UPDLOCK, ROWLOCK)
            WHERE i.iid_id = @p1;
        `
		argsSel = []interface{}{IidID}
	} else {
		// ถ้ามีโอกาสหลายแถวต่อ mit_id แนะนำเพิ่ม TOP 1 + ORDER BY ให้ชัดเจน
		querySel = `
            SELECT TOP 1 i.iid_qty, i.iid_status_flag
            FROM info_inventory_detail i WITH (UPDLOCK, ROWLOCK)
            WHERE i.mit_id = @p1
            ORDER BY i.iid_id DESC;
        `
		argsSel = []interface{}{MitID}
	}

	if err := tx.QueryRow(querySel, argsSel...).Scan(&beforeQty, &beforeFlag); err != nil {
		_ = tx.Rollback()
		if err == sql.ErrNoRows {
			return fmt.Errorf("no inventory detail found to update")
		}
		return fmt.Errorf("failed to read current quantity: %v", err)
	}

	// ถ้าค่าเดิมเป็น 9 ให้คง 9; ถ้าไม่ใช่ ให้เป็น 2 (ตาม logic เดิมของคุณ)
	targetFlag := 2
	if beforeFlag.Valid && beforeFlag.Int64 == 9 {
		targetFlag = 9
	}

	// อัปเดต
	if IidID > 0 {
		const updateQ = `
            UPDATE info_inventory_detail
            SET iid_qty = @p1,
                iid_updated_date = GETDATE(),
                iid_updated_by = @p2,
                iid_status_flag = @p3
            WHERE iid_id = @p4;
        `
		if _, err := tx.Exec(updateQ,
			sql.Named("p1", qty),
			sql.Named("p2", employee),
			sql.Named("p3", targetFlag),
			sql.Named("p4", IidID),
		); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to update item quantity: %v", err)
		}
	} else {
		const updateQ = `
            UPDATE info_inventory_detail
            SET iid_qty = @p1,
                iid_updated_date = GETDATE(),
                iid_updated_by = @p2,
                iid_status_flag = @p3
            WHERE mit_id = @p4;
        `
		if _, err := tx.Exec(updateQ,
			sql.Named("p1", qty),
			sql.Named("p2", employee),
			sql.Named("p3", targetFlag),
			sql.Named("p4", MitID),
		); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to update item quantity: %v", err)
		}
	}

	// หา iid_id สำหรับ log
	lidChannel := 2
	var iidIDForLog int
	if IidID > 0 {
		iidIDForLog = IidID
	} else {
		var tmpIid sql.NullInt64
		// ใช้ TOP 1 ให้ตรงกับการเลือกด้านบน
		err = tx.QueryRow(`
            SELECT TOP 1 iid_id 
            FROM info_inventory_detail WITH (UPDLOCK, ROWLOCK) 
            WHERE mit_id = @p1 
            ORDER BY iid_id DESC;`,
			MitID,
		).Scan(&tmpIid)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to determine iid_id for logging: %v", err)
		}
		if tmpIid.Valid {
			iidIDForLog = int(tmpIid.Int64)
		} else {
			_ = tx.Rollback()
			return fmt.Errorf("iid_id is null for mit_id %d", MitID)
		}
	}

	// insert log; ใส่ lid_status_flag ตาม targetFlag (ถ้าต้องการให้สอดคล้อง)
	insertLogQuery := `
        INSERT INTO log_inventory_detail
            (iid_id, lid_before_qty, lid_after_qty, lid_channel, lid_created_date, lid_created_by, lid_status_flag)
        VALUES
            (@p1, @p2, @p3, @p4, GETDATE(), @p5, @p6);
    `
	beforeVal := 0.0
	if beforeQty.Valid {
		beforeVal = beforeQty.Float64
	}

	if _, err := tx.Exec(insertLogQuery,
		sql.Named("p1", iidIDForLog),
		sql.Named("p2", beforeVal),
		sql.Named("p3", qty),
		sql.Named("p4", lidChannel),
		sql.Named("p5", employee),
		sql.Named("p6", targetFlag), // เดิมเป็น 2; ถ้าอยากตรึงไว้ที่ 2 ก็เปลี่ยนกลับได้
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert log entry: %v", err)
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

func UpdateReconfirmQty(db *sql.DB, qty float64, employee string, itemCode string, warehouse string, iidID int64) error {
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// 1) ดึง iid_qty ปัจจุบัน พร้อมล็อกแถวกันชนกัน
	//    (สำหรับ SQL Server แนะนำ WITH (UPDLOCK, ROWLOCK) ลดโอกาส deadlock)
	var beforeQty sql.NullFloat64
	selectQ := `
        SELECT i.iid_qty
        FROM info_inventory_detail i WITH (UPDLOCK, ROWLOCK)
        WHERE i.iid_id = @p1;
    `
	if err := tx.QueryRowContext(ctx, selectQ, iidID).Scan(&beforeQty); err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return fmt.Errorf("iid_id %d not found", iidID)
		}
		return fmt.Errorf("failed to select current iid_qty: %v", err)
	}

	// 2) insert log (before/after)
	insertLogQ := `
        INSERT INTO log_inventory_detail
            (iid_id, lid_before_qty, lid_after_qty, lid_channel, lid_created_date, lid_created_by, lid_status_flag)
        VALUES
            (@p1, @p2, @p3, @p4, GETDATE(), @p5, 9);
    `
	if _, err := tx.ExecContext(ctx, insertLogQ,
		iidID,
		func() float64 {
			if beforeQty.Valid {
				return beforeQty.Float64
			}
			return 0
		}(),
		qty,
		2,
		employee,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert log_inventory_detail: %v", err)
	}

	// 3) update detail เป็นค่าล่าสุด + สถานะ 9
	updateQ := `
        UPDATE info_inventory_detail
        SET iid_qty = @p1,
            iid_updated_date = GETDATE(),
            iid_updated_by = @p2,
            iid_status_flag = 9
        WHERE iid_id = @p3;
    `
	if _, err := tx.ExecContext(ctx, updateQ, qty, employee, iidID); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update item quantity: %v", err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

func UpdateNoconfirmQty(db *sql.DB, qty float64, employee string, itemCode string, warehouse string, iidID int64) error {
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// 1) ดึง iid_qty ปัจจุบัน พร้อมล็อกแถวกันชนกัน
	//    (สำหรับ SQL Server แนะนำ WITH (UPDLOCK, ROWLOCK) ลดโอกาส deadlock)
	var beforeQty sql.NullFloat64
	selectQ := `
        SELECT i.iid_qty
        FROM info_inventory_detail i WITH (UPDLOCK, ROWLOCK)
        WHERE i.iid_id = @p1;
    `
	if err := tx.QueryRowContext(ctx, selectQ, iidID).Scan(&beforeQty); err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return fmt.Errorf("iid_id %d not found", iidID)
		}
		return fmt.Errorf("failed to select current iid_qty: %v", err)
	}

	// 3) update detail เป็นค่าล่าสุด + สถานะ 9
	updateQ := `
        UPDATE info_inventory_detail
        SET iid_qty = @p1,
            iid_updated_date = GETDATE(),
            iid_updated_by = @p2,
            iid_status_flag = 9
        WHERE iid_id = @p3;
    `
	if _, err := tx.ExecContext(ctx, updateQ, qty, employee, iidID); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update item quantity: %v", err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

func ConfirmEditQtyAdjust(db *sql.DB, qty float64, employee string, itemCode string, warehouse string, iidID int64) error {
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// 1) ดึง iid_qty ปัจจุบัน (ล็อกแถวเพื่อลด race)
	var beforeQty sql.NullFloat64
	selectQ := `
		SELECT i.iid_qty
		FROM info_inventory_detail i WITH (UPDLOCK, ROWLOCK)
		WHERE i.iid_id = @p1;
	`
	if err := tx.QueryRowContext(ctx, selectQ, iidID).Scan(&beforeQty); err != nil {
		_ = tx.Rollback()
		if err == sql.ErrNoRows {
			return fmt.Errorf("iid_id %d not found", iidID)
		}
		return fmt.Errorf("failed to select current iid_qty: %v", err)
	}

	// 2) บันทึก log before/after
	insertLogQ := `
		INSERT INTO log_inventory_detail
			(iid_id, lid_before_qty, lid_after_qty, lid_channel, lid_created_date, lid_created_by,lid_status_flag)
		VALUES
			(@p1, @p2, @p3, @p4, GETDATE(), @p5, 9);
	`
	if _, err := tx.ExecContext(ctx, insertLogQ,
		iidID,
		func() float64 {
			if beforeQty.Valid {
				return beforeQty.Float64
			}
			return 0
		}(),
		qty,
		2,
		employee,
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert log_inventory_detail: %v", err)
	}

	// 3) อัปเดต/สร้าง info_inventory_audit
	//    - อัปเดต: set iia_reconf_qty = qty, iia_conf_flag = 1, iia_reconf_flag = 1
	updateAuditQ := `
		UPDATE info_inventory_audit
		SET
			iia_reconf_qty   = @p1,
			iia_reconf_flag  = 1,
			iia_conf_flag    = 1,
			iia_updated_date = GETDATE(),
			iia_updated_by   = @p2
		WHERE iid_id = @p3;
	`
	res, err := tx.ExecContext(ctx, updateAuditQ, qty, employee, iidID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update info_inventory_audit: %v", err)
	}
	aff, _ := res.RowsAffected()

	if aff == 0 {
		// ถ้าไม่มีแถว audit เดิม → INSERT ใหม่ โดยเก็บ conf_qty = before, reconf_qty = after
		insertAuditQ := `
			INSERT INTO info_inventory_audit
				(iid_id, iia_conf_qty, iia_reconf_qty, iia_conf_flag, iia_reconf_flag, iia_channel,
				 iia_created_date, iia_created_by, iia_updated_date, iia_updated_by)
			VALUES
				(@p1, @p2, @p3, 1, 1, @p4, GETDATE(), @p5, NULL, NULL);
		`
		if _, err := tx.ExecContext(ctx, insertAuditQ,
			iidID,
			func() float64 {
				if beforeQty.Valid {
					return beforeQty.Float64
				}
				return 0
			}(),
			qty, // iia_reconf_qty (หลังแก้)
			2,   // iia_channel (ปรับได้)
			employee,
		); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to insert info_inventory_audit: %v", err)
		}
	}

	// 4) อัปเดต detail เป็นค่าล่าสุด + สถานะ 9
	updateDetailQ := `
		UPDATE info_inventory_detail
		SET iid_qty = @p1,
			iid_updated_date = GETDATE(),
			iid_updated_by   = @p2,
			iid_status_flag  = 9
		WHERE iid_id = @p3;
	`
	if _, err := tx.ExecContext(ctx, updateDetailQ, qty, employee, iidID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update item quantity: %v", err)
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

func InsertInventoryAdjustment(db *sql.DB, iidID int, confQty float64, employee string, channel int) (err error) {
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	// ถ้ามี error ที่ไหนระหว่างทาง ให้ rollback
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1) อ่านจำนวนเดิม (ล็อกแถวไว้ใน tx เดียวกัน)
	var beforeQty float64
	selectQ := `
		SELECT ISNULL(i.iid_qty, 0)
		FROM info_inventory_detail i WITH (UPDLOCK, ROWLOCK)
		WHERE i.iid_id = @p1;
	`
	if err = tx.QueryRowContext(ctx, selectQ, iidID).Scan(&beforeQty); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("iid_id %d not found", iidID)
		}
		return fmt.Errorf("select current iid_qty: %w", err)
	}

	// 2) อัปเดตสถานะ
	const updateQ = `
		UPDATE info_inventory_detail
		SET iid_status_flag = 9
		WHERE iid_id = @p1;
	`
	if _, err = tx.ExecContext(ctx, updateQ, iidID); err != nil {
		return fmt.Errorf("update inventory detail: %w", err)
	}

	// 3) เขียน log
	const insertLogQ = `
		INSERT INTO log_inventory_detail
			(iid_id, lid_before_qty, lid_after_qty, lid_channel, lid_created_date, lid_created_by, lid_status_flag)
		VALUES
			(@p1, @p2, @p3, @p4, GETDATE(), @p5, 9);
	`
	if _, err = tx.ExecContext(ctx, insertLogQ,
		iidID,
		beforeQty, // ก่อนหน้า
		confQty,   // หลังปรับ
		channel,   // ใช้พารามิเตอร์ channel ที่ส่งเข้า function
		employee,
	); err != nil {
		return fmt.Errorf("insert log_inventory_detail: %w", err)
	}

	// 4) คอมมิตปิดท้าย
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func InsertAdjustInventory(db *sql.DB, iidID int, confQty float64, employee string, channel int) (err error) {
	ctx := context.Background()

	// เปิด TX ครั้งเดียวและกำหนด isolation ที่เหมาะสม
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	// ถ้ามี err ระหว่างทาง จะ rollback อัตโนมัติ
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1) อ่าน qty เดิม พร้อมล็อกแถว ลด race
	var beforeQty sql.NullFloat64
	const selectQ = `
		SELECT i.iid_qty
		FROM info_inventory_detail i WITH (UPDLOCK, ROWLOCK)
		WHERE i.iid_id = @p1;
	`
	if err = tx.QueryRowContext(ctx, selectQ, iidID).Scan(&beforeQty); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("iid_id %d not found", iidID)
		}
		return fmt.Errorf("select current iid_qty: %w", err)
	}

	// 2) บันทึกลงตาราง audit
	const insertAuditQ = `
		INSERT INTO info_inventory_audit
			(iid_id, iia_conf_qty, iia_reconf_qty, iia_conf_flag, iia_reconf_flag, iia_channel,
			 iia_created_date, iia_created_by, iia_updated_date, iia_updated_by)
		VALUES
			(@p1, @p2, 0, 1, 0, @p3, GETDATE(), @p4, NULL, NULL);
	`
	if _, err = tx.ExecContext(ctx, insertAuditQ, iidID, confQty, channel, employee); err != nil {
		return fmt.Errorf("insert info_inventory_audit: %w", err)
	}

	// 3) เขียน log before/after
	const insertLogQ = `
		INSERT INTO log_inventory_detail
			(iid_id, lid_before_qty, lid_after_qty, lid_channel, lid_created_date, lid_created_by, lid_status_flag)
		VALUES
			(@p1, @p2, @p3, @p4, GETDATE(), @p5, 5);
	`
	before := 0.0
	if beforeQty.Valid {
		before = beforeQty.Float64
	}
	if _, err = tx.ExecContext(ctx, insertLogQ, iidID, before, confQty, channel, employee); err != nil {
		return fmt.Errorf("insert log_inventory_detail: %w", err)
	}

	// 4) อัปเดตสถานะรายการหลัก
	const updateQ = `
		UPDATE info_inventory_detail
		SET iid_status_flag = 5
		WHERE iid_id = @p1;
	`
	if _, err = tx.ExecContext(ctx, updateQ, iidID); err != nil {
		return fmt.Errorf("update info_inventory_detail: %w", err)
	}

	// 5) ปิดธุรกรรม
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

type LocationDetail struct {
	Location string  `json:"location"`
	Qty      float64 `json:"qty"`
}

type InventoryItem struct {
	Code      string           `json:"code"`
	Warehouse string           `json:"warehouse"` // ✅ เพิ่มฟิลด์นี้
	Total     float64          `json:"total"`
	Details   []LocationDetail `json:"details"`
}

func GetInventorySummary(db *sql.DB, search string) ([]InventoryItem, error) {
	// ค้นหาอย่างเดียว: ว่าง = ไม่ค้นหา
	if len(search) == 0 {
		return []InventoryItem{}, nil
	}
	const q = `
				 SELECT
						m.mit_item_code,
						ISNULL(m.mit_location, '') AS mit_location,
						CAST(SUM(ISNULL(i.iid_qty, 0)) AS DECIMAL(18,3)) AS qty, -- ล็อก scale ชัดเจน
						ISNULL(m.mit_warehouse_code, '') AS mit_warehouse_code
					FROM info_inventory_detail i
					JOIN mst_inventory_tag m ON i.mit_id = m.mit_id
					WHERE m.mit_item_code = @p1
					GROUP BY m.mit_item_code, m.mit_location, m.mit_warehouse_code
					ORDER BY SUM(ISNULL(i.iid_qty, 0)) DESC;
				`
	rows, err := db.QueryContext(context.Background(), q, sql.Named("p1", search))
	if err != nil {
		return nil, fmt.Errorf("query inventory summary: %w", err)
	}
	defer rows.Close()

	var result []InventoryItem

	index := map[string]int{}

	for rows.Next() {
		var code sql.NullString
		var location sql.NullString
		var qty sql.NullFloat64
		var warehouse sql.NullString

		if err := rows.Scan(&code, &location, &qty, &warehouse); err != nil {
			return nil, fmt.Errorf("scan inventory summary: %w", err)
		}

		c := ""
		if code.Valid {
			c = code.String
			// ถ้าต้องการตัดความยาว 25 ตัว + trim:
			if len(c) > 25 {
				c = strings.TrimSpace(c[:25])
			} else {
				c = strings.TrimSpace(c)
			}
		}
		loc := ""
		if location.Valid {
			loc = location.String
		}
		qtyVal := 0.0
		if qty.Valid {
			qtyVal = qty.Float64
		}
		warehouseCode := ""
		if warehouse.Valid {
			warehouseCode = warehouse.String
		}

		key := c + "|" + warehouseCode
		if pos, ok := index[key]; ok {
			result[pos].Details = append(result[pos].Details, LocationDetail{Location: loc, Qty: qtyVal})
			result[pos].Total += qtyVal
		} else {
			item := InventoryItem{
				Code:      c,
				Warehouse: warehouseCode,
				Total:     qtyVal,
				Details:   []LocationDetail{{Location: loc, Qty: qtyVal}},
			}
			index[key] = len(result)
			result = append(result, item)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err inventory summary: %w", err)
	}

	return result, nil
}

func GetInventorySummaryAll(db *sql.DB, employee string) ([]InventoryItem, error) {
	baseQuery := `
					SELECT
					m.mit_item_code,
						ISNULL(m.mit_location, '')                                  AS mit_location,
						CAST(SUM(ISNULL(i.iid_qty, 0)) AS DECIMAL(18,3))            AS qty,
						ISNULL(m.mit_warehouse_code,'')                             AS mit_warehouse_code
					FROM info_inventory_detail i
					JOIN mst_inventory_tag m ON i.mit_id = m.mit_id
					JOIN info_inventory_audit a ON a.iid_id = i.iid_id AND a.iia_created_by = @p1
					GROUP BY m.mit_item_code, m.mit_location, m.mit_warehouse_code
					ORDER BY m.mit_item_code;
				`

	// always pass employee (may be empty string)
	args := []interface{}{employee}

	rows, err := db.QueryContext(context.Background(), baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query inventory summary all: %w", err)
	}
	defer rows.Close()

	var result []InventoryItem
	index := map[string]int{} // key: code|warehouse -> index

	for rows.Next() {
		var (
			code      sql.NullString
			location  sql.NullString
			qty       sql.NullFloat64
			warehouse sql.NullString
		)

		if err := rows.Scan(&code, &location, &qty, &warehouse); err != nil {
			return nil, fmt.Errorf("scan inventory summary all: %w", err)
		}

		c := strings.TrimSpace(code.String)
		loc := location.String
		q := 0.0
		if qty.Valid {
			q = qty.Float64
		}

		w := warehouse.String

		key := c + "|" + w
		if pos, ok := index[key]; ok {
			result[pos].Details = append(result[pos].Details, LocationDetail{Location: loc, Qty: q})
			result[pos].Total += q // รวมตามต้องการ (ถ้าจะรวม qty_adjust ก็ + qa)
		} else {
			item := InventoryItem{
				Code:      c,
				Warehouse: w,
				Total:     q, // หรือ q + qa ถ้าต้องการ
				Details:   []LocationDetail{{Location: loc, Qty: q}},
			}
			index[key] = len(result)
			result = append(result, item)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err inventory summary all: %w", err)
	}

	return result, nil
}

type FlatInventoryRow struct {
	IidID        int64
	ItemCD       string
	Warehouse    string
	Location     string
	QtyBefore    float64
	QtyAuditor   float64
	QtyConfirm   float64
	QtyAdjust    float64
	ReConfirmFlg int64
}

func GetInventoryRowsAll(db *sql.DB, employee string) ([]FlatInventoryRow, error) {
	const q = `
				SELECT
					i.iid_id,
					m.mit_item_code,
					ISNULL(m.mit_warehouse_code,'') AS mit_warehouse_code,
					ISNULL(m.mit_location,'')       AS mit_location,
					a.iia_conf_qty AS qty_adjust,
					i.iid_qty AS qty_before
				FROM info_inventory_detail i
				JOIN mst_inventory_tag m ON i.mit_id = m.mit_id
				JOIN info_inventory_audit a ON a.iid_id = i.iid_id
				WHERE a.iia_conf_flag = 1
				AND i.iid_status_flag = 5
				AND i.iid_updated_by = @p1
				ORDER BY i.iid_id DESC;`

	rows, err := db.QueryContext(context.Background(), q, employee)
	if err != nil {
		return nil, fmt.Errorf("query inventory rows: %w", err)
	}
	defer rows.Close()

	var out []FlatInventoryRow
	for rows.Next() {
		var iid sql.NullInt64
		var item sql.NullString
		var wh sql.NullString
		var loc sql.NullString
		var qtyA sql.NullFloat64
		var qtyB sql.NullFloat64
		if err := rows.Scan(&iid, &item, &wh, &loc, &qtyA, &qtyB); err != nil {
			return nil, fmt.Errorf("scan inventory rows: %w", err)
		}

		it := FlatInventoryRow{
			IidID:     int64(iid.Int64),
			ItemCD:    strings.TrimSpace(item.String),
			Warehouse: wh.String,
			Location:  loc.String,
			QtyAdjust: func() float64 {
				if qtyA.Valid {
					return qtyA.Float64
				}
				return 0
			}(),
			QtyBefore: func() float64 {
				if qtyB.Valid {
					return qtyB.Float64
				}
				return 0
			}(),
		}
		out = append(out, it)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err inventory rows: %w", err)
	}
	return out, nil
}

func GetInventoryHistory(db *sql.DB, employee string) ([]FlatInventoryRow, error) {
	const q = `
SELECT 
    m.mit_id,
    m.mit_item_code,
    ISNULL(m.mit_warehouse_code,'') AS mit_warehouse_code,
    ISNULL(m.mit_location,'')       AS mit_location,
    l.lid_before_qty AS qty_before,
    a.iia_conf_qty   AS qty_auditor,
    a.iia_reconf_qty AS qty_confirm,
    a.iia_reconf_flag
FROM info_inventory_detail i
LEFT JOIN mst_inventory_tag m ON i.mit_id = m.mit_id
LEFT JOIN info_inventory_audit a ON a.iid_id = i.iid_id
OUTER APPLY (
    SELECT TOP 1 lid_before_qty, lid_after_qty
    FROM log_inventory_detail l
    WHERE l.iid_id = i.iid_id
    ORDER BY l.lid_id DESC
) l
WHERE a.iia_conf_flag = 1
  AND i.iid_status_flag = 9
  AND i.iid_updated_by = @p1
ORDER BY m.mit_item_code;
`

	rows, err := db.QueryContext(context.Background(), q, employee)
	if err != nil {
		return nil, fmt.Errorf("query inventory rows: %w", err)
	}
	defer rows.Close()

	var out []FlatInventoryRow
	for rows.Next() {
		var (
			iid          sql.NullInt64
			item         sql.NullString
			wh           sql.NullString
			loc          sql.NullString
			qtyBefore    sql.NullFloat64
			qtyAuditor   sql.NullFloat64
			qtyConfirm   sql.NullFloat64
			ReConfirmFlg sql.NullInt64 // iia_reconf_flag (ไม่ใช้ แต่ต้อง Scan ให้ครบคอลัมน์)
		)

		if err := rows.Scan(
			&iid, &item, &wh, &loc,
			&qtyBefore, &qtyAuditor, &qtyConfirm, &ReConfirmFlg,
		); err != nil {
			return nil, fmt.Errorf("scan inventory rows: %w", err)
		}

		ff := func(n sql.NullFloat64) float64 {
			if n.Valid {
				return n.Float64
			}
			return 0
		}

		out = append(out, FlatInventoryRow{
			IidID:        iid.Int64,
			ItemCD:       strings.TrimSpace(item.String),
			Warehouse:    strings.TrimSpace(wh.String),
			Location:     strings.TrimSpace(loc.String),
			QtyBefore:    ff(qtyBefore),
			QtyAuditor:   ff(qtyAuditor),
			QtyConfirm:   ff(qtyConfirm),
			ReConfirmFlg: ReConfirmFlg.Int64,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err inventory rows: %w", err)
	}
	return out, nil
}

func GetInventoryReconfirmCount(db *sql.DB, employee string) (int, error) {
	const q = `SELECT COUNT(*) 
		FROM info_inventory_detail i
		JOIN mst_inventory_tag m ON i.mit_id = m.mit_id
		JOIN info_inventory_audit a ON a.iid_id = i.iid_id
		WHERE a.iia_conf_flag = 1
		AND i.iid_status_flag = 5
		AND i.iid_updated_by = @p1
	
	`
	var count int
	err := db.QueryRowContext(context.Background(), q, employee).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query inventory reconfirm count: %w", err)
	}
	return count, nil
}
