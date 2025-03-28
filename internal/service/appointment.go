package service

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/internal/storage"
)

type appointmentService struct {
	strg storage.StoragePG
	cfg  *config.Config
}

func NewAppointmentService(strg storage.StoragePG, cfg *config.Config) *appointmentService {
	return &appointmentService{strg: strg, cfg: cfg}
}

func (r *appointmentService) Create(ctx context.Context, req models.AppointmentCreateInput) (*models.CommonGetByIDResponse, error) {
	return nil, nil
}
