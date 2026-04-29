package requests

import "github.com/go-playground/validator/v10"

type CreateRefreshToken struct {
	UserId    int    `json:"user_id" validate:"required,min=1"`
	Token     string `json:"token" validate:"required,min=1"`
	ExpiresAt string `json:"expires_at" validate:"required,min=1"`
}

type UpdateRefreshToken struct {
	UserId    int    `json:"user_id" validate:"required,min=1"`
	Token     string `json:"token" validate:"required,min=1"`
	ExpiresAt string `json:"expires_at" validate:"required,min=1"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=1"`
}

func (r *CreateRefreshToken) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}

func (r *UpdateRefreshToken) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}

func (r *RefreshTokenRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	if err != nil {
		return err
	}

	return nil
}
