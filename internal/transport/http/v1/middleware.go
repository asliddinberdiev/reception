package v1

import (
	"errors"
	"strings"

	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/pkg/auth"
	"github.com/asliddinberdiev/reception/pkg/helper"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) middlewareJWT(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	tks := strings.Split(token, " ")

	if len(tks) < 2 || tks[0] != "Bearer" {
		return fiber.NewError(fiber.StatusForbidden, "invalid token")
	}

	var claims models.UserClaims
	err := auth.ParseToken(tks[1], h.cfg.Auth.SecretKey, &claims)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	c.Request().Header.Set("UserID", claims.ID)

	return c.Next()
}

func (h *Handler) otpMiddleware(c *fiber.Ctx) error {

	token := c.Get("Authorization")

	tks := strings.Split(token, " ")

	if len(tks) < 2 || tks[0] != "Bearer" {
		return fiber.NewError(fiber.StatusForbidden, "invalid token")
	}

	c.Request().Header.Add("Token", tks[1])

	return c.Next()
}

func (h *Handler) MwGetUserID(ctx *fiber.Ctx) (string, error) {
	userID := ctx.Get("UserID")
	if userID == "" {
		return "", errors.New("user id not found")
	}
	return userID, nil
}

func (h *Handler) MwGetToken(ctx *fiber.Ctx) (string, error) {
	token := ctx.Get("Token")
	if token == "" {
		return "", errors.New("token not found")
	}
	return token, nil
}

func MakeRequestSearch(c *fiber.Ctx) *models.CommonGetALL {
	var query models.CommonGetALL

	query.Limit = uint32(helper.ParseInt(c.Query("limit"), 10))
	query.Page = uint32(helper.ParseInt(c.Query("page"), 1))
	query.Search = c.Query("search", "")

	return &query
}
