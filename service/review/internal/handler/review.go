package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type reviewHandleGrpc struct {
	pb.UnimplementedReviewServiceServer
	reviewQueryService   service.ReviewQueryService
	reviewCommandService service.ReviewCommandService
	mapping              protomapper.ReviewProtoMapper
}

func NewReviewHandleGrpc(service *service.Service) *reviewHandleGrpc {
	return &reviewHandleGrpc{
		reviewQueryService:   service.ReviewQuery,
		reviewCommandService: service.ReviewCommand,
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

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	review, totalRecords, err := s.reviewQueryService.FindAllReviews(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationReview(paginationMeta, "success", "Successfully fetched categories", review)
	return so, nil
}

func (s *reviewHandleGrpc) FindByProduct(ctx context.Context, request *pb.FindAllReviewProductRequest) (*pb.ApiResponsePaginationReviewDetail, error) {
	product_id := int(request.GetProductId())
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	reqService := requests.FindAllReviewByProduct{
		ProductID: product_id,
		Page:      page,
		PageSize:  pageSize,
		Search:    search,
	}

	review, totalRecords, err := s.reviewQueryService.FindByProduct(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationReviewsDetail(paginationMeta, "success", "Successfully fetched review product", review)
	return so, nil
}

func (s *reviewHandleGrpc) FindByMerchant(ctx context.Context, request *pb.FindAllReviewMerchantRequest) (*pb.ApiResponsePaginationReviewDetail, error) {
	merchant_id := int(request.GetMerchantId())
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	reqService := requests.FindAllReviewByMerchant{
		MerchantID: merchant_id,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	review, totalRecords, err := s.reviewQueryService.FindByMerchant(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationReviewsDetail(paginationMeta, "success", "Successfully fetched review merchant", review)
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

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.reviewQueryService.FindByActive(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationReviewDeleteAt(paginationMeta, "success", "Successfully fetched active reviews", users)

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

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.reviewQueryService.FindByTrashed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationReviewDeleteAt(paginationMeta, "success", "Successfully fetched trashed categories", users)

	return so, nil
}

func (s *reviewHandleGrpc) Create(ctx context.Context, request *pb.CreateReviewRequest) (*pb.ApiResponseReview, error) {
	req := &requests.CreateReviewRequest{
		UserID:    int(request.GetUserId()),
		ProductID: int(request.GetProductId()),
		Rating:    int(request.GetRating()),
		Comment:   request.GetComment(),
	}

	if err := req.Validate(); err != nil {
		return nil, review_errors.ErrGrpcValidateCreateReview
	}

	review, err := s.reviewCommandService.CreateReview(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseReview("success", "Successfully created review", review), nil
}

func (s *reviewHandleGrpc) Update(ctx context.Context, request *pb.UpdateReviewRequest) (*pb.ApiResponseReview, error) {
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

	review, err := s.reviewCommandService.UpdateReview(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseReview("success", "Successfully updated review", review), nil
}

func (s *reviewHandleGrpc) TrashedReview(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, review_errors.ErrGrpcInvalidID
	}

	Review, err := s.reviewCommandService.TrashedReview(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseReviewDeleteAt("success", "Successfully trashed Review", Review)

	return so, nil
}

func (s *reviewHandleGrpc) RestoreReview(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, review_errors.ErrGrpcInvalidID
	}

	Review, err := s.reviewCommandService.RestoreReview(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseReviewDeleteAt("success", "Successfully restored Review", Review)

	return so, nil
}

func (s *reviewHandleGrpc) DeleteReviewPermanent(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, review_errors.ErrGrpcInvalidID
	}

	_, err := s.reviewCommandService.DeleteReviewPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseReviewDelete("success", "Successfully deleted Review permanently")

	return so, nil
}

func (s *reviewHandleGrpc) RestoreAllReview(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	_, err := s.reviewCommandService.RestoreAllReviews()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseReviewAll("success", "Successfully restore all Review")

	return so, nil
}

func (s *reviewHandleGrpc) DeleteAllReviewPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	_, err := s.reviewCommandService.DeleteAllReviewsPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseReviewAll("success", "Successfully delete Review permanen")

	return so, nil
}
