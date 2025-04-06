package handlers

import (
	"time"

	"github.com/MWT-proger/time-tracking/internal/domain"
	"github.com/MWT-proger/time-tracking/internal/service"
	"github.com/MWT-proger/time-tracking/pkg/config"
	"github.com/MWT-proger/time-tracking/pkg/logger"
	"github.com/manifoldco/promptui"
)

type SystrayHandler interface {
	SetTracking(project string, start *time.Time)
	StopTrayTicker()
}

// Handlers - структура для обработчиков приложения
type Handlers struct {
	ProjectService  *service.ProjectService
	TrackingService *service.TrackingService
	SystrayHandler  SystrayHandler
	Logger          logger.Logger
	Config          *config.Config
	Projects        map[string]*domain.Project
}

// NewHandlers - создание новых обработчиков
func NewHandlers(
	projectService *service.ProjectService,
	trackingService *service.TrackingService,
	systrayHandler SystrayHandler,
	logger logger.Logger,
	config *config.Config,
) *Handlers {
	return &Handlers{
		ProjectService:  projectService,
		TrackingService: trackingService,
		SystrayHandler:  systrayHandler,
		Logger:          logger,
		Config:          config,
	}
}

// SetProjects - установка проектов
func (h *Handlers) SetProjects(projects map[string]*domain.Project) {
	h.Projects = projects
}

// StartTracking - начало отслеживания времени через системный трей
func (h *Handlers) StartTracking() {
	// Выбор проекта
	projectName := h.ChooseProject()
	if projectName == "" {
		return
	}

	h.StartTrackingForProject(projectName)
}

// StopTracking - остановка отслеживания времени через системный трей
func (h *Handlers) StopTracking() {
	// Создаем список активных проектов
	var activeProjects []string

	for name, project := range h.Projects {
		if project.StartTime != nil {
			activeProjects = append(activeProjects, name)
		}
	}

	if len(activeProjects) == 0 {
		return
	}

	if len(activeProjects) == 1 {
		h.StopTrackingForProject(activeProjects[0])
		return
	}

	// Если активных проектов несколько, предлагаем выбрать
	prompt := promptui.Select{
		Label: "Выберите проект для остановки",
		Items: activeProjects,
	}

	_, projectName, err := prompt.Run()
	if err != nil {
		return
	}

	h.StopTrackingForProject(projectName)
}
