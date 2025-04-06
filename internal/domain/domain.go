package domain

import "time"

// TimeEntry - запись о затраченном времени
type TimeEntry struct {
	Date        string `json:"date"`
	TimeSpent   int    `json:"time_spent"`
	Description string `json:"description"`
}

// Sprint - этап проекта
type Sprint struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	StartDate   string               `json:"start_date,omitempty"`
	EndDate     string               `json:"end_date,omitempty"`
	Entries     map[string]TimeEntry `json:"entries,omitempty"`
	IsActive    bool                 `json:"is_active"`
}

// Project - проект
type Project struct {
	StartTime    *time.Time         `json:"start_time,omitempty"`
	Entries      []TimeEntry        `json:"entries,omitempty"`
	Sprints      map[string]*Sprint `json:"sprints,omitempty"`
	ActiveSprint string             `json:"active_sprint,omitempty"`
	Archived     bool               `json:"archived,omitempty"`
}

// Entry - структура записи времени
type Entry struct {
	TimeSpent   int    `json:"time_spent"`
	Description string `json:"description"`
	Date        string `json:"date"`
}
