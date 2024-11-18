package main

import (
	"fmt"

	"github.com/monzim/uiuBot/bot"
	"github.com/monzim/uiuBot/models"
	"gorm.io/gorm"
)

func SendToUnizimTest(db *gorm.DB) {
	var notices []models.Notice
	if err := db.Find(&notices).Error; err != nil {
		fmt.Println("Error fetching notices from the database:", err)
	}

	for _, notice := range notices {
		payload := &bot.NoticePayload{
			DocumentID: "unique()",
			Data: bot.NoticeBody{
				HashID:     notice.ID,
				Title:      notice.Title,
				Summary:    notice.Summary,
				Image:      notice.Image,
				Link:       notice.Link,
				Date:       notice.Date.Format("2006-01-02"),
				Department: string(notice.Department),
			},
		}

		if err := bot.NoticePushToUnizim(payload, true); err != nil {
			fmt.Println("Error:", err)
		}
	}

}
