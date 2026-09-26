package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"Hospital-Midderware/internal/application"
	httpDelivery "Hospital-Midderware/internal/delivery/http"
	"Hospital-Midderware/internal/infrastructure"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Initialize Database
	db := initDB()
	defer db.Close()

	// 2. Setup Dependencies (Infrastructure -> Application -> Delivery)
	jwtSecret := getEnv("JWT_SECRET", "super-secret-key-5678")

	staffRepo := infrastructure.NewPostgresStaffRepository(db)
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider(jwtSecret)

	authService := application.NewAuthService(staffRepo, hasher, tokenProvider)
	staffService := application.NewStaffService(staffRepo, hasher)
	staffHandler := httpDelivery.NewStaffHandler(authService, &staffService)

	log.Println("ประกอบร่าง Dependencies เรียบร้อย")

	// 3. Register Routes ด้วย Gin Router
	r := gin.Default()

	// Auth Routes
	r.POST("/staff/login", staffHandler.Login)

	// API v1 Grouping
	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/login", staffHandler.Login)
		v1.POST("/staff/create", staffHandler.CreateStaff) // 🟢 ผูก CreateStaff Route
	}

	// 4. Start HTTP Server
	port := getEnv("SERVER_PORT", ":8080")
	log.Printf("HTTP Server (Gin) กำลังทำงานที่พอร์ต %s ...\n", port)

	// Gin ใช้ r.Run(port) ในการเริ่มเซิร์ฟเวอร์
	if err := r.Run(port); err != nil {
		log.Fatalf("Server ทำงานผิดพลาด: %v", err)
	}
}

// initDB ทำหน้าที่เชื่อมต่อและตั้งค่า Connection Pool ให้ PostgreSQL
func initDB() *sql.DB {
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "KhonNaRak5555")
	dbName := getEnv("DB_NAME", "hospital_middleware")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("ไม่สามารถเปิด Database Connection ได้: %v", err)
	}

	// กำหนด Connection Pool Configurations สำหรับ Production
	db.SetMaxOpenConns(25)                 // จำนวน Connection สูงสุดที่เปิดพร้อมกันได้
	db.SetMaxIdleConns(25)                 // จำนวน Connection สำรองที่รอทำงาน
	db.SetConnMaxLifetime(5 * time.Minute) // อายุสูงสุดของ Connection ก่อนจะถูกปิดและสร้างใหม่

	if err := db.Ping(); err != nil {
		log.Fatalf("ไม่สามารถเชื่อมต่อ Database ได้ (Ping failed): %v", err)
	}

	log.Println("เชื่อมต่อ Database PostgreSQL สำเร็จเรียบร้อย")
	return db
}

// getEnv Helper ฟังก์ชันอ่านค่าจาก Environment Variables ถ้าไม่มีให้ใช้ค่า default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
