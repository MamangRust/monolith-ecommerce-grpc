package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type SliderQueryCache interface {
	GetSliderAllCache(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersRow, *int, bool)
	SetSliderAllCache(ctx context.Context, req *requests.FindAllSlider, data []*db.GetSlidersRow, total *int)

	GetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersActiveRow, *int, bool)
	SetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider, data []*db.GetSlidersActiveRow, total *int)

	GetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersTrashedRow, *int, bool)
	SetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider, data []*db.GetSlidersTrashedRow, total *int)

	GetSliderCache(ctx context.Context, slider_id int) (*db.GetSliderByIDRow, bool)
	SetSliderCache(ctx context.Context, data *db.GetSliderByIDRow)
}

type SliderCommandCache interface {
	DeleteSliderCache(ctx context.Context, slider_id int)
	InvalidateSliderCache(ctx context.Context)
}
