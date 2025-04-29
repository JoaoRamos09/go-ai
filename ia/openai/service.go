package openai

import (
	"context"
	"log"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/joaoramos09/go-ai/ia"
	"fmt"
)

type Service struct {
	apiKey string
}

func NewService(apiKey string) *Service {
	return &Service{apiKey: apiKey}
}

func (s *Service) Invoke(ctx context.Context, input string, model string) (string, error) {
	if model == "" {
		model = openai.ChatModelGPT4oMini
	}
	client := openai.NewClient(option.WithAPIKey(s.apiKey))
	cc, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(input),
		},
		Model: model,
	})

	if err != nil {
		log.Printf("[OPENAI] Error invoking OpenAI: %v", err)
		return "", err
	}

	log.Printf("[OPENAI] Successfully invoked OpenAI")
	return cc.Choices[0].Message.Content, nil
}

func (s *Service) InvokeWithSystemPrompt(ctx context.Context, input string, documents []ia.Document, model string) (string, error) {
	if model == "" {
		model = openai.ChatModelGPT4oMini
	}

	var documentsString string

	for _, document := range documents {
		documentsString += fmt.Sprintf("\n\nCategory: %s\n\nContent: %s", document.Category, document.Text)
	}

	systemPrompt := `
	You are a helpful assistant that can answer questions and help with tasks, 
	you are a expert in the field of AI and technology. 
	Analyze the following documents and answer the user's question based on the documents: ` + documentsString

	client := openai.NewClient(option.WithAPIKey(s.apiKey))
	cc, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(input),
		},
		Model: model,
	})

	if err != nil {
		log.Printf("[OPENAI] Error invoking OpenAI: %v", err)
		return "", err
	}

	log.Printf("[OPENAI] Successfully invoked OpenAI")
	return cc.Choices[0].Message.Content, nil
}
