package models

import (
	"time"

	"gorm.io/gorm"
)

type CronMessage struct {
	gorm.Model
	Message   string    `json:"message"`
	ChannelID string    `json:"channel_id"`
	UserID    string    `json:"user_id"`
	Done      bool      `json:"done"`
	SendTime  time.Time `json:"send_time"`
}
