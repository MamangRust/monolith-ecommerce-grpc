package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantBusinessCommandService struct {
	ctx                               context.Context
	trace                             trace.Tracer
	merchantQueryRepository           repository.MerchantQueryRepository
	merchantBusinessCommandRepository repository.MerchantBusinessCommandRepository
	logger                            logger.LoggerInterface
	mapping                           response_service.MerchantBusinessResponseMapper
	requestCounter                    *prometheus.CounterVec
	requestDuration                   *prometheus.HistogramVec
}

func NewMerchantBusinessCommandService(
	ctx context.Context,
	merchantQueryRepository repository.MerchantQueryRepository,
	merchantBusinessCommandRepository repository.MerchantBusinessCommandRepository,
	logger logger.LoggerInterface,
	mapping response_service.MerchantBusinessResponseMapper,
) *merchantBusinessCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_business_command_service_requests_total",
			Help: "Total number of requests to the MerchantBusinessCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_business_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantBusinessCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantBusinessCommandService{
		ctx:                               ctx,
		trace:                             otel.Tracer("merchant-business-command-service"),
		merchantQueryRepository:           merchantQueryRepository,
		merchantBusinessCommandRepository: merchantBusinessCommandRepository,
		logger:                            logger,
		mapping:                           mapping,
		requestCounter:                    requestCounter,
		requestDuration:                   requestDuration,
	}
}

func (s *merchantBusinessCommandService) CreateMerchant(req *requests.CreateMerchantBusinessInformationRequest) (*response.MerchantBusinessResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new merchant")

	merchant, err := s.merchantBusinessCommandRepository.CreateMerchantBusiness(req)

	if err != nil {
		s.logger.Error("Failed to create new merchant",
			zap.Error(err),
			zap.Any("request", req))

		return nil, merchantbusiness_errors.ErrFailedCreateMerchantBusiness
	}

	return s.mapping.ToMerchantBusinessResponse(merchant), nil
}

func (s *merchantBusinessCommandService) UpdateMerchant(req *requests.UpdateMerchantBusinessInformationRequest) (*response.MerchantBusinessResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating merchant", zap.Int("merchantID", *req.MerchantBusinessInfoID))

	merchant, err := s.merchantBusinessCommandRepository.UpdateMerchantBusiness(req)

	if err != nil {
		s.logger.Error("Failed to update merchant",
			zap.Error(err),
			zap.Any("request", req))

		return nil, merchantbusiness_errors.ErrFailedUpdateMerchantBusiness
	}

	return s.mapping.ToMerchantBusinessResponse(merchant), nil
}

func (s *merchantBusinessCommandService) TrashedMerchant(merchantID int) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantBusinessCommandRepository.TrashedMerchantBusiness(merchantID)

	if err != nil {
		s.logger.Error("Failed to move merchant to trash",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, merchantbusiness_errors.ErrFailedTrashedMerchantBusiness
	}

	return s.mapping.ToMerchantBusinessResponseDeleteAt(merchant), nil
}

func (s *merchantBusinessCommandService) RestoreMerchant(merchantID int) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantBusinessCommandRepository.RestoreMerchantBusiness(merchantID)

	if err != nil {
		s.logger.Error("Failed to restore merchant from trash",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return nil, merchantbusiness_errors.ErrFailedRestoreMerchantBusiness
	}

	return s.mapping.ToMerchantBusinessResponseDeleteAt(merchant), nil
}

func (s *merchantBusinessCommandService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting merchant permanently", zap.Int("merchantID", merchantID))

	success, err := s.merchantBusinessCommandRepository.DeleteMerchantBusinessPermanent(merchantID)

	if err != nil {
		s.logger.Error("Failed to permanently delete merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID))

		return false, merchantbusiness_errors.ErrFailedDeleteMerchantBusinessPermanent
	}

	return success, nil
}

func (s *merchantBusinessCommandService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed merchants")

	success, err := s.merchantBusinessCommandRepository.RestoreAllMerchantBusiness()

	if err != nil {
		s.logger.Error("Failed to restore all trashed merchants",
			zap.Error(err))

		return false, merchantbusiness_errors.ErrFailedRestoreAllMerchantBusiness
	}

	return success, nil
}

func (s *merchantBusinessCommandService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all merchants")

	success, err := s.merchantBusinessCommandRepository.DeleteAllMerchantBusinessPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed merchants",
			zap.Error(err))

		return false, merchantbusiness_errors.ErrFailedDeleteAllMerchantBusinessPermanent
	}

	return success, nil
}
