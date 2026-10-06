package handler

import (
	"encoding/json"
	"log"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pbreview_detail "github.com/MamangRust/monolith-ecommerce-pb/review_detail"
)

func (h *Handler) mapToReviewDetailResponse(data interface{}) interface{} {
	switch v := data.(type) {
	case *db.GetReviewDetailRow:
		return &pbreview_detail.ReviewDetailsResponse{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.GetReviewDetailsRow:
		return &pbreview_detail.ReviewDetailsResponse{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.GetReviewDetailsActiveRow:
		var deletedAt string
		if v.DeletedAt.Valid {
			deletedAt = v.DeletedAt.Time.Format("2006-01-02")
		}
		return &pbreview_detail.ReviewDetailsResponseDeleteAt{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		}
	case *db.GetReviewDetailsTrashedRow:
		var deletedAt string
		if v.DeletedAt.Valid {
			deletedAt = v.DeletedAt.Time.Format("2006-01-02")
		}
		return &pbreview_detail.ReviewDetailsResponseDeleteAt{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		}
	case *db.CreateReviewDetailRow:
		return &pbreview_detail.ReviewDetailsResponse{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.UpdateReviewDetailRow:
		return &pbreview_detail.ReviewDetailsResponse{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.ReviewDetail:
		var deletedAt string
		if v.DeletedAt.Valid {
			deletedAt = v.DeletedAt.Time.Format("2006-01-02")
		}
		return &pbreview_detail.ReviewDetailsResponseDeleteAt{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		}
	default:
		log.Printf("Unknown type for mapping: %T", v)
		return nil
	}
}

func (h *Handler) mapToPayload(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
