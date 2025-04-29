package database

import (
	"context"
	"log"
	"github.com/pinecone-io/go-pinecone/v3/pinecone"
)


func NewPinecone(ctx context.Context, indexName string, apiKey string, namespace string) *pinecone.IndexConnection {
	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: apiKey,
	})
	if err != nil {
		log.Fatal(err)
	}

	idxModel, err := pc.DescribeIndex(ctx, indexName)

	if err != nil {
		log.Fatal(err)
	}

	idxConnection, err := pc.Index(pinecone.NewIndexConnParams{Host: idxModel.Host, Namespace: namespace})

	if err != nil {
		log.Fatal(err)
	}


	return idxConnection
}
