package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	IdentityError      IdentityErrorHandler
	KafkaError         KafkaErrorHandler
	LoginError         LoginErrorHandler
	MarshalError       MarshalErrorHandler
	PasswordError      PasswordErrorHandler
	PasswordResetError PasswordResetErrorHandler
	RandomString       RandomStringErrorHandler
	RegisterError      RegisterErrorHandler
	TokenError         TokenErrorHandler
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		IdentityError:      NewIdentityError(logger),
		KafkaError:         NewKafkaError(logger),
		LoginError:         NewLoginError(logger),
		MarshalError:       NewMarshalError(logger),
		PasswordError:      NewPasswordError(logger),
		PasswordResetError: NewPasswordResetError(logger),
		RandomString:       NewRandomStringError(logger),
		RegisterError:      NewRegisterError(logger),
		TokenError:         NewTokenError(logger),
	}
}
