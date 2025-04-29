package user

import (
	"context"

	"github.com/joaoramos09/go-ai/internal/errs"
)

type contextKey string
const UserIDKey = contextKey("user_id")

type UseCase interface {
	Get(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id int) error
	GetFromContext(ctx context.Context) (*User, error)
}

type Service struct {
	userRepository Repository
}

func NewService(userRepository Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) Get(ctx context.Context, id int) (*User, error) {
	if id == 0 || id < 0 {
		return nil, errs.ErrInvalidUserID
	}
	return s.userRepository.Get(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.userRepository.GetByEmail(ctx, email)
}

func (s *Service) Create(ctx context.Context, user *User) (*User, error) {

	_, err := s.GetByEmail(ctx, user.Email)
	if err == nil {
		return nil, errs.ErrUserAlreadyExists
	}

	u, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id == 0 || id < 0 {
		return errs.ErrInvalidUserID
	}
	return s.userRepository.Delete(ctx, id)
}

func (s *Service) GetFromContext(ctx context.Context) (*User, error) {
	userID := ctx.Value(UserIDKey).(int)
	return s.Get(ctx, userID)
}
