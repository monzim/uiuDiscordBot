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
}

type UserActivity struct {
	UserID           string `json:"user_id" gorm:"primaryKey"`
	CommandsExecuted int    `json:"commands_executed"`
	LastActivity     string `json:"last_activity"`
}

type UserDetails struct {
	UserID    string    `json:"user_id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	JoinedAt  time.Time `json:"joined_at"`

	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	Locale        string `json:"locale"`
	Discriminator string `json:"discriminator"`
	Token         string `json:"token"`
	Verified      bool   `json:"verified"`
	MFAEnabled    bool   `json:"mfa_enabled"`
	Banner        string `json:"banner"`
	AccentColor   int    `json:"accent_color"`
	Bot           bool   `json:"bot"`
	PublicFlags   int    `json:"public_flags"`
	PremiumType   int    `json:"premium_type"`
	System        bool   `json:"system"`
	Flags         int    `json:"flags"`
}

type EventLog struct {
	gorm.Model
	EventType        string `json:"event_type"`
	EventDescription string `json:"event_description"`
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
}
