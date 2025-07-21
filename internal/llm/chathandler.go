package llm

import (
	"encoding/json"
	"fmt"
	"log"
	"me-ai/internal/middleware"
	"me-ai/internal/models"
	"net/http"
	"time"
)

type ChatHandler struct {
	llmService *LLMService
}

func NewChatHandler(llmService *LLMService) *ChatHandler {
	return &ChatHandler{
		llmService: llmService,
	}
}

func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	email := middleware.GetUserEmail(r)
	userRepo := &models.UserRepository{}
	user, err := userRepo.FindByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	var req struct {
		Message        string `json:"message"`
		ConversationID int    `json:"conversation_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	if req.Message == "" || req.ConversationID == 0 {
		http.Error(w, "Сообщение и conversation_id обязательны", http.StatusBadRequest)
		return
	}

	msgRepo := &models.MessageRepository{}
	userMsg := &models.Message{
		ConversationID: req.ConversationID,
		UserID:         user.ID,
		Content:        req.Message,
		Role:           "user",
		Timestamp:      time.Now(),
	}
	_, err = msgRepo.Create(userMsg)
	if err != nil {
		log.Printf("Ошибка сохранения сообщения пользователя: %v", err)
	}

	response, err := h.llmService.GenerateResponse(r.Context(), req.Message, req.ConversationID)
	if err != nil {
		log.Printf("Ошибка получения ответа от LLM: %v", err)
		http.Error(w, "Ошибка генерации ответа", http.StatusInternalServerError)
		return
	}

	llmMsg := &models.Message{
		ConversationID: req.ConversationID,
		UserID:         user.ID,
		Content:        response,
		Role:           "assistant",
		Timestamp:      time.Now(),
	}
	_, err = msgRepo.Create(llmMsg)
	if err != nil {
		log.Printf("Ошибка сохранения сообщения LLM: %v", err)
	}

	chatResponse := models.ChatResponse{
		Message:   response,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatResponse)
}

func (h *ChatHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	email := middleware.GetUserEmail(r)
	userRepo := &models.UserRepository{}
	user, err := userRepo.FindByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	repo := &models.ConversationRepository{}
	convos, err := repo.ListByUser(user.ID)
	if err != nil {
		http.Error(w, "Failed to get conversations", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(convos)
}

func (h *ChatHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	email := middleware.GetUserEmail(r)
	userRepo := &models.UserRepository{}
	user, err := userRepo.FindByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	var req models.CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	convo := &models.Conversation{
		UserID: user.ID,
		Title:  req.Title,
	}

	repo := &models.ConversationRepository{}
	created, err := repo.Create(convo)
	if err != nil {
		http.Error(w, "Failed to create conversation", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

func (h *ChatHandler) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	email := middleware.GetUserEmail(r)
	userRepo := &models.UserRepository{}
	user, err := userRepo.FindByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	var req models.DeleteConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	repo := &models.ConversationRepository{}
	if err := repo.Delete(req.ID, user.ID); err != nil {
		http.Error(w, "Failed to delete conversation", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ChatHandler) RenameConversation(w http.ResponseWriter, r *http.Request) {
	email := middleware.GetUserEmail(r)
	userRepo := &models.UserRepository{}
	user, err := userRepo.FindByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	var req struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.ID == 0 || req.Title == "" {
		http.Error(w, "id и title обязательны", http.StatusBadRequest)
		return
	}
	repo := &models.ConversationRepository{}
	if err := repo.UpdateTitle(req.ID, user.ID, req.Title); err != nil {
		log.Printf("Ошибка обновления названия чата: %v", err)
		http.Error(w, "Failed to rename conversation", http.StatusInternalServerError)
		return
	}
	log.Printf("Переименование чата: user_id=%d, chat_id=%d, title=%s", user.ID, req.ID, req.Title)
	w.WriteHeader(http.StatusNoContent)
}

func (h *ChatHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	convoID := r.URL.Query().Get("conversation_id")
	if convoID == "" {
		http.Error(w, "Missing conversation_id", http.StatusBadRequest)
		return
	}
	var id int
	if _, err := fmt.Sscanf(convoID, "%d", &id); err != nil {
		http.Error(w, "Invalid conversation_id", http.StatusBadRequest)
		return
	}
	repo := &models.MessageRepository{}
	msgs, err := repo.ListByConversation(id)
	if err != nil {
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(msgs)
}

func (h *ChatHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	email := middleware.GetUserEmail(r)
	userRepo := &models.UserRepository{}
	user, err := userRepo.FindByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	var req models.DeleteMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	repo := &models.MessageRepository{}
	if err := repo.Delete(req.ID, user.ID); err != nil {
		http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
