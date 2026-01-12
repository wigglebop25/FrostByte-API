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

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name            string   `json:"name"`
		Description     string   `json:"description"`
		Price           float64  `json:"price"`
		ProductImageURI string   `json:"product_image_uri"`
		CategoryNames   []string `json:"categories"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product := domain.Product{
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		ProductImageURI: req.ProductImageURI,
	}

	if err := h.service.CreateProductWithCategories(&product, req.CategoryNames); err != nil {
		if strings.Contains(err.Error(), "1062") {
			http.Error(w, "The product name you are trying to update already exists. Try again.", http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "record not found") {
			http.Error(w, "You are trying to update a product categories with non existing, please try again", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name            string   `json:"name"`
		Description     string   `json:"description"`
		Price           float64  `json:"price"`
		ProductImageURI string   `json:"product_image_uri"`
		CategoryNames   []string `json:"categories"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product := domain.Product{
		ProductID:       uint(id),
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		ProductImageURI: req.ProductImageURI,
	}

	if err := h.service.UpdateProductWithCategories(&product, req.CategoryNames); err != nil {
		if strings.Contains(err.Error(), "1062") {
			http.Error(w, "The product name you are trying to update already exists. Try again.", http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "record not found") {
			http.Error(w, "You are trying to update a product categories with non existing, please try again", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := "Product successfully updated"
	var updates []string
	if req.Name != "" {
		updates = append(updates, "name")
	}
	if req.Description != "" {
		updates = append(updates, "description")
	}
	if req.Price != 0 {
		updates = append(updates, "price")
	}
	if req.ProductImageURI != "" {
		updates = append(updates, "image")
	}
	if req.CategoryNames != nil {
		updates = append(updates, "categories")
	}

	if len(updates) == 1 {
		msg = "Product " + updates[0] + " successfully updated"
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": msg})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteProduct(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
