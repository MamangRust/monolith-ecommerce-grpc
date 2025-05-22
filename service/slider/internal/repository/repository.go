package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	SliderQuery   SliderQueryRepository
	SliderCommand SliderCommandRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps Deps) *Repositories {
	mapper := recordmapper.NewSliderRecordMapper()

	return &Repositories{
		SliderQuery:   NewSliderQueryRepository(deps.DB, deps.Ctx, mapper),
		SliderCommand: NewSliderCommandRepository(deps.DB, deps.Ctx, mapper),
	}
}
