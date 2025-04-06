package domain

import "time"

// Project - структура данных о проекте
type Project struct {
	Entries   []Entry    `json:"entries"`
	StartTime *time.Time `json:"start_time,omitempty"`
}

// Entry - структура записи времени
type Entry struct {
	TimeSpent   int    `json:"time_spent"`
	Description string `json:"description"`
	Date        string `json:"date"`
}
