package repository

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/models"
)

type UserPgI interface {
	GetAllDoctors(ctx context.Context, req models.GetALLRequest) (*models.GetAllProfileShort, error)
}
