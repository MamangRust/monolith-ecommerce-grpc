package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchant_social_link_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_social_link_errors"
)

type merchantSocialLinkCommandRepository struct {
	db *db.Queries
}

func NewMerchantSocialLinkCommandRepository(db *db.Queries) *merchantSocialLinkCommandRepository {
	return &merchantSocialLinkCommandRepository{
		db: db,
	}
}

func (r *merchantSocialLinkCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantSocialRequest) (*db.CreateMerchantSocialMediaLinkRow, error) {
	params := db.CreateMerchantSocialMediaLinkParams{
		MerchantDetailID: int32(*req.MerchantDetailID),
		Platform:         req.Platform,
		Url:              req.Url,
	}

	res, err := r.db.CreateMerchantSocialMediaLink(ctx, params)
	if err != nil {
		return nil, merchant_social_link_errors.ErrCreateMerchantSocialLink.WithInternal(err)
	}

	return res, nil
}

func (r *merchantSocialLinkCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (*db.UpdateMerchantSocialMediaLinkRow, error) {
	params := db.UpdateMerchantSocialMediaLinkParams{
		MerchantSocialID: int32(req.ID),
		Platform:         req.Platform,
		Url:              req.Url,
	}

	res, err := r.db.UpdateMerchantSocialMediaLink(ctx, params)
	if err != nil {
		return nil, merchant_social_link_errors.ErrUpdateMerchantSocialLink.WithInternal(err)
	}

	return res, nil
}

func (r *merchantSocialLinkCommandRepository) Trash(ctx context.Context, socialID int) (bool, error) {
	_, err := r.db.TrashMerchantSocialMediaLink(ctx, int32(socialID))
	if err != nil {
		return false, merchant_social_link_errors.ErrTrashMerchantSocialLink.WithInternal(err)
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) Restore(ctx context.Context, socialID int) (bool, error) {
	_, err := r.db.RestoreMerchantSocialMediaLink(ctx, int32(socialID))
	if err != nil {
		return false, merchant_social_link_errors.ErrRestoreMerchantSocialLink.WithInternal(err)
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) DeletePermanent(ctx context.Context, socialID int) (bool, error) {
	err := r.db.DeleteMerchantSocialMediaLinkPermanently(ctx, int32(socialID))
	if err != nil {
		return false, merchant_social_link_errors.ErrDeletePermanentMerchantSocialLink.WithInternal(err)
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantSocialMediaLinks(ctx)
	if err != nil {
		return false, merchant_social_link_errors.ErrRestoreAllMerchantSocialLinks.WithInternal(err)
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllMerchantSocialMediaLinksPermanently(ctx)
	if err != nil {
		return false, merchant_social_link_errors.ErrDeleteAllPermanentMerchantSocialLinks.WithInternal(err)
	}

	return true, nil
}
