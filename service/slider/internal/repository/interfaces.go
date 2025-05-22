package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type SliderQueryRepository interface {
	FindAllSlider(req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error)
	FindByActive(req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error)
	FindByTrashed(req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error)
}

type SliderCommandRepository interface {
	CreateSlider(request *requests.CreateSliderRequest) (*record.SliderRecord, error)
	UpdateSlider(request *requests.UpdateSliderRequest) (*record.SliderRecord, error)
	TrashSlider(slider_id int) (*record.SliderRecord, error)
	RestoreSlider(slider_id int) (*record.SliderRecord, error)
	DeleteSliderPermanently(slider_id int) (bool, error)
	RestoreAllSlider() (bool, error)
	DeleteAllPermanentSlider() (bool, error)
}
