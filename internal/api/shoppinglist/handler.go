package shoppinglist

import (
	"encoding/json"
	"errors"
	"localiza-compra/backend/internal/api/middleware"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type shoppingHandler struct {
	service Service
}

func NewHandler(s Service) *shoppingHandler {
	return &shoppingHandler{
		service: s,
	}
}

func (h *shoppingHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	listIDParam := chi.URLParam(r, "listID")
	listID, err := strconv.ParseInt(listIDParam, 10, 64)
	if err != nil {
		http.Error(w, "ID da lista inválido", http.StatusBadRequest)
		return
	}

	userIDCtx := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDCtx.(int64)

	if !ok {
		http.Error(w, "ID do usuário inválido no contexto", http.StatusInternalServerError)
		return
	}

	var req CreateShoppingListItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	itemToCreate := ShoppingListItem{
		ShoppingListID: listID,
		ProductID:      req.ProductID,
		Quantity:       req.Quantity,
	}

	createdList, err := h.service.CreateItem(r.Context(), userID, itemToCreate)

	if err != nil {
		if errors.Is(err, ErrShoppingListNotFound) {
			http.Error(w, "Lista não encontrada", http.StatusNotFound)
			return
		}
		if err.Error() == "não autorizado: você não é o dono desta lista" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		log.Printf("Erro ao criar item: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdList)
}

func (h *shoppingHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)

	if !ok {
		http.Error(w, "ID do usuário inválido no contexto", http.StatusInternalServerError)
		return
	}

	var req CreateShoppingListRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo da requisição inválida", http.StatusBadRequest)
		return
	}

	listToCreate := ShoppingList{
		UserID: userID,
		Name:   req.Name,
	}

	createdList, err := h.service.CreateList(r.Context(), listToCreate)
	if err != nil {
		log.Printf("Erro ao criar lista de compras: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdList)
}

func (h *shoppingHandler) GetAllByUserID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)

	if !ok {
		http.Error(w, "ID do usuário inválido no contexto", http.StatusInternalServerError)
		return
	}

	lists, err := h.service.GetAllByUserID(r.Context(), userID)
	if err != nil {
		log.Printf("Erro ao buscar lists: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lists)
}

func (h *shoppingHandler) GetAllItemsByListID(w http.ResponseWriter, r *http.Request) {
	listIDParam := chi.URLParam(r, "listID")
	listID, err := strconv.ParseInt(listIDParam, 10, 64)
	if err != nil {
		http.Error(w, "ID da lista inválido", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "ID do usuário inválido no contexto", http.StatusInternalServerError)
		return
	}

	list, err := h.service.GetAllItemsByListID(r.Context(), userID, listID)
	if err != nil {
		log.Printf("Erro ao buscar lists: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (h *shoppingHandler) UpdateItemStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "ID do usuário inválido no contexto", http.StatusInternalServerError)
		return
	}

	listIDParam := chi.URLParam(r, "listID")
	listID, err := strconv.ParseInt(listIDParam, 10, 64)
	if err != nil {
		http.Error(w, "ID da lista inválido", http.StatusBadRequest)
		return
	}

	itemIDParam := chi.URLParam(r, "itemID")
	itemID, err := strconv.ParseInt(itemIDParam, 10, 64)
	if err != nil {
		http.Error(w, "ID do produto inválido", http.StatusBadRequest)
		return
	}

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateItemStatus(r.Context(), userID, listID, itemID, req.IsChecked)
	if err != nil {
		log.Printf("Erro ao buscar lists: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
