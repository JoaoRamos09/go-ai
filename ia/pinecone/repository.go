package pinecone

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/joaoramos09/go-ai/ia"
	"github.com/pinecone-io/go-pinecone/v3/pinecone"
)

type Repository struct {
	IndexConnection *pinecone.IndexConnection
}

func NewRepository(IndexConnection *pinecone.IndexConnection) *Repository {
	return &Repository{IndexConnection: IndexConnection}
}

func (r *Repository) Insert(ctx context.Context, documents []ia.Document) error {

	integratedRecords := make([]*pinecone.IntegratedRecord, len(documents))

	for i, document := range documents {
		integratedRecords[i] = &pinecone.IntegratedRecord{
			"id":       uuid.New().String(),
			"text":     document.Text,
			"category": document.Category,
		}
	}

	err := r.IndexConnection.UpsertRecords(ctx, integratedRecords)
	if err != nil {
		log.Printf("Error inserting documents: %v", err)
		return err
	}

	return nil
}

func (r *Repository) Retrieve(ctx context.Context, input string) ([]ia.Document, error) {
	res, err := r.IndexConnection.SearchRecords(ctx, &pinecone.SearchRecordsRequest{
		Query: pinecone.SearchRecordsQuery{
			TopK: 3,
			Inputs: &map[string]interface{}{
				"text": input,
			},
		},
		Fields: &[]string{"text", "category"},
	})
	if err != nil {
		log.Printf("Error retrieving documents: %v", err)
		return nil, err
	}

	documents := make([]ia.Document, len(res.Result.Hits))
	for i, hit := range res.Result.Hits {
		documents[i] = ia.Document{
			ID:       hit.Id,
			Text: hit.Fields["text"].(string),
			Category: hit.Fields["category"].(string),
		}
	}

	return documents, nil
}
