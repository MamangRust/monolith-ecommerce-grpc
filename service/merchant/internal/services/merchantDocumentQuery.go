package services

import (
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
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
	errorhandler                    errorhandler.MerchantDocumentQueryErrorHandler
	mencache                        mencache.MerchantDocumentQueryCache
	trace                           trace.Tracer
	merchantDocumentQueryRepository repository.MerchantDocumentQueryRepository
	logger                          logger.LoggerInterface
	mapping                         response_service.MerchantDocumentResponseMapper
	requestCounter                  *prometheus.CounterVec
	requestDuration                 *prometheus.HistogramVec
}

func NewMerchantDocumentQueryService(
	errorhandler errorhandler.MerchantDocumentQueryErrorHandler,
	mencache mencache.MerchantDocumentQueryCache,
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
		errorhandler:                    errorhandler,
		mencache:                        mencache,
		trace:                           otel.Tracer("merchant-document-query-service"),
		merchantDocumentQueryRepository: merchantDocumentQueryRepository,
		logger:                          logger,
		mapping:                         mapping,
		requestCounter:                  requestCounter,
		requestDuration:                 requestDuration,
	}
}

func (s *merchantDocumentQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantDocuments(ctx, req); found {
		logSuccess("Successfully fetched merchant documents from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindAllDocuments(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_MERCHANT_DOCUMENT", span, &status, zap.Error(err))
	}

	merchantResponse := s.mapping.ToMerchantDocumentsResponse(merchantDocuments)

	s.mencache.SetCachedMerchantDocuments(ctx, req, merchantResponse, total)

	logSuccess("Successfully fetched merchant documents from database", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return merchantResponse, total, nil
}

func (s *merchantDocumentQueryService) FindById(ctx context.Context, merchant_id int) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("merchantDocument.id", merchant_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMerchantDocument(ctx, merchant_id); found {
		logSuccess("Successfully fetched merchant document from cache", zap.Int("merchantDocument.id", merchant_id))

		return data, nil
	}

	merchantDocument, err := s.merchantDocumentQueryRepository.FindById(ctx, merchant_id)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_DOCUMENT_BY_ID", span, &status, merchantdocument_errors.ErrMerchantDocumentNotFoundRes, zap.Int("merchantDocument.id", merchant_id))
	}

	merchantResponse := s.mapping.ToMerchantDocumentResponse(merchantDocument)

	s.mencache.SetCachedMerchantDocument(ctx, merchantResponse)

	logSuccess("Successfully fetched merchant document from database", zap.Int("merchantDocument.id", merchant_id))

	return merchantResponse, nil
}

func (s *merchantDocumentQueryService) FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantDocuments(ctx, req); found {
		logSuccess("Successfully fetched merchant documents from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_BY_ACTIVE_MERCHANT_DOCUMENTS", span, &status, zap.Error(err))
	}

	merchantResponse := s.mapping.ToMerchantDocumentsResponse(merchantDocuments)

	s.mencache.SetCachedMerchantDocuments(ctx, req, merchantResponse, total)

	logSuccess("Successfully fetched merchant documents from database", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return merchantResponse, total, nil
}

func (s *merchantDocumentQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantDocumentsTrashed(ctx, req); found {
		logSuccess("Successfully fetched merchant documents from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_MERCHANT_DOCUMENTS", span, &status, merchantdocument_errors.ErrFailedFindTrashedMerchantDocuments, zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
	}

	merchantResponse := s.mapping.ToMerchantDocumentsResponseDeleteAt(merchantDocuments)

	s.mencache.SetCachedMerchantDocumentsTrashed(ctx, req, merchantResponse, total)

	logSuccess("Successfully fetched merchant documents from database", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return merchantResponse, total, nil
}

func (s *merchantDocumentQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	ctx, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return ctx, span, end, status, logSuccess
}

func (s *merchantDocumentQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *merchantDocumentQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
