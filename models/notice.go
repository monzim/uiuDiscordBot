package models

import (
	"time"

	uiuscraper "github.com/monzim/uiu-notice-scraper"
)

type Department uiuscraper.Department

type Notice struct {
	ID         string     `gorm:"primaryKey" json:"id"`
	Title      string     `gorm:"index" json:"title"`
	Summary    string     `json:"summary"`
	Image      string     `json:"image"`
	Date       time.Time  `json:"date"`
	Link       string     `json:"link"`
	Notified   bool       `json:"notified"`
	Department Department `json:"department"`
	TimeCommon
}
