package seeder

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"go.uber.org/zap"
)

type merchantBusinessSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantBusinessSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *merchantBusinessSeeder {
	return &merchantBusinessSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *merchantBusinessSeeder) Seed() error {
	businessInfos := []db.CreateMerchantBusinessInformationParams{
		{
			MerchantID:        1,
			BusinessType:      toStringPtr("Toko Elektronik"),
			TaxID:             toStringPtr("01.234.567.8-999.000"),
			EstablishedYear:   toInt32Ptr(2010),
			NumberOfEmployees: toInt32Ptr(15),
			WebsiteUrl:        toStringPtr("https://technostore.com"),
		},
		{
			MerchantID:        2,
			BusinessType:      toStringPtr("Produk Kecantikan"),
			TaxID:             toStringPtr("02.345.678.9-888.000"),
			EstablishedYear:   toInt32Ptr(2015),
			NumberOfEmployees: toInt32Ptr(10),
			WebsiteUrl:        toStringPtr("https://glowbeauty.id"),
		},
		{
			MerchantID:        3,
			BusinessType:      toStringPtr("Toko Makanan Organik"),
			TaxID:             toStringPtr("03.456.789.0-777.000"),
			EstablishedYear:   toInt32Ptr(2012),
			NumberOfEmployees: toInt32Ptr(20),
			WebsiteUrl:        toStringPtr("https://dapsehat.id"),
		},
		{
			MerchantID:        4,
			BusinessType:      toStringPtr("Retail Gadget"),
			TaxID:             toStringPtr("04.567.890.1-666.000"),
			EstablishedYear:   toInt32Ptr(2018),
			NumberOfEmployees: toInt32Ptr(8),
			WebsiteUrl:        toStringPtr("https://gadgethub.com"),
		},
		{
			MerchantID:        5,
			BusinessType:      toStringPtr("Produk Ibu & Bayi"),
			TaxID:             toStringPtr("05.678.901.2-555.000"),
			EstablishedYear:   toInt32Ptr(2019),
			NumberOfEmployees: toInt32Ptr(6),
			WebsiteUrl:        toStringPtr("https://bayiceria.id"),
		},
		{
			MerchantID:        6,
			BusinessType:      toStringPtr("Peralatan Olahraga"),
			TaxID:             toStringPtr("06.789.012.3-444.000"),
			EstablishedYear:   toInt32Ptr(2016),
			NumberOfEmployees: toInt32Ptr(12),
			WebsiteUrl:        toStringPtr("https://tokosehat.id"),
		},
		{
			MerchantID:        7,
			BusinessType:      toStringPtr("Gaming Store"),
			TaxID:             toStringPtr("07.890.123.4-333.000"),
			EstablishedYear:   toInt32Ptr(2020),
			NumberOfEmployees: toInt32Ptr(5),
			WebsiteUrl:        toStringPtr("https://gameworld.com"),
		},
		{
			MerchantID:        8,
			BusinessType:      toStringPtr("Aksesori Otomotif"),
			TaxID:             toStringPtr("08.901.234.5-222.000"),
			EstablishedYear:   toInt32Ptr(2013),
			NumberOfEmployees: toInt32Ptr(9),
			WebsiteUrl:        toStringPtr("https://otomotifmart.com"),
		},
	}

	for _, info := range businessInfos {
		if _, err := r.db.CreateMerchantBusinessInformation(r.ctx, info); err != nil {
			r.logger.Error("failed to seed merchant business info", zap.Error(err))
			return err
		}
	}

	r.logger.Info("merchant business successfully seeded")

	return nil
}
