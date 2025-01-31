package http

import (
	"encoding/json"
	"errors"
	"github.com/ehenko97/apirest/internal/entity"
	"github.com/ehenko97/apirest/internal/service"
	"net/http"
	"strconv"
	"strings"
)

// ProductController структура контроллера для обработки запросов продуктов.
type ProductController struct {
	productService service.ProductService
	userService    service.UserService
}

// NewProductController создает новый контроллер продуктов.
func NewProductController(productService service.ProductService, userService service.UserService) *ProductController {
	return &ProductController{
		productService: productService,
		userService:    userService, // Передаем userService
	}
}

// Структура для объединения данных о пользователе и продуктах
type UserProductsResponse struct {
	User     entity.User      `json:"user"`
	Products []entity.Product `json:"products"`
}

// CreateProduct создает новый продукт.
func (pc *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product entity.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdProduct, err := pc.productService.Create(r.Context(), product)
	if err != nil {
		http.Error(w, "Failed to create product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdProduct); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetAllProducts возвращает список всех продуктов.
func (pc *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := pc.productService.FindAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetProduct возвращает продукт по ID.
func (pc *ProductController) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := extractProductIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid product ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	product, err := pc.productService.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Product not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdateProduct обновляет продукт по ID.
func (pc *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := extractProductIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid product ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var updatedProduct entity.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	updatedProduct.ID = id

	product, err := pc.productService.Update(r.Context(), updatedProduct)
	if err != nil {
		http.Error(w, "Failed to update product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DeleteProduct удаляет продукт по ID.
func (pc *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := extractProductIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid product ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = pc.productService.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (pc *ProductController) GetUserProducts(w http.ResponseWriter, r *http.Request) {
	// Проверка, что метод GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлечение userID из URL
	userID, err := extractUserIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем информацию о пользователе через userService
	user, err := pc.userService.FindByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем все продукты для указанного пользователя
	products, err := pc.productService.FindByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Формируем ответ с информацией о пользователе и продуктах
	response := UserProductsResponse{
		User:     user,
		Products: products,
	}

	// Возвращаем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// extractProductIDFromURL извлекает ID продукта из URL
func extractProductIDFromURL(path string) (int, error) {
	// Убираем только префикс "/products/"
	parts := strings.Split(strings.TrimPrefix(path, "/products/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		return 0, errors.New("missing product ID")
	}

	// Преобразуем ID в число
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, errors.New("invalid product ID format")
	}

	return id, nil
}
