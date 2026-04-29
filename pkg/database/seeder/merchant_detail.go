package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"go.uber.org/zap"
)

type merchantDetailSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantDetailSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *merchantDetailSeeder {
	return &merchantDetailSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *merchantDetailSeeder) Seed() error {
	details := []db.CreateMerchantDetailParams{
		{
			MerchantID:       1,
			DisplayName:      toStringPtr("Techno Store"),
			CoverImageUrl:    toStringPtr("cover/techno.jpg"),
			LogoUrl:          toStringPtr("logo/techno.png"),
			ShortDescription: toStringPtr("Pusat elektronik terpercaya sejak 2010"),
			WebsiteUrl:       toStringPtr("https://technostore.com"),
		},
		{
			MerchantID:       2,
			DisplayName:      toStringPtr("Glow Beauty"),
			CoverImageUrl:    toStringPtr("cover/beauty.jpg"),
			LogoUrl:          toStringPtr("logo/beauty.png"),
			ShortDescription: toStringPtr("Produk kecantikan alami dan aman"),
			WebsiteUrl:       toStringPtr("https://glowbeauty.id"),
		},
		{
			MerchantID:       3,
			DisplayName:      toStringPtr("Dapur Sehat"),
			CoverImageUrl:    toStringPtr("cover/dapur.jpg"),
			LogoUrl:          toStringPtr("logo/dapur.png"),
			ShortDescription: toStringPtr("Makanan sehat dan organik"),
			WebsiteUrl:       toStringPtr("https://dapsehat.id"),
		},
		{
			MerchantID:       4,
			DisplayName:      toStringPtr("Gadget Hub"),
			CoverImageUrl:    toStringPtr("cover/gadget.jpg"),
			LogoUrl:          toStringPtr("logo/gadget.png"),
			ShortDescription: toStringPtr("Semua tentang gadget terbaru"),
			WebsiteUrl:       toStringPtr("https://gadgethub.com"),
		},
		{
			MerchantID:       5,
			DisplayName:      toStringPtr("Bayi Ceria"),
			CoverImageUrl:    toStringPtr("cover/bayi.jpg"),
			LogoUrl:          toStringPtr("logo/bayi.png"),
			ShortDescription: toStringPtr("Produk terbaik untuk si kecil"),
			WebsiteUrl:       toStringPtr("https://bayiceria.id"),
		},
		{
			MerchantID:       6,
			DisplayName:      toStringPtr("Toko Sehat"),
			CoverImageUrl:    toStringPtr("cover/sehat.jpg"),
			LogoUrl:          toStringPtr("logo/sehat.png"),
			ShortDescription: toStringPtr("Peralatan olahraga lengkap"),
			WebsiteUrl:       toStringPtr("https://tokosehat.id"),
		},
		{
			MerchantID:       7,
			DisplayName:      toStringPtr("Game World"),
			CoverImageUrl:    toStringPtr("cover/game.jpg"),
			LogoUrl:          toStringPtr("logo/game.png"),
			ShortDescription: toStringPtr("Konsol dan game terbaik"),
			WebsiteUrl:       toStringPtr("https://gameworld.com"),
		},
		{
			MerchantID:       8,
			DisplayName:      toStringPtr("Otomotif Mart"),
			CoverImageUrl:    toStringPtr("cover/otomotif.jpg"),
			LogoUrl:          toStringPtr("logo/otomotif.png"),
			ShortDescription: toStringPtr("Aksesori kendaraan terpercaya"),
			WebsiteUrl:       toStringPtr("https://otomotifmart.com"),
		},
	}

	for i, detail := range details {
		_, err := r.db.CreateMerchantDetail(r.ctx, detail)
		if err != nil {
			r.logger.Error("failed to seed merchant detail", zap.Error(err))
			return err
		}

		merchantDetailID := int32(i + 1)
		socialMedia := []db.CreateMerchantSocialMediaLinkParams{
			{
				MerchantDetailID: merchantDetailID,
				Platform:         "Facebook",
				Url:              "https://www.facebook.com/merchant" + fmt.Sprint(merchantDetailID),
			},
			{
				MerchantDetailID: merchantDetailID,
				Platform:         "Instagram",
				Url:              "https://www.instagram.com/merchant" + fmt.Sprint(merchantDetailID),
			},
			{
				MerchantDetailID: merchantDetailID,
				Platform:         "Twitter",
				Url:              "https://www.twitter.com/merchant" + fmt.Sprint(merchantDetailID),
			},
		}

		for _, sm := range socialMedia {
			if _, err := r.db.CreateMerchantSocialMediaLink(r.ctx, sm); err != nil {
				r.logger.Error("failed to seed merchant social media link", zap.Error(err))
				return err
			}
		}
	}

	r.logger.Info("merchant detail & merchant social link successfully seeded")

	return nil
}
