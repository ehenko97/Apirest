package projectapirest

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

// Server представляет собой сервер, который может запускать HTTP и gRPC серверы.
type Server struct {
	httpServer *http.Server
	grpcServer *grpc.Server
}

// NewServer создаёт новый экземпляр Server.
func NewServer() *Server {
	return &Server{}
}

// RunHTTP запускает HTTP сервер на указанном порту.
// Принимает порт и HTTP-обработчик.
// Возвращает ошибку, если сервер не может быть запущен.
func (s *Server) RunHTTP(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    30 * time.Second, // Таймаут для неактивных соединений
	}

	log.Printf("Запуск HTTP сервера на порту %s...\n", port)
	return s.httpServer.ListenAndServe()
}

// RunGRPC запускает gRPC сервер на указанном порту.
// Принимает порт и функцию для регистрации gRPC-сервисов.
// Возвращает ошибку, если сервер не может быть запущен.
func (s *Server) RunGRPC(port string, register func(*grpc.Server)) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("не удалось прослушать порт %s: %v", port, err)
	}

	s.grpcServer = grpc.NewServer()
	register(s.grpcServer) // Регистрируем gRPC-сервисы

	log.Printf("Запуск gRPC сервера на порту %s...\n", port)
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("ошибка запуска gRPC сервера: %v", err)
	}
	return nil
}

// Shutdown корректно завершает работу HTTP и gRPC серверов.
// Принимает контекст для управления таймаутом.
// Возвращает ошибку, если завершение работы не удалось.
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Завершение работы серверов...")

	var err error

	// Graceful shutdown для HTTP сервера
	if s.httpServer != nil {
		if shutdownErr := s.httpServer.Shutdown(ctx); shutdownErr != nil {
			log.Printf("Ошибка завершения HTTP сервера: %v\n", shutdownErr)
			err = shutdownErr
		}
	}

	// Graceful shutdown для gRPC сервера
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	if err != nil {
		return fmt.Errorf("ошибка завершения работы серверов: %v", err)
	}

	log.Println("Серверы успешно остановлены")
	return nil
}
