package auth

import (
	"log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joaoramos09/go-ai/internal/errs"
	"golang.org/x/crypto/bcrypt"
)

type UseCase interface {
	EncodePassword(password string) ([]byte, error)
	MatchPassword(password string, hash []byte) bool
	TokenAuth
}

type TokenAuth interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type Service struct {
	TokenAuth
}

func NewService(TokenAuth TokenAuth) *Service {
	return &Service{
		TokenAuth: TokenAuth,
	}
}

func (s *Service) EncodePassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error encoding password: %v", err)
		return nil, errs.ErrDecodePassword
	}
	return hash, nil
}

func (s *Service) MatchPassword(password string, hash []byte) bool {

	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		log.Printf("Error matching password: %v", err)
		return false
	}
	return true
}
