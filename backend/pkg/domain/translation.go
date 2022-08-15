package domain

import (
	"github.com/google/uuid"
	"time"
)

type Translation struct {
	Id            string `json:"id"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
	Transcription string `json:"transcription"`
	Translation   string `json:"translation"`
	Text          string `json:"text"`
	Example       string `json:"example"`
}

type TranslationRequest struct {
	Transcription string `json:"transcription"`
	Translation   string `json:"translation"`
	Text          string `json:"text"`
	Example       string `json:"example"`
}

func NewTranslation(request TranslationRequest) *Translation {
	now := time.Now().Unix()
	return &Translation{
		Id:            uuid.New().String(),
		CreatedAt:     now,
		UpdatedAt:     now,
		Translation:   request.Translation,
		Transcription: request.Transcription,
		Text:          request.Text,
		Example:       request.Example,
	}
}

func (t *Translation) ApplyChanges(request TranslationRequest) {
	t.UpdatedAt = time.Now().Unix()
	t.Transcription = request.Transcription
	t.Text = request.Text
	t.Translation = request.Translation
	t.Example = request.Example
}
