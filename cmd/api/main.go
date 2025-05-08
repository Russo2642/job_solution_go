package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"job_solition/internal/config"
	"job_solition/internal/server"
)

// @title JobSolution API
// @version 1.0
// @description API для платформы отзывов сотрудников о компаниях JobSolution
// @BasePath /api
// @description API также доступно по пути /api/v1 для обратной совместимости

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Используйте JWT токен с префиксом "Bearer ". Пример: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	srv := server.NewServer(cfg)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()

	log.Printf("Сервер запущен на порту %s в режиме %s", cfg.Server.Port, cfg.Server.Mode)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Выключение сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	log.Println("Сервер успешно остановлен")
}
