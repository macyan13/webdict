package server

import "time"

type translationRequest struct {
	Transcription string   `json:"transcription"`
	Translation   string   `json:"translation"`
	Text          string   `json:"text"`
	Example       string   `json:"example"`
	TagIds        []string `json:"tag_ids"`
}

type tagRequest struct {
	Tag string `json:"tag"`
}

type translationResponse struct {
	ID            string        `json:"id"`
	CreatedAt     time.Time     `json:"created_at"`
	Transcription string        `json:"transcription"`
	Translation   string        `json:"translation"`
	Text          string        `json:"text"`
	Example       string        `json:"example"`
	Tags          []tagResponse `json:"tags"`
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
