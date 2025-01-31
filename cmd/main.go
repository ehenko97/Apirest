package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ehenko97/apirest"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ehenko97/apirest/api/routes"
	"github.com/ehenko97/apirest/internal/cache"
	mygrpc "github.com/ehenko97/apirest/internal/controller/grpc"
	httpController "github.com/ehenko97/apirest/internal/controller/http"
	"github.com/ehenko97/apirest/internal/repository"
	"github.com/ehenko97/apirest/internal/service"
	pb "github.com/ehenko97/apirest/pkg/proto"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

// Статические ошибки
var (
	ErrMissingEnvVar    = errors.New("необходимо указать переменную окружения")
	ErrUnknownCacheType = errors.New("неизвестный тип кэша")
)

func main() {
	// Настройка логгера
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("Ошибка", "err", err)
		os.Exit(1) // Завершение программы с кодом ошибки
	}
}

func run() error {
	// Загрузка переменных окружения из .env
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("ошибка загрузки .env файла: %w", err)
	}

	// Проверка обязательных переменных окружения
	requiredEnvVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE", "CACHE_TYPE"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("%w: %s", ErrMissingEnvVar, envVar)
		}
	}

	// Чтение конфигурации из переменных окружения
	dbConfig := repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Подключение к базе данных
	db, err := repository.NewPostgresDB(dbConfig)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}
	defer db.Close()

	slog.Info("Подключение к базе данных установлено")

	// Создание репозиториев
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Чтение конфигурации кэша из переменных окружения
	cacheType := os.Getenv("CACHE_TYPE")

	var cacheService cache.Cache

	switch cacheType {
	case "redis":
		cacheService = cache.NewRedisCache(os.Getenv("REDIS_ADDR"))
	case "inmemory":
		cacheService = cache.NewInMemoryCache()
	default:
		return fmt.Errorf("%w: %s", ErrUnknownCacheType, cacheType)
	}

	// Создание сервисов
	userService := service.NewUserService(userRepo, cacheService)
	productService := service.NewProductService(productRepo, cacheService)

	// Инициализация HTTP контроллеров
	userController := httpController.NewUserController(userService)
	productController := httpController.NewProductController(productService, userService)

	// Создание маршрутизатора
	router := routes.NewRouter(userController, productController)

	// Создание сервера
	server := projectapirest.NewServer()

	// Запуск HTTP сервера в горутине
	go func() {
		slog.Info("Запуск HTTP сервера на порту 8080...")
		if err := server.RunHTTP("8080", router); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Ошибка запуска HTTP сервера", "err", err)
		}
	}()

	// Запуск gRPC сервера в горутине
	go func() {
		slog.Info("Запуск gRPC сервера на порту 50051...")
		if err := server.RunGRPC("50051", func(srv *grpc.Server) {
			// Регистрируем UserService
			userServiceServer := mygrpc.NewUserController(userService)
			pb.RegisterUserServiceServer(srv, userServiceServer)

			// Регистрируем ProductService
			productServiceServer := mygrpc.NewProductController(productService)
			pb.RegisterProductServiceServer(srv, productServiceServer)
		}); err != nil {
			slog.Error("Ошибка запуска gRPC сервера", "err", err)
		}
	}()

	// Ожидание сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Получен сигнал для завершения работы...")

	// Корректное завершение работы серверов
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("ошибка завершения работы серверов: %w", err)
	}

	slog.Info("Серверы успешно остановлены")
	return nil
}
