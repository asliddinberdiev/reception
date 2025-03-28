package service

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/internal/storage"
)

type userService struct {
	strg storage.StoragePG
	cfg  *config.Config
}

func NewUserService(strg storage.StoragePG, cfg *config.Config) *userService {
	return &userService{strg: strg, cfg: cfg}
}

func (s *userService) GetAllDoctors(ctx context.Context, req models.GetALLRequest) (*models.GetAllProfileShort, error) {
	return s.strg.User().GetAllDoctors(ctx, req)
}
