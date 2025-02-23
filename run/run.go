package run

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"pet-store/config"
	"pet-store/internal/db"
	"pet-store/internal/infrastructure/component"
	migrations "pet-store/internal/infrastructure/db/migrate"
	"pet-store/internal/infrastructure/errors"
	"pet-store/internal/infrastructure/responder"
	"pet-store/internal/infrastructure/router"
	"pet-store/internal/infrastructure/server"
	"pet-store/internal/infrastructure/tools/cryptography"
	"pet-store/internal/modules"
	"pet-store/internal/storages"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/sync/errgroup"

	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

// App - структура приложения
type App struct {
	conf     config.AppConf
	logger   *zap.Logger
	srv      server.Server
	Sig      chan os.Signal
	Storages *storages.Storages
	Servises *modules.Services
}

// Runner - интерфейс запуска приложения
type Runner interface {
	Run() int
}

// NewApp - конструктор приложения
func NewApp(conf config.AppConf, logger *zap.Logger) *App {
	return &App{conf: conf, logger: logger, Sig: make(chan os.Signal, 1)}
}

func (a *App) Run() int {
	// создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	errGroup, ctx := errgroup.WithContext(ctx)

	// запускаем горутину для graceful shutdown
	// при получении сигнала SIGINT
	// вызываем cancel для контекста
	errGroup.Go(func() error {
		sigInt := <-a.Sig
		a.logger.Info("signal interrupt recieved", zap.Stringer("os_signal", sigInt))
		cancel()
		return nil
	})

	errGroup.Go(func() error {
		err := a.srv.Serve(ctx)
		if err != nil && err != http.ErrServerClosed {
			a.logger.Error("app: server error", zap.Error(err))
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return errors.GeneralError
	}

	return errors.NoError
}

// Bootstrap - инициализация приложения
func (a *App) Bootstrap(options ...interface{}) Runner {
	// инициализация менеджера токенов
	tokenManager := cryptography.NewTokenJWT(a.conf.Token)
	// инициализация декодера
	decoder := godecoder.NewDecoder(jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		DisallowUnknownFields:  true,
	})
	// инициализация менеджера ответов сервера
	responseManager := responder.NewResponder(decoder, a.logger)
	// инициализация генератора uuid
	uuID := cryptography.NewUUIDGenerator()
	// инициализация хешера
	hash := cryptography.NewHash(uuID)
	// инициализация компонентов
	components := component.NewComponents(a.conf, tokenManager, responseManager, decoder, hash, a.logger)
	// инициализация базы данных sql и его адаптера
	_, sqlAdapter, err := db.NewSqlDB(a.conf.DB, a.logger)
	if err != nil {
		a.logger.Fatal("error init db", zap.Error(err))
	}

	// инициализация хранилищ
	newStorages := storages.NewStorages(sqlAdapter)
	a.Storages = newStorages
	// инициализация сервисов
	services := modules.NewServices(newStorages, components)
	a.Servises = services
	controllers := modules.NewControllers(services, components)
	// инициализация роутера
	r := router.NewRouter(controllers, components)
	//Применение миграций если они ещё не применены
	migrations.MigrationInit(a.conf, a.logger)
	// конфигурация сервера
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.conf.Server.Port),
		Handler: r,
	}
	// инициализация сервера
	a.srv = server.NewHttpServer(a.conf.Server, srv, a.logger)
	// возвращаем приложение
	return a
}
