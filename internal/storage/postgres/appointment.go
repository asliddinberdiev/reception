package postgres

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/internal/storage/repository"
	"github.com/asliddinberdiev/reception/pkg/helper"
	"github.com/asliddinberdiev/reception/pkg/logger"
	"github.com/pkg/errors"
)

type appointmentRepo struct {
	db  models.DB
	log logger.Logger
}

func NewAppointmentRepo(db models.DB, log logger.Logger) repository.AppointmentPgI {
	return &appointmentRepo{db: db, log: log}
}

func (r *appointmentRepo) Create(ctx context.Context, req models.AppointmentCreateInput) (*models.CommonGetByID, error) {
	id := helper.NewV7ID()

	query := `
		INSERT INTO appointments (
			id,
			patient_id,
			doctor_id,
			appointment_date,
			appointment_time
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	if _, err := r.db.Exec(
		ctx,
		query,
		id,
		req.UserID,
		req.DoctorID,
		req.AppointmentDate,
		req.AppointmentTime,
	); err != nil {
		return nil, errors.Wrap(err, "failed to create appointment")
	}

	return &models.CommonGetByID{ID: id}, nil
}

func (r *appointmentRepo) GetByID(ctx context.Context, req models.CommonGetByID) (*models.Appointment, error) {
	query := `
		SELECT
			a.id,
			a.appointment_date,
			a.appointment_time,
			a.status,
			a.doctor_approval,
			a.created_at,
			a.updated_at,
			d.id,
			d.first_name,
			d.last_name,
			p.id,
			p.first_name,
			p.last_name
		FROM "appointments" a
		JOIN "users" d ON d.id = a.doctor_id
		JOIN "patients" p ON p.id = a.patient_id
		WHERE a.id = $1 
	`

	var data models.Appointment
	if err := r.db.QueryRow(
		ctx,
		query,
		req.ID,
	).Scan(
		&data.ID,
		&data.AppointmentDate,
		&data.AppointmentTime,
		&data.Status,
		&data.DoctorApproval,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.Doctor.ID,
		&data.Doctor.Firstname,
		&data.Doctor.Lastname,
		&data.Patient.ID,
		&data.Patient.Firstname,
		&data.Patient.Lastname,
	); err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *appointmentRepo) GetByRangeTime(ctx context.Context, req models.AppointmentRangeTime) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1 
            FROM appointments 
            WHERE doctor_id = $1 
            AND appointment_date = $2 
            AND NOT (appointment_time >= $4 OR appointment_time + duration <= $3)
            AND status NOT IN ('cancelled', 'rejected')
        )
    `

	var exists bool
	err := r.db.QueryRow(
		ctx,
		query,
		req.DoctorID,
		req.Date,
		req.StartTime,
		req.EndTime,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return !exists, nil
}

func (r *appointmentRepo) UpdateStatus(ctx context.Context, req models.CommonGetByID) error {
	query := `
		UPDATE "appointments"
		SET status = 'cancelled'
		WHERE id = $1 
	`

	if _, err := r.db.Exec(ctx, query, req.ID); err != nil {
		return errors.Wrap(err, "failed to status update")
	}

	return nil
}
