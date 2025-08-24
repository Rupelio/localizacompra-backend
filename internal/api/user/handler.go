package user

import (
	"encoding/json"
	"errors"
	"localiza-compra/backend/internal/api/middleware"
	"log"
	"net/http"
	"time"
)

type userHandler struct {
	service Service
}

func NewHandler(s Service) *userHandler {
	return &userHandler{
		service: s,
	}
}

func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	userToCreate := User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
		Phone:        req.Phone,
	}

	createdUser, err := h.service.Create(r.Context(), userToCreate)
	if err != nil {
		log.Printf("Erro ao criar usuário: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	tokenString, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("Erro ao entrar na conta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *userHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userIDCtx := r.Context().Value(middleware.UserIDKey)

	userID, ok := userIDCtx.(int64)

	if !ok {
		http.Error(w, "ID do usuário inválido no contexto", http.StatusInternalServerError)
		return
	}

	user, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			http.Error(w, "Usuário não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *userHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusNoContent)
}
