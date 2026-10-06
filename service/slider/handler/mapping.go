package handler

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pbslider "github.com/MamangRust/monolith-ecommerce-pb/slider"
)

func MapToSliderResponse(slider *db.Slider) *pbslider.SliderResponse {
	return &pbslider.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func MapToSliderResponseGetSlidersRow(slider *db.GetSlidersRow) *pbslider.SliderResponse {
	return &pbslider.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func MapToSliderResponseGetSliderByIDRow(slider *db.GetSliderByIDRow) *pbslider.SliderResponse {
	return &pbslider.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func MapToSliderResponseCreateSliderRow(slider *db.CreateSliderRow) *pbslider.SliderResponse {
	return &pbslider.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func MapToSliderResponseUpdateSliderRow(slider *db.UpdateSliderRow) *pbslider.SliderResponse {
	return &pbslider.SliderResponse{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func MapToSliderResponseDeleteAt(slider *db.Slider) *pbslider.SliderResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if slider.DeletedAt.Valid {
		deletedAt = &wrapperspb.StringValue{Value: slider.DeletedAt.Time.Format("2006-01-02")}
	}

	return &pbslider.SliderResponseDeleteAt{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: deletedAt,
	}
}

func MapToSliderResponseDeleteAtGetSlidersActiveRow(slider *db.GetSlidersActiveRow) *pbslider.SliderResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if slider.DeletedAt.Valid {
		deletedAt = &wrapperspb.StringValue{Value: slider.DeletedAt.Time.Format("2006-01-02")}
	}

	return &pbslider.SliderResponseDeleteAt{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: deletedAt,
	}
}

func MapToSliderResponseDeleteAtGetSlidersTrashedRow(slider *db.GetSlidersTrashedRow) *pbslider.SliderResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if slider.DeletedAt.Valid {
		deletedAt = &wrapperspb.StringValue{Value: slider.DeletedAt.Time.Format("2006-01-02")}
	}

	return &pbslider.SliderResponseDeleteAt{
		Id:        int32(slider.SliderID),
		Name:      slider.Name,
		Image:     slider.Image,
		CreatedAt: slider.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: slider.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: deletedAt,
	}
}
