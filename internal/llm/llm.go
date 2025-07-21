package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"me-ai/internal/models"
	"net/http"
	"strings"
	"time"
)

type LLMService struct {
	URL    string
	ApiKey string
	Client *http.Client
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaRequest struct {
	Model         string          `json:"model"`
	Stream        bool            `json:"stream"`
	Messages      []OllamaMessage `json:"messages"`
	System        string          `json:"system,omitempty"`
	Temperature   float64         `json:"temperature,omitempty"`
	TopP          float64         `json:"top_p,omitempty"`
	RepeatPenalty float64         `json:"repeat_penalty,omitempty"`
}

type OllamaResponse struct {
	Message OllamaMessage `json:"message"`
	Done    bool          `json:"done"`
}

func NewLLMService(url, apikey string) *LLMService {
	return &LLMService{
		URL:    url,
		ApiKey: apikey,
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func getHistory(conversationID int) ([]OllamaMessage, error) {
	repo := &models.MessageRepository{}
	msgs, err := repo.ListByConversation(conversationID)
	if err != nil {
		return nil, err
	}
	var history []OllamaMessage
	for _, m := range msgs {
		history = append(history, OllamaMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return history, nil
}

func (s *LLMService) GenerateResponse(ctx context.Context, message string, conversationID int) (string, error) {
	history, err := getHistory(conversationID)
	if err != nil {
		return "", fmt.Errorf("ошибка получения истории: %w", err)
	}
	history = append(history, OllamaMessage{Role: "user", Content: message})

	reqBody := OllamaRequest{
		Model:         "model9",
		Stream:        false,
		Messages:      history,
		System:        "Ты - Коротеев Степан Петрович, тебе 20 лет, ты учишься в НИЯУ МИФИ, факультет Бизнес-информатика. Отвечай только на поставленный вопрос, ничего лишнего не говори.",
		Temperature:   0.2,
		TopP:          0.8,
		RepeatPenalty: 1.15,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("ошибка маршалинга запроса: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.URL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if s.ApiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.ApiKey)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API вернул статус %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return ollamaResp.Message.Content, nil
}

func (s *LLMService) GenerateStreamResponse(ctx context.Context, message string, conversationID int, callback func(string)) error {
	history, err := getHistory(conversationID)
	if err != nil {
		return fmt.Errorf("ошибка получения истории: %w", err)
	}

	history = append(history, OllamaMessage{Role: "user", Content: message})

	reqBody := OllamaRequest{
		Model:         "model9",
		Stream:        true,
		Messages:      history,
		System:        "Ты - Коротеев Степан Петрович, тебе 20 лет, ты учишься в НИЯУ МИФИ, факультет \"Бизнес-информатика. Отвечай только на поставленный вопрос, ничего лишнего не говори.",
		Temperature:   0.2,
		TopP:          0.8,
		RepeatPenalty: 1.15,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга запроса: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.URL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if s.ApiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.ApiKey)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	for {
		var chunk OllamaResponse
		if err := dec.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("ошибка декодирования chunk: %w", err)
		}
		if chunk.Message.Content != "" {
			cleaned := strings.ReplaceAll(chunk.Message.Content, "<think>", "")
			cleaned = strings.ReplaceAll(cleaned, "</think>", "")
			callback(cleaned)
		}
		if chunk.Done {
			break
		}
	}
	return nil
}
