package repository

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/models"
)

type PatientPgI interface {
	Create(ctx context.Context, inp models.PatientCreateInput) (*models.CommonGetByID, error)
	GetByPhone(ctx context.Context, phoneNumber string) (*models.Patient, error)
	SetAsVerified(ctx context.Context, phoneNumber string) (*models.Patient, error)
}
