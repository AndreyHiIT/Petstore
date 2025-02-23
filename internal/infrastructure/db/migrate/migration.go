package migrations

import (
	"fmt"
	"pet-store/config"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres" // Импорт драйвера PostgreSQL
	_ "github.com/golang-migrate/migrate/source/file"       // Импорт драйвера для файлов миграций
	"go.uber.org/zap"
)

func MigrationInit(config config.AppConf, logger *zap.Logger) {
	db := config.DB
	// Формирование строки подключения с использованием переменных окружения
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db.User, db.Password, db.Host, db.Port, db.Name)
	// Создаем объект миграции
	m, err := migrate.New(
		"file://internal/infrastructure/db/migrate", // Путь к миграциям
		dbURL, // URL для подключения к базе
	)
	if err != nil {
		logger.Fatal("Error", zap.Error(err))
	}
	// Получаем текущую версию миграций
	version, _, err := m.Version()
	if err != nil && err.Error() != "no migration" {
		logger.Fatal("Error", zap.Error(err))
	}
	if version == 0 {
		// Если миграции ещё не применены, выполняем их
		logger.Info("No migrations applied yet. Applying migrations...")
		err = m.Up()
		if err != nil && err.Error() != "no change" {
			logger.Fatal("Error", zap.Error(err))
		}
		logger.Info("Migrations applied successfully.")
	} else {
		logger.Info("Migrations already applied", zap.Uint("Current version", version))
	}
}
