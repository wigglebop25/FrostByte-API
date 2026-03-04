package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"frostbyte-api/internal/domain"
	"frostbyte-api/internal/service"
	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCategory(&category); err != nil {
		if strings.Contains(err.Error(), "1062") {
			http.Error(w, "Category already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	category, err := h.service.GetCategoryByID(uint(id))
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	category.CategoryID = uint(id)

	if err := h.service.UpdateCategory(&category); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCategory(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CategoryID uint   `json:"category_id"`
		ProductID  uint   `json:"product_id"`
		Category   string `json:"category"`
		Product    string `json:"product"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var catID, prodID uint = req.CategoryID, req.ProductID

	// Resolve Category ID by Name if provided
	if catID == 0 && req.Category != "" {
		cat, err := h.service.GetCategoryByName(req.Category)
		if err != nil {
			http.Error(w, "Category not found: "+req.Category, http.StatusNotFound)
			return
		}
		catID = cat.CategoryID
	}

	// Resolve Product ID by Name if provided
	if prodID == 0 && req.Product != "" {
		// We need to inject ProductService or use a helper.
		// For now, I'll assume CategoryService can handle this lookup or I'll add it to CategoryService.
		prod, err := h.service.GetProductByName(req.Product)
		if err != nil {
			http.Error(w, "Product not found: "+req.Product, http.StatusNotFound)
			return
		}
		prodID = prod.ProductID
	}

	if catID == 0 || prodID == 0 {
		http.Error(w, "category and product (names or ids) are required", http.StatusBadRequest)
		return
	}

	if err := h.service.AddProductToCategory(catID, prodID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CategoryHandler) RemoveProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CategoryID uint   `json:"category_id"`
		ProductID  uint   `json:"product_id"`
		Category   string `json:"category"`
		Product    string `json:"product"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var catID, prodID uint = req.CategoryID, req.ProductID

	if catID == 0 && req.Category != "" {
		cat, err := h.service.GetCategoryByName(req.Category)
		if err != nil {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		catID = cat.CategoryID
	}

	if prodID == 0 && req.Product != "" {
		prod, err := h.service.GetProductByName(req.Product)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		prodID = prod.ProductID
	}

	if catID == 0 || prodID == 0 {
		http.Error(w, "category and product (names or ids) are required", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveProductFromCategory(catID, prodID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
