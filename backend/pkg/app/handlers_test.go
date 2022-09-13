package app

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/repository"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const v1TranslationApi = "/v1/api/translations"

func TestServer_CreateTranslation(t *testing.T) {
	s := initTestServer()
	transcription := "[CreateTranscription]"
	tr := "CreateTranslation"
	text := "CreateText"
	example := "CreateExample"

	request := translation.Request{
		Transcription: transcription,
		Translation:   tr,
		Text:          text,
		Example:       example,
	}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationApi, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	records := getExistingRecords(s)
	created := records[0]

	assert.Equal(t, tr, created.Translation)
	assert.Equal(t, text, created.Text)
	assert.Equal(t, transcription, created.Transcription)
	assert.Equal(t, example, created.Example)
}

func TestServer_DeleteTranslationById(t *testing.T) {
	s := initTestServer()

	request := translation.Request{}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationApi, bytes.NewBuffer(jsonValue))
	s.router.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingRecords(s)[0].Id
	req, _ = http.NewRequest("DELETE", v1TranslationApi+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Zero(t, len(getExistingRecords(s)))
}

func TestServer_UpdateTranslation(t *testing.T) {
	s := initTestServer()
	request := translation.Request{}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationApi, bytes.NewBuffer(jsonValue))
	s.router.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingRecords(s)[0].Id

	transcription := "[updateTranscription]"
	tr := "UpdateTranslation"
	text := "UpdateText"
	example := "UpdateExample"

	request = translation.Request{
		Transcription: transcription,
		Translation:   tr,
		Text:          text,
		Example:       example,
	}
	jsonValue, _ = json.Marshal(request)
	req, _ = http.NewRequest("PUT", v1TranslationApi+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", v1TranslationApi+"/"+id, nil)
	w = httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var record translation.Translation
	json.Unmarshal(w.Body.Bytes(), &record)

	assert.Equal(t, tr, record.Translation)
	assert.Equal(t, text, record.Text)
	assert.Equal(t, transcription, record.Transcription)
	assert.Equal(t, example, record.Example)
}

func getExistingRecords(s *Server) []translation.Translation {
	req, _ := http.NewRequest("GET", v1TranslationApi, nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var records []translation.Translation
	json.Unmarshal(w.Body.Bytes(), &records)
	return records
}

func initTestServer() *Server {
	router := gin.Default()
	translationService := translation.NewTranslationService(repository.NewTranslationRepository())
	s := NewServer(router, *translationService)
	s.BuildRoutes()
	return s
}
