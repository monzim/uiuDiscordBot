package models

import (
	"gorm.io/gorm"
)

type CommandLog struct {
	gorm.Model
	UserID       string `json:"user_id"`
	Command      string `json:"command"`
	Parameters   string `json:"parameters"`
	ResponseTime string `json:"response_time"`
}

type UserActivity struct {
	UserID           string `json:"user_id" gorm:"primaryKey"`
	CommandsExecuted int    `json:"commands_executed"`
	LastActivity     string `json:"last_activity"`
}

type ErrorLog struct {
	gorm.Model
	ErrorMessage string `json:"error_message"`
}

type ServerStats struct {
	ServerID      string `json:"server_id" gorm:"primaryKey"`
	MembersCount  int    `json:"members_count"`
	ChannelsCount int    `json:"channels_count"`
	CreatedAt     string `json:"created_at"`
}

type EventLog struct {
	gorm.Model
	EventType        string `json:"event_type"`
	EventDescription string `json:"event_description"`
}

type UserDetails struct {
	UserID    string `json:"user_id" gorm:"primaryKey"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	JoinedAt  string `json:"joined_at"`
}
