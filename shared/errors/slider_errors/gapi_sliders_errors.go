package slider_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	"net/http"
)

var (
	ErrGrpcInvalidID = errors.NewGrpcError("invalid ID", http.StatusBadRequest)

	ErrGrpcValidateCreateSlider = errors.NewGrpcError("validation failed: invalid create slider request", http.StatusBadRequest)
	ErrGrpcValidateUpdateSlider = errors.NewGrpcError("validation failed: invalid update slider request", http.StatusBadRequest)
)
