package ia

import (
	"context"

	"github.com/joaoramos09/go-ai/internal/errs"
)

type UseCase interface {
	LLMUseCase
	Insert(ctx context.Context, documents []Document) error
	Query(ctx context.Context, input string) ([]Document, error)
}

type LLMUseCase interface {
	Invoke(ctx context.Context, input string, model string) (string, error)
	InvokeWithSystemPrompt(ctx context.Context, input string, documents []Document, model string) (string, error)
}

type Service struct {
	LLMUseCase
	repo Repository
}

func NewService(llm LLMUseCase, repo Repository) *Service {
	return &Service{LLMUseCase: llm, repo: repo}
}

func (s *Service) Insert(ctx context.Context, documents []Document) error {
	if err := s.repo.Insert(ctx, documents); err != nil {
		return errs.ErrAIInsert
	}
	return nil
}

func (s *Service) Query(ctx context.Context, input string) ([]Document, error) {
	documents, err := s.repo.Retrieve(ctx, input)
	if err != nil {
		return nil, errs.ErrAISearch
	}

	return documents, nil
}












