package repository

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/models"
)

type AppointmentPgI interface {
	Create(ctx context.Context, req models.AppointmentCreateInput) (*models.CommonGetByIDResponse, error)
}
