package server

import "time"

type translationRequest struct {
	Text          string   `json:"text"`
	Transcription string   `json:"transcription"`
	Meaning       string   `json:"meaning"`
	Example       string   `json:"example"`
	TagIds        []string `json:"tag_ids"`
}

type tagRequest struct {
	Tag string `json:"tag"`
}

type translationResponse struct {
	ID            string        `json:"id"`
	Text          string        `json:"text"`
	Transcription string        `json:"transcription"`
	Meaning       string        `json:"meaning"`
	Example       string        `json:"example"`
	Tags          []tagResponse `json:"tags"`
	CreatedAt     time.Time     `json:"created_at"`
}

type tagResponse struct {
	ID  string `json:"id"`
	Tag string `json:"tag"`
}

type idResponse struct {
	ID string `json:"id"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthTokenResponse struct {
	AccessToken string `json:"accessToken"`
	Type        string `json:"type"`
}
