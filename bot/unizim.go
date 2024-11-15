package bot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

type NoticeBoyd struct {
	HashID     string `json:"hash_id"`
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	Image      string `json:"image"`
	Date       string `json:"date"`
	Department string `json:"department"`
}

type NoticePayload struct {
	DocumentID string     `json:"documentId"`
	Data       NoticeBoyd `json:"data"`
}

func NoticePushToUnizim(payload *NoticePayload) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", os.Getenv("UNIZIM_API_URL"), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("X-Appwrite-Project", os.Getenv("APPWRITE_PROJECT_ID"))
	req.Header.Set("X-Appwrite-Key", os.Getenv("APPWRITE_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Info().Msgf("Status code: %d", resp.StatusCode)
	log.Info().Msgf("Notice sent to Unizim app server: %s", payload.Data.HashID)
	return nil
}
