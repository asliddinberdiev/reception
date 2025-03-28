package service

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/internal/storage"
)

type User interface {
	GetAllDoctors(ctx context.Context, req models.GetALLRequest) (*models.GetAllProfileShort, error)
}

type Service struct {
	User
}

func NewService(strg storage.StoragePG, cfg *config.Config) *Service {
	return &Service{
		User: NewUserService(strg, cfg),
	}
}
