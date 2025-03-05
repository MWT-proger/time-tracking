package app

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
)

var (
	trackedProject   string
	trackingStart    *time.Time
	updateTrayTicker *time.Ticker
)

func onReady() {
	systray.SetIcon(nil)
	systray.SetTitle("Таймер")
	systray.SetTooltip("Учет времени")

	startItem := systray.AddMenuItem("Начать отслеживание", "Начать отслеживание времени")
	stopItem := systray.AddMenuItem("Остановить отслеживание", "Остановить отслеживание времени")
	exitItem := systray.AddMenuItem("Выход", "Выход из приложения")

	go func() {
		for {
			select {
			case <-startItem.ClickedCh:
				// Логика для начала отслеживания
			case <-stopItem.ClickedCh:
				// Логика для остановки отслеживания
			case <-exitItem.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	if updateTrayTicker != nil {
		updateTrayTicker.Stop()
	}
}

func startTrayTicker() {
	if updateTrayTicker != nil {
		return
	}
	updateTrayTicker = time.NewTicker(1 * time.Second)
	go func() {
		for range updateTrayTicker.C {
			if trackingStart != nil {
				elapsed := time.Since(*trackingStart)
				systray.SetTitle(fmt.Sprintf("%s: %v", trackedProject, elapsed.Round(time.Second)))
			}
		}
	}()
}

func stopTrayTicker() {
	if updateTrayTicker != nil {
		updateTrayTicker.Stop()
		updateTrayTicker = nil
	}
	systray.SetTitle("Таймер")
}
