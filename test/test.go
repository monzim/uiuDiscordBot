package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	uiuscraper "github.com/monzim/uiu-notice-scraper"
	db "github.com/monzim/uiuBot/database"
	"github.com/monzim/uiuBot/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func getSemester(date time.Time) string {
	year := date.Year()
	springStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	springMid := time.Date(year, time.March, 15, 12, 0, 0, 0, time.UTC)
	springEnd := time.Date(year, time.May, 31, 23, 59, 59, 999999999, time.UTC)
	summerStart := time.Date(year, time.June, 1, 0, 0, 0, 0, time.UTC)
	summerMid := time.Date(year, time.August, 15, 12, 0, 0, 0, time.UTC)
	summerEnd := time.Date(year, time.September, 30, 23, 59, 59, 999999999, time.UTC)
	fallStart := time.Date(year, time.October, 1, 0, 0, 0, 0, time.UTC)
	fallMid := time.Date(year, time.November, 15, 12, 0, 0, 0, time.UTC)
	fallEnd := time.Date(year, time.December, 31, 23, 59, 59, 999999999, time.UTC)

	if date.After(springStart) && date.Before(springMid) {
		return fmt.Sprintf("%d_SPRING_MID", year)
	} else if date.After(springMid) && date.Before(springEnd) {
		return fmt.Sprintf("%d_SPRING_FINAL", year)
	} else if date.After(summerStart) && date.Before(summerMid) {
		return fmt.Sprintf("%d_SUMMER_MID", year)
	} else if date.After(summerMid) && date.Before(summerEnd) {
		return fmt.Sprintf("%d_SUMMER_FINAL", year)
	} else if date.After(fallStart) && date.Before(fallMid) {
		return fmt.Sprintf("%d_FALL_MID", year)
	} else if date.After(fallMid) && date.Before(fallEnd) {
		return fmt.Sprintf("%d_FALL_FINAL", year)
	} else {
		return "Invalid Date"
	}
}

func TestAllMonths(year int) {
	// Iterate over all months
	for month := time.January; month <= time.December; month++ {
		// Get the first day of the month
		firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		// Test the first day of the month
		semester := getSemester(firstDayOfMonth)
		fmt.Printf("%s: %s\n", firstDayOfMonth.Format("2006-01-02"), semester)

		// Get a mid-day of the month for additional testing
		midDayOfMonth := time.Date(year, month, 15, 12, 0, 0, 0, time.UTC)
		// Test the mid-day of the month
		semester = getSemester(midDayOfMonth)
		fmt.Printf("%s: %s\n", midDayOfMonth.Format("2006-01-02"), semester)
	}
}

func Oldmain() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}

	postgres, err := db.NewDatabaseConnection(
		// os.Getenv("LOG_DATABASE_URI"),
		os.Getenv("DATABASE_URI"),
	)
	if err != nil {
		log.Error().Err(err).Msg("Error initializing the database connection")

	}

	err = postgres.AutoMigrate(&models.Exam{}, &models.Notice{})
	if err != nil {
		log.Error().Err(err).Msg("Error migrating the database")
	}

	var notices []models.Notice

	res := postgres.Where("department = ?", "").Select("id, department, date").Find(&notices)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("Error fetching notices")
	}

	for i, notice := range notices {
		notice.Department = models.Department(uiuscraper.DepartmentAll)
		res := postgres.Save(&notice)
		if res.Error != nil {
			log.Error().Err(res.Error).Msgf("%d Error updating the notice", i)
		}

	}

	// get the count of notices per department
	cases := []struct {
		testName    string
		department  uiuscraper.Department
		allowDomain string
		noticeSite  string
	}{
		{"UIU Notices", uiuscraper.DepartmentAll, uiuscraper.AllowDomainUIU, uiuscraper.Notice_Site_UIU},
		{"CSE Notices", uiuscraper.DepartmentCSE, uiuscraper.AllowDomainCSE, uiuscraper.Notice_Site_CSE},
		{"EEE Notices", uiuscraper.DepartmentEEE, uiuscraper.AllowDomainEEE, uiuscraper.Notice_Site_EEE},
		{"CE Notices", uiuscraper.DepartmentCivil, uiuscraper.AllowDomainCE, uiuscraper.Notice_Site_CE},
		{"Pharmacy Notices", uiuscraper.DepartmentPharmacy, uiuscraper.AllowDomainPharmacy, uiuscraper.Notice_Site_Pharmacy},
	}

	for _, tc := range cases {
		var count int64
		postgres.Model(&models.Notice{}).Where("department = ?", tc.department).Count(&count)
		log.Info().Msgf("Total notices for %s: %d", tc.testName, count)

		config := setupConfig(tc.department, tc.allowDomain, tc.noticeSite)

		notices := uiuscraper.ScrapNotice(config)
		for _, notice := range notices {

			var n models.Notice
			res := postgres.Where("id = ?", notice.ID).First(&n)
			if res.Error != nil {
				log.Error().Err(res.Error).Msg("Error fetching the notice")
				log.Info().Msgf("---->> Notice: %s", notice)

				// create the notice
				noticeModel := models.Notice{
					ID:         notice.ID,
					Title:      notice.Title,
					Department: models.Department(notice.Department),
					Notified:   true,
					Summary:    notice.Summary,
					Image:      notice.Image,
					Date:       notice.Date,
					Link:       notice.Link,
				}

				res := postgres.Create(&noticeModel)
				if res.Error != nil {
					log.Error().Err(res.Error).Msg("Error creating the notice")
				}

			}

			if n.ID == notice.ID {
				n.Summary = notice.Summary
				res := postgres.Save(&n)
				if res.Error != nil {
					log.Error().Err(res.Error).Msg("Error updating the notice")
				}

				continue
			}

		}
	}

}

func setupConfig(department uiuscraper.Department, allowDomain string, noticeSite string) *uiuscraper.NoticeScrapConfig {
	config := uiuscraper.NoticeScrapConfig{
		LastNoticeId: nil,
		Department:   department,
		AllowDomain:  allowDomain,
		NOTICE_SITE:  noticeSite,
	}
	return &config
}
