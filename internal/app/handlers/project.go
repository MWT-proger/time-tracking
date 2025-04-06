package handlers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
)

// CreateProject - создание нового проекта
func (h *Handlers) CreateProject() {
	prompt := promptui.Prompt{
		Label: "Название проекта",
		Validate: func(input string) error {
			// Проверка на пустое имя
			if input == "" {
				return fmt.Errorf("имя проекта не может быть пустым")
			}

			// Проверка на уникальность имени
			if _, exists := h.Projects[input]; exists {
				return fmt.Errorf("проект с именем '%s' уже существует", input)
			}

			return nil
		},
	}

	name, err := prompt.Run()
	if err != nil {
		h.Logger.Warnf("Отмена создания проекта: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	h.Logger.Infof("Попытка создания проекта: %s", name)

	err = h.ProjectService.CreateProject(h.Projects, name)
	if err != nil {
		h.Logger.Errorf("Ошибка создания проекта: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	h.Logger.Infof("Проект создан: %s", name)
	fmt.Printf("Проект '%s' успешно создан\n", name)
}

// ChooseProject - выбор проекта из списка
func (h *Handlers) ChooseProject() string {
	// Создаем список проектов с информацией о статусе
	var activeProjects []string
	var inactiveProjects []string

	for name, project := range h.Projects {
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
		h.Logger.Warnf("Ошибка выбора проекта: %v", err)
		return ""
	}

	// Если выбрана опция "Назад" (индекс 0), возвращаем пустую строку
	if idx == 0 {
		return ""
	}

	// Удаляем префикс статуса из имени проекта
	return strings.TrimPrefix(strings.TrimPrefix(result, "▶ "), "  ")
}

// SelectAndManageProject - выбор и управление проектом
func (h *Handlers) SelectAndManageProject() {
	// Выбор проекта
	projectName := h.ChooseProject()
	if projectName == "" {
		return
	}

	// Меню управления проектом
	for {
		prompt := promptui.Select{
			Label: fmt.Sprintf("Проект: %s", projectName),
			Items: []string{
				"Начать отслеживание",
				"Остановить отслеживание",
				"Управление спринтами",
				"Статистика проекта",
				"Назад в главное меню",
			},
		}
		_, cmd, _ := prompt.Run()

		switch cmd {
		case "Начать отслеживание":
			h.StartTrackingForProject(projectName)
		case "Остановить отслеживание":
			h.StopTrackingForProject(projectName)
		case "Управление спринтами":
			h.ManageSprintsForProject(projectName)
		case "Статистика проекта":
			h.ShowProjectStatistics(projectName)
		case "Назад в главное меню":
			return
		}
	}
}

// ShowProjectStatistics - вывод статистики по конкретному проекту
func (h *Handlers) ShowProjectStatistics(projectName string) {
	project := h.Projects[projectName]

	fmt.Printf("\nСтатистика проекта \"%s\":\n", projectName)

	// Общее время по проекту
	var totalProject int
	for _, entry := range project.Entries {
		totalProject += entry.TimeSpent
	}

	fmt.Printf("  Общее время: %v сек\n", totalProject)

	// Если есть спринты, показываем статистику по ним
	if project.Sprints != nil && len(project.Sprints) > 0 {
		fmt.Println("  Спринты:")

		sprints, _ := h.ProjectService.GetProjectSprints(h.Projects, projectName)
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
