package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantsociallink_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_social_link_errors"
)

type merchantSocialLinkCommandRepository struct {
	db *db.Queries
}

func NewMerchantSocialLinkCommandRepository(db *db.Queries) *merchantSocialLinkCommandRepository {
	return &merchantSocialLinkCommandRepository{
		db: db,
	}
}

func (r *merchantSocialLinkCommandRepository) CreateSocialLink(ctx context.Context, req *requests.CreateMerchantSocialRequest) (bool, error) {
	params := db.CreateMerchantSocialMediaLinkParams{
		MerchantDetailID: int32(*req.MerchantDetailID),
		Platform:         req.Platform,
		Url:              req.Url,
	}

	_, err := r.db.CreateMerchantSocialMediaLink(ctx, params)
	if err != nil {
		return false, merchantsociallink_errors.ErrCreateMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) UpdateSocialLink(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (bool, error) {
	params := db.UpdateMerchantSocialMediaLinkParams{
		MerchantSocialID: int32(req.ID),
		Platform:         req.Platform,
		Url:              req.Url,
	}

	_, err := r.db.UpdateMerchantSocialMediaLink(ctx, params)
	if err != nil {
		return false, merchantsociallink_errors.ErrUpdateMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) TrashSocialLink(ctx context.Context, socialID int) (bool, error) {
	_, err := r.db.TrashMerchantSocialMediaLink(ctx, int32(socialID))
	if err != nil {
		return false, merchantsociallink_errors.ErrTrashMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) RestoreSocialLink(ctx context.Context, socialID int) (bool, error) {
	_, err := r.db.RestoreMerchantSocialMediaLink(ctx, int32(socialID))
	if err != nil {
		return false, merchantsociallink_errors.ErrRestoreMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) DeletePermanentSocialLink(ctx context.Context, socialID int) (bool, error) {
	err := r.db.DeleteMerchantSocialMediaLinkPermanently(ctx, int32(socialID))
	if err != nil {
		return false, merchantsociallink_errors.ErrDeletePermanentMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) RestoreAllSocialLink(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantSocialMediaLinks(ctx)
	if err != nil {
		return false, merchantsociallink_errors.ErrRestoreAllMerchantSocialLinks
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) DeleteAllPermanentSocialLink(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllMerchantSocialMediaLinksPermanently(ctx)
	if err != nil {
		return false, merchantsociallink_errors.ErrDeleteAllPermanentMerchantSocialLinks
	}

	return true, nil
}
