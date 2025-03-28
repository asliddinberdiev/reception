package postgres

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/internal/storage/repository"
	"github.com/asliddinberdiev/reception/pkg/helper"
	"github.com/asliddinberdiev/reception/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type patientRepo struct {
	db  models.DB
	log logger.Logger
}

func NewPatientRepo(db models.DB, log logger.Logger) repository.PatientPgI {
	return &patientRepo{db: db, log: log}
}

func (r *patientRepo) Create(ctx context.Context, inp models.PatientCreateInput) (*models.CommonGetByID, error) {
	id := helper.NewV7ID()
	query := `
		INSERT INTO patients (
			id,
			phone_number,
			first_name,
			last_name,
			gender,
			password
		) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	if _, err := r.db.Exec(ctx, query, id, inp.PhoneNumber, inp.Firstname, inp.Lastname, inp.Gender, inp.Password); err != nil {
		return nil, errors.Wrap(err, "failed to create patient")
	}

	return &models.CommonGetByID{ID: id}, nil
}

func (r *patientRepo) GetByPhone(ctx context.Context, phoneNumber string) (*models.Patient, error) {
	var res models.Patient

	query := `
		SELECT
			id,
			phone_number,
			is_validated,
			first_name,
			last_name,
			gender,
			password,
			created_at,
			updated_at
		FROM patients
		WHERE phone_number = $1 AND deleted_at = 0
	`
	if err := r.db.QueryRow(ctx, query, phoneNumber).
		Scan(
			&res.ID,
			&res.PhoneNumber,
			&res.IsVerified,
			&res.Firstname,
			&res.Lastname,
			&res.Gender,
			&res.Password,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
		return nil, errors.Wrap(err, "failed to scan")
	}

	return &res, nil
}

func (r *patientRepo) SetAsVerified(ctx context.Context, phoneNumber string) (*models.Patient, error) {
	query := `
		UPDATE "patients"
		SET is_validated=TRUE
		WHERE phone_number = $1
	`
	res, err := r.db.Exec(ctx, query, phoneNumber)
	if err != nil {
		return nil, err
	}

	if i := res.RowsAffected(); i == 0 {
		return nil, pgx.ErrNoRows
	}

	return r.GetByPhone(ctx, phoneNumber)
}
