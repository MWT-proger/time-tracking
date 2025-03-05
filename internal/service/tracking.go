package service

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/MWT-proger/time-tracking/internal/domain"
)

type TrackingService struct {
	projectService *ProjectService
}

func NewTrackingService() *TrackingService {
	return &TrackingService{
		projectService: NewProjectService(),
	}
}

func (s *TrackingService) StartTracking(name string) {
	data := s.projectService.GetProjects()
	project, exists := data[name]
	if !exists || project.StartTime != nil {
		return
	}
	start := time.Now()
	project.StartTime = &start
	s.projectService.saveData(data)
}

func (s *TrackingService) StopTracking(name, description string) {
	data := s.projectService.GetProjects()
	project, exists := data[name]
	if !exists || project.StartTime == nil {
		return
	}
	elapsed := time.Since(*project.StartTime)
	project.Entries = append(project.Entries, domain.Entry{
		TimeSpent:   int(elapsed.Seconds()),
		Description: description,
		Date:        time.Now().Format("2006-01-02 15:04:05"),
	})
	project.StartTime = nil
	s.projectService.saveData(data)
}

func notify(project string) {
	exec.Command("notify-send", "Оповещение", fmt.Sprintf("Вы работаете над проектом '%s' уже 25 минут! Время сделать перерыв.", project)).Run()
}
