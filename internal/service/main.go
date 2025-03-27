package service

import (
	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/asliddinberdiev/reception/internal/storage"
)

type Service struct{}

func NewService(strg storage.StoragePG, cfg *config.Config) *Service {
	return &Service{}
}
