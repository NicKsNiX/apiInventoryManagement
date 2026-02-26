package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"inventory-management/database"
	"inventory-management/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// แปลงค่าอะไรก็ได้ให้เป็น string แบบปลอดภัย (ถ้า nil คืน "")
func asString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		// JSON number → float64
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case nil:
		return ""
	default:
		return fmt.Sprint(t)
	}
}

func Login(c *fiber.Ctx) error {
	db := database.DB

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).SendString("Invalid request")
	}

	loginURL := "http://192.168.161.102:9999/login"
	loginPayload := map[string]string{
		"username": input.Username,
		"password": input.Password,
	}
	jsonData, _ := json.Marshal(loginPayload)

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(500).SendString("Error during login request: " + err.Error())
	}
	defer resp.Body.Close()

	var loginResult map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResult); err != nil {
		return c.Status(500).SendString("Error decoding login response: " + err.Error())
	}

	// ถ้า status code จาก service อื่นไม่ใช่ 200 ให้ส่งต่อข้อความ error กลับไป
	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).JSON(loginResult)
	}

	// ดึง "user" แบบปลอดภัย
	userObj, ok := loginResult["user"].(map[string]interface{})
	if !ok || userObj == nil {
		log.Printf("unexpected login response structure: %+v", loginResult)
		return c.Status(500).JSON(fiber.Map{
			"error": "Unexpected login response: 'user' not found",
			"code":  "USER_FIELD_MISSING",
		})
	}

	// อ่านฟิลด์ต่าง ๆ แบบไม่ panic
	empCode := asString(userObj["employeeID"])
	department := asString(userObj["division"])
	if empCode == "" {
		return c.Status(400).SendString("emp_code is invalid or not found")
	}
	firstname := asString(userObj["name"])
	lastname := asString(userObj["surname"])
	email := asString(userObj["email"])
	position := asString(userObj["position"])

	// หา mpu_id เฉพาะเมื่อมี position
	var mpuID int
	if position != "" {
		var err error
		mpuID, err = models.GetMpuIDByName(db, position)
		if err != nil {
			return c.Status(500).SendString("Error getting mpu_id: " + err.Error())
		}
	} else {
		mpuID = 0 // หรือกำหนด default ตามที่ระบบต้องการ
	}

	var departmentID int
	if department != "" {
		id, err := models.GetDepartment(db, department)
		if err != nil {
			return c.Status(500).SendString("Error getting department ID: " + err.Error())
		}
		departmentID = id
	} else {
		departmentID = 0 // default
	}

	// ดึงกลุ่มสิทธิ์เริ่มต้น "Member"
	webSpgID, appSpgID, err := models.GetPermissionGroupIDsByName(db, "Member")
	if err != nil {
		return c.Status(500).SendString("Error getting permission group: " + err.Error())
	}

	// ตรวจสอบว่ามีบัญชีอยู่แล้วหรือยัง
	exists, perm, err := models.GetEmpPermission(db, empCode)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":  "check exists failed",
			"code":   "CHECK_EXISTS_FAILED",
			"detail": err.Error(),
		})
	}

	if exists {
		// ถ้ามีอยู่แล้วและ sa_permission_app == 0 ให้ปรับเป็น 3 (หรือค่าที่ต้องการ)
		if perm == 0 {
			if err := models.UpdatePermissionApp(db, empCode, 3); err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error":  "failed to update permission",
					"code":   "UPDATE_PERMISSION_FAILED",
					"detail": err.Error(),
				})
			}
		}
	} else {
		// ยังไม่มีบัญชี → insert ใหม่
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":  "failed to hash password",
				"code":   "HASH_FAILED",
				"detail": err.Error(),
			})
		}

		// ตั้งค่ากลุ่มสิทธิ์ (ถ้าไม่มีให้ 0)
		var saPermWeb, saPermApp int
		if webSpgID > 0 {
			saPermWeb = webSpgID
		}
		if appSpgID > 0 {
			saPermApp = appSpgID
		}

		newAcc := models.SysAccount{
			SaPermissionWeb: saPermWeb,
			SaPermissionApp: saPermApp,
			MpuID:           mpuID,
			SaUsernameAD:    input.Username,
			SaEmpCode:       empCode,
			SaPassword:      string(hashed),
			SaFirstName:     firstname,
			SaLastName:      lastname,
			SaEmail:         email,
			SaStatusFlag:    1,
			SaCreatedDate:   time.Now(),
			SaCreatedBy:     "SYSTEM",
			SaUpdatedDate:   time.Now(),
			SaUpdatedBy:     "SYSTEM",
			SdID:            departmentID,
		}

		if err := models.InsertAccount(db, newAcc); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":  "Error inserting new account",
				"code":   "INSERT_FAILED",
				"detail": err.Error(),
			})
		}
	}

	var SaPermissionApp int
	if err := db.QueryRow("SELECT sa_permission_app FROM dbo.sys_account WHERE sa_emp_code = @p1", empCode).Scan(&SaPermissionApp); err != nil {
		log.Printf("failed to read sa_permission_app for %s: %v", empCode, err)
		SaPermissionApp = 0
	}
	var SaStatusFlg int
	if err := db.QueryRow("SELECT sa_status_flag FROM dbo.sys_account WHERE sa_emp_code = @p1", empCode).Scan(&SaStatusFlg); err != nil {
		log.Printf("failed to read sa_permission_app for %s: %v", empCode, err)
		SaStatusFlg = 0
	}
	// log.Printf("DEBUG: SaStatusFlg=%d", SaStatusFlg)

	var spgName string
	if SaPermissionApp > 0 {
		if err := db.QueryRow("SELECT spg_name FROM dbo.sys_permission_group WHERE spg_id = @p1", SaPermissionApp).Scan(&spgName); err != nil {
			log.Printf("failed to read spg_name for spg_id=%d: %v", SaPermissionApp, err)
			spgName = ""
		}
	}

	var respHash string
	if h, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost); err == nil {
		respHash = string(h)
	}

	return c.JSON(fiber.Map{
		"message":           "Login successful",
		"data":              loginResult,
		"sa_permission_app": SaPermissionApp,
		"spg_name":          spgName,
		"sa_status_flg":     SaStatusFlg,
		"credentials": fiber.Map{
			"password_hash": respHash,
		},
	})
}
