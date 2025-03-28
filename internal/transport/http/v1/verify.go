package v1

import (
	"time"

	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/pkg/auth"
	"github.com/asliddinberdiev/reception/pkg/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func (h *Handler) verifyOtp(c *fiber.Ctx) error {
	var (
		claims models.OtpClaims
		req    models.VerifyInput
	)

	token := c.Get("Token", "")
	if token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "token not found")
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	if err := h.valid.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := auth.ParseToken(token, req.Otp, &claims); err != nil {
		if helper.ErrorIs(err, "token is expired") {
			return fiber.NewError(fiber.StatusBadRequest, "token is expired")
		}
		return fiber.NewError(fiber.StatusBadRequest, "invalid otp")
	}

	patient, err := h.sve.Patient.SetAsVerified(c.Context(), claims.PhoneNumber)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	accessClaims := models.PatientClaims{
		ID:             patient.ID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Duration(h.cfg.Auth.AccessTTL)).Unix()},
	}
	accessToken, err := auth.GenerateToken(accessClaims, h.cfg.Auth.SecretKey)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshClaims := models.PatientClaims{
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
