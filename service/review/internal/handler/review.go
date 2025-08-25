package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type reviewHandleGrpc struct {
	pb.UnimplementedReviewServiceServer
	reviewQueryService   service.ReviewQueryService
	reviewCommandService service.ReviewCommandService
	logger               logger.LoggerInterface
	mapping              protomapper.ReviewProtoMapper
}

func NewReviewHandleGrpc(service *service.Service, logger logger.LoggerInterface) pb.ReviewServiceServer {
	return &reviewHandleGrpc{
		reviewQueryService:   service.ReviewQuery,
		reviewCommandService: service.ReviewCommand,
		logger:               logger,
		mapping:              protomapper.NewReviewProtoMapper(),
	}
}

func (s *reviewHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReview, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all reviews",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviews, totalRecords, err := s.reviewQueryService.FindAllReviews(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all reviews",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched all reviews",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_reviews_count", len(reviews)),
	)

	so := s.mapping.ToProtoResponsePaginationReview(paginationMeta, "success", "Successfully fetched reviews", reviews)
	return so, nil
}

func (s *reviewHandleGrpc) FindByProduct(ctx context.Context, request *pb.FindAllReviewProductRequest) (*pb.ApiResponsePaginationReviewDetail, error) {
	productID := int(request.GetProductId())
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching reviews by product",
		zap.Int("product_id", productID),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllReviewByProduct{
		ProductID: productID,
		Page:      page,
		PageSize:  pageSize,
		Search:    search,
	}

	reviews, totalRecords, err := s.reviewQueryService.FindByProduct(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch reviews by product",
			zap.Int("product_id", productID),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched reviews by product",
		zap.Int("product_id", productID),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int("fetched_reviews_count", len(reviews)),
	)

	so := s.mapping.ToProtoResponsePaginationReviewsDetail(paginationMeta, "success", "Successfully fetched reviews for product", reviews)
	return so, nil
}

func (s *reviewHandleGrpc) FindByMerchant(ctx context.Context, request *pb.FindAllReviewMerchantRequest) (*pb.ApiResponsePaginationReviewDetail, error) {
	merchantID := int(request.GetMerchantId())
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching reviews by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllReviewByMerchant{
		MerchantID: merchantID,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	reviews, totalRecords, err := s.reviewQueryService.FindByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch reviews by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched reviews by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int("fetched_reviews_count", len(reviews)),
	)

	so := s.mapping.ToProtoResponsePaginationReviewsDetail(paginationMeta, "success", "Successfully fetched reviews for merchant", reviews)
	return so, nil
}

func (s *reviewHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active reviews",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviews, totalRecords, err := s.reviewQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active reviews",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched active reviews",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_reviews_count", len(reviews)),
	)

	so := s.mapping.ToProtoResponsePaginationReviewDeleteAt(paginationMeta, "success", "Successfully fetched active reviews", reviews)
	return so, nil
}

func (s *reviewHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed reviews",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviews, totalRecords, err := s.reviewQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed reviews",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched trashed reviews",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_reviews_count", len(reviews)),
	)

	so := s.mapping.ToProtoResponsePaginationReviewDeleteAt(paginationMeta, "success", "Successfully fetched trashed reviews", reviews)
	return so, nil
}

func (s *reviewHandleGrpc) Create(ctx context.Context, request *pb.CreateReviewRequest) (*pb.ApiResponseReview, error) {
	s.logger.Info("Creating new review",
		zap.Int("user_id", int(request.GetUserId())),
		zap.Int("product_id", int(request.GetProductId())),
		zap.Int("rating", int(request.GetRating())),
	)

	req := &requests.CreateReviewRequest{
		UserID:    int(request.GetUserId()),
		ProductID: int(request.GetProductId()),
		Rating:    int(request.GetRating()),
		Comment:   request.GetComment(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on review creation",
			zap.Int("user_id", int(request.GetUserId())),
			zap.Int("product_id", int(request.GetProductId())),
			zap.Error(err),
		)
		return nil, review_errors.ErrGrpcValidateCreateReview
	}

	review, err := s.reviewCommandService.CreateReview(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create review",
			zap.Int("user_id", int(request.GetUserId())),
			zap.Int("product_id", int(request.GetProductId())),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Review created successfully",
		zap.Int("review_id", int(review.ID)),
		zap.Int("rating", review.Rating),
		zap.String("product_name", review.Name),
	)

	return s.mapping.ToProtoResponseReview("success", "Successfully created review", review), nil
}

func (s *reviewHandleGrpc) Update(ctx context.Context, request *pb.UpdateReviewRequest) (*pb.ApiResponseReview, error) {
	id := int(request.GetReviewId())

	if id == 0 {
		s.logger.Error("Invalid review ID provided for update", zap.Int("review_id", id))
		return nil, review_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Updating review", zap.Int("review_id", id))

	req := &requests.UpdateReviewRequest{
		ReviewID: &id,
		Name:     request.GetName(),
		Rating:   int(request.GetRating()),
		Comment:  request.GetComment(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on review update",
			zap.Int("review_id", id),
			zap.Error(err),
		)
		return nil, review_errors.ErrGrpcValidateUpdateReview
	}

	review, err := s.reviewCommandService.UpdateReview(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update review",
			zap.Int("review_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Review updated successfully",
		zap.Int("review_id", id),
		zap.Int("rating", review.Rating),
	)

	return s.mapping.ToProtoResponseReview("success", "Successfully updated review", review), nil
}

func (s *reviewHandleGrpc) TrashedReview(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid review ID for trashing", zap.Int("review_id", id))
		return nil, review_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Moving review to trash", zap.Int("review_id", id))

	review, err := s.reviewCommandService.TrashedReview(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash review",
			zap.Int("review_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Review moved to trash successfully",
		zap.Int("review_id", id),
		zap.String("product_name", review.Name),
		zap.Int("user_id", int(review.UserID)),
	)

	so := s.mapping.ToProtoResponseReviewDeleteAt("success", "Successfully trashed review", review)
	return so, nil
}

func (s *reviewHandleGrpc) RestoreReview(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid review ID for restore", zap.Int("review_id", id))
		return nil, review_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Restoring review from trash", zap.Int("review_id", id))

	review, err := s.reviewCommandService.RestoreReview(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore review",
			zap.Int("review_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Review restored successfully",
		zap.Int("review_id", id),
		zap.String("product_name", review.Name),
	)

	so := s.mapping.ToProtoResponseReviewDeleteAt("success", "Successfully restored review", review)
	return so, nil
}

func (s *reviewHandleGrpc) DeleteReviewPermanent(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid review ID for permanent deletion", zap.Int("review_id", id))
		return nil, review_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Permanently deleting review", zap.Int("review_id", id))

	_, err := s.reviewCommandService.DeleteReviewPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete review",
			zap.Int("review_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Review permanently deleted", zap.Int("review_id", id))

	so := s.mapping.ToProtoResponseReviewDelete("success", "Successfully deleted review permanently")
	return so, nil
}

func (s *reviewHandleGrpc) RestoreAllReview(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	s.logger.Info("Restoring all trashed reviews")

	_, err := s.reviewCommandService.RestoreAllReviews(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all reviews", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All reviews restored successfully")

	so := s.mapping.ToProtoResponseReviewAll("success", "Successfully restored all reviews")
	return so, nil
}

func (s *reviewHandleGrpc) DeleteAllReviewPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	s.logger.Info("Permanently deleting all trashed reviews")

	_, err := s.reviewCommandService.DeleteAllReviewsPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all reviews", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All reviews permanently deleted")

	so := s.mapping.ToProtoResponseReviewAll("success", "Successfully deleted all reviews permanently")
	return so, nil
}
