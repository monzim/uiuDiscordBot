package models

type MessageLog struct {
	ID        string `json:"id" gorm:"primaryKey"`
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	Message   string `json:"message"`
	Data      []byte `gorm:"type:jsonb" json:"data"`
}
