package merchant_document_errors

import (
	"github.com/labstack/echo/v4"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

func ErrApiFailedFindAllMerchantDocuments(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrInternal.WithMessage("Failed to find all merchant documents"), "")
}

func ErrApiInvalidMerchantDocumentID(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrBadRequest.WithMessage("Invalid merchant document ID"), "")
}

func ErrApiFailedFindByIdMerchantDocument(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrInternal.WithMessage("Failed to find merchant document by ID"), "")
}

func ErrApiFailedFindAllActiveMerchantDocuments(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrInternal.WithMessage("Failed to find active merchant documents"), "")
}

func ErrApiFailedFindAllTrashedMerchantDocuments(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrInternal.WithMessage("Failed to find trashed merchant documents"), "")
}

func ErrApiValidateCreateMerchantDocument(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrBadRequest.WithMessage("Validation failed for creating merchant document"), "")
}

func ErrApiFailedCreateMerchantDocument(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrInternal.WithMessage("Failed to create merchant document"), "")
}

func ErrApiValidateUpdateMerchantDocument(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrBadRequest.WithMessage("Validation failed for updating merchant document"), "")
}

func ErrApiFailedUpdateMerchantDocument(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrInternal.WithMessage("Failed to update merchant document"), "")
}

func ErrApiBindUpdateMerchantDocumentStatus(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrBadRequest.WithMessage("Failed to bind update merchant document status"), "")
}

func ErrApiFailedTrashMerchantDocument(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrInternal.WithMessage("Failed to trash merchant document"), "")
}

func ErrApiFailedRestoreMerchantDocument(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrInternal.WithMessage("Failed to restore merchant document"), "")
}

func ErrApiFailedDeleteMerchantDocumentPermanent(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrInternal.WithMessage("Failed to delete merchant document permanently"), "")
}

func ErrApiFailedRestoreAllMerchantDocuments(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrInternal.WithMessage("Failed to restore all merchant documents"), "")
}

func ErrApiFailedDeleteAllMerchantDocumentsPermanent(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrInternal.WithMessage("Failed to delete all merchant documents permanently"), "")
}

func ErrApiInvalidMerchantId(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrBadRequest.WithMessage("Invalid merchant ID"), "")
}

func ErrApiDocumentTypeRequired(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrBadRequest.WithMessage("Document type is required"), "")
}

func ErrApiDocumentFileRequired(c echo.Context) error {
	return errors.HandleApiError(c, errors.ErrBadRequest.WithMessage("Document file is required"), "")
}
