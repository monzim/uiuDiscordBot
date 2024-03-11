package models

import (
	"time"

	"gorm.io/gorm"
)

type CommandLog struct {
	gorm.Model
	UserID       string `json:"user_id"`
	Command      string `json:"command"`
	Parameters   string `json:"parameters"`
	ResponseTime string `json:"response_time"`
	ServerID     string `json:"server_id"`
}

type UserActivity struct {
	UserID           string    `json:"user_id" gorm:"primaryKey"`
	CommandsExecuted int       `json:"commands_executed"`
	LastActivity     string    `json:"last_activity"`
	TimeStamps       time.Time `json:"time_stamps"`
	ServerID         string    `json:"server_id"`
}

type EventLog struct {
	gorm.Model
	EventType        string `json:"event_type"`
	EventDescription string `json:"event_description"`
	DM               bool   `json:"dm"`
	ServerID         string `json:"server_id"`
}

type DMLog struct {
	gorm.Model
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	UserData []byte `gorm:"type:jsonb" json:"user_data"`
	Data     []byte `gorm:"type:jsonb" json:"data"`
}

type ServerStats struct {
	ServerID      string `json:"server_id" gorm:"primaryKey"`
	MembersCount  int    `json:"members_count"`
	ChannelsCount int    `json:"channels_count"`
	CreatedAt     string `json:"created_at"`
}

type ExamTimeLog struct {
	gorm.Model
	UserID       string `json:"user_id"`
	Department   string `json:"department"`
	CourseCode   string `json:"course_code"`
	Section      string `json:"section"`
	ResponseTime string `json:"response_time"`
	ServerID     string `json:"server_id"`
}
