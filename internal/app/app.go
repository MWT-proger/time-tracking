package app

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/manifoldco/promptui"

	"github.com/MWT-proger/time-tracking/internal/app/handlers"
	"github.com/MWT-proger/time-tracking/internal/domain"
	"github.com/MWT-proger/time-tracking/internal/service"
	"github.com/MWT-proger/time-tracking/pkg/config"
	"github.com/MWT-proger/time-tracking/pkg/logger"
)

// App - основной класс приложения
type App struct {
	ProjectService  *service.ProjectService
	TrackingService *service.TrackingService
	SystrayHandler  *SystrayHandler
	Projects        map[string]*domain.Project
	Logger          logger.Logger
	Config          *config.Config
	Handlers        *handlers.Handlers
}

// NewApp - создание нового экземпляра приложения
func NewApp(cfg *config.Config, log logger.Logger) *App {
	projectService := service.NewProjectService(log, cfg.DataFile)
	trackingService := service.NewTrackingService(projectService, log, cfg)
	systrayHandler := NewSystrayHandler()

	app := &App{
		ProjectService:  projectService,
		TrackingService: trackingService,
		SystrayHandler:  systrayHandler,
		Logger:          log,
		Config:          cfg,
	}

	// Инициализируем обработчики после создания App
	app.Handlers = handlers.NewHandlers(app.ProjectService, app.TrackingService, app.SystrayHandler, app.Logger, app.Config)

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

		systray.Run(a.onReady, a.onExit)
	}()

	for {
		prompt := promptui.Select{
			Label: "Главное меню",
			Items: []string{
				"Выбрать проект",
				"Создать проект",
				"Сводка по всем проектам",
				"Выход",
			},
		}
		_, cmd, _ := prompt.Run()

		switch cmd {
		case "Выбрать проект":
			a.Handlers.SelectAndManageProject()
		case "Создать проект":
			a.Handlers.CreateProject()
		case "Сводка по всем проектам":
			a.Handlers.ShowSummary()
		case "Выход":
			a.ProjectService.SaveData(a.Projects)
			systray.Quit()
			return
		}
	}
}

// onReady - обработчик готовности системного трея
func (a *App) onReady() {
	// Загружаем иконку из встроенных ресурсов
	iconData, err := GetIcon()
	if err != nil {
		// Если не удалось загрузить иконку, используем встроенную иконку
		a.Logger.Warnf("Не удалось загрузить иконку: %v, используется встроенная иконка", err)
	}

	systray.SetIcon(iconData)
	systray.SetTitle("Таймер")
	systray.SetTooltip("Учет времени")

	// Создаем пункты меню
	mStart := systray.AddMenuItem("Начать отслеживание", "Начать отслеживание времени")
	mStop := systray.AddMenuItem("Остановить отслеживание", "Остановить отслеживание времени")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Выход", "Выход из приложения")

	// Обработка событий меню
	go func() {
		for {
			select {
			case <-mStart.ClickedCh:
				a.Handlers.StartTracking()
			case <-mStop.ClickedCh:
				a.Handlers.StopTracking()
			case <-mQuit.ClickedCh:
				a.Logger.Info("Выход из приложения через системный трей")
				systray.Quit()
				return
			}
		}
	}()
}

// onExit - обработчик выхода из системного трея
func (a *App) onExit() {
	// Логика при выходе
}
