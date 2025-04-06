package app

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
)

// SystrayHandler - обработчик системного трея
type SystrayHandler struct {
	TrackedProject   string
	TrackingStart    *time.Time
	UpdateTrayTicker *time.Ticker
}

// NewSystrayHandler - создание нового обработчика системного трея
func NewSystrayHandler() *SystrayHandler {
	return &SystrayHandler{}
}

// StartTrayTicker - запуск тикера обновления системного трея
func (h *SystrayHandler) StartTrayTicker() {
	if h.UpdateTrayTicker != nil {
		return
	}
	h.UpdateTrayTicker = time.NewTicker(1 * time.Second)
	go func() {
		for range h.UpdateTrayTicker.C {
			if h.TrackingStart != nil {
				elapsed := time.Since(*h.TrackingStart)
				systray.SetTitle(fmt.Sprintf("%s: %v", h.TrackedProject, elapsed.Round(time.Second)))
			}
		}
	}()
}

// StopTrayTicker - остановка тикера обновления системного трея
func (h *SystrayHandler) StopTrayTicker() {
	if h.UpdateTrayTicker != nil {
		h.UpdateTrayTicker.Stop()
		h.UpdateTrayTicker = nil
	}
	systray.SetTitle("Таймер")
}

// SetTracking - установка отслеживаемого проекта
func (h *SystrayHandler) SetTracking(project string, start *time.Time) {
	h.TrackedProject = project
	h.TrackingStart = start
	if start != nil {
		h.StartTrayTicker()
	} else {
		h.StopTrayTicker()
	}
}
