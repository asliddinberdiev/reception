package service

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/internal/storage"
	"github.com/asliddinberdiev/reception/pkg/helper"
	"github.com/jackc/pgx/v5"
)

type patientService struct {
	strg storage.StoragePG
	cfg  *config.Config
}

func NewPatientService(strg storage.StoragePG, cfg *config.Config) *patientService {
	return &patientService{strg: strg, cfg: cfg}
}

func (s *patientService) Create(ctx context.Context, inp models.PatientCreateInput) (*models.CommonGetByID, error) {
	return s.strg.Patient().Create(ctx, inp)
}

func (s *patientService) GetPatientByPhone(ctx context.Context, phoneNumber string) (*models.Patient, error) {
	res, err := s.strg.Patient().GetByPhone(ctx, phoneNumber)
	if err != nil {
		if helper.ErrorIs(err, pgx.ErrNoRows.Error()) {
			return &models.Patient{}, nil
		}

		return nil, err
	}
	return res, nil
}

func (s *patientService) SetAsVerified(ctx context.Context, phoneNumber string) (*models.Patient, error) {
	tx, err := s.strg.WithTransaction(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()
	res, err := tx.Patient().SetAsVerified(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}

	return res, nil
}
