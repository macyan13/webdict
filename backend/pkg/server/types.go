package server

import "time"

type translationRequest struct {
	Source        string   `json:"source"`
	Transcription string   `json:"transcription"`
	Target        string   `json:"target"`
	Example       string   `json:"example"`
	TagIds        []string `json:"tag_ids"`
	LangID        string   `json:"lang_id"`
}

type tagRequest struct {
	Tag string `json:"tag"`
}

type userRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type updateUserRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	Role            int    `json:"role"`
	DefaultLangID   string `json:"default_lang_id"`
}

type updateProfileRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	DefaultLangID   string `json:"default_lang_id"`
}

type langRequest struct {
	Name string `json:"name"`
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type translationResponse struct {
	ID            string        `json:"id"`
	Source        string        `json:"source"`
	Transcription string        `json:"transcription"`
	Target        string        `json:"target"`
	Example       string        `json:"example"`
	Tags          []tagResponse `json:"tags"`
	CreatedAt     time.Time     `json:"created_at"`
	Lang          langResponse  `json:"lang"`
}

type lastTranslationsResponse struct {
	Translations []translationResponse `json:"translations"`
	TotalRecords int                   `json:"total_records"`
}

type tagResponse struct {
	ID  string `json:"id"`
	Tag string `json:"tag"`
}

type langResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type userResponse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	Role        int          `json:"role"`
	DefaultLang langResponse `json:"default_lang"`
}

type idResponse struct {
	ID string `json:"id"`
}

type AuthTokenResponse struct {
	AccessToken string `json:"accessToken"`
	Type        string `json:"type"`
}
