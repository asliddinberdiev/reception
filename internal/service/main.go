package service

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/internal/storage"
)

type User interface {
	GetAllDoctors(ctx context.Context, req models.CommonGetALL) (*models.GetAllProfileShort, error)
}

type Patient interface {
	Create(ctx context.Context, inp models.PatientCreateInput) (*models.CommonGetByID, error)
	GetPatientByPhone(ctx context.Context, phoneNumber string) (*models.Patient, error)
	SetAsVerified(ctx context.Context, phoneNumber string) (*models.Patient, error)
}

type Appointment interface {
	Create(ctx context.Context, req models.AppointmentCreateInput) (*models.CommonGetByID, error)
	GetByID(ctx context.Context, req models.CommonGetByID) (*models.Appointment, error)
	GetByRangeTime(ctx context.Context, req models.AppointmentRangeTime) (bool, error)
	UpdateStatus(ctx context.Context, req models.CommonGetByID) error
}

type Service struct {
	User
	Patient
	Appointment
}

func NewService(strg storage.StoragePG, cfg *config.Config) *Service {
	return &Service{
		User:        NewUserService(strg, cfg),
		Patient:     NewPatientService(strg, cfg),
		Appointment: NewAppointmentService(strg, cfg),
	}
}
