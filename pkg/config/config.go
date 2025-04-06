package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config - структура конфигурации приложения
type Config struct {
	// Путь к файлу данных
	DataFile string

	// Директория для логов
	LogDir string

	// Уровень логирования
	LogLevel string

	// Время для уведомления в секундах
	NotificationTime int

	// Версия приложения
	Version string

	// Показать справку и выйти
	ShowHelp bool
}

// DefaultConfig - конфигурация по умолчанию
func DefaultConfig() *Config {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir, _ = os.UserHomeDir()
	}

	return &Config{
		DataFile:         filepath.Join(homeDir, "учет_времени.json"),
		LogDir:           filepath.Join(homeDir, ".time-tracker", "logs"),
		LogLevel:         "info",
		NotificationTime: 1500, // 25 минут в секундах
		ShowHelp:         false,
	}
}

// ParseFlags - парсинг флагов командной строки
func ParseFlags(version string) *Config {
	config := DefaultConfig()
	config.Version = version

	// Определение флагов
	flag.StringVar(&config.DataFile, "data", config.DataFile, "Путь к файлу данных")
	flag.StringVar(&config.LogDir, "log-dir", config.LogDir, "Директория для логов")
	flag.StringVar(&config.LogLevel, "log-level", config.LogLevel, "Уровень логирования (debug, info, warn, error, fatal)")
	flag.IntVar(&config.NotificationTime, "notify-time", config.NotificationTime, "Время для уведомления в секундах")
	flag.BoolVar(&config.ShowHelp, "help", false, "Показать справку и выйти")
	flag.BoolVar(&config.ShowHelp, "h", false, "Показать справку и выйти (сокращение)")

	// Переопределяем стандартный обработчик справки
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Трекер времени v%s\n\n", version)
		fmt.Fprintf(os.Stderr, "Использование: %s [флаги]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Флаги:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nПримеры:\n")
		fmt.Fprintf(os.Stderr, "  %s -data /path/to/data.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -log-level debug\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -notify-time 1800\n", os.Args[0])
	}

	// Парсинг флагов
	flag.Parse()

	// Если запрошена справка, показываем её и выходим
	if config.ShowHelp {
		flag.Usage()
		os.Exit(0)
	}

	return config
}

// GetLogFile - получение пути к файлу логов
func (c *Config) GetLogFile() string {
	return filepath.Join(c.LogDir, "time-tracker.log")
}
