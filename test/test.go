package main

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	db "github.com/monzim/uiuBot/database"
	"github.com/monzim/uiuBot/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func trimRoomSpaces(s string) string {
	return strings.ReplaceAll(s, "  ", "")
}

func generateID(courseCode, section, room string) string {
	hashInput := fmt.Sprintf("%s-%s-%s", courseCode, section, room)
	hash := sha256.Sum256([]byte(hashInput))
	return fmt.Sprintf("%x", hash)
}

func csvToJSON(csvFilePath, jsonFilePath string) ([]models.Exam, error) {
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var exams []models.Exam

	for _, row := range csvData[1:] {
		exam := models.Exam{
			ID:          generateID(row[1], row[3], trimRoomSpaces(row[7])),
			Department:  strings.ToLower(strings.TrimSpace(row[0])),
			CourseCode:  strings.ToLower(strings.TrimSpace(row[1])),
			CourseTitle: strings.ToLower(strings.TrimSpace(row[2])),
			Section:     strings.ToLower(strings.TrimSpace(row[3])),
			Teacher:     row[4],
			ExamDate:    row[5],
			ExamTime:    row[6],
			Room:        trimRoomSpaces(row[7]),
		}
		exams = append(exams, exam)
	}

	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonEncoder := json.NewEncoder(jsonFile)
	jsonEncoder.SetIndent("", "    ") // Optional: Set indentation for better readability
	if err := jsonEncoder.Encode(exams); err != nil {
		return nil, err
	}

	fmt.Println("Conversion successful!")
	return exams, nil
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}

	postgres, err := db.NewDatabaseConnection(&db.DatabaseConfig{
		// Host:     os.Getenv("DB_HOST"),
		// Port:     os.Getenv("DB_PORT"),
		// DBname:   os.Getenv("DB_NAME"),
		// User:     os.Getenv("DB_USER"),
		// Password: os.Getenv("DB_PASSWORD"),
		// SSlMode:  os.Getenv("DB_SSL_MODE"),
		Host:     os.Getenv("LOG_DB_HOST"),
		Port:     os.Getenv("LOG_DB_PORT"),
		DBname:   os.Getenv("LOG_DB_NAME"),
		User:     os.Getenv("LOG_DB_USER"),
		Password: os.Getenv("LOG_DB_PASSWORD"),
		SSlMode:  os.Getenv("LOG_DB_SSL_MODE"),
	})
	if err != nil {
		log.Error().Err(err).Msg("Error initializing the database connection")

	}

	err = postgres.AutoMigrate(&models.Exam{})
	if err != nil {
		log.Error().Err(err).Msg("Error migrating the database")
	}

	// get all the unique department names
	var departments []string
	postgres.Table("exams").Distinct("department").Pluck("department", &departments)

	log.Info().Msgf("Unique departments: %v", departments)
	log.Info().Msgf("Total departments: %v", len(departments))

	// var courses []string
	// postgres.Table("exams").Distinct("course_code").Pluck("course_code", &courses)

	// log.Info().Msgf("Unique course codes: %v", courses)

	// var exams []models.Exam

	// var department = "Bba"
	// courseCode := "4204 "
	// section := "A"

	// // res := postgres.Where("department = ? AND section = ?",
	// // 	strings.ToLower(department),
	// // 	strings.ToLower(section),
	// // ).Where("course_code LIKE ?", "%"+courseCode+"%").Find(&exams)

	// res := postgres.Where(models.Exam{
	// 	Department: strings.ToLower(department),
	// 	Section:    strings.ToLower(section),
	// }).Where(
	// 	"course_code LIKE ?", "%"+strings.ToLower(strings.TrimSpace(courseCode))+"%",
	// ).Find(&exams)

	// if res.Error != nil {
	// 	log.Error().Err(res.Error).Msg("Error fetching exams")
	// }

	// log.Info().Msgf("Total exams: %v", len(exams))
	// return

	csvFils := []string{
		"data/BSCSE-BSDS-BSEEE-BBA-BBA_in_AIS-&-BSECO-BSCE-BSSEDS-MSCSE-BA-in-English-MBA:EMBA.csv",
	}

	for _, csvFile := range csvFils {
		jsonFilePath := strings.Replace(csvFile, ".csv", ".json", 1)

		json, err := csvToJSON(csvFile, jsonFilePath)
		if err != nil {
			fmt.Println(err)
		}

		log.Info().Msgf("Data from %s converted to %s", csvFile, jsonFilePath)

		for i, exam := range json {
			var e models.Exam
			result := postgres.Where("id = ?", exam.ID).First(&e)
			if result.Error == nil {
				log.Info().Msgf("%d. Exam already exists", i+1)
				continue
			}

			postgres.Create(&exam)
			log.Info().Msgf("%d. Data inserted", i+1)
		}
	}

	log.Info().Msg("Data inserted successfully!")

}
