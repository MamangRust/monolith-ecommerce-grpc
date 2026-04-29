package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenExpired = errors.New("token expired")

//go:generate mockgen -source=token.go -destination=mocks/token.go
type TokenManager interface {
	GenerateToken(userId int, audience string) (string, error)
	ValidateToken(tokenString string) (string, error)
}

type Manager struct {
	secretKey []byte
}

func NewManager(secretKey string) (*Manager, error) {
	if secretKey == "" {
		return nil, errors.New("empty secret key")
	}
	return &Manager{secretKey: []byte(secretKey)}, nil
}

func (m *Manager) GenerateToken(userId int, audience string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(12 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expireTime),
		Subject:   strconv.Itoa(userId),
		Audience:  []string{audience},
	})

	return token.SignedString([]byte(m.secretKey))
}

func (m *Manager) ValidateToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrTokenExpired
		}
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil

}
