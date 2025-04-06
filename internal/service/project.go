package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/MWT-proger/time-tracking/internal/domain"
	"github.com/MWT-proger/time-tracking/pkg/logger"
	"github.com/google/uuid"
)

// ProjectService - сервис для работы с проектами
type ProjectService struct {
	DataFile string
	Logger   logger.Logger
}

// NewProjectService - создание нового сервиса проектов
func NewProjectService(log logger.Logger, dataFile string) *ProjectService {
	return &ProjectService{
		DataFile: dataFile,
		Logger:   log,
	}
}

// LoadData - загрузка данных из файла
func (s *ProjectService) LoadData() (map[string]*domain.Project, error) {
	s.Logger.Debug("Загрузка данных из файла:", s.DataFile)
	data := make(map[string]*domain.Project)

	// Создаем директорию для файла данных, если она не существует
	dir := filepath.Dir(s.DataFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		s.Logger.Errorf("Ошибка создания директории для данных: %v", err)
		return nil, err
	}

	file, err := os.Open(s.DataFile)
	if err != nil && !os.IsNotExist(err) {
		s.Logger.Errorf("Ошибка открытия файла данных: %v", err)
		return nil, err
	}
	if os.IsNotExist(err) {
		s.Logger.Info("Файл данных не существует, будет создан новый")
		return data, nil
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		s.Logger.Errorf("Ошибка декодирования данных: %v", err)
		return data, err
	}
	return data, nil
}

// SaveData - сохранение данных в файл
func (s *ProjectService) SaveData(data map[string]*domain.Project) error {
	s.Logger.Debug("Сохранение данных в файл:", s.DataFile)

	// Создаем директорию для файла данных, если она не существует
	dir := filepath.Dir(s.DataFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		s.Logger.Errorf("Ошибка создания директории для данных: %v", err)
		return err
	}

	file, err := os.Create(s.DataFile)
	if err != nil {
		s.Logger.Errorf("Ошибка создания файла данных: %v", err)
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(data)
}

// CreateProject - создание нового проекта
func (s *ProjectService) CreateProject(data map[string]*domain.Project, name string) error {
	s.Logger.Infof("Создание проекта: %s", name)

	// Проверка на пустое имя
	if name == "" {
		s.Logger.Warn("Попытка создать проект с пустым именем")
		return fmt.Errorf("имя проекта не может быть пустым")
	}

	// Проверка на уникальность имени
	if _, exists := data[name]; exists {
		s.Logger.Warnf("Попытка создать проект с существующим именем: %s", name)
		return fmt.Errorf("проект с именем '%s' уже существует", name)
	}

	data[name] = &domain.Project{}
	return s.SaveData(data)
}

// GetProjectNames - получение списка имен проектов
func (s *ProjectService) GetProjectNames(data map[string]*domain.Project, includeArchived bool) []string {
	s.Logger.Debug("Получение списка имен проектов")

	names := make([]string, 0, len(data))

	for name, project := range data {
		if !project.Archived || includeArchived {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	return names
}

// CreateSprint - создание нового спринта проекта
func (s *ProjectService) CreateSprint(data map[string]*domain.Project, projectName, sprintName, description string) error {
	s.Logger.Infof("Создание спринта '%s' для проекта '%s'", sprintName, projectName)

	// Проверка на существование проекта
	project, exists := data[projectName]
	if !exists {
		s.Logger.Warnf("Попытка создать спринт для несуществующего проекта: %s", projectName)
		return fmt.Errorf("проект '%s' не существует", projectName)
	}

	// Проверка на пустое имя спринта
	if sprintName == "" {
		s.Logger.Warn("Попытка создать спринт с пустым именем")
		return fmt.Errorf("имя спринта не может быть пустым")
	}

	// Инициализация карты спринтов, если она не существует
	if project.Sprints == nil {
		project.Sprints = make(map[string]*domain.Sprint)
	}

	// Проверка на уникальность имени спринта
	for _, sprint := range project.Sprints {
		if sprint.Name == sprintName {
			s.Logger.Warnf("Попытка создать спринт с существующим именем: %s", sprintName)
			return fmt.Errorf("спринт с именем '%s' уже существует в проекте '%s'", sprintName, projectName)
		}
	}

	// Создание уникального ID для спринта
	sprintID := uuid.New().String()

	// Создание нового спринта
	project.Sprints[sprintID] = &domain.Sprint{
		ID:          sprintID,
		Name:        sprintName,
		Description: description,
		StartDate:   time.Now().Format("2006-01-02"),
		Entries:     make(map[string]domain.TimeEntry),
		IsActive:    true,
	}

	// Установка спринта как активного для проекта
	project.ActiveSprint = sprintID

	return s.SaveData(data)
}

// GetProjectSprints - получение списка спринтов проекта
func (s *ProjectService) GetProjectSprints(data map[string]*domain.Project, projectName string) ([]*domain.Sprint, error) {
	project, exists := data[projectName]
	if !exists {
		return nil, fmt.Errorf("проект '%s' не существует", projectName)
	}

	if project.Sprints == nil || len(project.Sprints) == 0 {
		return []*domain.Sprint{}, nil
	}

	sprints := make([]*domain.Sprint, 0, len(project.Sprints))
	for _, sprint := range project.Sprints {
		sprints = append(sprints, sprint)
	}

	// Сортировка спринтов: сначала активные, затем по имени
	sort.Slice(sprints, func(i, j int) bool {
		if sprints[i].IsActive != sprints[j].IsActive {
			return sprints[i].IsActive
		}
		return sprints[i].Name < sprints[j].Name
	})

	return sprints, nil
}

// SetActiveSprint - установка активного спринта для проекта
func (s *ProjectService) SetActiveSprint(data map[string]*domain.Project, projectName, sprintID string) error {
	project, exists := data[projectName]
	if !exists {
		return fmt.Errorf("проект '%s' не существует", projectName)
	}

	if project.Sprints == nil || len(project.Sprints) == 0 {
		return fmt.Errorf("у проекта '%s' нет спринтов", projectName)
	}

	if _, exists := project.Sprints[sprintID]; !exists {
		return fmt.Errorf("спринт с ID '%s' не существует в проекте '%s'", sprintID, projectName)
	}

	// Сначала деактивируем все спринты
	for _, sprint := range project.Sprints {
		sprint.IsActive = false
	}

	// Активируем выбранный спринт
	project.Sprints[sprintID].IsActive = true
	project.ActiveSprint = sprintID

	return s.SaveData(data)
}

// ArchiveProject - архивирование проекта
func (s *ProjectService) ArchiveProject(data map[string]*domain.Project, name string) error {
	s.Logger.Infof("Архивирование проекта: %s", name)

	project, exists := data[name]
	if !exists {
		s.Logger.Warnf("Попытка архивировать несуществующий проект: %s", name)
		return fmt.Errorf("проект '%s' не существует", name)
	}

	// Проверяем, не запущено ли отслеживание для проекта
	if project.StartTime != nil {
		s.Logger.Warnf("Попытка архивировать проект с запущенным отслеживанием: %s", name)
		return fmt.Errorf("невозможно архивировать проект с запущенным отслеживанием")
	}

	project.Archived = true

	return s.SaveData(data)
}

// RestoreProject - восстановление проекта из архива
func (s *ProjectService) RestoreProject(data map[string]*domain.Project, name string) error {
	s.Logger.Infof("Восстановление проекта из архива: %s", name)

	project, exists := data[name]
	if !exists {
		s.Logger.Warnf("Попытка восстановить несуществующий проект: %s", name)
		return fmt.Errorf("проект '%s' не существует", name)
	}

	if !project.Archived {
		s.Logger.Warnf("Попытка восстановить неархивированный проект: %s", name)
		return fmt.Errorf("проект '%s' не находится в архиве", name)
	}

	project.Archived = false

	return s.SaveData(data)
}
