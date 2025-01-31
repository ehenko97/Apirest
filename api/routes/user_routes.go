package routes

import (
	http2 "github.com/ehenko97/apirest/internal/controller/http"
	"github.com/gorilla/mux"
)

// SetupUserRoutes регистрирует маршруты для работы с пользователями
func SetupUserRoutes(r *mux.Router, userController *http2.UserController) {
	// GET получение всех пользователей
	r.HandleFunc("/users", userController.GetAllUsers).Methods("GET")

	// POST создание нового пользователя
	r.HandleFunc("/users", userController.CreateUser).Methods("POST")

	// GET  получение пользователя по ID
	r.HandleFunc("/users/{id}", userController.GetUser).Methods("GET")

	// PUT обновление данных пользователя по ID
	r.HandleFunc("/users/{id}", userController.UpdateUser).Methods("PUT")

	// DELETE  удаление пользователя по ID
	r.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE")
}
