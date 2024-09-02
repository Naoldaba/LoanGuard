package main

import (
	"LoanGuard/internal/delivery/controllers"
	"LoanGuard/internal/delivery/routers"
	"LoanGuard/internal/infrastructures/database"
	"LoanGuard/internal/infrastructures/middlewares"
	"LoanGuard/internal/infrastructures/services"
	"LoanGuard/internal/infrastructures/services/email_service"
	"LoanGuard/internal/usecases"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"LoanGuard/internal/repository/implementations"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbName := os.Getenv("DB_NAME")
	accessSecretKey := os.Getenv("ACCESS_SECRET_KEY")
	refreshSecretKey := os.Getenv("REFRESH_SECRET_KEY")
	verificationSecretKey := os.Getenv("VERIFICATION_SECRET_KEY")

	smtpPortStr := os.Getenv("SMTP_PORT")
	userName := os.Getenv("USERNAME")
	smtpHost := os.Getenv("SMTP_HOST")
	passWord := os.Getenv("PASSWORD")
	cachePort := os.Getenv("CACHE_PORT")
	cacheHost := os.Getenv("CACHE_HOST")
	
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}

	dbClient, err := database.NewMongoDB(context.Background(), dbName)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB!", dbClient.Db.Name())

	//services
	emailSvc := email_service.NewEmailService(smtpHost, smtpPort, userName, passWord)
	passSvc := services.NewPasswordService()
	validationSvc := services.NewValidationService()
	jwtSvc := services.NewJWTService(accessSecretKey, refreshSecretKey, verificationSecretKey)
	cacheSvc := services.NewCacheService(cacheHost + ":" + cachePort , "", 0)
	cloudSvc := services.NewCloudinaryService(os.Getenv("CLOUDINARY_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"), os.Getenv("CLOUDINARY_UPLOAD_FOLDER"),)

	//repo implementations
	userRepo := implementations.NewMongoUserRepository(dbClient.Db, cacheSvc)
	otpRepo := implementations.NewMongoOtpRepository(dbClient.Db)
	loanRepo := implementations.NewMongoLoanRepository(dbClient.Db)
	logRepo := implementations.NewMongoLogRepository(dbClient.Db)

	//middlewares
	authMiddleware := middlewares.NewAuthMiddleware(jwtSvc, cacheSvc)


	//usecases
	userUsecase := usecases.NewUserUsecase(userRepo, passSvc, validationSvc, emailSvc, jwtSvc, cloudSvc, "http://localhost:8080")
	otpUsecase := usecases.NewOtpUseCase(otpRepo, userRepo, emailSvc, passSvc, "http://localhost:8080", validationSvc)
	loanUsecase := usecases.NewLoanUsecase(loanRepo, logRepo)
	adminUsecase := usecases.NewAdminUsecase(loanRepo, logRepo)	

	// controllers
	userController := controllers.NewUserController(userUsecase)
	otpController := controllers.NewOTPController(otpUsecase)
	loanController := controllers.NewLoanController(loanUsecase)
	adminController := controllers.NewAdminController(userUsecase, adminUsecase)
	

	//gin engine initialization
	router := gin.New()
	router.Use(gin.Logger())


	// routers
	routers.CreateUserRouter(router, userController, otpController, authMiddleware)
	routers.CreateAdminRouter(router, adminController, authMiddleware)
	routers.CreateLoanRouter(router, loanController, authMiddleware)

	if err := router.Run(":" + os.Getenv("PORT")); err!= nil{
		log.Fatal(err)
	}
}
