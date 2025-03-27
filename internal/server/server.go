package server

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App
	cfg *config.Config
}

func NewServer(cfg *config.Config, app *fiber.App) *Server {
	return &Server{app: app, cfg: cfg}
}

func (s *Server) Run() error {
	return s.app.Listen(s.cfg.GetAddress())
}

func (s *Server) Stop(ctx context.Context) error {
	return s.app.Shutdown()
}
