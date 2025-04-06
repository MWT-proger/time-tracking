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
			Items: []string{
				"Создать проект",
				"Начать отслеживание",
				"Остановить отслеживание",
				"Управление спринтами",
				"Сводка",
				"Выход",
			},
		}
		_, cmd, _ := prompt.Run()

		switch cmd {
		case "Создать проект":
			a.createProject()
		case "Начать отслеживание":
			a.startTracking()
		case "Остановить отслеживание":
			a.stopTracking()
		case "Управление спринтами":
			a.manageSprints()
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
		Validate: func(input string) error {
			// Проверка на пустое имя
			if input == "" {
				return fmt.Errorf("имя проекта не может быть пустым")
			}

			// Проверка на уникальность имени
			if _, exists := a.Projects[input]; exists {
				return fmt.Errorf("проект с именем '%s' уже существует", input)
			}

			return nil
		},
	}

	name, err := prompt.Run()
	if err != nil {
		a.Logger.Warnf("Отмена создания проекта: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	a.Logger.Infof("Попытка создания проекта: %s", name)

	err = a.ProjectService.CreateProject(a.Projects, name)
	if err != nil {
		a.Logger.Errorf("Ошибка создания проекта: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	a.Logger.Infof("Проект создан: %s", name)
	fmt.Printf("Проект '%s' успешно создан\n", name)
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

// manageSprints - управление спринтами проектов
func (a *App) manageSprints() {
	for {
		prompt := promptui.Select{
			Label: "Управление спринтами",
			Items: []string{
				"Создать спринт",
				"Выбрать активный спринт",
				"Просмотреть спринты проекта",
				"Назад",
			},
		}
		_, cmd, _ := prompt.Run()

		switch cmd {
		case "Создать спринт":
			a.createSprint()
		case "Выбрать активный спринт":
			a.setActiveSprint()
		case "Просмотреть спринты проекта":
			a.viewProjectSprints()
		case "Назад":
			return
		}
	}
}

// createSprint - создание нового спринта проекта
func (a *App) createSprint() {
	// Выбор проекта
	projectName := a.chooseProject()
	if projectName == "" {
		return
	}

	// Ввод имени спринта
	namePrompt := promptui.Prompt{
		Label: "Название спринта",
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("имя спринта не может быть пустым")
			}

			// Проверка на уникальность имени спринта
			project := a.Projects[projectName]
			if project.Sprints != nil {
				for _, sprint := range project.Sprints {
					if sprint.Name == input {
						return fmt.Errorf("спринт с именем '%s' уже существует", input)
					}
				}
			}

			return nil
		},
	}

	sprintName, err := namePrompt.Run()
	if err != nil {
		a.Logger.Warnf("Отмена создания спринта: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	// Ввод описания спринта
	descPrompt := promptui.Prompt{
		Label: "Описание спринта",
	}

	description, _ := descPrompt.Run()

	// Создание спринта
	err = a.ProjectService.CreateSprint(a.Projects, projectName, sprintName, description)
	if err != nil {
		a.Logger.Errorf("Ошибка создания спринта: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	a.Logger.Infof("Спринт '%s' создан для проекта '%s'", sprintName, projectName)
	fmt.Printf("Спринт '%s' успешно создан для проекта '%s'\n", sprintName, projectName)
}

// setActiveSprint - установка активного спринта для проекта
func (a *App) setActiveSprint() {
	// Выбор проекта
	projectName := a.chooseProject()
	if projectName == "" {
		return
	}

	project := a.Projects[projectName]

	// Проверка наличия спринтов
	if project.Sprints == nil || len(project.Sprints) == 0 {
		fmt.Printf("У проекта '%s' нет спринтов\n", projectName)
		return
	}

	// Получение списка спринтов
	sprints, err := a.ProjectService.GetProjectSprints(a.Projects, projectName)
	if err != nil {
		a.Logger.Errorf("Ошибка получения спринтов: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	// Создание списка для выбора
	var options []string
	options = append(options, "← Назад")

	for _, sprint := range sprints {
		prefix := "  "
		if sprint.IsActive {
			prefix = "▶ "
		}
		options = append(options, prefix+sprint.Name)
	}

	// Выбор спринта
	prompt := promptui.Select{
		Label: "Выберите спринт",
		Items: options,
	}

	idx, result, err := prompt.Run()
	if err != nil || idx == 0 {
		return
	}

	// Получение ID выбранного спринта
	selectedName := strings.TrimPrefix(strings.TrimPrefix(result, "▶ "), "  ")
	var selectedID string

	for _, sprint := range sprints {
		if sprint.Name == selectedName {
			selectedID = sprint.ID
			break
		}
	}

	// Установка активного спринта
	err = a.ProjectService.SetActiveSprint(a.Projects, projectName, selectedID)
	if err != nil {
		a.Logger.Errorf("Ошибка установки активного спринта: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	a.Logger.Infof("Спринт '%s' установлен как активный для проекта '%s'", selectedName, projectName)
	fmt.Printf("Спринт '%s' установлен как активный для проекта '%s'\n", selectedName, projectName)
}

// viewProjectSprints - просмотр спринтов проекта
func (a *App) viewProjectSprints() {
	// Выбор проекта
	projectName := a.chooseProject()
	if projectName == "" {
		return
	}

	project := a.Projects[projectName]

	// Проверка наличия спринтов
	if project.Sprints == nil || len(project.Sprints) == 0 {
		fmt.Printf("У проекта '%s' нет спринтов\n", projectName)
		return
	}

	// Получение списка спринтов
	sprints, err := a.ProjectService.GetProjectSprints(a.Projects, projectName)
	if err != nil {
		a.Logger.Errorf("Ошибка получения спринтов: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("\nСпринты проекта '%s':\n", projectName)

	for _, sprint := range sprints {
		status := "Неактивный"
		if sprint.IsActive {
			status = "Активный"
		}

		fmt.Printf("\n%s - %s\n", sprint.Name, status)
		if sprint.Description != "" {
			fmt.Printf("  Описание: %s\n", sprint.Description)
		}
		fmt.Printf("  Дата начала: %s\n", sprint.StartDate)

		// Подсчет времени по спринту
		var total int
		for _, entry := range sprint.Entries {
			total += entry.TimeSpent
		}

		fmt.Printf("  Общее время: %v сек\n", total)

		// Вывод записей спринта
		if len(sprint.Entries) > 0 {
			fmt.Println("  Записи:")
			for _, entry := range sprint.Entries {
				fmt.Printf("    %s - %v сек: %s\n", entry.Date, entry.TimeSpent, entry.Description)
			}
		}
	}
}

// showSummary - вывод сводки по проектам
func (a *App) showSummary() {
	if len(a.Projects) == 0 {
		fmt.Println("Нет данных для отображения.")
		return
	}

	for name, project := range a.Projects {
		fmt.Printf("\nПроект \"%s\":\n", name)

		// Общее время по проекту
		var totalProject int
		for _, entry := range project.Entries {
			totalProject += entry.TimeSpent
		}

		fmt.Printf("  Общее время: %v сек\n", totalProject)

		// Если есть спринты, показываем статистику по ним
		if project.Sprints != nil && len(project.Sprints) > 0 {
			fmt.Println("  Спринты:")

			sprints, _ := a.ProjectService.GetProjectSprints(a.Projects, name)
			for _, sprint := range sprints {
				var totalSprint int
				for _, entry := range sprint.Entries {
					totalSprint += entry.TimeSpent
				}

				status := ""
				if sprint.IsActive {
					status = " (Активный)"
				}

				fmt.Printf("    %s%s: %v сек\n", sprint.Name, status, totalSprint)
			}
		}

		// Показываем записи проекта
		fmt.Println("  Записи:")
		for _, entry := range project.Entries {
			fmt.Printf("    %s - %v сек: %s\n", entry.Date, entry.TimeSpent, entry.Description)
		}
	}
}
