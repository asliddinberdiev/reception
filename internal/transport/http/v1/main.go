package v1

import (
	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/asliddinberdiev/reception/internal/service"
	"github.com/asliddinberdiev/reception/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	log     logger.Logger
	cfg     *config.Config
	valid   *validator.Validate
	service *service.Service
}

func NewHandler(log logger.Logger, cfg *config.Config, services *service.Service) *Handler {
	return &Handler{
		log:     log,
		cfg:     cfg,
		valid:   validator.New(),
		service: services,
	}
}

func (h *Handler) Init(router fiber.Router) {
	_ = router.Group("/v1")
}
