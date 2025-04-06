package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MWT-proger/time-tracking/internal/domain"
	"github.com/MWT-proger/time-tracking/pkg/logger"
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
func (s *ProjectService) GetProjectNames(data map[string]*domain.Project) []string {
	names := []string{}
	for name := range data {
		names = append(names, name)
	}
	return names
}
