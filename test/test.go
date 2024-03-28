package main

import (
	"os"

	"github.com/joho/godotenv"
	uiuscraper "github.com/monzim/uiu-notice-scraper"
	db "github.com/monzim/uiuBot/database"
	"github.com/monzim/uiuBot/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// func trimRoomSpaces(s string) string {
// 	return strings.ReplaceAll(s, "  ", "")
// }

// func generateID(courseCode, section, room string) string {
// 	hashInput := fmt.Sprintf("%s-%s-%s", courseCode, section, room)
// 	hash := sha256.Sum256([]byte(hashInput))
// 	return fmt.Sprintf("%x", hash)
// }

// func csvToJSON(csvFilePath, jsonFilePath string) ([]models.Exam, error) {
// 	csvFile, err := os.Open(csvFilePath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer csvFile.Close()

// 	csvReader := csv.NewReader(csvFile)
// 	csvData, err := csvReader.ReadAll()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var exams []models.Exam

// 	for _, row := range csvData[1:] {
// 		exam := models.Exam{
// 			ID:          generateID(row[1], row[3], trimRoomSpaces(row[7])),
// 			Department:  strings.ToLower(strings.TrimSpace(row[0])),
// 			CourseCode:  strings.ToLower(strings.TrimSpace(row[1])),
// 			CourseTitle: strings.ToLower(strings.TrimSpace(row[2])),
// 			Section:     strings.ToLower(strings.TrimSpace(row[3])),
// 			Teacher:     row[4],
// 			ExamDate:    row[5],
// 			ExamTime:    row[6],
// 			Room:        trimRoomSpaces(row[7]),
// 		}
// 		exams = append(exams, exam)
// 	}

// 	jsonFile, err := os.Create(jsonFilePath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer jsonFile.Close()

// 	jsonEncoder := json.NewEncoder(jsonFile)
// 	jsonEncoder.SetIndent("", "    ") // Optional: Set indentation for better readability
// 	if err := jsonEncoder.Encode(exams); err != nil {
// 		return nil, err
// 	}

// 	fmt.Println("Conversion successful!")
// 	return exams, nil
// }

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}

	postgres, err := db.NewDatabaseConnection(&db.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBname:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSlMode:  os.Getenv("DB_SSL_MODE"),
		// Host:     os.Getenv("LOG_DB_HOST"),
		// Port:     os.Getenv("LOG_DB_PORT"),
		// DBname:   os.Getenv("LOG_DB_NAME"),
		// User:     os.Getenv("LOG_DB_USER"),
		// Password: os.Getenv("LOG_DB_PASSWORD"),
		// SSlMode:  os.Getenv("LOG_DB_SSL_MODE"),
	})
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

		// 	config := setupConfig(tc.department, tc.allowDomain, tc.noticeSite)

		// 	notices := uiuscraper.ScrapNotice(config)
		// 	for _, notice := range notices {

		// 		var n models.Notice
		// 		res := postgres.Where("id = ?", notice.ID).First(&n)
		// 		if res.Error != nil {
		// 			log.Error().Err(res.Error).Msg("Error fetching the notice")
		// 			log.Info().Msgf("---->> Notice: %s", notice)

		// 			// create the notice
		// 			noticeModel := models.Notice{
		// 				ID:         notice.ID,
		// 				Title:      notice.Title,
		// 				Department: models.Department(notice.Department),
		// 				Notified:   true,
		// 				Summary:    notice.Summary,
		// 				Image:      notice.Image,
		// 				Date:       notice.Date,
		// 				Link:       notice.Link,
		// 			}

		// 			res := postgres.Create(&noticeModel)
		// 			if res.Error != nil {
		// 				log.Error().Err(res.Error).Msg("Error creating the notice")
		// 			}

		// 		}

		// 		if n.ID == notice.ID {
		// 			n.Summary = notice.Summary
		// 			res := postgres.Save(&n)
		// 			if res.Error != nil {
		// 				log.Error().Err(res.Error).Msg("Error updating the notice")
		// 			}

		// 			continue
		// 		}

		// 	}
	}

}

// func setupConfig(department uiuscraper.Department, allowDomain string, noticeSite string) *uiuscraper.NoticeScrapConfig {
// 	config := uiuscraper.NoticeScrapConfig{
// 		LastNoticeId: nil,
// 		Department:   department,
// 		AllowDomain:  allowDomain,
// 		NOTICE_SITE:  noticeSite,
// 	}
// 	return &config
// }
