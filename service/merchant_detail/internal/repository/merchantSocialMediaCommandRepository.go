package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantsociallink_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_social_link_errors"
)

type merchantSocialLinkCommandRepository struct {
	db  *db.Queries
	ctx context.Context
}

func NewMerchantSocialLinkCommandRepository(db *db.Queries, ctx context.Context) *merchantSocialLinkCommandRepository {
	return &merchantSocialLinkCommandRepository{
		db:  db,
		ctx: ctx,
	}
}

func (r *merchantSocialLinkCommandRepository) CreateSocialLink(req *requests.CreateMerchantSocialRequest) (bool, error) {
	params := db.CreateMerchantSocialMediaLinkParams{
		MerchantDetailID: int32(*req.MerchantDetailID),
		Platform:         req.Platform,
		Url:              req.Url,
	}

	_, err := r.db.CreateMerchantSocialMediaLink(r.ctx, params)
	if err != nil {
		return false, merchantsociallink_errors.ErrCreateMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) UpdateSocialLink(req *requests.UpdateMerchantSocialRequest) (bool, error) {
	params := db.UpdateMerchantSocialMediaLinkParams{
		MerchantSocialID: int32(req.ID),
		Platform:         req.Platform,
		Url:              req.Url,
	}

	_, err := r.db.UpdateMerchantSocialMediaLink(r.ctx, params)
	if err != nil {
		return false, merchantsociallink_errors.ErrUpdateMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) TrashSocialLink(socialID int) (bool, error) {
	_, err := r.db.TrashMerchantSocialMediaLink(r.ctx, int32(socialID))
	if err != nil {
		return false, merchantsociallink_errors.ErrTrashMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) RestoreSocialLink(socialID int) (bool, error) {
	_, err := r.db.RestoreMerchantSocialMediaLink(r.ctx, int32(socialID))
	if err != nil {
		return false, merchantsociallink_errors.ErrRestoreMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) DeletePermanentSocialLink(socialID int) (bool, error) {
	err := r.db.DeleteMerchantSocialMediaLinkPermanently(r.ctx, int32(socialID))
	if err != nil {
		return false, merchantsociallink_errors.ErrDeletePermanentMerchantSocialLink
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) RestoreAllSocialLink() (bool, error) {
	err := r.db.RestoreAllMerchantSocialMediaLinks(r.ctx)
	if err != nil {
		return false, merchantsociallink_errors.ErrRestoreAllMerchantSocialLinks
	}

	return true, nil
}

func (r *merchantSocialLinkCommandRepository) DeleteAllPermanentSocialLink() (bool, error) {
	err := r.db.DeleteAllMerchantSocialMediaLinksPermanently(r.ctx)
	if err != nil {
		return false, merchantsociallink_errors.ErrDeleteAllPermanentMerchantSocialLinks
	}

	return true, nil
}
