package postgres

import (
	"context"
	"fmt"

	"github.com/asliddinberdiev/reception/internal/models"
	"github.com/asliddinberdiev/reception/internal/storage/repository"
	"github.com/asliddinberdiev/reception/pkg/logger"
	"github.com/pkg/errors"
)

type userRepo struct {
	db  models.DB
	log logger.Logger
}

func NewUserRepo(db models.DB, log logger.Logger) repository.UserPgI {
	return &userRepo{db: db, log: log}
}

func (r *userRepo) GetAllDoctors(ctx context.Context, req models.GetALLRequest) (*models.GetAllProfileShort, error) {
	list := &models.GetAllProfileShort{
		Data: make([]models.ProfileShort, 0),
	}

	query := `
		SELECT 
			u.id,
			u.first_name,
			u.last_name,
			u.specialty,
			u.description,
			json_agg(json_build_object(
				'week_day', uwt.week_day,
				'start_time', uwt.start_time,
				'finish_time', uwt.finish_time
			) ORDER BY uwt.week_day),
			COUNT(u.id) OVER()
		FROM "users" u
		JOIN "user_work_times" uwt ON uwt.user_id = u.id  
		JOIN "user_roles" r ON r.user_id = u.id AND r.role_id = '0195d786-5a00-71f1-be49-b7cd48c150d2'
		WHERE u.deleted_at = 0
	`

	args := []any{}
	placeholder := 1

	if req.Search != "" {
		query += fmt.Sprintf(` AND (u.first_name LIKE $%d OR u.specialty LIKE $%d)`, placeholder, placeholder+1)
		searchPattern := "%" + req.Search + "%"
		args = append(args, searchPattern, searchPattern)
		placeholder += 2
	}

	query += fmt.Sprintf(`
		GROUP BY 
			u.id,
			u.first_name,
			u.last_name,
			u.specialty,
			u.description
		ORDER BY u.id DESC
		LIMIT $%d OFFSET $%d
	`, placeholder, placeholder+1)
	args = append(args, req.Limit, (req.Page-1)*req.Limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all doctors")
	}
	defer rows.Close()

	for rows.Next() {
		var p models.ProfileShort
		if err := rows.Scan(
			&p.ID,
			&p.FirstName,
			&p.LastName,
			&p.Specialty,
			&p.Description,
			&p.WorkTime,
			&list.Total,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan doctors")
		}
		list.Data = append(list.Data, p)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to get all doctors")
	}

	return list, nil
}
