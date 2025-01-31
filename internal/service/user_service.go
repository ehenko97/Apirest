package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ehenko97/apirest/internal/cache"
	"github.com/ehenko97/apirest/internal/entity"
	"log"
	"time"
)

// UsersService реализует логику работы с пользователями
type UsersService struct {
	repo  UserRepository
	cache cache.Cache // Используем интерфейс Cache для работы с кэшем
}

// NewUserService создает новый сервис пользователей
func NewUserService(repo UserRepository, cache cache.Cache) *UsersService {
	return &UsersService{
		repo:  repo,
		cache: cache,
	}
}

// Вспомогательная функция для генерации ключа кэша
func userCacheKey(id int) string {
	return fmt.Sprintf("user:%d", id)
}

// Create добавляет нового пользователя.
func (s *UsersService) Create(ctx context.Context, user entity.User) (entity.User, error) {
	if user.Name == "" {
		return entity.User{}, errors.New("имя пользователя не может быть пустым")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	// Сохраняем пользователя в кэше
	_ = s.cache.Set(ctx, userCacheKey(createdUser.ID), createdUser, 10*time.Minute)

	return createdUser, nil
}

// FindByID возвращает пользователя по ID.
func (s *UsersService) FindByID(ctx context.Context, id int) (entity.User, error) {
	if id <= 0 {
		return entity.User{}, errors.New("некорректный ID пользователя")
	}

	// Попробуем найти пользователя в кэше
	var user entity.User
	err := s.cache.Get(ctx, userCacheKey(id), &user)
	if err == nil {
		return user, nil // Если пользователь найден в кэше, возвращаем его
	}

	// Если в кэше пользователя нет, получаем его из базы данных
	user, err = s.repo.FindByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	// Сохраняем пользователя в кэше
	_ = s.cache.Set(ctx, userCacheKey(id), user, 10*time.Minute)

	return user, nil
}

// Update обновляет информацию о пользователе.
func (s *UsersService) Update(ctx context.Context, user entity.User) (entity.User, error) {
	if user.ID <= 0 || user.Name == "" {
		return entity.User{}, errors.New("некорректные данные для обновления пользователя")
	}

	user.UpdatedAt = time.Now()

	// Обновляем пользователя в базе данных
	err := s.repo.Update(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	// Обновляем пользователя в кэше
	if err := s.cache.Set(ctx, userCacheKey(user.ID), user, 10*time.Minute); err != nil {
		log.Printf("Ошибка обновления в кэше пользователя %d: %v", user.ID, err)
	}

	return user, nil
}

// Delete удаляет пользователя.
func (s *UsersService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("некорректный ID пользователя")
	}

	// Удаляем пользователя из базы данных
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Удаляем пользователя из кэша
	_ = s.cache.Delete(ctx, userCacheKey(id))

	return nil
}

// FindAll возвращает всех пользователей.
func (s *UsersService) FindAll(ctx context.Context) ([]entity.User, error) {
	return s.repo.FindAll(ctx)
}
