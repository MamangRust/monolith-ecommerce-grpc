package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-user/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	UserQuery   UserQueryHandler
	UserCommand UserCommandHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		UserQuery:   NewUserQueryHandler(deps.Service.UserQuery, deps.Logger),
		UserCommand: NewUserCommandHandler(deps.Service.UserCommand, deps.Logger),
	}
}
