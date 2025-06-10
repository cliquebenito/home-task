package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"testowoe/internal/models"
)

type Repository struct {
	log *slog.Logger
	*pgxpool.Pool
}

type Storage interface {
	SaveStatistics(ctx context.Context, bannerID int) error
	LoadStats(ctx context.Context, params models.StatsQueryParams) ([]models.StatsResponse, error)
	CreateBanner(ctx context.Context, name string) error
}

func NewRepository(db *pgxpool.Pool, lg *slog.Logger) Storage {
	return &Repository{
		log:  lg,
		Pool: db,
	}
}

const saveuQuery = `
  INSERT INTO stats (banner_id, timestamp, count)
  VALUES ($1, $2, 1)
  ON CONFLICT (banner_id, timestamp)
  DO UPDATE SET count = stats.count + 1;`

func (r *Repository) SaveStatistics(ctx context.Context, bannerID int) error {
	now := time.Now().Truncate(time.Minute)
	log := r.log.With(
		slog.String("op", "Repository.SaveStatistics"),
		slog.Int("banner_id", bannerID),
		slog.Time("timestamp", now),
	)

	_, err := r.Exec(ctx, saveuQuery, bannerID, now)
	if err != nil {
		log.Error("insert failed", slog.Any("err", err))
		return fmt.Errorf("%w: %v", ErrSaveStatisticsFailed, err)
	}

	log.Info("statistics saved")
	return nil
}

const loadStatsQuery = `SELECT timestamp, count
		FROM stats
		WHERE banner_id = $1 AND timestamp >= $2 AND timestamp < $3
		ORDER BY timestamp;`

func (r *Repository) LoadStats(ctx context.Context, params models.StatsQueryParams) ([]models.StatsResponse, error) {
	log := r.log.With(
		slog.String("op", "Repository.LoadStats"),
		slog.Int("banner_id", params.BannerID),
	)

	rows, err := r.Query(ctx, loadStatsQuery, params.BannerID, params.From, params.To)
	if err != nil {
		log.Error("query failed", slog.Any("err", err))
		return nil, fmt.Errorf("%w: %v", ErrStatsQueryFailed, err)
	}
	defer rows.Close()

	var result []models.StatsResponse
	for rows.Next() {
		var row models.StatsResponse
		if err = rows.Scan(&row.TS, &row.V); err != nil {
			log.Error("row scan failed", slog.Any("err", err))
			return nil, fmt.Errorf("%w: %v", ErrStatsScanFailed, err)
		}
		result = append(result, row)
	}

	log.Info("stats loaded", slog.Int("rows", len(result)))
	return result, nil
}

const createBannerQuery = `INSERT INTO banners (name) VALUES ($1);`

func (r *Repository) CreateBanner(ctx context.Context, name string) error {
	log := r.log.With(
		slog.String("op", "Repository.CreateBanner"),
		slog.String("name", name),
	)
	_, err := r.Exec(ctx, createBannerQuery, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "banners_name_unique" {
				log.Warn("banner name already exists")
				return fmt.Errorf("%w: banner name already exists", ErrBannerNameExists)
			}
		}
		log.Error("insert failed", slog.Any("err", err))
		return fmt.Errorf("%w: %v", ErrCreateBannerFailed, err)
	}
	log.Info("banner created")
	return nil
}
