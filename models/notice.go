package models

import "time"

type Notice struct {
	ID       string    `gorm:"primaryKey" json:"id"`
	Title    string    `gorm:"index" json:"title"`
	Image    string    `json:"image"`
	Date     time.Time `json:"date"`
	Link     string    `json:"link"`
	Notified bool      `json:"notified"`
	TimeCommon
}
