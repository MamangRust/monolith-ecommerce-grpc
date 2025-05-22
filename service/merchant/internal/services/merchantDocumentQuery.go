package services

import (
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantdocument_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_document_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type merchantDocumentQueryService struct {
	ctx                             context.Context
	trace                           trace.Tracer
	merchantDocumentQueryRepository repository.MerchantDocumentQueryRepository
	logger                          logger.LoggerInterface
	mapping                         response_service.MerchantDocumentResponseMapper
	requestCounter                  *prometheus.CounterVec
	requestDuration                 *prometheus.HistogramVec
}

func NewMerchantDocumentQueryService(
	ctx context.Context,
	merchantDocumentQueryRepository repository.MerchantDocumentQueryRepository,
	logger logger.LoggerInterface,
	mapping response_service.MerchantDocumentResponseMapper,
) *merchantDocumentQueryService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_document_query_request_count",
		Help: "Number of merchant document query requests MerchantDocumentQueryService",
	}, []string{"status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_document_query_request_duration_seconds",
		Help:    "The duration of requests MerchantDocumentQueryService",
		Buckets: prometheus.DefBuckets,
	}, []string{"status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantDocumentQueryService{
		ctx:                             ctx,
		trace:                           otel.Tracer("merchant-document-query-service"),
		merchantDocumentQueryRepository: merchantDocumentQueryRepository,
		logger:                          logger,
		mapping:                         mapping,
		requestCounter:                  requestCounter,
		requestDuration:                 requestDuration,
	}
}

func (s *merchantDocumentQueryService) FindAll(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindAll")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching all merchant document records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindAllDocuments(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_MERCHANT_DOCUMENTS")

		s.logger.Error("Failed to find all merchant documents",
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find all merchant documents")

		status = "failed_to_find_all_merchant_documents"

		return nil, nil, merchantdocument_errors.ErrFailedFindAllMerchantDocuments
	}

	merchantResponse := s.mapping.ToMerchantDocumentsResponse(merchantDocuments)

	s.logger.Debug("Merchant document records found",
		zap.Int("total", *total),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	return merchantResponse, total, nil
}

func (s *merchantDocumentQueryService) FindById(merchant_id int) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("Finding merchant document by ID", zap.Int("merchant_id", merchant_id))

	merchantDocument, err := s.merchantDocumentQueryRepository.FindById(merchant_id)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT_DOCUMENT_BY_ID")

		s.logger.Error("Failed to find merchant document by ID",
			zap.Error(err),
			zap.Int("merchant_id", merchant_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find merchant document by ID")

		status = "failed_to_find_merchant_document_by_id"

		return nil, merchantdocument_errors.ErrFailedFindMerchantDocumentById
	}

	merchantResponse := s.mapping.ToMerchantDocumentResponse(merchantDocument)

	s.logger.Debug("Merchant document found by ID", zap.Int("merchant_id", merchant_id))

	return merchantResponse, nil
}

func (s *merchantDocumentQueryService) FindByActive(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByActive")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching all merchant document active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindByActive(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_BY_ACTIVE_MERCHANT_DOCUMENTS")

		s.logger.Error("Failed to find by active merchant documents",
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find by active merchant documents")

		status = "failed_to_find_by_active_merchant_documents"

		return nil, nil, merchantdocument_errors.ErrFailedFindActiveMerchantDocuments
	}

	merchantResponse := s.mapping.ToMerchantDocumentsResponse(merchantDocuments)

	s.logger.Debug("Merchant document records found",
		zap.Int("total", *total),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	return merchantResponse, total, nil
}

func (s *merchantDocumentQueryService) FindByTrashed(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashed", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByTrashed")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching fetched trashed merchant documents",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindByTrashed(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_BY_TRASHED_MERCHANT_DOCUMENTS")

		s.logger.Error("Failed to find by trashed merchant documents",
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find by trashed merchant documents")

		status = "failed_to_find_by_trashed_merchant_documents"

		return nil, nil, merchantdocument_errors.ErrFailedFindTrashedMerchantDocuments
	}

	merchantResponse := s.mapping.ToMerchantDocumentsResponseDeleteAt(merchantDocuments)

	s.logger.Debug("Merchant document records found",
		zap.Int("total", *total),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	return merchantResponse, total, nil
}

func (s *merchantDocumentQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
