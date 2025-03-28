package postgres

import (
	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/internal/storage/repository"
	"github.com/asliddinberdiev/reception/pkg/logger"
)

type appointmentRepo struct {
	db  models.DB
	log logger.Logger
}

func NewAppointmentRepo(db models.DB, log logger.Logger) repository.AppointmentPgI {
	return &userRepo{db: db, log: log}
}
