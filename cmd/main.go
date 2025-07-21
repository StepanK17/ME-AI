package main

import (
	"fmt"
	"me-ai/configs"
	"me-ai/internal/llm"

	"me-ai/internal/auth"
	"me-ai/internal/middleware"
	"me-ai/pkg/db"
	"net/http"
)

func main() {
	cfg := configs.LoadConfig()
	if err := db.Init(); err != nil {
		panic(err)
	}

	llmService := llm.NewLLMService(cfg.LLM.URL, cfg.LLM.ApiKey)
	chatHandler := llm.NewChatHandler(llmService)
	wsHandler := llm.NewWebSocketHandler(llmService)

	router := http.NewServeMux()

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{Config: cfg})

	jwtMw := middleware.JWTAuth(cfg.Auth.Secret)
	corsMw := middleware.CORS

	protected := http.NewServeMux()
	protected.HandleFunc("/api/chat", chatHandler.HandleChat)
	protected.HandleFunc("/api/ws", wsHandler.HandleWebSocket)
	protected.HandleFunc("/api/conversations", chatHandler.ListConversations)         // GET
	protected.HandleFunc("/api/conversations/create", chatHandler.CreateConversation) // POST
	protected.HandleFunc("/api/conversations/delete", chatHandler.DeleteConversation) // POST
	protected.HandleFunc("/api/conversations/rename", chatHandler.RenameConversation) // POST
	protected.HandleFunc("/api/messages", chatHandler.ListMessages)                   // GET
	protected.HandleFunc("/api/messages/delete", chatHandler.DeleteMessage)           // POST

	router.Handle("/api/", corsMw(jwtMw(protected)))

	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
