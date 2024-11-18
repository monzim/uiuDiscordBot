package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type Notice struct {
	HashID    string    `json:"hash_id"`
	ID        string    `json:"$id"`
	CreatedAt time.Time `json:"$createdAt"`
	UpdatedAt time.Time `json:"$updatedAt"`
}

type NoticeResponse struct {
	Total     int      `json:"total"`
	Documents []Notice `json:"documents"`
}

type NoticeBody struct {
	HashID     string `json:"hash_id"`
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	Image      string `json:"image"`
	Date       string `json:"date"`
	Link       string `json:"link"`
	Department string `json:"department"`
}

type NoticePayload struct {
	DocumentID string     `json:"documentId"`
	Data       NoticeBody `json:"data"`
}

func NoticePushToUnizim(payload *NoticePayload, duplicationCheck bool) error {
	if duplicationCheck {
		notice, err := FetchNoticeByHashID(payload.Data.HashID)
		if err == nil {
			log.Info().Msgf("Notice already exists in the database: %s", notice.ID)
			return nil
		}
	}

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

	log.Info().Msgf("Notice sent to Unizim app server: %s", payload.Data.HashID)
	return nil
}

func FetchNoticeByHashID(hashID string) (*Notice, error) {
	params := url.Values{}
	queries := []string{
		fmt.Sprintf(`{"method":"equal","attribute":"hash_id","values":["%s"]}`, hashID),
		`{"method":"limit","values":[1]}`,
	}

	for i, q := range queries {
		params.Add(fmt.Sprintf("queries[%d]", i), q)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", os.Getenv("UNIZIM_API_URL"), params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Appwrite-Project", os.Getenv("APPWRITE_PROJECT_ID"))
	req.Header.Set("X-Appwrite-Key", os.Getenv("APPWRITE_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var noticeResp NoticeResponse
	if err := json.NewDecoder(resp.Body).Decode(&noticeResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if len(noticeResp.Documents) < 1 {
		return nil, fmt.Errorf("no notice found with hash ID: %s", hashID)
	}

	return &noticeResp.Documents[0], nil
}
