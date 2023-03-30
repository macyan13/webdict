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

const v1TranslationAPI = "/v1/api/translations"

func TestServer_CreateTranslation(t *testing.T) {
	s := initTestServer()
	transcription := "CreateTranscription"
	tr := "CreateTranslation"
	text := "CreateText"
	example := "CreateExample"
	tg := "testTag"

	user, err := s.userRepo.GetByEmail(s.opts.Admin.AdminEmail)
	assert.Nil(t, err)

	_, err = s.app.Commands.AddTag.Handle(command.AddTag{
		Tag:      tg,
		AuthorID: user.ID(),
	})
	assert.Nil(t, err)

	tags, err := s.app.Queries.AllTags.Handle(query.AllTags{AuthorID: user.ID()})
	assert.Nil(t, err)
	tagID := tags[0].ID

	request := translationRequest{
		Transcription: transcription,
		Target:        tr,
		Source:        text,
		Example:       example,
		TagIds:        []string{tagID},
	}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	records := getExistingTranslations(t, s)
	created := records[0]

	assert.Equal(t, tr, created.Target)
	assert.Equal(t, text, created.Source)
	assert.Equal(t, transcription, created.Transcription)
	assert.Equal(t, example, created.Example)
	assert.Equal(t, tagID, created.Tags[0].ID)
}

func TestServer_CreateTranslationUnauthorised(t *testing.T) {
	s := initTestServer()
	transcription := "CreateTranscription"
	tr := "CreateTranslation"
	text := "CreateText"
	example := "CreateExample"

	request := translationRequest{
		Transcription: transcription,
		Target:        tr,
		Source:        text,
		Example:       example,
		TagIds:        []string{},
	}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationAPI, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Zero(t, len(getExistingTranslations(t, s)))
}

func TestServer_DeleteTranslationById(t *testing.T) {
	s := initTestServer()

	jsonValue, _ := json.Marshal(translationRequest{Source: "test", Target: "test"})
	req, _ := http.NewRequest("POST", v1TranslationAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingTranslations(t, s)[0].ID
	req, _ = http.NewRequest("DELETE", v1TranslationAPI+"/"+id, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Zero(t, len(getExistingTranslations(t, s)))
}

func TestServer_DeleteTranslationByIdUnauthorised(t *testing.T) {
	s := initTestServer()

	jsonValue, _ := json.Marshal(translationRequest{Source: "test", Target: "test"})
	req, _ := http.NewRequest("POST", v1TranslationAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingTranslations(t, s)[0].ID
	req, _ = http.NewRequest("DELETE", v1TranslationAPI+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, 1, len(getExistingTranslations(t, s)))
}

func TestServer_UpdateTranslation(t *testing.T) {
	s := initTestServer()

	jsonValue, _ := json.Marshal(translationRequest{Source: "test", Target: "test"})
	req, _ := http.NewRequest("POST", v1TranslationAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingTranslations(t, s)[0].ID

	transcription := "[updateTranscription]"
	tr := "UpdateTranslation"
	source := "UpdateText"
	example := "UpdateExample"

	request := translationRequest{
		Transcription: transcription,
		Target:        tr,
		Source:        source,
		Example:       example,
	}
	jsonValue, _ = json.Marshal(request)
	req, _ = http.NewRequest("PUT", v1TranslationAPI+"/"+id, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", v1TranslationAPI+"/"+id, http.NoBody)
	setAuthToken(s, req)
	w = httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record translationResponse
	err := json.Unmarshal(w.Body.Bytes(), &record)
	assert.Nil(t, err)

	assert.Equal(t, tr, record.Target)
	assert.Equal(t, source, record.Source)
	assert.Equal(t, transcription, record.Transcription)
	assert.Equal(t, example, record.Example)
}

func TestServer_UpdateTranslationUnauthorised(t *testing.T) {
	s := initTestServer()
	originalTranslation := "originalTranslation"
	request := translationRequest{Target: originalTranslation, Source: "test"}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingTranslations(t, s)[0].ID

	transcription := "[updateTranscription]"
	tr := "UpdateTranslation"
	text := "UpdateText"
	example := "UpdateExample"

	request = translationRequest{
		Transcription: transcription,
		Target:        tr,
		Source:        text,
		Example:       example,
	}
	jsonValue, _ = json.Marshal(request)
	req, _ = http.NewRequest("PUT", v1TranslationAPI+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	req, _ = http.NewRequest("GET", v1TranslationAPI+"/"+id, http.NoBody)
	setAuthToken(s, req)
	w = httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record translationResponse
	err := json.Unmarshal(w.Body.Bytes(), &record)
	assert.Nil(t, err)

	assert.Equal(t, originalTranslation, record.Target)
}

func getExistingTranslations(t *testing.T, s *HTTPServer) []translationResponse {
	req, _ := http.NewRequest("GET", v1TranslationAPI+"/last?pageSize=10", http.NoBody)
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var response lastTranslationsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	return response.Translations
}
