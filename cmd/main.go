package main

import (
	"awesomeProject2/cmd/config"
	"awesomeProject2/cmd/db"
	"awesomeProject2/cmd/handler"
	"awesomeProject2/cmd/service"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Не удалось создать логгер: %v", err)
	}
	defer logger.Sync()

	env := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.GetOrDefault("DB_HOST", "db"),
		config.GetOrDefault("DB_PORT", "5432"),
		config.GetOrDefault("DB_USER", "admin"),
		config.GetOrDefault("DB_PASSWORD", "3228"),
		config.GetOrDefault("DB_NAME", "mydb"),
		config.GetOrDefault("SSL_MODE", "disable"),
	)

	db, err := sqlx.Open("postgres", env)
	if err != nil {
		logger.Fatal("Не удалось подключиться к БД", zap.Any("env", env))
	}
	logger.Info("Приложение успешно стартовало")
	boardStore := storage.NewBoardStorage(db)
	listStore := storage.NewListStorage(db)
	cardStore := storage.NewCardStorage(db)
	boardService := service.NewBoardService(boardStore, logger)
	listService := service.NewListService(listStore, logger)
	cardService := service.NewCardService(cardStore, logger)
	boardHandler := handler.NewBoardHandler(boardService, logger)
	listHandler := handler.NewListHandler(listService, logger)
	cardHandler := handler.NewCardHandler(cardService, logger)
	http.HandleFunc("/boards", boardHandler.HandleBoards)
	http.HandleFunc("/lists", listHandler.HandleLists)
	http.HandleFunc("/cards", cardHandler.HandleCards)
	http.ListenAndServe(":8080", nil)
}
