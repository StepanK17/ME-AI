package llm

import (
	"log"
	"net/http"

	"me-ai/internal/middleware"
	"me-ai/internal/models"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	llmService *LLMService
}

func NewWebSocketHandler(llmService *LLMService) *WebSocketHandler {
	return &WebSocketHandler{
		llmService: llmService,
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка обновления соединения: %v", err)
		return
	}
	defer conn.Close()

	email := middleware.GetUserEmail(r)
	userRepo := &models.UserRepository{}
	user, err := userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("User not found: %v", err)
		return
	}

	log.Println("Новое WebSocket соединение")

	for {
		var msg struct {
			Type           string `json:"type"`
			Content        string `json:"content"`
			Role           string `json:"role"`
			ConversationID int    `json:"conversation_id"`
		}
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Ошибка WebSocket: %v", err)
			}
			break
		}

		if msg.Type == "user_message" {
			if msg.ConversationID == 0 || msg.Content == "" {
				conn.WriteJSON(models.WebSocketMessage{
					Type:    "error",
					Content: "conversation_id и content обязательны",
					Role:    "system",
				})
				continue
			}

			msgRepo := &models.MessageRepository{}
			userMsg := &models.Message{
				ConversationID: msg.ConversationID,
				UserID:         user.ID,
				Content:        msg.Content,
				Role:           "user",
				Timestamp:      time.Now(),
			}
			_, err := msgRepo.Create(userMsg)
			if err != nil {
				log.Printf("Ошибка сохранения сообщения пользователя: %v", err)
			}

			userMsgOut := models.WebSocketMessage{
				Type:    "user_message",
				Content: msg.Content,
				Role:    "user",
			}
			conn.WriteJSON(userMsgOut)

			go h.handleLLMResponse(conn, msg.Content, msg.ConversationID, user.ID)
		}
	}
}

func (h *WebSocketHandler) handleLLMResponse(conn *websocket.Conn, message string, conversationID int, userID int) {

	typingMsg := models.WebSocketMessage{
		Type:    "typing",
		Content: "LLM думает...",
		Role:    "assistant",
	}
	conn.WriteJSON(typingMsg)

	var fullResponse string
	err := h.llmService.GenerateStreamResponse(nil, message, conversationID, func(chunk string) {
		fullResponse += chunk
		streamMsg := models.WebSocketMessage{
			Type:    "assistant_chunk",
			Content: chunk,
			Role:    "assistant",
		}
		conn.WriteJSON(streamMsg)
	})

	if err != nil {
		log.Printf("Ошибка генерации ответа: %v", err)
		errorMsg := models.WebSocketMessage{
			Type:    "error",
			Content: "Извините, произошла ошибка при генерации ответа",
			Role:    "system",
		}
		conn.WriteJSON(errorMsg)
		return
	}

	msgRepo := &models.MessageRepository{}
	llmMsg := &models.Message{
		ConversationID: conversationID,
		UserID:         userID,
		Content:        fullResponse,
		Role:           "assistant",
		Timestamp:      time.Now(),
	}
	_, err = msgRepo.Create(llmMsg)
	if err != nil {
		log.Printf("Ошибка сохранения сообщения LLM: %v", err)
	}

	finalMsg := models.WebSocketMessage{
		Type:    "assistant_complete",
		Content: fullResponse,
		Role:    "assistant",
	}
	conn.WriteJSON(finalMsg)
}
