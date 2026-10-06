package handler

import (
	pbreview "github.com/MamangRust/monolith-ecommerce-pb/review"
)

type ReviewHandleGrpc interface {
	pbreview.ReviewQueryServiceServer
	pbreview.ReviewCommandServiceServer
}

type ReviewQueryHandler interface {
	pbreview.ReviewQueryServiceServer
}

type ReviewCommandHandler interface {
	pbreview.ReviewCommandServiceServer
}
