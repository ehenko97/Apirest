package routes

import (
	http2 "github.com/ehenko97/apirest/internal/controller/http"
	"github.com/gorilla/mux"
)

// NewRouter создаёт и возвращает маршрутизатор с зарегистрированными маршрутами
func NewRouter(
	userController *http2.UserController,
	productController *http2.ProductController,
) *mux.Router {
	// Создаем новый маршрутизатор gorilla/mux
	router := mux.NewRouter()

	// Настройка маршрутов пользователей
	SetupUserRoutes(router, userController)

	// Настройка маршрутов продуктов
	SetupProductRoutes(router, productController)

	// Возвращаем маршрутизатор
	return router
}
