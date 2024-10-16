package main

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	db "github.com/monzim/uiuBot/database"
	"github.com/monzim/uiuBot/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// SUMMER_24_MID = "24_SUMMER_MID"
	FINAL_SUMMER_24 = "24_SUMMER_FINAL"
)

func AddExamsToDatabase() {
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

	// convert the exam time csv to json
	exams, err := CsvToJSON(FINAL_SUMMER_24, fmt.Sprintf("%s.csv", FINAL_SUMMER_24), fmt.Sprintf("%s.json", FINAL_SUMMER_24))
	if err != nil {
		log.Error().Err(err).Msg("Error converting the csv to json")
	}

	wg := sync.WaitGroup{}
	for i := 0; i < len(exams); i += 100 {
		end := i + 100
		if end > len(exams) {
			end = len(exams)
		}

		wg.Add(1)
		go func(exams []models.Exam) {
			for _, exam := range exams {
				res := postgres.Create(&exam)
				if res.Error != nil {
					log.Error().Err(res.Error).Msg("Error creating the exam")
				}

				log.Info().Msgf("Exam: %s inserted", exam.ID)
			}

			wg.Done()
		}(exams[i:end])
	}

	wg.Wait()

}

func trimRoomSpaces(s string) string {
	return strings.ReplaceAll(s, "  ", "")
}

func generateID(semisterCode, courseCode, section, room string) string {
	hashInput := fmt.Sprintf("%s-%s-%s-%s", semisterCode, courseCode, section, room)
	hash := sha256.Sum256([]byte(hashInput))
	return fmt.Sprintf("%x", hash)
}

func CsvToJSON(semisterCode, csvFilePath, jsonFilePath string) ([]models.Exam, error) {
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
			ID:          generateID(semisterCode, row[1], row[3], trimRoomSpaces(row[7])),
			TrimsterID:  semisterCode,
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
	jsonEncoder.SetIndent("", "    ")
	if err := jsonEncoder.Encode(exams); err != nil {
		return nil, err
	}

	fmt.Println("Conversion successful!")
	return exams, nil
}
