package postgres

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/internal/storage/repository"
	"github.com/asliddinberdiev/reception/pkg/logger"
)

type appointmentRepo struct {
	db  models.DB
	log logger.Logger
}

func NewAppointmentRepo(db models.DB, log logger.Logger) repository.AppointmentPgI {
	return &appointmentRepo{db: db, log: log}
}

func (r *appointmentRepo) Create(ctx context.Context, req models.AppointmentCreateInput) (*models.CommonGetByIDResponse, error) {
	return nil, nil
}
