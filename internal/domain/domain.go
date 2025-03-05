package domain

import "time"

type Project struct {
	Entries   []Entry    `json:"entries"`
	StartTime *time.Time `json:"start_time,omitempty"`
}

type Entry struct {
	TimeSpent   int    `json:"time_spent"`
	Description string `json:"description"`
	Date        string `json:"date"`
}
