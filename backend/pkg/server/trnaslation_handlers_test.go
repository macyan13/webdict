package server

import (
	"bytes"
	"encoding/json"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const v1TranslationApi = "/v1/api/translations"

func TestServer_CreateTranslation(t *testing.T) {
	s := initTestServer()
	transcription := "CreateTranscription"
	tr := "CreateTranslation"
	text := "CreateText"
	example := "CreateExample"
	tg := "testTag"

	user, err := s.userRepo.GetByEmail(s.opts.Admin.AdminEmail)
	assert.Nil(t, err)

	err = s.app.Commands.AddTag.Handle(command.AddTag{
		Tag:      tg,
		AuthorId: user.Id(),
	})
	assert.Nil(t, err)

	tags, err := s.app.Queries.AllTags.Handle(query.AllTags{AuthorId: user.Id()})
	assert.Nil(t, err)
	tagId := tags[0].Id

	request := translationRequest{
		Transcription: transcription,
		Translation:   tr,
		Text:          text,
		Example:       example,
		TagIds:        []string{tagId},
	}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationApi, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	records := getExistingTranslations(s)
	created := records[0]

	assert.Equal(t, tr, created.Translation)
	assert.Equal(t, text, created.Text)
	assert.Equal(t, transcription, created.Transcription)
	assert.Equal(t, example, created.Example)
	assert.Equal(t, tagId, created.Tags[0].Id)
}

func TestServer_CreateTranslationUnauthorised(t *testing.T) {
	s := initTestServer()
	transcription := "CreateTranscription"
	tr := "CreateTranslation"
	text := "CreateText"
	example := "CreateExample"

	request := translationRequest{
		Transcription: transcription,
		Translation:   tr,
		Text:          text,
		Example:       example,
		TagIds:        []string{},
	}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationApi, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Zero(t, len(getExistingTranslations(s)))
}

func TestServer_DeleteTranslationById(t *testing.T) {
	s := initTestServer()

	request := translationRequest{}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationApi, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingTranslations(s)[0].Id
	req, _ = http.NewRequest("DELETE", v1TranslationApi+"/"+id, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Zero(t, len(getExistingTranslations(s)))
}

func TestServer_DeleteTranslationByIdUnauthorised(t *testing.T) {
	s := initTestServer()

	request := translationRequest{}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationApi, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingTranslations(s)[0].Id
	req, _ = http.NewRequest("DELETE", v1TranslationApi+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, 1, len(getExistingTranslations(s)))
}

func TestServer_UpdateTranslation(t *testing.T) {
	s := initTestServer()
	request := translationRequest{}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationApi, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingTranslations(s)[0].Id

	transcription := "[updateTranscription]"
	tr := "UpdateTranslation"
	text := "UpdateText"
	example := "UpdateExample"

	request = translationRequest{
		Transcription: transcription,
		Translation:   tr,
		Text:          text,
		Example:       example,
	}
	jsonValue, _ = json.Marshal(request)
	req, _ = http.NewRequest("PUT", v1TranslationApi+"/"+id, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", v1TranslationApi+"/"+id, nil)
	setAuthToken(s, req)
	w = httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record translationResponse
	json.Unmarshal(w.Body.Bytes(), &record)

	assert.Equal(t, tr, record.Translation)
	assert.Equal(t, text, record.Text)
	assert.Equal(t, transcription, record.Transcription)
	assert.Equal(t, example, record.Example)
}

func TestServer_UpdateTranslationUnauthorised(t *testing.T) {
	s := initTestServer()
	originalTranslation := "originalTranslation"
	request := translationRequest{Translation: originalTranslation}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationApi, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingTranslations(s)[0].Id

	transcription := "[updateTranscription]"
	tr := "UpdateTranslation"
	text := "UpdateText"
	example := "UpdateExample"

	request = translationRequest{
		Transcription: transcription,
		Translation:   tr,
		Text:          text,
		Example:       example,
	}
	jsonValue, _ = json.Marshal(request)
	req, _ = http.NewRequest("PUT", v1TranslationApi+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	req, _ = http.NewRequest("GET", v1TranslationApi+"/"+id, nil)
	setAuthToken(s, req)
	w = httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record translationResponse
	json.Unmarshal(w.Body.Bytes(), &record)

	assert.Equal(t, originalTranslation, record.Translation)
}

func getExistingTranslations(s *HttpServer) []translationResponse {
	req, _ := http.NewRequest("GET", v1TranslationApi+"/last?limit=10", nil)
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var records []translationResponse
	json.Unmarshal(w.Body.Bytes(), &records)
	return records
}
