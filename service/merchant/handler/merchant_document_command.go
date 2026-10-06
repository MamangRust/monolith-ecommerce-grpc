package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-merchant/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	"google.golang.org/protobuf/types/known/emptypb"

	pbmerchant_document "github.com/MamangRust/monolith-ecommerce-pb/merchant_document"
)

type merchantDocumentCommandHandler struct {
	pbmerchant_document.UnimplementedMerchantDocumentCommandServiceServer
	merchantDocumentCommand service.MerchantDocumentCommandService
	logger                  logger.LoggerInterface
}

func NewMerchantDocumentCommandHandler(svc service.MerchantDocumentCommandService, logger logger.LoggerInterface) pbmerchant_document.MerchantDocumentCommandServiceServer {
	return &merchantDocumentCommandHandler{
		merchantDocumentCommand: svc,
		logger:                  logger,
	}
}

func (s *merchantDocumentCommandHandler) Create(ctx context.Context, req *pbmerchant_document.CreateMerchantDocumentRequest) (*pbmerchant_document.ApiResponseMerchantDocument, error) {
	request := requests.CreateMerchantDocumentRequest{
		MerchantID:   int(req.GetMerchantId()),
		DocumentType: req.GetDocumentType(),
		DocumentUrl:  req.GetDocumentUrl(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateCreateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.Create(ctx, &request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_document.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully created merchant document",
		Data:    mapToProtoMerchantDocumentResponse(document),
	}, nil
}

func (s *merchantDocumentCommandHandler) Update(ctx context.Context, req *pbmerchant_document.UpdateMerchantDocumentRequest) (*pbmerchant_document.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	request := requests.UpdateMerchantDocumentRequest{
		DocumentID:   &id,
		MerchantID:   int(req.GetMerchantId()),
		DocumentType: req.GetDocumentType(),
		DocumentUrl:  req.GetDocumentUrl(),
		Status:       req.GetStatus(),
		Note:         req.GetNote(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcFailedUpdateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.Update(ctx, &request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_document.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully updated merchant document",
		Data:    mapToProtoMerchantDocumentResponse(document),
	}, nil
}

func (s *merchantDocumentCommandHandler) UpdateStatus(ctx context.Context, req *pbmerchant_document.UpdateMerchantDocumentStatusRequest) (*pbmerchant_document.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	request := requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: &id,
		MerchantID: int(req.GetMerchantId()),
		Status:     req.GetStatus(),
		Note:       req.GetNote(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcFailedUpdateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.UpdateStatus(ctx, &request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_document.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully updated merchant document status",
		Data:    mapToProtoMerchantDocumentResponse(document),
	}, nil
}

func (s *merchantDocumentCommandHandler) Trashed(ctx context.Context, req *pbmerchant_document.TrashedMerchantDocumentRequest) (*pbmerchant_document.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_document.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully trashed merchant document",
		Data:    mapToProtoMerchantDocumentResponse(document),
	}, nil
}

func (s *merchantDocumentCommandHandler) Restore(ctx context.Context, req *pbmerchant_document.RestoreMerchantDocumentRequest) (*pbmerchant_document.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_document.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully restored merchant document",
		Data:    mapToProtoMerchantDocumentResponse(document),
	}, nil
}

func (s *merchantDocumentCommandHandler) DeletePermanent(ctx context.Context, req *pbmerchant_document.DeleteMerchantDocumentPermanentRequest) (*pbmerchant_document.ApiResponseMerchantDocumentDelete, error) {
	id := int(req.GetDocumentId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	_, err := s.merchantDocumentCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_document.ApiResponseMerchantDocumentDelete{
		Status:  "success",
		Message: "Successfully permanently deleted merchant document",
	}, nil
}

func (s *merchantDocumentCommandHandler) RestoreAll(ctx context.Context, _ *emptypb.Empty) (*pbmerchant_document.ApiResponseMerchantDocumentAll, error) {
	_, err := s.merchantDocumentCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_document.ApiResponseMerchantDocumentAll{
		Status:  "success",
		Message: "Successfully restored all merchant documents",
	}, nil
}

func (s *merchantDocumentCommandHandler) DeleteAllPermanent(ctx context.Context, _ *emptypb.Empty) (*pbmerchant_document.ApiResponseMerchantDocumentAll, error) {
	_, err := s.merchantDocumentCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_document.ApiResponseMerchantDocumentAll{
		Status:  "success",
		Message: "Successfully permanently deleted all merchant documents",
	}, nil
}
