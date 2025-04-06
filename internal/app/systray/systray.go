package systray

import (
	"time"
)

// Handler - интерфейс для работы с системным треем
type Handler interface {
	SetTracking(projectName string, startTime *time.Time)
}

// SystrayHandler - обработчик системного трея
type SystrayHandler struct {
	CurrentProject string
	StartTime      *time.Time
}

// NewSystrayHandler - создание нового обработчика системного трея
func NewSystrayHandler() *SystrayHandler {
	return &SystrayHandler{}
}

// SetTracking - установка информации об отслеживании
func (h *SystrayHandler) SetTracking(projectName string, startTime *time.Time) {
	h.CurrentProject = projectName
	h.StartTime = startTime
}
