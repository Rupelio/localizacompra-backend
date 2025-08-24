package category

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type categoryHandler struct {
	service Service
}

func NewHandler(s Service) *categoryHandler {
	return &categoryHandler{
		service: s,
	}
}

func (h *categoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	categoryToCreate := Category{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	createdCategory, err := h.service.Create(r.Context(), categoryToCreate)
	if err != nil {
		log.Printf("Erro ao criar categoria: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCategory)
}

func (h *categoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "ID da categoria inválido", http.StatusBadRequest)
		return
	}
	category, err := h.service.GetByID(r.Context(), categoryID)
	if err != nil {
		log.Printf("Erro ao buscar categoria: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func (h *categoryHandler) PartialUpdate(w http.ResponseWriter, r *http.Request) {
	var req UpdateCategoryRequest

	categoryID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "ID da categoria inválido", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	err = h.service.PartialUpdate(r.Context(), categoryID, req)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *categoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "ID da categoria inválido", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), categoryID)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro ao deletar produto: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *categoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll(r.Context())
	if err != nil {
		log.Printf("Erro ao buscar categorias: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}
