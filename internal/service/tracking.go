package service

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/MWT-proger/time-tracking/internal/domain"
	"github.com/MWT-proger/time-tracking/pkg/config"
	"github.com/MWT-proger/time-tracking/pkg/logger"
	"github.com/google/uuid"
)

// TrackingService - сервис для отслеживания времени
type TrackingService struct {
	ProjectService   *ProjectService
	Logger           logger.Logger
	NotificationTime int
}

// NewTrackingService - создание нового сервиса отслеживания
func NewTrackingService(projectService *ProjectService, log logger.Logger, cfg *config.Config) *TrackingService {
	return &TrackingService{
		ProjectService:   projectService,
		Logger:           log,
		NotificationTime: cfg.NotificationTime,
	}
}

// StartTracking - начало отслеживания времени
func (s *TrackingService) StartTracking(data map[string]*domain.Project, name string) error {
	s.Logger.Debugf("Попытка начать отслеживание для проекта: %s", name)
	project, exists := data[name]
	if !exists {
		return fmt.Errorf("проект '%s' не существует", name)
	}

	if project.StartTime != nil {
		return fmt.Errorf("отслеживание для проекта '%s' уже запущено", name)
	}

	// Проверяем, есть ли у проекта активный этап
	if project.Sprints != nil && len(project.Sprints) > 0 && project.ActiveSprint != "" {
		if _, exists := project.Sprints[project.ActiveSprint]; !exists {
			// Если активный этап не существует, сбрасываем его
			project.ActiveSprint = ""
		}
	}

	now := time.Now()
	project.StartTime = &now

	return s.ProjectService.SaveData(data)
}

// StopTracking - остановка отслеживания времени
func (s *TrackingService) StopTracking(data map[string]*domain.Project, name string, description string) (time.Duration, error) {
	s.Logger.Debugf("Попытка остановить отслеживание для проекта: %s", name)
	project, exists := data[name]
	if !exists || project.StartTime == nil {
		return 0, fmt.Errorf("отслеживание для проекта '%s' не запущено", name)
	}

	elapsed := time.Since(*project.StartTime)
	seconds := int(elapsed.Seconds())

	entry := domain.TimeEntry{
		Date:        time.Now().Format("2006-01-02 15:04:05"),
		TimeSpent:   seconds,
		Description: description,
	}

	// Если у проекта есть активный этап, добавляем запись к нему
	if project.Sprints != nil && project.ActiveSprint != "" {
		if sprint, exists := project.Sprints[project.ActiveSprint]; exists {
			if sprint.Entries == nil {
				sprint.Entries = make(map[string]domain.TimeEntry)
			}
			entryID := uuid.New().String()
			sprint.Entries[entryID] = entry
		}
	}

	// В любом случае добавляем запись к проекту для обратной совместимости
	if project.Entries == nil {
		project.Entries = []domain.TimeEntry{}
	}
	project.Entries = append(project.Entries, entry)

	project.StartTime = nil

	if err := s.ProjectService.SaveData(data); err != nil {
		return 0, err
	}

	return elapsed, nil
}

// Notify - отправка уведомления
func (s *TrackingService) Notify(project string) error {
	s.Logger.Debugf("Отправка уведомления для проекта: %s", project)
	return exec.Command("notify-send", "Оповещение", fmt.Sprintf("Вы работаете над проектом '%s' уже %d минут! Время сделать перерыв.", project, s.NotificationTime/60)).Run()
}

// Summary - вывод сводки по проектам
func (s *TrackingService) Summary(data map[string]*domain.Project) map[string]int {
	result := make(map[string]int)

	for name, project := range data {
		var total int
		for _, entry := range project.Entries {
			total += entry.TimeSpent
		}
		result[name] = total
	}

	return result
}
