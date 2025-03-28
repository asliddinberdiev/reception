package repository

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/models"
)

type AppointmentPgI interface {
	Create(ctx context.Context, req models.AppointmentCreateInput) (*models.CommonGetByID, error)
	GetByID(ctx context.Context, req models.CommonGetByID) (*models.Appointment, error)
	GetByRangeTime(ctx context.Context, req models.AppointmentRangeTime) (bool, error)
}
