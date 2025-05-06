package auth_jwt

import (
	"fmt"
	"time"
	"log"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secretKey   string
	tokenExpiry time.Duration
	tokenIssuer string
	tokenAudience string
}

func NewService(secretKey string, tokenExpiry time.Duration, tokenIssuer string, tokenAudience string) *Service {
	return &Service{
		secretKey:   secretKey,
		tokenExpiry: tokenExpiry,
		tokenIssuer: tokenIssuer,
		tokenAudience: tokenAudience,
	}
}

func (s *Service) Generate(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"aud": s.tokenAudience,
		"exp": time.Now().Add(s.tokenExpiry).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": s.tokenIssuer,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.secretKey))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", err
	}
	return token, nil
}

func (s *Service) Validate(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(s.tokenIssuer),
		jwt.WithIssuer(s.tokenIssuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
