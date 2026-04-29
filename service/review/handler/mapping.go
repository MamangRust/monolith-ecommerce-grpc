package handler

import (
	"encoding/json"
	"log"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (h *reviewHandleGrpc) mapResponse(data interface{}) interface{} {
	switch v := data.(type) {
	case *db.Review:
		return h.mapReview(v)
	case *db.GetReviewsRow:
		return h.mapGetReviewsRow(v)
	case *db.GetReviewsActiveRow:
		return h.mapGetReviewsActiveRow(v)
	case *db.GetReviewsTrashedRow:
		return h.mapGetReviewsTrashedRow(v)
	case *db.GetReviewByProductIdRow:
		return h.mapGetReviewByProductIdRow(v)
	case *db.GetReviewByMerchantIdRow:
		return h.mapGetReviewByMerchantIdRow(v)
	case *db.GetReviewByIDRow:
		return h.mapGetReviewByIDRow(v)
	case *db.CreateReviewRow:
		return h.mapCreateReviewRow(v)
	case *db.UpdateReviewRow:
		return h.mapUpdateReviewRow(v)
	case []*db.GetReviewsRow:
		res := make([]*pb.ReviewResponse, len(v))
		for i, r := range v {
			res[i] = h.mapGetReviewsRow(r)
		}
		return res
	case []*db.GetReviewsActiveRow:
		res := make([]*pb.ReviewResponseDeleteAt, len(v))
		for i, r := range v {
			res[i] = h.mapGetReviewsActiveRow(r)
		}
		return res
	case []*db.GetReviewsTrashedRow:
		res := make([]*pb.ReviewResponseDeleteAt, len(v))
		for i, r := range v {
			res[i] = h.mapGetReviewsTrashedRow(r)
		}
		return res
	case []*db.GetReviewByProductIdRow:
		res := make([]*pb.ReviewsDetailResponse, len(v))
		for i, r := range v {
			res[i] = h.mapGetReviewByProductIdRow(r)
		}
		return res
	case []*db.GetReviewByMerchantIdRow:
		res := make([]*pb.ReviewsDetailResponse, len(v))
		for i, r := range v {
			res[i] = h.mapGetReviewByMerchantIdRow(r)
		}
		return res
	default:
		return nil
	}
}

func (h *reviewHandleGrpc) mapReview(v *db.Review) *pb.ReviewResponseDeleteAt {
	var deletedAt string
	if v.DeletedAt.Valid {
		deletedAt = v.DeletedAt.Time.Format("2006-01-02")
	}

	return &pb.ReviewResponseDeleteAt{
		Id:        int32(v.ReviewID),
		UserId:    int32(v.UserID),
		ProductId: int32(v.ProductID),
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    int32(v.Rating),
		CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
	}
}

func (h *reviewHandleGrpc) mapGetReviewsRow(v *db.GetReviewsRow) *pb.ReviewResponse {
	return &pb.ReviewResponse{
		Id:        int32(v.ReviewID),
		UserId:    int32(v.UserID),
		ProductId: int32(v.ProductID),
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    int32(v.Rating),
		CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func (h *reviewHandleGrpc) mapGetReviewsActiveRow(v *db.GetReviewsActiveRow) *pb.ReviewResponseDeleteAt {
	var deletedAt string
	if v.DeletedAt.Valid {
		deletedAt = v.DeletedAt.Time.Format("2006-01-02")
	}

	return &pb.ReviewResponseDeleteAt{
		Id:        int32(v.ReviewID),
		UserId:    int32(v.UserID),
		ProductId: int32(v.ProductID),
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    int32(v.Rating),
		CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
	}
}

func (h *reviewHandleGrpc) mapGetReviewsTrashedRow(v *db.GetReviewsTrashedRow) *pb.ReviewResponseDeleteAt {
	var deletedAt string
	if v.DeletedAt.Valid {
		deletedAt = v.DeletedAt.Time.Format("2006-01-02")
	}

	return &pb.ReviewResponseDeleteAt{
		Id:        int32(v.ReviewID),
		UserId:    int32(v.UserID),
		ProductId: int32(v.ProductID),
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    int32(v.Rating),
		CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
	}
}

func (h *reviewHandleGrpc) mapGetReviewByProductIdRow(v *db.GetReviewByProductIdRow) *pb.ReviewsDetailResponse {
	res := &pb.ReviewsDetailResponse{
		Id:        int32(v.ReviewID),
		UserId:    int32(v.UserID),
		ProductId: int32(v.ProductID),
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    int32(v.Rating),
		CreatedAt: v.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
	}

	if v.DeletedAt.Valid {
		res.DeletedAt = v.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
	}

	if v.ReviewDetails != nil {
		h.mapReviewDetails(v.ReviewDetails, res)
	}

	return res
}

func (h *reviewHandleGrpc) mapGetReviewByMerchantIdRow(v *db.GetReviewByMerchantIdRow) *pb.ReviewsDetailResponse {
	res := &pb.ReviewsDetailResponse{
		Id:        int32(v.ReviewID),
		UserId:    int32(v.UserID),
		ProductId: int32(v.ProductID),
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    int32(v.Rating),
		CreatedAt: v.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
	}

	if v.DeletedAt.Valid {
		res.DeletedAt = v.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
	}

	if v.ReviewDetails != nil {
		h.mapReviewDetails(v.ReviewDetails, res)
	}

	return res
}

func (h *reviewHandleGrpc) mapGetReviewByIDRow(v *db.GetReviewByIDRow) *pb.ReviewResponse {
	return &pb.ReviewResponse{
		Id:        int32(v.ReviewID),
		UserId:    int32(v.UserID),
		ProductId: int32(v.ProductID),
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    int32(v.Rating),
		CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func (h *reviewHandleGrpc) mapCreateReviewRow(v *db.CreateReviewRow) *pb.ReviewResponse {
	return &pb.ReviewResponse{
		Id:        int32(v.ReviewID),
		UserId:    int32(v.UserID),
		ProductId: int32(v.ProductID),
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    int32(v.Rating),
		CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func (h *reviewHandleGrpc) mapUpdateReviewRow(v *db.UpdateReviewRow) *pb.ReviewResponse {
	return &pb.ReviewResponse{
		Id:        int32(v.ReviewID),
		UserId:    int32(v.UserID),
		ProductId: int32(v.ProductID),
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    int32(v.Rating),
		CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func (h *reviewHandleGrpc) mapReviewDetails(reviewDetails interface{}, res *pb.ReviewsDetailResponse) {
	var details []struct {
		DetailID  int    `json:"detail_id"`
		Type      string `json:"type"`
		URL       string `json:"url"`
		Caption   string `json:"caption"`
		CreatedAt string `json:"created_at"`
	}

	detailsBytes, err := json.Marshal(reviewDetails)
	if err != nil {
		log.Printf("Error marshaling review details: %v", err)
		return
	}

	err = json.Unmarshal(detailsBytes, &details)
	if err != nil {
		log.Printf("Error unmarshaling review details: %v", err)
		return
	}

	if len(details) > 0 {
		firstDetail := details[0]
		res.ReviewDetail = &pb.ReviewDetailResponse{
			Id:        int32(firstDetail.DetailID),
			Type:      firstDetail.Type,
			Url:       firstDetail.URL,
			Caption:   firstDetail.Caption,
			CreatedAt: firstDetail.CreatedAt,
		}
	}
}
