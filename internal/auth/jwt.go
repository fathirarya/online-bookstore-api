package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey      string
	issuer         string
	expireDuration time.Duration
}

type JWTCustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWTService
func NewJWTService(cfg *JWTConfig) *JWTService {
	return &JWTService{
		secretKey:      cfg.SecretKey,
		issuer:         cfg.Issuer,
		expireDuration: cfg.ExpireDuration,
	}
}

// GenerateToken generates a JWT token with custom claims
func (j *JWTService) GenerateToken(userID int) (string, error) {
	claims := JWTCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken validates the JWT token string and returns the claims if valid
func (j *JWTService) ValidateToken(tokenString string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (j *JWTService) ExpireDuration() time.Duration {
	return j.expireDuration
}
