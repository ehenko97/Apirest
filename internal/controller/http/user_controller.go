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

// UserController структура контроллера для обработки запросов пользователей.
type UserController struct {
	userService service.UserService
}

// NewUserController создает новый контроллер пользователей.
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser создает нового пользователя.
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	// Десериализация запроса
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация данных пользователя
	if user.Name == "" || user.Email == "" {
		http.Error(w, "User name and email are required", http.StatusBadRequest)
		return
	}

	createdUser, err := uc.userService.Create(r.Context(), user)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сериализация ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// GetAllUsers возвращает всех пользователей.
func (uc *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uc.userService.FindAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сериализация ответа
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// GetUser возвращает пользователя по ID.
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractUserIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid user ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uc.userService.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Сериализация ответа
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// UpdateUser обновляет пользователя по ID.
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractUserIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid user ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var updatedUser entity.User
	// Десериализация запроса
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	updatedUser.ID = id

	// Валидация данных пользователя
	if updatedUser.Name == "" || updatedUser.Email == "" {
		http.Error(w, "User name and email are required", http.StatusBadRequest)
		return
	}

	user, err := uc.userService.Update(r.Context(), updatedUser)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сериализация ответа
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// DeleteUser удаляет пользователя по ID.
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractUserIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid user ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = uc.userService.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// extractUserIDFromURL извлекает ID пользователя из URL
func extractUserIDFromURL(path string) (int, error) {
	// Убираем только префикс "/users/"
	parts := strings.Split(strings.TrimPrefix(path, "/users/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		return 0, errors.New("missing user ID")
	}

	// Преобразуем ID в число
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, errors.New("invalid user ID format")
	}

	return id, nil
}
