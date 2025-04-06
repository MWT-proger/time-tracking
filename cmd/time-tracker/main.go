package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MWT-proger/time-tracking/internal/app"
	"github.com/MWT-proger/time-tracking/pkg/config"
	"github.com/MWT-proger/time-tracking/pkg/logger"
)

func main() {
	// Парсим флаги командной строки
	cfg := config.ParseFlags(app.Version)

	// Создаем директорию для логов
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		fmt.Printf("Ошибка создания директории для логов: %v\n", err)
	}

	// Создаем директорию для данных
	dataDir := filepath.Dir(cfg.DataFile)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("Ошибка создания директории для данных: %v\n", err)
	}

	// Создаем файловый логгер
	logFile := cfg.GetLogFile()
	fileLogger, err := logger.NewFileLogger(cfg.LogLevel, logFile)
	if err != nil {
		fmt.Printf("Ошибка создания файлового логгера: %v\n", err)
		// Если не удалось создать файловый логгер, используем стандартный
		fileLogger = logger.NewLogger(cfg.LogLevel, nil)
	} else {
		// Логируем запуск приложения
		fileLogger.Infof("Запуск трекера времени v%s (сборка: %s, коммит: %s)",
			app.Version, app.BuildDate, app.GitCommit)
	}

	fmt.Printf("Трекер времени v%s (сборка: %s, коммит: %s)\n",
		app.Version, app.BuildDate, app.GitCommit)
	fmt.Printf("Логи сохраняются в: %s\n", logFile)
	fmt.Printf("Файл данных: %s\n", cfg.DataFile)

	// Создаем и инициализируем приложение
	application := app.NewApp(cfg, fileLogger)

	if err := application.Initialize(); err != nil {
		fmt.Printf("Ошибка инициализации: %v\n", err)
		os.Exit(1)
	}

	application.Run()
}
