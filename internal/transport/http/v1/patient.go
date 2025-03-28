package v1

import (
	"fmt"
	"time"

	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/pkg/auth"
	"github.com/asliddinberdiev/reception/pkg/helper"
	"github.com/asliddinberdiev/reception/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) initPatientRoutes(router fiber.Router) {
	patient := router.Group("patients")
	{
		patient.Post("register", h.registerPatient)
		patient.Post("login", h.loginPatient)
		patient.Post("verify", h.middlewareOTP, h.verifyOtp)
	}
}

func (h *Handler) registerPatient(c *fiber.Ctx) error {
	var req models.PatientCreateInput

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	if err := h.valid.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	patient, err := h.sve.Patient.GetPatientByPhone(c.Context(), req.PhoneNumber)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if patient.IsVerified {
		return fiber.NewError(fiber.StatusOK, "already exists")
	}

	if len(patient.ID) == 0 && !patient.IsVerified {
		hashedPassword, err := helper.PasswordHash(req.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		req.Password = hashedPassword

		if _, err := h.sve.Patient.Create(c.Context(), req); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	otp := helper.RandNumberStringRunes(h.cfg.Auth.CodeLength)

	if h.cfg.App.Environment == "dev" {
		otp = "12345678"
	}

	msg := fmt.Sprintf("Your secret key: %s", otp)
	h.log.Info("OTP message", logger.Any("msg", msg))

	if h.cfg.App.Environment != "dev" {
		// send sms func
	}

	claims := models.OtpClaims{
		PhoneNumber:    req.PhoneNumber,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Duration(h.cfg.Auth.OtpTTL)).Unix()},
	}

	token, err := auth.GenerateToken(claims, otp)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func (h *Handler) loginPatient(c *fiber.Ctx) error {
	var req models.LoginInput

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	if err := h.valid.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	patient, err := h.sve.Patient.GetPatientByPhone(c.Context(), req.PhoneNumber)
	if err != nil {
		if helper.ErrorIs(err, pgx.ErrNoRows.Error()) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !helper.PasswordCompare(patient.Password, req.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	accessClaims := models.UserClaims{
		ID:             patient.ID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Duration(h.cfg.Auth.AccessTTL)).Unix()},
	}
	accessToken, err := auth.GenerateToken(accessClaims, h.cfg.Auth.SecretKey)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshClaims := models.UserClaims{
		ID:             patient.ID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Duration(h.cfg.Auth.RefreshTTL)).Unix()},
	}
	refreshToken, err := auth.GenerateToken(refreshClaims, h.cfg.Auth.SecretKey)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	tokenResponse := models.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": tokenResponse})
}
