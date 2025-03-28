package storage

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/storage/postgres"
	"github.com/asliddinberdiev/reception/internal/storage/repository"
	"github.com/asliddinberdiev/reception/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repos struct {
	userRepo        repository.UserPgI
	patientRepo     repository.PatientPgI
	appointmentRepo repository.AppointmentPgI
}

type repoIs interface {
	User() repository.UserPgI
	Patient() repository.PatientPgI
	Appointment() repository.AppointmentPgI
}

type storage struct {
	db  *pgxpool.Pool
	log logger.Logger
	repos
}

type storageTr struct {
	tx   pgx.Tx
	conn *pgxpool.Conn
	log  logger.Logger
	repos
}

type StorageTrI interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	repoIs
}

type StoragePG interface {
	WithTransaction(ctx context.Context) (StorageTrI, error)
	repoIs
}

func NewStoragePg(db *pgxpool.Pool, log logger.Logger) StoragePG {
	return &storage{
		db:  db,
		log: log,
		repos: repos{
			userRepo:        postgres.NewUserRepo(db, log),
			patientRepo:     postgres.NewPatientRepo(db, log),
			appointmentRepo: postgres.NewAppointmentRepo(db, log),
		},
	}
}

func (s *storage) WithTransaction(ctx context.Context) (StorageTrI, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		s.log.Error("failed to acquire db connection", logger.Error(err))
		return nil, err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		conn.Release()
		s.log.Error("failed to begin db transaction", logger.Error(err))
		return nil, err
	}

	return &storageTr{
		tx:   tx,
		conn: conn,
		repos: repos{
			userRepo:        postgres.NewUserRepo(tx, s.log),
			patientRepo:     postgres.NewPatientRepo(tx, s.log),
			appointmentRepo: postgres.NewAppointmentRepo(tx, s.log),
		},
	}, nil
}

func (s *storageTr) Commit(ctx context.Context) error {
	if err := s.tx.Commit(ctx); err != nil {
		s.log.Error("failed to commit tx", logger.Error(err))
		return err
	}

	s.conn.Release()
	return nil
}

func (s *storageTr) Rollback(ctx context.Context) error {
	if err := s.tx.Rollback(ctx); err != nil {
		s.log.Error("failed to rollback tx", logger.Error(err))
		return err
	}

	s.conn.Release()
	return nil
}

func (s *repos) User() repository.UserPgI {
	return s.userRepo
}

func (s *repos) Patient() repository.PatientPgI {
	return s.patientRepo
}

func (s *repos) Appointment() repository.AppointmentPgI {
	return s.appointmentRepo
}
