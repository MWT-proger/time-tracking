package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/getlantern/systray"
	"github.com/manifoldco/promptui"
)

// Путь к файлу для хранения данных
var dataFile = filepath.Join(os.Getenv("HOME"), "учет_времени.json")

// Время для уведомления (25 минут в секундах)
const notificationTime = 1500 // 25 минут

// Структура данных о проекте
type Project struct {
	Entries   []Entry    `json:"entries"`
	StartTime *time.Time `json:"start_time,omitempty"`
}

// Структура записи времени
type Entry struct {
	TimeSpent   int    `json:"time_spent"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

// Глобальные переменные для системного трея
var (
	trackedProject   string
	trackingStart    *time.Time
	updateTrayTicker *time.Ticker
)

// Загрузка данных из файла
func loadData() (map[string]*Project, error) {
	data := make(map[string]*Project)
	file, err := os.Open(dataFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func notify(project string) {
	exec.Command("notify-send", "Оповещение", fmt.Sprintf("Вы работаете над проектом '%s' уже 25 минут! Время сделать перерыв.", project)).Run()
}

// Сохранение данных в файл
func saveData(data map[string]*Project) error {
	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(data)
}

// Создание нового проекта
func createProject(data map[string]*Project, name string) {
	if _, exists := data[name]; exists {
		fmt.Println("Проект уже существует.")
		return
	}
	data[name] = &Project{}
	saveData(data)
	fmt.Println("Проект создан:", name)
}

// Начало отслеживания времени
func startTracking(data map[string]*Project, name string) {
	project, exists := data[name]
	if !exists {
		fmt.Println("Проект не найден.")
		return
	}
	if project.StartTime != nil {
		fmt.Println("Отслеживание уже запущено.")
		return
	}
	start := time.Now()
	project.StartTime = &start
	saveData(data)
	fmt.Println("Начато отслеживание для проекта:", name)

	trackedProject = name
	trackingStart = &start
	startTrayTicker()

	go func() {
		time.Sleep(notificationTime * time.Second)
		notify(name)
	}()
}

// Остановка отслеживания времени
func stopTracking(data map[string]*Project, name string, description string) {
	project, exists := data[name]
	if !exists || project.StartTime == nil {
		fmt.Println("Проект не найден или неактивен.")
		return
	}
	elapsed := time.Since(*project.StartTime)
	project.Entries = append(project.Entries, Entry{
		TimeSpent:   int(elapsed.Seconds()),
		Description: description,
		Date:        time.Now().Format("2006-01-02 15:04:05"),
	})
	project.StartTime = nil
	saveData(data)
	fmt.Printf("Отслеживание остановлено для проекта %s. Время: %v\n", name, elapsed)

	stopTrayTicker()
}

// Вывод сводки по проектам
func summary(data map[string]*Project) {
	if len(data) == 0 {
		fmt.Println("Нет данных для отображения.")
		return
	}

	for name, project := range data {
		fmt.Printf("\nПроект \"%s\":\n", name)
		var total int
		for _, entry := range project.Entries {
			total += entry.TimeSpent
			fmt.Printf("  %s - %v сек: %s\n", entry.Date, entry.TimeSpent, entry.Description)
		}
		fmt.Printf("  Общее время: %v сек\n", total)
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

func onReady() {
	systray.SetIcon(nil) // Добавьте путь к иконке или используйте заглушку
	systray.SetTitle("Таймер")
	systray.SetTooltip("Учет времени")

	// Элементы меню
	startItem := systray.AddMenuItem("Начать отслеживание", "Начать отслеживание времени")
	stopItem := systray.AddMenuItem("Остановить отслеживание", "Остановить отслеживание времени")
	exitItem := systray.AddMenuItem("Выход", "Выход из приложения")

	go func() {
		for {
			select {
			case <-startItem.ClickedCh:
				// Добавьте логику выбора проекта
			case <-stopItem.ClickedCh:
				// Добавьте логику остановки отслеживания
			case <-exitItem.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}

func onExit() {
	// Логика при выходе
}

func main() {
	data, _ := loadData()
	go systray.Run(onReady, onExit)

	for {
		prompt := promptui.Select{
			Label: "Выберите команду",
			Items: []string{"Создать проект", "Начать отслеживание", "Остановить отслеживание", "Сводка", "Выход"},
		}
		_, cmd, _ := prompt.Run()

		switch cmd {
		case "Создать проект":
			prompt := promptui.Prompt{
				Label: "Название проекта",
			}
			name, _ := prompt.Run()
			createProject(data, name)
		case "Начать отслеживание":
			projectName := chooseProject(data)
			startTracking(data, projectName)
		case "Остановить отслеживание":
			projectName := chooseProject(data)
			prompt := promptui.Prompt{
				Label: "Что сделано",
			}
			description, _ := prompt.Run()
			stopTracking(data, projectName, description)
		case "Сводка":
			summary(data)
		case "Выход":
			saveData(data)
			return
		}
	}
}

func chooseProject(data map[string]*Project) string {
	names := []string{}
	for name := range data {
		names = append(names, name)
	}
	prompt := promptui.Select{
		Label: "Выберите проект",
		Items: names,
	}
	_, project, _ := prompt.Run()
	return project
}
