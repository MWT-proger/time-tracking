package handlers

import (
	"fmt"
	"time"

	"github.com/MWT-proger/time-tracking/internal/domain"
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
	fmt.Printf("Отслеживание остановлено для проекта %s. Время: %s\n", projectName, h.FormatDuration(elapsed))
	h.SystrayHandler.SetTracking("", nil)
}

// FormatTimeSpent - форматирует время в виде "Xh Ym Zs"
func (h *Handlers) FormatTimeSpent(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	result := ""
	if hours > 0 {
		result += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 || hours > 0 {
		result += fmt.Sprintf("%dm ", minutes)
	}
	result += fmt.Sprintf("%ds", secs)

	return result
}

// FormatDuration - форматирует time.Duration в виде "Xh Ym Zs"
func (h *Handlers) FormatDuration(duration time.Duration) string {
	seconds := int(duration.Seconds())
	return h.FormatTimeSpent(seconds)
}

// ShowSummary - вывод сводки по проектам
func (h *Handlers) ShowSummary() {
	h.Logger.Debug("Отображение сводки по проектам")

	if len(h.Projects) == 0 {
		h.Logger.Debug("Нет данных для отображения")
		fmt.Println("Нет данных для отображения.")
		return
	}

	// Разделяем проекты на активные и архивные
	activeProjects := make(map[string]*domain.Project)
	archivedProjects := make(map[string]*domain.Project)

	for name, project := range h.Projects {
		if project.Archived {
			archivedProjects[name] = project
		} else {
			activeProjects[name] = project
		}
	}

	h.Logger.Debugf("Найдено активных проектов: %d, архивных проектов: %d",
		len(activeProjects), len(archivedProjects))

	// Выводим активные проекты
	if len(activeProjects) > 0 {
		fmt.Println("\nАктивные проекты:")
		for name, project := range activeProjects {
			fmt.Printf("\nПроект \"%s\":\n", name)

			// Общее время по проекту
			var totalProject int
			for _, entry := range project.Entries {
				totalProject += entry.TimeSpent
			}

			fmt.Printf("  Общее время: %s\n", h.FormatTimeSpent(totalProject))

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

					fmt.Printf("    %s%s: %s\n", sprint.Name, status, h.FormatTimeSpent(totalSprint))
				}
			}

			// Показываем записи проекта
			if len(project.Entries) > 0 {
				fmt.Println("  Записи:")
				for _, entry := range project.Entries {
					fmt.Printf("    %s - %s: %s\n", entry.Date, h.FormatTimeSpent(entry.TimeSpent), entry.Description)
				}
			}
		}
	}

	// Выводим архивные проекты
	if len(archivedProjects) > 0 {
		h.Logger.Debug("Отображение архивных проектов")
		fmt.Println("\nАрхивные проекты:")
		for name, project := range archivedProjects {
			var totalProject int
			for _, entry := range project.Entries {
				totalProject += entry.TimeSpent
			}

			fmt.Printf("  %s: %s\n", name, h.FormatTimeSpent(totalProject))
		}
	}
}
