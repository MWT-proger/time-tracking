package handlers

import (
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
)

// StartTrackingForProject - начало отслеживания времени для конкретного проекта
func (h *Handlers) StartTrackingForProject(projectName string) {
	project := h.Projects[projectName]

	// Если у проекта есть спринты, предлагаем выбрать активный спринт
	if project.Sprints != nil && len(project.Sprints) > 0 {
		prompt := promptui.Select{
			Label: "Выберите спринт для отслеживания",
			Items: []string{
				"Использовать текущий активный спринт",
				"Выбрать другой спринт",
			},
		}
		_, choice, _ := prompt.Run()

		if choice == "Выбрать другой спринт" {
			h.SetActiveSprintForProject(projectName)
		}
	}

	h.Logger.Infof("Попытка начать отслеживание для проекта: %s", projectName)
	err := h.TrackingService.StartTracking(h.Projects, projectName)
	if err != nil {
		h.Logger.Errorf("Ошибка начала отслеживания: %v", err)
		fmt.Println(err)
		return
	}

	h.Logger.Infof("Начато отслеживание для проекта: %s", projectName)
	fmt.Println("Начато отслеживание для проекта:", projectName)

	project = h.Projects[projectName]
	h.SystrayHandler.SetTracking(projectName, project.StartTime)

	go func() {
		time.Sleep(time.Duration(h.Config.NotificationTime) * time.Second)
		h.Logger.Infof("Отправка уведомления о перерыве для проекта: %s", projectName)
		h.TrackingService.Notify(projectName)
	}()
}

// StopTrackingForProject - остановка отслеживания времени для конкретного проекта
func (h *Handlers) StopTrackingForProject(projectName string) {
	project := h.Projects[projectName]

	// Проверяем, запущено ли отслеживание для этого проекта
	if project.StartTime == nil {
		fmt.Printf("Отслеживание для проекта '%s' не запущено\n", projectName)
		return
	}

	prompt := promptui.Prompt{
		Label: "Что сделано",
	}
	description, _ := prompt.Run()

	h.Logger.Infof("Попытка остановить отслеживание для проекта: %s", projectName)
	elapsed, err := h.TrackingService.StopTracking(h.Projects, projectName, description)
	if err != nil {
		h.Logger.Errorf("Ошибка остановки отслеживания: %v", err)
		fmt.Println(err)
		return
	}

	h.Logger.Infof("Отслеживание остановлено для проекта %s. Время: %v", projectName, elapsed)
	fmt.Printf("Отслеживание остановлено для проекта %s. Время: %v\n", projectName, elapsed)
	h.SystrayHandler.SetTracking("", nil)
}

// ShowSummary - вывод сводки по проектам
func (h *Handlers) ShowSummary() {
	if len(h.Projects) == 0 {
		fmt.Println("Нет данных для отображения.")
		return
	}

	for name, project := range h.Projects {
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

			sprints, _ := h.ProjectService.GetProjectSprints(h.Projects, name)
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
