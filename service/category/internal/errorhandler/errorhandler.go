package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	CategoryQueryError           CategoryQueryError
	CategoryCommandError         CategoryCommandError
	CategoryStatsError           CategoryStatsError
	CategoryStatsByIdError       CategoryStatsByIdError
	CategoryStatsByMerchantError CategoryStatsByMerchantError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		CategoryQueryError:           NewCategoryQueryError(logger),
		CategoryCommandError:         NewCategoryCommandError(logger),
		CategoryStatsError:           NewCategoryStatsError(logger),
		CategoryStatsByIdError:       NewCategoryStatsByIdError(logger),
		CategoryStatsByMerchantError: NewCategoryStatsByMerchantError(logger),
	}
}
