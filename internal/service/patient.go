package service

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/internal/storage"
)

type patientService struct {
	strg storage.StoragePG
	cfg  *config.Config
}

func NewPatientService(strg storage.StoragePG, cfg *config.Config) *patientService {
	return &patientService{strg: strg, cfg: cfg}
}

func (s *patientService) Create(ctx context.Context, inp models.PatientCreateInput) (*models.CommonGetByIDResponse, error) {
	return s.strg.Patient().Create(ctx, inp)
}

func (s *patientService) HasPatientByPhone(ctx context.Context, phoneNumber string) (bool, error) {
	if _, err := s.strg.Patient().GetByPhone(ctx, phoneNumber); err != nil {
		return false, err
	}

	return true, nil
}

func (s *patientService) GetPatientByPhone(ctx context.Context, phoneNumber string) (*models.Patient, error) {
	return s.strg.Patient().GetByPhone(ctx, phoneNumber)
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
