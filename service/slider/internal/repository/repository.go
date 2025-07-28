package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	SliderQuery   SliderQueryRepository
	SliderCommand SliderCommandRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapper := recordmapper.NewSliderRecordMapper()

	return &Repositories{
		SliderQuery:   NewSliderQueryRepository(DB, mapper),
		SliderCommand: NewSliderCommandRepository(DB, mapper),
	}
}
