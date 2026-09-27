package main

import (
	"database/sql"
	"log"
	"time"

	"Hospital-Middleware/internal/application"
	"Hospital-Middleware/internal/config"
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
	if err := r.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

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
	cfg := config.Load()

	// 1. Initialize Database
	db := initDB(cfg)
	defer db.Close()

	// 2. Setup Dependencies
	staffRepo := infrastructure.NewPostgresStaffRepository(db)
	hasher := infrastructure.NewBcryptHasher()
	tokenProvider := infrastructure.NewJWTTokenProvider(cfg.JWTSecret)

	hospitalAAdapter := infrastructure.NewHospitalAAPIAdapter(cfg.HospitalABaseURL)
	hospitalBAdapter := infrastructure.NewHospitalAAPIAdapter(cfg.HospitalBBaseURL)
	hospitalResolver := infrastructure.NewHospitalResolver(hospitalAAdapter, hospitalBAdapter)

	authService := application.NewAuthService(staffRepo, hasher, tokenProvider)
	staffService := application.NewStaffService(staffRepo, hasher)
	searchPatientUseCase := application.NewSearchPatient(hospitalResolver)

	log.Println("ประกอบร่าง Dependencies เรียบร้อย")

	// 3. Register Routes ด้วย Gin Router
	r := setupRouter(authService, &staffService, searchPatientUseCase, tokenProvider)

	// 4. Start HTTP Server
	log.Printf("HTTP Server (Gin) กำลังทำงานที่พอร์ต %s ...\n", cfg.ServerPort)

	if err := r.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Server ทำงานผิดพลาด: %v", err)
	}
}

func initDB(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.DSN())
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
