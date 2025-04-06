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

// ArchiveProject - архивирование проекта
func (h *Handlers) ArchiveProject(projectName string) {
	h.Logger.Infof("Попытка архивировать проект: %s", projectName)

	// Проверяем, запущено ли отслеживание для проекта
	project := h.Projects[projectName]
	if project.StartTime != nil {
		h.Logger.Warnf("Невозможно архивировать проект '%s' с запущенным отслеживанием", projectName)
		fmt.Printf("Невозможно архивировать проект '%s' с запущенным отслеживанием\n", projectName)
		return
	}

	// Запрашиваем подтверждение с помощью Select вместо Prompt
	prompt := promptui.Select{
		Label: fmt.Sprintf("Вы уверены, что хотите архивировать проект '%s'?", projectName),
		Items: []string{"Да", "Нет"},
	}

	idx, _, err := prompt.Run()
	if err != nil {
		h.Logger.Warnf("Отмена архивирования проекта из-за ошибки: %v", err)
		fmt.Println("Архивирование отменено")
		return
	}

	// Если выбран второй вариант (Нет), отменяем операцию
	if idx != 0 {
		h.Logger.Infof("Пользователь отменил архивирование проекта '%s'", projectName)
		fmt.Println("Архивирование отменено")
		return
	}

	err = h.ProjectService.ArchiveProject(h.Projects, projectName)
	if err != nil {
		h.Logger.Errorf("Ошибка архивирования проекта: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	h.Logger.Infof("Проект '%s' архивирован", projectName)
	fmt.Printf("Проект '%s' успешно архивирован\n", projectName)
}

// RestoreProject - восстановление проекта из архива
func (h *Handlers) RestoreProject(projectName string) {
	h.Logger.Infof("Попытка восстановить проект из архива: %s", projectName)

	// Запрашиваем подтверждение с помощью Select вместо Prompt
	prompt := promptui.Select{
		Label: fmt.Sprintf("Вы уверены, что хотите восстановить проект '%s' из архива?", projectName),
		Items: []string{"Да", "Нет"},
	}

	idx, _, err := prompt.Run()
	if err != nil {
		h.Logger.Warnf("Отмена восстановления проекта из-за ошибки: %v", err)
		fmt.Println("Восстановление отменено")
		return
	}

	// Если выбран второй вариант (Нет), отменяем операцию
	if idx != 0 {
		h.Logger.Infof("Пользователь отменил восстановление проекта '%s'", projectName)
		fmt.Println("Восстановление отменено")
		return
	}

	err = h.ProjectService.RestoreProject(h.Projects, projectName)
	if err != nil {
		h.Logger.Errorf("Ошибка восстановления проекта: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	h.Logger.Infof("Проект '%s' восстановлен из архива", projectName)
	fmt.Printf("Проект '%s' успешно восстановлен из архива\n", projectName)
}

// ChooseProject - выбор проекта из списка
func (h *Handlers) ChooseProject() string {
	h.Logger.Debug("Выбор проекта из списка")

	// Создаем списки проектов с информацией о статусе
	var activeProjects []string
	var inactiveProjects []string
	var archivedProjects []string

	for name, project := range h.Projects {
		if project.Archived {
			archivedProjects = append(archivedProjects, "📦 "+name)
		} else if project.StartTime != nil {
			activeProjects = append(activeProjects, "▶ "+name)
		} else {
			inactiveProjects = append(inactiveProjects, "⏹ "+name)
		}
	}

	// Сортируем проекты по алфавиту
	sort.Strings(activeProjects)
	sort.Strings(inactiveProjects)
	sort.Strings(archivedProjects)

	// Объединяем списки: сначала "Назад", затем активные, затем неактивные, затем разделитель, затем архивные
	options := append([]string{"← Назад"}, activeProjects...)
	options = append(options, inactiveProjects...)

	if len(archivedProjects) > 0 {
		h.Logger.Debug("Добавление архивных проектов в список выбора")
		options = append(options, "───────────────")
		options = append(options, archivedProjects...)
	}

	if len(options) == 1 { // Только опция "Назад"
		h.Logger.Debug("Нет доступных проектов для выбора")
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
		h.Logger.Debug("Выбрана опция 'Назад'")
		return ""
	}

	// Если выбран разделитель, возвращаем пустую строку
	if result == "───────────────" {
		h.Logger.Debug("Выбран разделитель")
		return ""
	}

	// Удаляем префикс статуса из имени проекта
	projectName := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(result, "▶ "), "⏹ "), "📦 ")
	h.Logger.Debugf("Выбран проект: %s", projectName)

	return projectName
}

// SelectAndManageProject - выбор и управление проектом
func (h *Handlers) SelectAndManageProject() {
	h.Logger.Debug("Запуск меню выбора и управления проектом")

	// Выбор проекта
	projectName := h.ChooseProject()
	if projectName == "" {
		h.Logger.Debug("Проект не выбран, возврат в главное меню")
		return
	}

	project := h.Projects[projectName]
	h.Logger.Debugf("Управление проектом: %s (архивирован: %v)", projectName, project.Archived)

	// Меню управления проектом
	for {
		var menuItems []string
		var projectLabel string

		if project.Archived {
			menuItems = []string{
				"Восстановить из архива",
				"Статистика проекта",
				"Назад в главное меню",
			}
			projectLabel = fmt.Sprintf("Проект: %s (в архиве)", projectName)
			h.Logger.Debugf("Отображение меню для архивированного проекта: %s", projectName)
		} else {
			menuItems = []string{
				"Начать отслеживание",
				"Остановить отслеживание",
				"Управление спринтами",
				"Статистика проекта",
				"Архивировать проект",
				"Назад в главное меню",
			}
			projectLabel = fmt.Sprintf("Проект: %s", projectName)
			h.Logger.Debugf("Отображение меню для активного проекта: %s", projectName)
		}

		prompt := promptui.Select{
			Label: projectLabel,
			Items: menuItems,
		}
		_, cmd, _ := prompt.Run()

		h.Logger.Debugf("Выбрана команда: %s для проекта %s", cmd, projectName)

		switch cmd {
		case "Начать отслеживание":
			h.StartTrackingForProject(projectName)
		case "Остановить отслеживание":
			h.StopTrackingForProject(projectName)
		case "Управление спринтами":
			h.ManageSprintsForProject(projectName)
		case "Статистика проекта":
			h.ShowProjectStatistics(projectName)
		case "Архивировать проект":
			h.ArchiveProject(projectName)
			// После архивирования возвращаемся в главное меню
			return
		case "Восстановить из архива":
			h.RestoreProject(projectName)
			// Обновляем проект после восстановления
			project = h.Projects[projectName]
			h.Logger.Debugf("Проект %s восстановлен из архива, обновление состояния", projectName)
		case "Назад в главное меню":
			h.Logger.Debugf("Возврат в главное меню из проекта %s", projectName)
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

	fmt.Printf("  Общее время: %s\n", h.FormatTimeSpent(totalProject))

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

			fmt.Printf("    %s%s: %s\n", sprint.Name, status, h.FormatTimeSpent(totalSprint))
		}
	}

	// Показываем записи проекта
	fmt.Println("  Записи:")
	for _, entry := range project.Entries {
		fmt.Printf("    %s - %s: %s\n", entry.Date, h.FormatTimeSpent(entry.TimeSpent), entry.Description)
	}
}
