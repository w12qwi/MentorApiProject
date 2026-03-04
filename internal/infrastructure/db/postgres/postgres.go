package postgres

import (
	"MentorApiProject/internal/config"
	"MentorApiProject/internal/domain/models"
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"time"

	_ "github.com/jackc/pgx/v5"
)

type CalculationsRepository struct {
	db *sql.DB
}

func NewCalculationsRepository() *CalculationsRepository {
	cfg := config.LoadPostgresConfig()
	println("POSTGRES CONFIG LOADED" + cfg.Dsn())

	db, err := sql.Open("pgx", cfg.Dsn())
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &CalculationsRepository{db: db}
}

func (r *CalculationsRepository) Close() error {
	return r.db.Close()
}

func (r *CalculationsRepository) SaveCalculation(ctx context.Context, calculation models.Calculation) error {

	query, args, err := squirrel.Insert("calculations").
		Columns("id", "num_a", "num_b", "sign", "result", "created_at").
		Values(calculation.Id, calculation.NumA, calculation.NumB, calculation.Sign, calculation.Result, calculation.CreatedAt).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *CalculationsRepository) GetById(ctx context.Context, id uuid.UUID) (models.Calculation, error) {

	query, args, err := squirrel.Select("id", "num_a", "num_b", "sign", "result", "created_at").
		From("calculations").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return models.Calculation{}, err
	}

	var calculation models.Calculation
	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&calculation.Id,
			&calculation.NumA,
			&calculation.NumB,
			&calculation.Sign,
			&calculation.Result,
			&calculation.CreatedAt)
	if err != nil {
		return models.Calculation{}, err
	}

	return calculation, nil
}

func (r *CalculationsRepository) GetAllCalculations(ctx context.Context) ([]models.Calculation, error) {

	query, args, err := squirrel.Select("id", "num_a", "num_b", "sign", "result", "created_at").
		From("calculations").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.Calculation{}

	for rows.Next() {
		calculation := models.Calculation{}

		err = rows.Scan(
			&calculation.Id,
			&calculation.NumA,
			&calculation.NumB,
			&calculation.Sign,
			&calculation.Result,
			&calculation.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, calculation)
	}

	return result, nil
}

func (r *CalculationsRepository) GetCalculationsByDate(ctx context.Context, date time.Time) ([]models.Calculation, error) {

	query, args, err := squirrel.Select("id", "num_a", "num_b", "sign", "result", "created_at").
		From("calculations").
		Where(squirrel.Eq{"created_at": date}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.Calculation{}

	for rows.Next() {
		calculation := models.Calculation{}

		err = rows.Scan(
			&calculation.Id,
			&calculation.NumA,
			&calculation.NumB,
			&calculation.Sign,
			&calculation.Result,
			&calculation.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, calculation)
	}

	return result, nil
}

func (r *CalculationsRepository) GetCalculationsByDateRange(ctx context.Context, from, to time.Time) ([]models.Calculation, error) {

	query, args, err := squirrel.Select("id", "num_a", "num_b", "sign", "result", "created_at").
		From("calculations").
		Where(squirrel.GtOrEq{"created_at": from}).
		Where(squirrel.LtOrEq{"created_at": to}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.Calculation{}

	for rows.Next() {
		calculation := models.Calculation{}

		err = rows.Scan(
			&calculation.Id,
			&calculation.NumA,
			&calculation.NumB,
			&calculation.Sign,
			&calculation.Result,
			&calculation.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, calculation)
	}

	return result, nil
}
