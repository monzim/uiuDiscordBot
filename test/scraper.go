package main

import (
	"fmt"

	uiuscraper "github.com/monzim/uiu-notice-scraper"
	"github.com/monzim/uiuBot/models"
	"gorm.io/gorm"
)

func ScrapData(db *gorm.DB) {
	var latestNotice models.Notice
	if err := db.Order("date desc").First(&latestNotice).Error; err != nil {
		fmt.Println("No latest notice found in the database")
	}

	config := uiuscraper.NoticeScrapConfig{
		LastNoticeId: &latestNotice.ID,
		Department:   uiuscraper.DepartmentAll,
		AllowDomain:  string(uiuscraper.DepartmentAll),
		NOTICE_SITE:  string(uiuscraper.Notice_Site_UIU),
	}

	notices := uiuscraper.ScrapNotice(&config)
	for _, notice := range notices {
		var n models.Notice = models.Notice{
			ID:         notice.ID,
			Title:      notice.Title,
			Image:      notice.Image,
			Date:       notice.Date,
			Link:       notice.Link,
			Department: models.Department(notice.Department),
		}

		if err := db.FirstOrCreate(&n).Error; err != nil {
			fmt.Println("Error creating notice:", err)
			continue
		}
	}

	fmt.Println("Total notices created:", len(notices))

}
