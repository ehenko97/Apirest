package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

// Настроим Redis и PostgreSQL для тестов
func setup() (*redis.Client, *pgx.Conn, func(), error) {
	// Создаем Redis клиента
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Используем имя сервиса Docker
	})

	// Проверяем подключение к Redis
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, nil, nil, fmt.Errorf("ошибка подключения к Redis: %v", err)
	}

	// Создаем подключение к PostgreSQL
	db, err := pgx.Connect(ctx, "postgresql://myuser:mypassword@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		redisClient.Close()
		return nil, nil, nil, fmt.Errorf("ошибка подключения к PostgreSQL: %v", err)
	}

	// Функция для очистки
	tearDown := func() {
		redisClient.Close()
		db.Close(ctx)
	}

	return redisClient, db, tearDown, nil
}

// Тест для проверки подключения к PostgreSQL
func TestPostgresConnection(t *testing.T) {
	// Подготовка
	_, db, tearDown, err := setup()
	if err != nil {
		t.Fatalf("Ошибка при настройке теста: %v", err)
	}
	defer tearDown()

	// Проверка подключения
	err = db.Ping(context.Background())
	assert.NoError(t, err, "Ошибка подключения к базе данных")
}

// Тест для проверки подключения к Redis
func TestRedisConnection(t *testing.T) {
	// Подготовка
	redisClient, _, tearDown, err := setup()
	if err != nil {
		t.Fatalf("Ошибка при настройке теста: %v", err)
	}
	defer tearDown()

	// Проверка подключения
	_, err = redisClient.Ping(context.Background()).Result()
	assert.NoError(t, err, "Ошибка подключения к Redis")
}
