package handler

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func MapToSliderResponse(slider *db.Slider) *pb.SliderResponse {
	return &pb.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func MapToSliderResponseGetSlidersRow(slider *db.GetSlidersRow) *pb.SliderResponse {
	return &pb.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func MapToSliderResponseGetSliderByIDRow(slider *db.GetSliderByIDRow) *pb.SliderResponse {
	return &pb.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func MapToSliderResponseCreateSliderRow(slider *db.CreateSliderRow) *pb.SliderResponse {
	return &pb.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func MapToSliderResponseUpdateSliderRow(slider *db.UpdateSliderRow) *pb.SliderResponse {
	return &pb.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func MapToSliderResponseDeleteAt(slider *db.Slider) *pb.SliderResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if slider.DeletedAt.Valid {
		deletedAt = &wrapperspb.StringValue{Value: slider.DeletedAt.Time.Format("2006-01-02")}
	}

	return &pb.SliderResponseDeleteAt{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: deletedAt,
	}
}

func MapToSliderResponseDeleteAtGetSlidersActiveRow(slider *db.GetSlidersActiveRow) *pb.SliderResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if slider.DeletedAt.Valid {
		deletedAt = &wrapperspb.StringValue{Value: slider.DeletedAt.Time.Format("2006-01-02")}
	}

	return &pb.SliderResponseDeleteAt{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: deletedAt,
	}
}

func MapToSliderResponseDeleteAtGetSlidersTrashedRow(slider *db.GetSlidersTrashedRow) *pb.SliderResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if slider.DeletedAt.Valid {
		deletedAt = &wrapperspb.StringValue{Value: slider.DeletedAt.Time.Format("2006-01-02")}
	}

	return &pb.SliderResponseDeleteAt{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: deletedAt,
	}
}
