package systray

import (
	"embed"
)

//go:embed assets/icons/clock.png
var iconFS embed.FS

// GetIcon - получение иконки из встроенных ресурсов
func GetIcon() ([]byte, error) {
	return iconFS.ReadFile("assets/icons/clock.png")
}
