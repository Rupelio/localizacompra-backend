package store

import (
	"encoding/json"
	"log"
	"net/http"
)

type storeHandler struct {
	service Service
}

func NewHandler(s Service) *storeHandler {
	return &storeHandler{
		service: s,
	}
}

func (h *storeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateStoreRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	storeToCreate := Store{
		Name:    req.Name,
		Address: req.Address,
		CNPJ:    req.CNPJ,
	}

	createdStore, err := h.service.Create(r.Context(), storeToCreate)
	if err != nil {
		log.Printf("Erro ao criar loja: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdStore)
}

func (h *storeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	stores, err := h.service.GetAll(r.Context())
	if err != nil {
		log.Printf("Erro ao buscar produtos: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stores)
}
