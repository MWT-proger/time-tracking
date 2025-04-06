package service

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/MWT-proger/time-tracking/internal/domain"
	"github.com/MWT-proger/time-tracking/pkg/config"
	"github.com/MWT-proger/time-tracking/pkg/logger"
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
		return fmt.Errorf("проект не найден")
	}
	if project.StartTime != nil {
		return fmt.Errorf("отслеживание уже запущено")
	}
	start := time.Now()
	project.StartTime = &start
	return s.ProjectService.SaveData(data)
}

// StopTracking - остановка отслеживания времени
func (s *TrackingService) StopTracking(data map[string]*domain.Project, name string, description string) (time.Duration, error) {
	s.Logger.Debugf("Попытка остановить отслеживание для проекта: %s", name)
	project, exists := data[name]
	if !exists || project.StartTime == nil {
		return 0, fmt.Errorf("проект не найден или неактивен")
	}
	elapsed := time.Since(*project.StartTime)
	project.Entries = append(project.Entries, domain.Entry{
		TimeSpent:   int(elapsed.Seconds()),
		Description: description,
		Date:        time.Now().Format("2006-01-02 15:04:05"),
	})
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
