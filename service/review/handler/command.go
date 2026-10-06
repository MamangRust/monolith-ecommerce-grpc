package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-review/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	"google.golang.org/protobuf/types/known/emptypb"

	pbreview "github.com/MamangRust/monolith-ecommerce-pb/review"
)

type reviewCommandHandler struct {
	pbreview.UnimplementedReviewCommandServiceServer
	reviewService service.ReviewCommandService
	logger        logger.LoggerInterface
}

func NewReviewCommandHandler(reviewService service.ReviewCommandService, logger logger.LoggerInterface) pbreview.ReviewCommandServiceServer {
	return &reviewCommandHandler{
		reviewService: reviewService,
		logger:        logger,
	}
}

func (h *reviewCommandHandler) Create(ctx context.Context, request *pbreview.CreateReviewRequest) (*pbreview.ApiResponseReview, error) {
	req := &requests.CreateReviewRequest{
		UserID:    int(request.GetUserId()),
		ProductID: int(request.GetProductId()),
		Rating:    int(request.GetRating()),
		Comment:   request.GetComment(),
	}

	if err := req.Validate(); err != nil {
		return nil, review_errors.ErrGrpcValidateCreateReview
	}

	review, err := h.reviewService.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var hMapping reviewHandleGrpc
	protoReview := hMapping.mapResponse(review).(*pbreview.ReviewResponse)

	return &pbreview.ApiResponseReview{
		Status:  "success",
		Message: "Successfully created review",
		Data:    protoReview,
	}, nil
}

func (h *reviewCommandHandler) Update(ctx context.Context, request *pbreview.UpdateReviewRequest) (*pbreview.ApiResponseReview, error) {
	id := int(request.GetReviewId())

	if id == 0 {
		return nil, review_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateReviewRequest{
		ReviewID: &id,
		Name:     request.GetName(),
		Rating:   int(request.GetRating()),
		Comment:  request.GetComment(),
	}

	if err := req.Validate(); err != nil {
		return nil, review_errors.ErrGrpcValidateUpdateReview
	}

	review, err := h.reviewService.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var hMapping reviewHandleGrpc
	protoReview := hMapping.mapResponse(review).(*pbreview.ReviewResponse)

	return &pbreview.ApiResponseReview{
		Status:  "success",
		Message: "Successfully updated review",
		Data:    protoReview,
	}, nil
}

func (h *reviewCommandHandler) TrashedReview(ctx context.Context, request *pbreview.FindByIdReviewRequest) (*pbreview.ApiResponseReviewDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, review_errors.ErrGrpcInvalidID
	}

	review, err := h.reviewService.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var hMapping reviewHandleGrpc
	protoReview := hMapping.mapResponse(review).(*pbreview.ReviewResponseDeleteAt)

	return &pbreview.ApiResponseReviewDeleteAt{
		Status:  "success",
		Message: "Successfully trashed review",
		Data:    protoReview,
	}, nil
}

func (h *reviewCommandHandler) RestoreReview(ctx context.Context, request *pbreview.FindByIdReviewRequest) (*pbreview.ApiResponseReviewDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, review_errors.ErrGrpcInvalidID
	}

	review, err := h.reviewService.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var hMapping reviewHandleGrpc
	protoReview := hMapping.mapResponse(review).(*pbreview.ReviewResponseDeleteAt)

	return &pbreview.ApiResponseReviewDeleteAt{
		Status:  "success",
		Message: "Successfully restored review",
		Data:    protoReview,
	}, nil
}

func (h *reviewCommandHandler) DeleteReviewPermanent(ctx context.Context, request *pbreview.FindByIdReviewRequest) (*pbreview.ApiResponseReviewDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, review_errors.ErrGrpcInvalidID
	}

	_, err := h.reviewService.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbreview.ApiResponseReviewDelete{
		Status:  "success",
		Message: "Successfully deleted review permanently",
	}, nil
}

func (h *reviewCommandHandler) RestoreAllReview(ctx context.Context, _ *emptypb.Empty) (*pbreview.ApiResponseReviewAll, error) {
	_, err := h.reviewService.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbreview.ApiResponseReviewAll{
		Status:  "success",
		Message: "Successfully restored all reviews",
	}, nil
}

func (h *reviewCommandHandler) DeleteAllReviewPermanent(ctx context.Context, _ *emptypb.Empty) (*pbreview.ApiResponseReviewAll, error) {
	_, err := h.reviewService.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbreview.ApiResponseReviewAll{
		Status:  "success",
		Message: "Successfully deleted all reviews permanently",
	}, nil
}
