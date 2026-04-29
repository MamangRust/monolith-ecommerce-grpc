package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-banner/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type bannerQueryHandler struct {
	pb.UnimplementedBannerQueryServiceServer
	BannerQuery service.BannerQueryService
	logger      logger.LoggerInterface
}

func NewBannerQueryHandler(svc service.BannerQueryService, logger logger.LoggerInterface) BannerQueryHandler {
	return &bannerQueryHandler{
		BannerQuery: svc,
		logger:      logger,
	}
}

func (s *bannerQueryHandler) FindAll(ctx context.Context, request *pb.FindAllBannerRequest) (*pb.ApiResponsePaginationBanner, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllBanner{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	banners, totalRecords, err := s.BannerQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoBanners := make([]*pb.BannerResponse, len(banners))
	for i, banner := range banners {
		protoBanners[i] = mapToProtoBannerResponse(banner)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationBanner{
		Status:     "success",
		Message:    "Successfully fetched banners",
		Data:       protoBanners,
		Pagination: paginationMeta,
	}, nil
}

func (s *bannerQueryHandler) FindById(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBanner, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, banner_errors.ErrGrpcBannerInvalidId
	}

	banner, err := s.BannerQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseBanner{
		Status:  "success",
		Message: "Successfully fetched banner",
		Data:    mapToProtoBannerResponse(banner),
	}, nil
}

func (s *bannerQueryHandler) FindByActive(ctx context.Context, request *pb.FindAllBannerRequest) (*pb.ApiResponsePaginationBannerDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllBanner{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	banners, totalRecords, err := s.BannerQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoBanners := make([]*pb.BannerResponseDeleteAt, len(banners))
	for i, banner := range banners {
		protoBanners[i] = mapToProtoBannerResponseDeleteAt(banner)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationBannerDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active banners",
		Data:       protoBanners,
		Pagination: paginationMeta,
	}, nil
}

func (s *bannerQueryHandler) FindByTrashed(ctx context.Context, request *pb.FindAllBannerRequest) (*pb.ApiResponsePaginationBannerDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllBanner{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	banners, totalRecords, err := s.BannerQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoBanners := make([]*pb.BannerResponseDeleteAt, len(banners))
	for i, banner := range banners {
		protoBanners[i] = mapToProtoBannerResponseDeleteAt(banner)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationBannerDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed banners",
		Data:       protoBanners,
		Pagination: paginationMeta,
	}, nil
}
