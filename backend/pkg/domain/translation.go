package domain

import (
	"github.com/google/uuid"
	"time"
)

type Translation struct {
	Id            string    `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Transcription string    `json:"transcription"`
	Translation   string    `json:"translation"`
	Text          string    `json:"text"`
	Example       string    `json:"example"`
}

type TranslationRequest struct {
	Transcription string `json:"transcription"`
	Translation   string `json:"translation"`
	Text          string `json:"text"`
	Example       string `json:"example"`
}

func NewTranslation(request TranslationRequest) *Translation {
	return &Translation{
		Id:            uuid.New().String(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Translation:   request.Translation,
		Transcription: request.Transcription,
		Text:          request.Text,
		Example:       request.Example,
	}
}

func (t *Translation) ApplyChanges(request TranslationRequest) {
	t.UpdatedAt = time.Now()
	t.Transcription = request.Transcription
	t.Text = request.Text
	t.Translation = request.Translation
	t.Example = request.Example
}
