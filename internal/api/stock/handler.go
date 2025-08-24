package stock

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type stockItemHandler struct {
	service Service
}

func NewHandler(s Service) *stockItemHandler {
	return &stockItemHandler{
		service: s,
	}
}

func (h *stockItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	storeIDParam := chi.URLParam(r, "storeID")
	storeID, err := strconv.ParseInt(storeIDParam, 10, 64)
	if err != nil {
		http.Error(w, "ID da loja inválido", http.StatusBadRequest)
		return
	}

	productIDParam := chi.URLParam(r, "productID")
	productID, err := strconv.ParseInt(productIDParam, 10, 64)
	if err != nil {
		http.Error(w, "ID do produto inválido", http.StatusBadRequest)
		return
	}

	var req CreateStockItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	stockItemToCreate := StockItem{
		StoreID:   storeID,
		ProductID: productID,
		Price:     req.Price,
		Quantity:  req.Quantity,
		Sector:    req.Sector,
	}

	createdStockItem, err := h.service.Create(r.Context(), stockItemToCreate)
	if err != nil {
		log.Printf("Erro ao criar estoque do item.")
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdStockItem)
}

func (h *stockItemHandler) GetAllByStoreId(w http.ResponseWriter, r *http.Request) {
	storeIDParam := chi.URLParam(r, "storeID")
	storeID, err := strconv.ParseInt(storeIDParam, 10, 64)
	if err != nil {
		http.Error(w, "ID da loja inválido", http.StatusBadRequest)
		return
	}

	productStore, err := h.service.GetAllByStoreId(r.Context(), storeID)
	if err != nil {
		log.Printf("Erro ao buscar produtos: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(productStore)
}
