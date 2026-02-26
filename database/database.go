// database/connect.go
package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb" // import driver
)

var DB *sql.DB

// ฟังก์ชันเชื่อมต่อฐานข้อมูล
func ConnectDB() error {
	user := "sa"
	password := "Te@m1nw"
	host := "192.168.161.101"
	port := "1433"
	database := "ims_db_dev"

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
		user, password, host, port, database)

	var err error
	DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		return err
	}

	// ตรวจสอบการเชื่อมต่อ
	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Connected to SQL Server " + database + " successfully!")
	return nil
}
