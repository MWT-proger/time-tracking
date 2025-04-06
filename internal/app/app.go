package app

import (
	"fmt"

	"github.com/MWT-proger/time-tracking/internal/app/handlers"
	"github.com/MWT-proger/time-tracking/internal/app/systray"
	"github.com/MWT-proger/time-tracking/internal/domain"
	"github.com/MWT-proger/time-tracking/internal/service"
	"github.com/MWT-proger/time-tracking/pkg/config"
	"github.com/MWT-proger/time-tracking/pkg/logger"
)

type SystrayHandler interface {
	Run()
	Quit()
}

// App - основной класс приложения
type App struct {
	ProjectService  *service.ProjectService
	TrackingService *service.TrackingService
	SystrayHandler  SystrayHandler
	Projects        map[string]*domain.Project
	Logger          logger.Logger
	Config          *config.Config
	Handlers        *handlers.Handlers
}

// NewApp - создание нового экземпляра приложения
func NewApp(cfg *config.Config, log logger.Logger) *App {
	projectService := service.NewProjectService(log, cfg.DataFile)
	trackingService := service.NewTrackingService(projectService, log, cfg)
	systrayHandler := systray.NewSystrayHandler(log)

	app := &App{
		ProjectService:  projectService,
		TrackingService: trackingService,
		SystrayHandler:  systrayHandler,
		Logger:          log,
		Config:          cfg,
	}

	// Инициализируем обработчики
	app.Handlers = handlers.NewHandlers(app.ProjectService, app.TrackingService, systrayHandler, app.Logger, app.Config)

	return app
}

// Initialize - инициализация приложения
func (a *App) Initialize() error {
	a.Logger.Info("Инициализация приложения")
	var err error
	a.Projects, err = a.ProjectService.LoadData()
	if err != nil {
		a.Logger.Errorf("Ошибка загрузки данных: %v", err)
		return err
	}

	// Передаем загруженные проекты в обработчики
	a.Handlers.SetProjects(a.Projects)

	a.Logger.Infof("Загружено проектов: %d", len(a.Projects))
	return nil
}

// Run - запуск приложения
func (a *App) Run() {
	a.Logger.Info("Запуск приложения")

	// Запускаем системный трей в отдельной горутине
	go func() {
		defer func() {
			if r := recover(); r != nil {
				a.Logger.Errorf("Ошибка при запуске системного трея: %v", r)
				fmt.Println("Не удалось запустить системный трей. Приложение будет работать только в режиме командной строки.")
			}
		}()

		a.SystrayHandler.Run()
	}()
	a.Handlers.GeneralMenu()
	a.SystrayHandler.Quit()
}
