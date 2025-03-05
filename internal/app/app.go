package app

import (
	"fmt"

	"github.com/MWT-proger/time-tracking/internal/service"
	"github.com/manifoldco/promptui"
)

type App struct {
	projectService  *service.ProjectService
	trackingService *service.TrackingService
}

func NewApp() *App {
	fmt.Printf("Time Tracker v%s\n", Version)
	fmt.Printf("Author: %s\n\n", Author)

	return &App{
		projectService:  service.NewProjectService(),
		trackingService: service.NewTrackingService(),
	}
}

func (a *App) Run() {
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
			a.summary()
		case "Выход":
			return
		}
	}
}

func (a *App) createProject() {
	prompt := promptui.Prompt{
		Label: "Название проекта",
	}
	name, _ := prompt.Run()
	a.projectService.CreateProject(name)
}

func (a *App) startTracking() {
	projectName := a.chooseProject()
	if projectName == "" {
		return
	}
	a.trackingService.StartTracking(projectName)
}

func (a *App) stopTracking() {
	projectName := a.chooseProject()
	if projectName == "" {
		return
	}
	prompt := promptui.Prompt{
		Label: "Что сделано",
	}
	description, _ := prompt.Run()
	a.trackingService.StopTracking(projectName, description)
}

func (a *App) summary() {
	projects := a.projectService.GetProjects()
	if len(projects) == 0 {
		fmt.Println("Нет данных для отображения.")
		return
	}

	for name, project := range projects {
		fmt.Printf("\nПроект \"%s\":\n", name)
		var total int
		for _, entry := range project.Entries {
			total += entry.TimeSpent
			fmt.Printf("  %s - %v сек: %s\n", entry.Date, entry.TimeSpent, entry.Description)
		}
		fmt.Printf("  Общее время: %v сек\n", total)
	}
}

func (a *App) chooseProject() string {
	projects := a.projectService.GetProjects()
	names := []string{}
	for name := range projects {
		names = append(names, name)
	}
	names = append(names, "Назад")

	prompt := promptui.Select{
		Label: "Выберите проект",
		Items: names,
	}
	_, project, _ := prompt.Run()

	if project == "Назад" {
		return ""
	}
	return project
}
