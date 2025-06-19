package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErroHandler struct {
	BannerQueryError   BannerQueryError
	BannerCommandError BannerCommandError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErroHandler {
	return &ErroHandler{
		BannerQueryError:   NewBannerQueryError(logger),
		BannerCommandError: NewBannerCommandError(logger),
	}
}
