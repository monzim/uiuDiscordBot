package models

import "time"

type UserDetails struct {
	ServerID      string     `json:"server_id" gorm:"primaryKey"`
	UserID        string     `json:"user_id" gorm:"primaryKey"`
	Username      string     `json:"username"`
	AvatarURL     string     `json:"avatar_url"`
	Department    Department `json:"department"`
	JoinedAt      time.Time  `json:"joined_at"`
	Email         string     `json:"email"`
	Avatar        string     `json:"avatar"`
	Locale        string     `json:"locale"`
	Discriminator string     `json:"discriminator"`
	Token         string     `json:"token"`
	Verified      bool       `json:"verified"`
	MFAEnabled    bool       `json:"mfa_enabled"`
	Banner        string     `json:"banner"`
	AccentColor   int        `json:"accent_color"`
	Bot           bool       `json:"bot"`
	PublicFlags   int        `json:"public_flags"`
	PremiumType   int        `json:"premium_type"`
	System        bool       `json:"system"`
	Flags         int        `json:"flags"`
}
