package service

import (
	"context"
	"errors"
	_ "github.com/ehenko97/apirest/internal/cache"
	"github.com/ehenko97/apirest/internal/entity"
	"github.com/ehenko97/apirest/internal/service/mocks"
	"github.com/ehenko97/apirest/internal/service/mocks/user"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

// Тесты для UsersService
func TestUsersService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockRepo := user.NewMockUserRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	userService := NewUserService(mockRepo, mockCache)

	ctx := context.Background()

	t.Run("Create User - Success", func(t *testing.T) {
		user := entity.User{ID: 1, Name: "Test User"}
		expectedUser := entity.User{ID: 1, Name: "Test User", CreatedAt: time.Now(), UpdatedAt: time.Now()}

		// Используем gomock.Any() для полей CreatedAt и UpdatedAt
		mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedUser, nil)
		mockCache.EXPECT().Set(ctx, userCacheKey(expectedUser.ID), expectedUser, 10*time.Minute).Return(nil)

		result, err := userService.Create(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("Create User - Error", func(t *testing.T) {
		user := entity.User{ID: 1, Name: "Test User"}
		expectedErr := errors.New("repository error")

		mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(entity.User{}, expectedErr)

		result, err := userService.Create(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, entity.User{}, result)
	})

	t.Run("FindByID User - Success (Cache Hit)", func(t *testing.T) {
		user := entity.User{ID: 1, Name: "Test User"}

		mockCache.EXPECT().Get(ctx, userCacheKey(1), gomock.Any()).SetArg(2, user).Return(nil)

		result, err := userService.FindByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("FindByID User - Success (Cache Miss)", func(t *testing.T) {
		user := entity.User{ID: 1, Name: "Test User"}

		mockCache.EXPECT().Get(ctx, userCacheKey(1), gomock.Any()).Return(errors.New("not found"))
		mockRepo.EXPECT().FindByID(ctx, 1).Return(user, nil)
		mockCache.EXPECT().Set(ctx, userCacheKey(1), user, 10*time.Minute).Return(nil)

		result, err := userService.FindByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("FindByID User - Error", func(t *testing.T) {
		expectedErr := errors.New("repository error")

		mockCache.EXPECT().Get(ctx, userCacheKey(1), gomock.Any()).Return(errors.New("not found"))
		mockRepo.EXPECT().FindByID(ctx, 1).Return(entity.User{}, expectedErr)

		result, err := userService.FindByID(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, entity.User{}, result)
	})

	t.Run("Update User - Success", func(t *testing.T) {
		user := entity.User{ID: 1, Name: "Updated User", UpdatedAt: time.Now()}

		mockRepo.EXPECT().Update(ctx, user).Return(nil)
		mockCache.EXPECT().Set(ctx, userCacheKey(user.ID), user, 10*time.Minute).Return(nil)

		result, err := userService.Update(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("Update User - Error", func(t *testing.T) {
		user := entity.User{ID: 1, Name: "Updated User", UpdatedAt: time.Now()}
		expectedErr := errors.New("repository error")

		mockRepo.EXPECT().Update(ctx, user).Return(expectedErr)

		result, err := userService.Update(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, entity.User{}, result)
	})

	t.Run("Delete User - Success", func(t *testing.T) {
		mockRepo.EXPECT().Delete(ctx, 1).Return(nil)
		mockCache.EXPECT().Delete(ctx, userCacheKey(1)).Return(nil)

		err := userService.Delete(ctx, 1)

		assert.NoError(t, err)
	})

	t.Run("Delete User - Error", func(t *testing.T) {
		expectedErr := errors.New("repository error")

		mockRepo.EXPECT().Delete(ctx, 1).Return(expectedErr)

		err := userService.Delete(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("FindAll Users - Success", func(t *testing.T) {
		users := []entity.User{
			{ID: 1, Name: "User 1"},
			{ID: 2, Name: "User 2"},
		}

		mockRepo.EXPECT().FindAll(ctx).Return(users, nil)

		result, err := userService.FindAll(ctx)

		assert.NoError(t, err)
		assert.Equal(t, users, result)
	})

	t.Run("FindAll Users - Error", func(t *testing.T) {
		expectedErr := errors.New("repository error")

		mockRepo.EXPECT().FindAll(ctx).Return(nil, expectedErr)

		result, err := userService.FindAll(ctx)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
	})
}

// Тесты для ProductServiceDB
func TestProductServiceDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockRepo := user.NewMockProductRepository(ctrl) // Убедитесь, что это правильный пакет
	mockCache := mocks.NewMockCache(ctrl)
	productService := NewProductService(mockRepo, mockCache)

	ctx := context.Background()

	t.Run("Create Product - Success", func(t *testing.T) {
		product := entity.Product{ID: 1, Name: "Test Product", Price: 100}
		expectedProduct := entity.Product{ID: 1, Name: "Test Product", Price: 100, CreatedAt: time.Now(), UpdatedAt: time.Now()}

		// Используем gomock.Any() для полей CreatedAt и UpdatedAt
		mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedProduct, nil)

		result, err := productService.Create(ctx, product)

		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, result)
	})

	t.Run("Create Product - Error", func(t *testing.T) {
		product := entity.Product{ID: 1, Name: "Test Product", Price: 100}
		expectedErr := errors.New("repository error")

		// Используем gomock.Any() для полей CreatedAt и UpdatedAt
		mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(entity.Product{}, expectedErr)

		result, err := productService.Create(ctx, product)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, entity.Product{}, result)
	})

	t.Run("FindByID Product - Success (Cache Hit)", func(t *testing.T) {
		product := entity.Product{ID: 1, Name: "Test Product", Price: 100}

		mockCache.EXPECT().Get(ctx, cacheKey(1), gomock.Any()).SetArg(2, product).Return(nil)

		result, err := productService.FindByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, product, result)
	})

	t.Run("FindByID Product - Success (Cache Miss)", func(t *testing.T) {
		product := entity.Product{ID: 1, Name: "Test Product", Price: 100}

		mockCache.EXPECT().Get(ctx, cacheKey(1), gomock.Any()).Return(errors.New("not found"))
		mockRepo.EXPECT().FindByID(ctx, 1).Return(product, nil)
		mockCache.EXPECT().Set(ctx, cacheKey(1), product, 10*time.Minute).Return(nil)

		result, err := productService.FindByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, product, result)
	})

	t.Run("FindByID Product - Error", func(t *testing.T) {
		expectedErr := errors.New("repository error")

		mockCache.EXPECT().Get(ctx, cacheKey(1), gomock.Any()).Return(errors.New("not found"))
		mockRepo.EXPECT().FindByID(ctx, 1).Return(entity.Product{}, expectedErr)

		result, err := productService.FindByID(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, entity.Product{}, result)
	})

	t.Run("Update Product - Success", func(t *testing.T) {
		product := entity.Product{ID: 1, Name: "Updated Product", Price: 200, UpdatedAt: time.Now()}

		mockRepo.EXPECT().Update(ctx, product).Return(product, nil)
		mockCache.EXPECT().Set(ctx, cacheKey(product.ID), product, 10*time.Minute).Return(nil)

		result, err := productService.Update(ctx, product)

		assert.NoError(t, err)
		assert.Equal(t, product, result)
	})

	t.Run("Update Product - Error", func(t *testing.T) {
		product := entity.Product{ID: 1, Name: "Updated Product", Price: 200, UpdatedAt: time.Now()}
		expectedErr := errors.New("repository error")

		mockRepo.EXPECT().Update(ctx, product).Return(entity.Product{}, expectedErr)

		result, err := productService.Update(ctx, product)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, entity.Product{}, result)
	})

	t.Run("Delete Product - Success", func(t *testing.T) {
		mockRepo.EXPECT().Delete(ctx, 1).Return(nil)
		mockCache.EXPECT().Delete(ctx, cacheKey(1)).Return(nil)

		err := productService.Delete(ctx, 1)

		assert.NoError(t, err)
	})

	t.Run("Delete Product - Error", func(t *testing.T) {
		expectedErr := errors.New("repository error")

		mockRepo.EXPECT().Delete(ctx, 1).Return(expectedErr)

		err := productService.Delete(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("FindAll Products - Success", func(t *testing.T) {
		products := []entity.Product{
			{ID: 1, Name: "Product 1", Price: 100},
			{ID: 2, Name: "Product 2", Price: 200},
		}

		mockRepo.EXPECT().FindAll(ctx).Return(products, nil)

		result, err := productService.FindAll(ctx)

		assert.NoError(t, err)
		assert.Equal(t, products, result)
	})

	t.Run("FindAll Products - Error", func(t *testing.T) {
		expectedErr := errors.New("repository error")

		mockRepo.EXPECT().FindAll(ctx).Return(nil, expectedErr)

		result, err := productService.FindAll(ctx)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
	})

	t.Run("FindByUserID Products - Success", func(t *testing.T) {
		products := []entity.Product{
			{ID: 1, Name: "Product 1", Price: 100, UserID: 1},
			{ID: 2, Name: "Product 2", Price: 200, UserID: 1},
		}

		mockRepo.EXPECT().FindByUserID(ctx, 1).Return(products, nil)

		result, err := productService.FindByUserID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, products, result)
	})

	t.Run("FindByUserID Products - Error", func(t *testing.T) {
		expectedErr := errors.New("repository error")

		mockRepo.EXPECT().FindByUserID(ctx, 1).Return(nil, expectedErr)

		result, err := productService.FindByUserID(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
	})
}
