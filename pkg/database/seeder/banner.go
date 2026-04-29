package seeder

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"go.uber.org/zap"
)

type bannerSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewBannerSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *bannerSeeder {
	return &bannerSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *bannerSeeder) Seed() error {
	banners := []db.CreateBannerParams{
		{
			Name:      "Banner 1",
			StartDate: toDate(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:   toDate(time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)),
			StartTime: toTime(parseTime("08:00")),
			EndTime:   toTime(parseTime("16:00")),
			IsActive:  toBoolPtr(true),
		},
		{
			Name:      "Banner 2",
			StartDate: toDate(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:   toDate(time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)),
			StartTime: toTime(parseTime("09:00")),
			EndTime:   toTime(parseTime("17:00")),
			IsActive:  toBoolPtr(true),
		},
		{
			Name:      "Banner 3",
			StartDate: toDate(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:   toDate(time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC)),
			StartTime: toTime(parseTime("10:00")),
			EndTime:   toTime(parseTime("18:00")),
			IsActive:  toBoolPtr(false),
		},
		{
			Name:      "Banner 4",
			StartDate: toDate(time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:   toDate(time.Date(2023, 4, 30, 0, 0, 0, 0, time.UTC)),
			StartTime: toTime(parseTime("07:00")),
			EndTime:   toTime(parseTime("15:00")),
			IsActive:  toBoolPtr(true),
		},
		{
			Name:      "Banner 5",
			StartDate: toDate(time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:   toDate(time.Date(2023, 5, 31, 0, 0, 0, 0, time.UTC)),
			StartTime: toTime(parseTime("06:00")),
			EndTime:   toTime(parseTime("14:00")),
			IsActive:  toBoolPtr(false),
		},
		{
			Name:      "Banner 6",
			StartDate: toDate(time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:   toDate(time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC)),
			StartTime: toTime(parseTime("12:00")),
			EndTime:   toTime(parseTime("20:00")),
			IsActive:  toBoolPtr(true),
		},
		{
			Name:      "Banner 7",
			StartDate: toDate(time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:   toDate(time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC)),
			StartTime: toTime(parseTime("08:30")),
			EndTime:   toTime(parseTime("16:30")),
			IsActive:  toBoolPtr(true),
		},
		{
			Name:      "Banner 8",
			StartDate: toDate(time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:   toDate(time.Date(2023, 8, 31, 0, 0, 0, 0, time.UTC)),
			StartTime: toTime(parseTime("09:30")),
			EndTime:   toTime(parseTime("17:30")),
			IsActive:  toBoolPtr(false),
		},
	}

	for _, banner := range banners {
		if _, err := r.db.CreateBanner(r.ctx, banner); err != nil {
			r.logger.Error("Failed to insert banner", zap.Error(err))
			return err
		}
	}

	r.logger.Info("banner successfully seeded")
	return nil
}

func parseTime(t string) time.Time {
	parsed, _ := time.Parse("15:04", t)
	return parsed
}
