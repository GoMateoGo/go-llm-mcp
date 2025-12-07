package main

import (
	"context"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type ChatOpenAI struct {
	Ctx          context.Context
	ModelName    string
	SystemPrompt string
	RagContext   string
	Tools        []mcp.Tool
	LLM          openai.Client
}

type LLMOption func(*ChatOpenAI)

func WithSystemPrompt(prompt string) LLMOption {
	return func(llm *ChatOpenAI) {
		llm.SystemPrompt = prompt
	}
}

func WithRagContext(context string) LLMOption {
	return func(llm *ChatOpenAI) {
		llm.RagContext = context
	}
}

func WithTools(tools []mcp.Tool) LLMOption {
	return func(llm *ChatOpenAI) {
		llm.Tools = tools
	}
}

func NewChatOpenAI(ctx context.Context, modelName string, opts ...LLMOption) *ChatOpenAI {
	if modelName == "" {
		panic("modelName is required")
	}

	var (
		apiKey  = os.Getenv("OPENAI_API")
		baseURL = os.Getenv("OPENAI_BASE_URL")
	)
	if apiKey == "" {
		panic("OPENAI_API is required")
	}

	clientOptions := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}

	if baseURL != "" {
		clientOptions = append(clientOptions, option.WithBaseURL(baseURL))
	}

	client := openai.NewClient(clientOptions...)

	llm := &ChatOpenAI{
		Ctx:       ctx,
		ModelName: modelName,
		LLM:       client,
	}

	for _, opt := range opts {
		opt(llm)
	}

	return llm
}

func main() {
	ctx := context.Background()
	_ = NewChatOpenAI(ctx, "gpt-3.5-turbo", WithSystemPrompt("You are a helpful assistant."))
}
