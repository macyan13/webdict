package query

import "time"

type Translation struct {
	Id            string
	CreatedAd     time.Time
	Transcription string
	Translation   string
	Text          string
	Example       string
	Tags          []Tag
}

type Tag struct {
	Id        string
	Tag       string
	CreatedAd time.Time
}
