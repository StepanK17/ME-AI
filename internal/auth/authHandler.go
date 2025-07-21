package auth

import (
	"me-ai/configs"
	"me-ai/internal/models"
	"me-ai/pkg/jwt"
	"me-ai/pkg/req"
	"me-ai/pkg/res"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
}

type AuthHandler struct {
	*configs.Config
	AuthService *AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	userRepo := &models.UserRepository{}
	service := NewAuthService(userRepo)
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: service,
	}
	router.HandleFunc("/api/auth/login", handler.Login())
	router.HandleFunc("/api/auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Login: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := LoginResponse{
			Token: token,
		}
		res.Json(w, data, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Login: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := RegisterResponse{
			Token: token,
		}
		res.Json(w, data, 200)
	}
}
