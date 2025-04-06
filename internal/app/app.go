package app

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/manifoldco/promptui"

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
}

// NewApp - создание нового экземпляра приложения
func NewApp(cfg *config.Config, log logger.Logger) *App {
	projectService := service.NewProjectService(log, cfg.DataFile)
	return &App{
		ProjectService:  projectService,
		TrackingService: service.NewTrackingService(projectService, log, cfg),
		SystrayHandler:  NewSystrayHandler(),
		Logger:          log,
		Config:          cfg,
	}
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
			Label: "Выберите команду",
			Items: []string{"Создать проект", "Начать отслеживание", "Остановить отслеживание", "Сводка", "Выход"},
		}
		_, cmd, _ := prompt.Run()

		switch cmd {
		case "Создать проект":
			a.createProject()
		case "Начать отслеживание":
			a.startTracking()
		case "Остановить отслеживание":
			a.stopTracking()
		case "Сводка":
			a.showSummary()
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
				a.startTracking()
			case <-mStop.ClickedCh:
				a.stopTracking()
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

// createProject - создание нового проекта
func (a *App) createProject() {
	prompt := promptui.Prompt{
		Label: "Название проекта",
	}
	name, _ := prompt.Run()
	a.Logger.Infof("Попытка создания проекта: %s", name)

	err := a.ProjectService.CreateProject(a.Projects, name)
	if err != nil {
		a.Logger.Errorf("Ошибка создания проекта: %v", err)
		fmt.Println(err)
		return
	}
	a.Logger.Infof("Проект создан: %s", name)
	fmt.Println("Проект создан:", name)
}

// chooseProject - выбор проекта из списка
func (a *App) chooseProject() string {
	// Создаем список проектов с информацией о статусе
	var activeProjects []string
	var inactiveProjects []string
	
	for name, project := range a.Projects {
		if project.StartTime != nil {
			activeProjects = append(activeProjects, "▶ "+name)
		} else {
			inactiveProjects = append(inactiveProjects, "  "+name)
		}
	}
	
	// Сортируем проекты по алфавиту
	sort.Strings(activeProjects)
	sort.Strings(inactiveProjects)
	
	// Объединяем списки: сначала "Назад", затем активные, затем неактивные
	options := append([]string{"← Назад"}, activeProjects...)
	options = append(options, inactiveProjects...)
	
	if len(options) == 1 { // Только опция "Назад"
		fmt.Println("Нет доступных проектов")
		return ""
	}
	
	prompt := promptui.Select{
		Label: "Выберите проект",
		Items: options,
	}
	idx, result, err := prompt.Run()
	
	if err != nil {
		a.Logger.Errorf("Ошибка при выборе проекта: %v", err)
		return ""
	}
	
	// Если выбрана опция "Назад" (индекс 0), возвращаем пустую строку
	if idx == 0 {
		return ""
	}
	
	// Удаляем префикс статуса из имени проекта
	return strings.TrimPrefix(strings.TrimPrefix(result, "▶ "), "  ")
}

// chooseActiveProject - выбор активного проекта из списка
func (a *App) chooseActiveProject() string {
	// Собираем имена активных проектов
	var activeProjects []string
	for name, project := range a.Projects {
		if project.StartTime != nil {
			activeProjects = append(activeProjects, name)
		}
	}
	
	if len(activeProjects) == 0 {
		fmt.Println("Нет активных проектов")
		return ""
	}
	
	// Добавляем опцию "Назад" в начало списка
	options := append([]string{"← Назад"}, activeProjects...)
	
	prompt := promptui.Select{
		Label: "Выберите активный проект",
		Items: options,
	}
	idx, result, err := prompt.Run()
	
	if err != nil {
		a.Logger.Errorf("Ошибка при выборе проекта: %v", err)
		return ""
	}
	
	// Если выбрана опция "Назад" (индекс 0), возвращаем пустую строку
	if idx == 0 {
		return ""
	}
	
	return result
}

// startTracking - начало отслеживания времени
func (a *App) startTracking() {
	projectName := a.chooseProject()
	if projectName == "" {
		a.Logger.Warn("Отмена начала отслеживания: проект не выбран")
		return
	}

	a.Logger.Infof("Попытка начать отслеживание для проекта: %s", projectName)
	err := a.TrackingService.StartTracking(a.Projects, projectName)
	if err != nil {
		a.Logger.Errorf("Ошибка начала отслеживания: %v", err)
		fmt.Println(err)
		return
	}

	a.Logger.Infof("Начато отслеживание для проекта: %s", projectName)
	fmt.Println("Начато отслеживание для проекта:", projectName)

	project := a.Projects[projectName]
	a.SystrayHandler.SetTracking(projectName, project.StartTime)

	go func() {
		time.Sleep(time.Duration(a.Config.NotificationTime) * time.Second)
		a.Logger.Infof("Отправка уведомления о перерыве для проекта: %s", projectName)
		a.TrackingService.Notify(projectName)
	}()
}

// stopTracking - остановка отслеживания времени
func (a *App) stopTracking() {
	projectName := a.chooseActiveProject()
	if projectName == "" {
		a.Logger.Warn("Отмена остановки отслеживания: проект не выбран")
		return
	}

	prompt := promptui.Prompt{
		Label: "Что сделано",
	}
	description, _ := prompt.Run()

	a.Logger.Infof("Попытка остановить отслеживание для проекта: %s", projectName)
	elapsed, err := a.TrackingService.StopTracking(a.Projects, projectName, description)
	if err != nil {
		a.Logger.Errorf("Ошибка остановки отслеживания: %v", err)
		fmt.Println(err)
		return
	}

	a.Logger.Infof("Отслеживание остановлено для проекта %s. Время: %v", projectName, elapsed)
	fmt.Printf("Отслеживание остановлено для проекта %s. Время: %v\n", projectName, elapsed)
	a.SystrayHandler.SetTracking("", nil)
}

// showSummary - вывод сводки по проектам
func (a *App) showSummary() {
	if len(a.Projects) == 0 {
		fmt.Println("Нет данных для отображения.")
		return
	}

	for name, project := range a.Projects {
		fmt.Printf("\nПроект \"%s\":\n", name)
		var total int
		for _, entry := range project.Entries {
			total += entry.TimeSpent
			fmt.Printf("  %s - %v сек: %s\n", entry.Date, entry.TimeSpent, entry.Description)
		}
		fmt.Printf("  Общее время: %v сек\n", total)
	}
}
