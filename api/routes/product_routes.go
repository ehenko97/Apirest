package routes

import (
	http2 "github.com/ehenko97/apirest/internal/controller/http"
	"github.com/gorilla/mux"
)

// SetupProductRoutes регистрирует маршруты для работы с продуктами
func SetupProductRoutes(r *mux.Router, productController *http2.ProductController) {
	// GET получение всех продуктов
	r.HandleFunc("/products", productController.GetAllProducts).Methods("GET")

	// POST создание нового продукта
	r.HandleFunc("/products", productController.CreateProduct).Methods("POST")

	// GET получение всех продуктов по ID
	r.HandleFunc("/products/{id}", productController.GetProduct).Methods("GET")

	// PUT обновление данных продукта по ID
	r.HandleFunc("/products/{id}", productController.UpdateProduct).Methods("PUT")

	// DELETE  удаление продукта по ID
	r.HandleFunc("/products/{id}", productController.DeleteProduct).Methods("DELETE")

	// GET получение пользователя по ID и продукта
	r.HandleFunc("/users/{id}/products", productController.GetUserProducts).Methods("GET")
}
