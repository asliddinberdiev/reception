package http

import (
	"encoding/json"
	"fmt"

	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/asliddinberdiev/reception/internal/service"
	v1 "github.com/asliddinberdiev/reception/internal/transport/http/v1"
	"github.com/asliddinberdiev/reception/pkg/helper"
	"github.com/asliddinberdiev/reception/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5"
)

var Validate = validator.New()

type Handler struct {
	log      logger.Logger
	cfg      *config.Config
	services *service.Service
}

func NewHandler(log logger.Logger, cfg *config.Config, services *service.Service) *Handler {
	return &Handler{
		log:      log,
		cfg:      cfg,
		services: services,
	}
}

// @title Reception
// @version 1.0
// @description API for Reception Application

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func (h *Handler) Init(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: h.errorHandler,
		JSONEncoder:  h.jsonEncoder,
		AppName:      cfg.App.ServiceName,
	})

	app.Use(flogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Authorization, Origin, Content-Type, Accept, Content-Language, Accept-Language, Access-Control-Allow-Headers",
	}))

	app.Get("/health", h.health)

	app.Get("/swagger/*any", swagger.HandlerDefault)

	app.Get("/swagger/*", swagger.New(swagger.Config{DeepLinking: true}))

	h.initAPI(app)

	return app
}

func (h *Handler) initAPI(app *fiber.App) {
	handlerV1 := v1.NewHandler(h.log, h.cfg, h.services)
	api := app.Group("api")
	{
		handlerV1.Init(api)
	}
}

func (h *Handler) errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "something went wrong"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	if code == fiber.StatusBadRequest {
		return ctx.Status(code).JSON(fiber.Map{"error": message})
	}

	if code == fiber.StatusNotFound {
		return ctx.Status(code).JSON(fiber.Map{"error": message})
	}

	if helper.ErrorIs(err, pgx.ErrNoRows.Error()) {
		code = fiber.StatusNotFound
		message = "not found"
	}

	if code == fiber.StatusUnauthorized {
		return ctx.Status(code).JSON(fiber.Map{"error": message})
	}

	if h.cfg.App.Environment == "prod" {
		h.log.Error("srv", logger.Any("error", message))
		return ctx.Status(code).JSON(fiber.Map{"error": message})
	}

	return ctx.Status(code).JSON(fiber.Map{"error": err.Error()})
}

func (h *Handler) jsonEncoder(v any) ([]byte, error) {
	r, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("error while json Marshaling: %w", err)
	}

	return r, nil
}

func (h *Handler) health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": true})
}
