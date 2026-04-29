package seeder

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"go.uber.org/zap"
)

type merchantSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *merchantSeeder {
	return &merchantSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *merchantSeeder) Seed() error {
	merchants := []db.CreateMerchantParams{
		{
			UserID:       1,
			Name:         "Elektronik Store",
			Description:  toStringPtr("Toko elektronik terpercaya dengan berbagai produk gadget dan aksesoris."),
			Address:      toStringPtr("Jl. Teknologi No.1, Jakarta"),
			ContactEmail: toStringPtr("support@elektronikstore.com"),
			ContactPhone: toStringPtr("081234567890"),
			Status:       "active",
		},
		{
			UserID:       2,
			Name:         "Kecantikan Sehat",
			Description:  toStringPtr("Produk skincare dan kesehatan pilihan."),
			Address:      toStringPtr("Jl. Kesehatan No.5, Bandung"),
			ContactEmail: toStringPtr("cs@kecantikansehat.com"),
			ContactPhone: toStringPtr("082345678901"),
			Status:       "active",
		},
		{
			UserID:       3,
			Name:         "Rumah Indah",
			Description:  toStringPtr("Peralatan rumah tangga berkualitas dan estetik."),
			Address:      toStringPtr("Jl. Rumah No.12, Surabaya"),
			ContactEmail: toStringPtr("info@rumahindah.com"),
			ContactPhone: toStringPtr("083456789012"),
			Status:       "active",
		},
		{
			UserID:       4,
			Name:         "Mom & Baby Care",
			Description:  toStringPtr("Semua kebutuhan ibu dan bayi ada di sini."),
			Address:      toStringPtr("Jl. Keluarga No.7, Depok"),
			ContactEmail: toStringPtr("support@momandbaby.com"),
			ContactPhone: toStringPtr("084567890123"),
			Status:       "active",
		},
		{
			UserID:       5,
			Name:         "Sport Zone",
			Description:  toStringPtr("Perlengkapan olahraga dan outdoor terlengkap."),
			Address:      toStringPtr("Jl. Atletik No.3, Yogyakarta"),
			ContactEmail: toStringPtr("halo@sportzone.com"),
			ContactPhone: toStringPtr("085678901234"),
			Status:       "active",
		},
		{
			UserID:       6,
			Name:         "Fresh Mart",
			Description:  toStringPtr("Toko makanan dan minuman segar dan kemasan."),
			Address:      toStringPtr("Jl. Pasar No.10, Semarang"),
			ContactEmail: toStringPtr("fresh@mart.com"),
			ContactPhone: toStringPtr("086789012345"),
			Status:       "active",
		},
		{
			UserID:       7,
			Name:         "Gamer Heaven",
			Description:  toStringPtr("Game, console, dan aksesori lengkap untuk gamers."),
			Address:      toStringPtr("Jl. Game No.8, Bekasi"),
			ContactEmail: toStringPtr("gamer@heaven.com"),
			ContactPhone: toStringPtr("087890123456"),
			Status:       "active",
		},
		{
			UserID:       8,
			Name:         "AutoParts Store",
			Description:  toStringPtr("Toko perlengkapan otomotif terpercaya."),
			Address:      toStringPtr("Jl. Otomotif No.6, Medan"),
			ContactEmail: toStringPtr("service@autoparts.com"),
			ContactPhone: toStringPtr("088901234567"),
			Status:       "active",
		},
	}

	for _, merchant := range merchants {
		if _, err := r.db.CreateMerchant(r.ctx, merchant); err != nil {
			r.logger.Error("failed to seed merchant", zap.Error(err))
			return err
		}
	}

	r.logger.Info("merchant succesfully seeded")

	return nil
}
