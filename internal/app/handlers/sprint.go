package handlers

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// ManageSprintsForProject - управление спринтами для конкретного проекта
func (h *Handlers) ManageSprintsForProject(projectName string) {
	for {
		prompt := promptui.Select{
			Label: fmt.Sprintf("Управление спринтами проекта: %s", projectName),
			Items: []string{
				"Создать спринт",
				"Выбрать активный спринт",
				"Просмотреть спринты",
				"Назад к управлению проектом",
			},
		}
		_, cmd, _ := prompt.Run()

		switch cmd {
		case "Создать спринт":
			h.CreateSprintForProject(projectName)
		case "Выбрать активный спринт":
			h.SetActiveSprintForProject(projectName)
		case "Просмотреть спринты":
			h.ViewSprintsForProject(projectName)
		case "Назад к управлению проектом":
			return
		}
	}
}

// CreateSprintForProject - создание нового спринта для конкретного проекта
func (h *Handlers) CreateSprintForProject(projectName string) {
	// Ввод имени спринта
	namePrompt := promptui.Prompt{
		Label: "Название спринта",
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("имя спринта не может быть пустым")
			}

			// Проверка на уникальность имени спринта
			project := h.Projects[projectName]
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
		h.Logger.Warnf("Отмена создания спринта: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	// Ввод описания спринта
	descPrompt := promptui.Prompt{
		Label: "Описание спринта",
	}

	description, _ := descPrompt.Run()

	// Создание спринта
	err = h.ProjectService.CreateSprint(h.Projects, projectName, sprintName, description)
	if err != nil {
		h.Logger.Errorf("Ошибка создания спринта: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	h.Logger.Infof("Спринт '%s' создан для проекта '%s'", sprintName, projectName)
	fmt.Printf("Спринт '%s' успешно создан для проекта '%s'\n", sprintName, projectName)
}

// SetActiveSprintForProject - установка активного спринта для конкретного проекта
func (h *Handlers) SetActiveSprintForProject(projectName string) {
	project := h.Projects[projectName]

	// Проверка наличия спринтов
	if project.Sprints == nil || len(project.Sprints) == 0 {
		fmt.Printf("У проекта '%s' нет спринтов\n", projectName)
		return
	}

	// Получение списка спринтов
	sprints, err := h.ProjectService.GetProjectSprints(h.Projects, projectName)
	if err != nil {
		h.Logger.Errorf("Ошибка получения спринтов: %v", err)
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
	err = h.ProjectService.SetActiveSprint(h.Projects, projectName, selectedID)
	if err != nil {
		h.Logger.Errorf("Ошибка установки активного спринта: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	h.Logger.Infof("Спринт '%s' установлен как активный для проекта '%s'", selectedName, projectName)
	fmt.Printf("Спринт '%s' установлен как активный для проекта '%s'\n", selectedName, projectName)
}

// ViewSprintsForProject - просмотр спринтов для конкретного проекта
func (h *Handlers) ViewSprintsForProject(projectName string) {
	project := h.Projects[projectName]

	// Проверка наличия спринтов
	if project.Sprints == nil || len(project.Sprints) == 0 {
		fmt.Printf("У проекта '%s' нет спринтов\n", projectName)
		return
	}

	// Получение списка спринтов
	sprints, err := h.ProjectService.GetProjectSprints(h.Projects, projectName)
	if err != nil {
		h.Logger.Errorf("Ошибка получения спринтов: %v", err)
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
