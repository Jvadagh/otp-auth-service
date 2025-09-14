package main

// @title OTP Auth Service API
// @version 1.0
// @description This is a sample OTP-based login and registration service.
// @host localhost:8080
// @BasePath /

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger" // swagger handler
	"github.com/jvadagh/otp-auth-service/internal/api"
	_ "github.com/jvadagh/otp-auth-service/internal/docs" // swagger docs
	"github.com/jvadagh/otp-auth-service/internal/middleware"
	"github.com/jvadagh/otp-auth-service/internal/repository"
	"github.com/jvadagh/otp-auth-service/internal/service"
	"github.com/jvadagh/otp-auth-service/pkg/config"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	db := repository.NewPostgres(cfg.PostgresDSN)
	rdb := repository.NewRedis(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)
	otpService := service.NewOTPService(rdb)
	authHandler := api.NewAuthHandler(db, otpService, cfg.JWTSecret)
	userHandler := api.NewUserHandler(repository.NewUserRepo(db))
	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Post("/auth/request-otp", authHandler.RequestOTP)
	app.Post("/auth/verify-otp", authHandler.VerifyOTP)
	users := app.Group("/users", middleware.JWTMiddleware(cfg.JWTSecret))
	users.Get("/", userHandler.ListUsers)
	users.Get("/:id", userHandler.GetUser)
	log.Println("Starting server on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
