package service

import (
	"context"

	"testowoe/internal/models"
	"testowoe/internal/repository"
)

type BannerService struct {
	repo repository.Storage
}

func NewBannerService(r repository.Storage) *BannerService {
	return &BannerService{
		repo: r,
	}
}

func (r *BannerService) SaveStatistics(ctx context.Context, bannerID int) error {
	return r.repo.SaveStatistics(ctx, bannerID)
}

func (r *BannerService) CreateBanner(ctx context.Context, name string) error {
	return r.repo.CreateBanner(ctx, name)
}
func (r *BannerService) LoadStats(ctx context.Context, opts models.StatsQueryParams) ([]models.StatsResponse, error) {
	return r.repo.LoadStats(ctx, opts)
}
