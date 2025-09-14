package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jvadagh/otp-auth-service/internal/repository"
	"github.com/jvadagh/otp-auth-service/internal/service"
	"github.com/jvadagh/otp-auth-service/pkg/utils"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type AuthHandler struct {
	userRepo   *repository.UserRepo
	otpService *service.OTPService
	jwtSecret  string
}

func NewAuthHandler(db *gorm.DB, otpService *service.OTPService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:   repository.NewUserRepo(db),
		otpService: otpService,
		jwtSecret:  jwtSecret,
	}
}

// RequestOTP godoc
// @Summary Request OTP
// @Description Generate a new OTP for the given phone number.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RequestOTPRequest true "Phone number"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/request-otp [post]
func (h *AuthHandler) RequestOTP(c *fiber.Ctx) error {
	var body struct {
		Phone string `json:"phone"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid body")
	}
	normalized, err := utils.NormalizePhone(body.Phone)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid phone number")
	}
	body.Phone = normalized
	code, err := h.otpService.GenerateOTP(body.Phone)
	if err != nil {
		if err.Error() == "rate limit exceeded" {
			return fiber.NewError(fiber.StatusTooManyRequests, "too many otp requests, try later")
		}
		return fiber.NewError(http.StatusInternalServerError, "failed to generate otp")
	}
	log.Printf("Generated OTP for %s: %s", body.Phone, code)
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "OTP sent"})
}

// VerifyOTP godoc
// @Summary Verify OTP
// @Description Verify the OTP for a phone number. Registers a new user if not existing, otherwise logs in.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "Phone number and OTP"
// @Success 200 {object} map[string]string "JWT token"
// @Failure 400 {object} map[string]string "Invalid body"
// @Failure 401 {object} map[string]string "Invalid or expired OTP"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var body struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid body")
	}
	normalized, err := utils.NormalizePhone(body.Phone)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid phone number")
	}
	body.Phone = normalized
	valid, err := h.otpService.ValidateOTP(body.Phone, body.OTP)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "error validating otp")
	}
	if !valid {
		return fiber.NewError(http.StatusUnauthorized, "invalid or expired otp")
	}
	user, err := h.userRepo.GetByPhone(body.Phone)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user, err = h.userRepo.Create(body.Phone)
			if err != nil {
				return fiber.NewError(http.StatusInternalServerError, "failed to create user")
			}
		} else {
			return fiber.NewError(http.StatusInternalServerError, "db error")
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed to sign token")
	}
	return c.JSON(fiber.Map{"token": signed})
}
