package v1

import (
	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initUserRoutes(router fiber.Router) {
	// public
	router.Get("doctors", h.getAllDoctors)
}

func (h *Handler) getAllDoctors(c *fiber.Ctx) error {
	params := h.MakeRequestSearch(c)

	res, err := h.sve.User.GetAllDoctors(
		c.Context(),
		models.CommonGetALL{
			Limit:  params.Limit,
			Page:   params.Page,
			Search: params.Search,
		},
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
