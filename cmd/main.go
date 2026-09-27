package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"Hospital-Middleware/internal/application"
	httpDelivery "Hospital-Middleware/internal/delivery/http"
	"Hospital-Middleware/internal/delivery/http/middleware"
	"Hospital-Middleware/internal/domain"
	"Hospital-Middleware/internal/infrastructure"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// setupRouter แยกฟังก์ชันประกอบร่าง Dependencies และ Route ทั้งหมดออกมาเพื่อความง่ายในการเขียน Unit Test
func setupRouter(
	authService *application.AuthService,
	staffService *application.StaffService,
	searchPatientUseCase *application.SearchPatient,
	tokenProvider domain.TokenProvider,
) *gin.Engine {
	r := gin.Default()

	staffHandler := httpDelivery.NewStaffHandler(authService, staffService)
	patientHandler := httpDelivery.NewPatientHandler(searchPatientUseCase)

	r.POST("/staff/login", staffHandler.Login)
	r.POST("/staff/create", staffHandler.CreateStaff)

	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware(tokenProvider))
	{
		protected.GET("/patient/search", patientHandler.SearchPatient)
	}

	return r
}

func main() {
	// 1. Initialize Database
	db := initDB()
	defer db.Close()

	// 2. Setup Dependencies
	jwtSecret := getEnv("JWT_SECRET", "super-secret-key-5678")

	staffRepo := infrastructure.NewPostgresStaffRepository(db)
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider(jwtSecret)

	hospitalABaseURL := getEnv("HOSPITAL_A_BASE_URL", "https://hospital-a.api.co.th")
	hospitalBBaseURL := getEnv("HOSPITAL_B_BASE_URL", "https://hospital-b.api.co.th")
	hospitalAAdapter := infrastructure.NewHospitalAAPIAdapter(hospitalABaseURL)
	hospitalBAdapter := infrastructure.NewHospitalAAPIAdapter(hospitalBBaseURL)
	hospitalResolver := infrastructure.NewHospitalResolver(hospitalAAdapter, hospitalBAdapter)

	authService := application.NewAuthService(staffRepo, hasher, tokenProvider)
	staffService := application.NewStaffService(staffRepo, hasher)
	searchPatientUseCase := application.NewSearchPatient(hospitalResolver)

	log.Println("ประกอบร่าง Dependencies เรียบร้อย")

	// 3. Register Routes ด้วย Gin Router
	r := setupRouter(authService, &staffService, searchPatientUseCase, tokenProvider)

	// 4. Start HTTP Server
	port := getEnv("SERVER_PORT", ":8080")
	log.Printf("HTTP Server (Gin) กำลังทำงานที่พอร์ต %s ...\n", port)

	if err := r.Run(port); err != nil {
		log.Fatalf("Server ทำงานผิดพลาด: %v", err)
	}
}

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

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("ไม่สามารถเชื่อมต่อ Database ได้ (Ping failed): %v", err)
	}

	log.Println("เชื่อมต่อ Database PostgreSQL สำเร็จเรียบร้อย")
	return db
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
