package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ehenko97/apirest/internal/cache"
	"github.com/ehenko97/apirest/internal/entity"
	"time"
)

// ProductServiceDB реализует логику работы с продуктами
type ProductServiceDB struct {
	repo  ProductRepository
	cache cache.Cache // Используем интерфейс Cache из пакета cache
}

// NewProductService создает новый сервис продуктов
func NewProductService(repo ProductRepository, cache cache.Cache) *ProductServiceDB {
	return &ProductServiceDB{
		repo:  repo,
		cache: cache,
	}
}

// Вспомогательная функция для генерации ключа кеша
func cacheKey(id int) string {
	return fmt.Sprintf("product:%d", id)
}

// Create добавляет новый продукт.
func (s *ProductServiceDB) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	if product.Name == "" || product.Price <= 0 {
		return entity.Product{}, errors.New("необходимо указать название и положительную цену")
	}

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	return s.repo.Create(ctx, product)
}

// FindByID возвращает продукт по ID.
func (s *ProductServiceDB) FindByID(ctx context.Context, id int) (entity.Product, error) {
	if id <= 0 {
		return entity.Product{}, errors.New("некорректный ID продукта")
	}

	// Попробуем найти продукт в кэше
	var product entity.Product
	err := s.cache.Get(ctx, cacheKey(id), &product)
	if err == nil {
		return product, nil // Если продукт найден в кэше, возвращаем его
	}

	// Если в кэше продукта нет, получаем его из базы данных
	product, err = s.repo.FindByID(ctx, id)
	if err != nil {
		return entity.Product{}, err
	}

	// Сохраняем продукт в кэше
	_ = s.cache.Set(ctx, cacheKey(id), product, 10*time.Minute) // Устанавливаем TTL 10 минут

	return product, nil
}

// Update обновляет информацию о продукте.
func (s *ProductServiceDB) Update(ctx context.Context, product entity.Product) (entity.Product, error) {
	if product.ID <= 0 || product.Name == "" || product.Price <= 0 {
		return entity.Product{}, errors.New("некорректные данные для обновления продукта")
	}

	product.UpdatedAt = time.Now()

	// Обновляем продукт в базе данных
	updatedProduct, err := s.repo.Update(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}

	// Обновляем кэш
	_ = s.cache.Set(ctx, cacheKey(product.ID), updatedProduct, 10*time.Minute)

	return updatedProduct, nil
}

// Delete удаляет продукт по ID.
func (s *ProductServiceDB) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("некорректный ID продукта")
	}

	// Удаляем продукт из базы данных
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Удаляем продукт из кэша
	_ = s.cache.Delete(ctx, cacheKey(id))

	return nil
}

// FindAll возвращает список всех продуктов.
func (s *ProductServiceDB) FindAll(ctx context.Context) ([]entity.Product, error) {
	return s.repo.FindAll(ctx)
}

// FindByUserID возвращает продукты, связанные с определенным userID.
func (s *ProductServiceDB) FindByUserID(ctx context.Context, userID int) ([]entity.Product, error) {
	if userID <= 0 {
		return nil, errors.New("некорректный ID пользователя")
	}
	return s.repo.FindByUserID(ctx, userID)
}
