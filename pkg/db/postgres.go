package db

import (
	"context"

	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/asliddinberdiev/reception/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg *config.Config, log logger.Logger, ctx context.Context) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, cfg.GetPostgresDSN())
	if err != nil {
		log.Error("failed to connect to postgres", logger.Error(err))
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		log.Error("failed to ping postgres", logger.Error(err))
		return nil, err
	}

	return db, nil
}
