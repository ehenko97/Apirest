package service

import (
	"context"
	"github.com/ehenko97/apirest/internal/entity"
)

// ProductService описывает методы управления продуктами.
type ProductService interface {
	Create(ctx context.Context, product entity.Product) (entity.Product, error)
	FindByID(ctx context.Context, id int) (entity.Product, error)
	Update(ctx context.Context, product entity.Product) (entity.Product, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]entity.Product, error)
	FindByUserID(ctx context.Context, userID int) ([]entity.Product, error)
}

// UserService описывает методы управления пользователями.
type UserService interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	FindByID(ctx context.Context, id int) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]entity.User, error)
}

// ProductRepository описывает методы для работы с продуктами.
type ProductRepository interface {
	Create(ctx context.Context, productID entity.Product) (entity.Product, error)
	FindByID(ctx context.Context, id int) (entity.Product, error)
	Update(ctx context.Context, product entity.Product) (entity.Product, error)
	Delete(ctx context.Context, productID int) error
	FindAll(ctx context.Context) ([]entity.Product, error)
	FindByUserID(ctx context.Context, userID int) ([]entity.Product, error)
}

// UserRepository описывает методы работы с пользователями.
type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	FindByID(ctx context.Context, id int) (entity.User, error)
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context) ([]entity.User, error)
}

type Repository struct {
	UserRepository
	ProductRepository
}
