package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

type Repositories struct {
	SliderQuery   SliderQueryRepository
	SliderCommand SliderCommandRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		SliderQuery:   NewSliderQueryRepository(DB),
		SliderCommand: NewSliderCommandRepository(DB),
	}
}
