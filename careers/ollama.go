package careers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model string `json:"model"`
	Stream bool `json:"stream"`
	Messages []Message `json:"messages"`
	Format json.RawMessage `json:"format"`
	Think bool `json:"think"`
}

type ChatResponse struct {
	Message Message `json:"message"`
}

const ollamaUrl = "http://localhost:11434"
const model = "gemma4:e2b"
const chatEndpoint = "/api/chat"

func Chat(userPrompt, systemPrompt string, schema json.RawMessage) (string, error) {
	req := ChatRequest {
		Model: model,
		Stream: false,
		Messages: []Message {
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Format: schema,
		Think: false,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal ChatRequest into json: %w", err)
	}

	res, err := http.Post(ollamaUrl + chatEndpoint, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return "", fmt.Errorf("Failed to reach ollama: %w", err)
	}
	defer res.Body.Close()

	var chatRes ChatResponse
	err = json.NewDecoder(res.Body).Decode(&chatRes)
	if err != nil {
		return "", fmt.Errorf("Failed to decode ChatResponse: %w", err)
	}

	return chatRes.Message.Content, nil
}
