package v1

import (
	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/internal/service"
	"github.com/asliddinberdiev/reception/pkg/helper"
	"github.com/asliddinberdiev/reception/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	log   logger.Logger
	cfg   *config.Config
	valid *validator.Validate
	sve   *service.Service
}

func NewHandler(log logger.Logger, cfg *config.Config, services *service.Service) *Handler {
	return &Handler{
		log:   log,
		cfg:   cfg,
		valid: validator.New(),
		sve:   services,
	}
}

func (h *Handler) Init(router fiber.Router) {
	v1 := router.Group("/v1")
	{
		h.initUserRoutes(v1)
		h.initPatientRoutes(v1)
		h.initAppointmentRoutes(v1)
	}
}

func (h *Handler) MakeRequestSearch(c *fiber.Ctx) *models.CommonGetALL {
	var query models.CommonGetALL

	query.Limit = uint32(helper.ParseInt(c.Query("limit"), 10))
	query.Page = uint32(helper.ParseInt(c.Query("page"), 1))
	query.Search = c.Query("search", "")

	return &query
}
