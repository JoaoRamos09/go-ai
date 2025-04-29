package ia

import (
	"context"
)

type Repository interface {
	Retrieve(ctx context.Context, input string) ([]Document, error)
	Insert(ctx context.Context, documents []Document) error
}

