package systray

import (
	"fmt"
	"sync"
	"time"

	"github.com/MWT-proger/time-tracking/pkg/logger"
	"github.com/getlantern/systray"
)

// SystrayHandler - обработчик системного трея
type SystrayHandler struct {
	mu               sync.RWMutex
	TrackedProject   string
	TrackingStart    *time.Time
	UpdateTrayTicker *time.Ticker
	Logger           logger.Logger
	StartTracking    func()
	StopTracking     func()
	OnExit           func()
}

// NewSystrayHandler - создание нового обработчика системного трея
func NewSystrayHandler(
	log logger.Logger,

) *SystrayHandler {
	return &SystrayHandler{
		Logger: log,
	}
}

// StartTrayTicker - запуск тикера обновления системного трея
func (h *SystrayHandler) StartTrayTicker() {
	h.mu.Lock()
	if h.UpdateTrayTicker != nil {
		h.mu.Unlock()
		return
	}
	h.UpdateTrayTicker = time.NewTicker(1 * time.Second)
	ticker := h.UpdateTrayTicker
	h.mu.Unlock()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.Logger.Errorf("Паника в горутине обновления таймера: %v", r)
				h.mu.Lock()
				if h.UpdateTrayTicker == ticker {
					h.UpdateTrayTicker = nil
				}
				h.mu.Unlock()
			}
		}()

		for range ticker.C {
			func() {
				defer func() {
					if r := recover(); r != nil {
						h.Logger.Errorf("Ошибка при обновлении заголовка таймера: %v", r)
					}
				}()

				h.mu.RLock()
				trackingStart := h.TrackingStart
				trackedProject := h.TrackedProject
				h.mu.RUnlock()

				if trackingStart != nil {
					elapsed := time.Since(*trackingStart)
					title := fmt.Sprintf("%s: %v", trackedProject, elapsed.Round(time.Second))
					systray.SetTitle(title)
				}
			}()
		}
	}()
}

// StopTrayTicker - остановка тикера обновления системного трея
func (h *SystrayHandler) StopTrayTicker() {
	h.mu.Lock()
	if h.UpdateTrayTicker != nil {
		h.UpdateTrayTicker.Stop()
		h.UpdateTrayTicker = nil
	}
	h.mu.Unlock()

	defer func() {
		if r := recover(); r != nil {
			h.Logger.Errorf("Ошибка при установке заголовка таймера: %v", r)
		}
	}()

	systray.SetTitle("Таймер")
}

// SetTracking - установка отслеживаемого проекта
func (h *SystrayHandler) SetTracking(project string, start *time.Time) {
	h.mu.Lock()
	wasTracking := h.TrackingStart != nil
	h.TrackedProject = project
	h.TrackingStart = start
	tickerExists := h.UpdateTrayTicker != nil
	h.mu.Unlock()

	if start != nil {
		if !tickerExists {
			h.StartTrayTicker()
		}
	} else {
		if wasTracking {
			h.StopTrayTicker()
		}
	}
}

// Run - запуск системного трея
func (h *SystrayHandler) Run() {
	systray.Run(h.onReady, h.onExit)
}

// onReady - обработчик готовности системного трея
func (h *SystrayHandler) onReady() {
	// Загружаем иконку из встроенных ресурсов
	iconData, err := GetIcon()
	if err != nil {
		// Если не удалось загрузить иконку, используем встроенную иконку
		h.Logger.Errorf("Не удалось загрузить иконку: %v", err)
		return
	}

	// Проверяем, что массив байтов не пустой
	if len(iconData) == 0 {
		h.Logger.Error("Получен пустой массив байтов иконки")
		return
	}

	systray.SetIcon(iconData)
	systray.SetTitle("Таймер")
	systray.SetTooltip("Учет времени")

	// Создаем пункты меню
	mStart := systray.AddMenuItem("Начать отслеживание", "Начать отслеживание времени")
	mStop := systray.AddMenuItem("Остановить отслеживание", "Остановить отслеживание времени")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Выход", "Выход из приложения")

	// Обработка событий меню
	go func() {
		for {
			select {
			case <-mStart.ClickedCh:
				if h.StartTracking != nil {
					h.StartTracking()
				}
			case <-mStop.ClickedCh:
				if h.StopTracking != nil {
					h.StopTracking()
				}
			case <-mQuit.ClickedCh:
				h.Logger.Info("Выход из приложения через системный трей")
				systray.Quit()
				return
			}
		}
	}()
}

// onExit - обработчик выхода из системного трея
func (h *SystrayHandler) onExit() {
	// Вызываем обработчик выхода, если он задан
	if h.OnExit != nil {
		h.OnExit()
	}
}

// Quit - завершение работы системного трея
func (h *SystrayHandler) Quit() {
	systray.Quit()
}
