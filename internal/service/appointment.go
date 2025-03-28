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

func (s *appointmentService) Create(ctx context.Context, req models.AppointmentCreateInput) (*models.CommonGetByID, error) {
	return s.strg.Appointment().Create(ctx, req)
}

func (s *appointmentService) GetByID(ctx context.Context, req models.CommonGetByID) (*models.Appointment, error) {
	return s.strg.Appointment().GetByID(ctx, req)
}

func (s *appointmentService) GetByRangeTime(ctx context.Context, req models.AppointmentRangeTime) (bool, error) {
	return s.strg.Appointment().GetByRangeTime(ctx, req)
}
