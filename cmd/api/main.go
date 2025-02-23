package main

import (
	"os"
	"pet-store/config"
	"pet-store/internal/infrastructure/logs"
	"pet-store/run"

	"github.com/joho/godotenv"
)

// @title Pet-Store
// @version 1.0
// @description Chtoto-gdeto

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	// Создаем конфигурацию приложения
	conf := config.NewAppConf()
	// Создаем логгер
	logger := logs.NewLogger(conf, os.Stdout)
	if err != nil {
		logger.Fatal("error loading .env file")
	}
	// Инициализируем конфигурацию приложения с логгером
	conf.Init(logger)

	// Создаем инстанс приложения
	app := run.NewApp(conf, logger)

	exitCode := app.
		// Инициализируем приложение
		Bootstrap().
		// Запускаем приложение
		Run()
	os.Exit(exitCode)
}
