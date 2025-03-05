package service

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/MWT-proger/time-tracking/internal/domain"
)

type ProjectService struct {
	dataFile string
}

func NewProjectService() *ProjectService {
	return &ProjectService{
		dataFile: filepath.Join(os.Getenv("HOME"), "учет_времени.json"),
	}
}

func (s *ProjectService) CreateProject(name string) {
	data, _ := s.loadData()
	if _, exists := data[name]; exists {
		return
	}
	data[name] = &domain.Project{}
	s.saveData(data)
}

func (s *ProjectService) GetProjects() map[string]*domain.Project {
	data, _ := s.loadData()
	return data
}

func (s *ProjectService) loadData() (map[string]*domain.Project, error) {
	data := make(map[string]*domain.Project)
	file, err := os.Open(s.dataFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func (s *ProjectService) saveData(data map[string]*domain.Project) error {
	file, err := os.Create(s.dataFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(data)
}
