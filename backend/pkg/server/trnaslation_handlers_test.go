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
	ln := "EN"

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
		Lang:          ln,
	}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	records := getExistingTranslations(t, s, ln)
	created := records[0]

	assert.Equal(t, tr, created.Target)
	assert.Equal(t, text, created.Source)
	assert.Equal(t, transcription, created.Transcription)
	assert.Equal(t, example, created.Example)
	assert.Equal(t, tagID, created.Tags[0].ID)
	assert.Equal(t, ln, created.Lang)
}

func TestServer_CreateTranslationUnauthorised(t *testing.T) {
	s := initTestServer()
	transcription := "CreateTranscription"
	tr := "CreateTranslation"
	text := "CreateText"
	example := "CreateExample"
	ln := "DE"

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
	assert.Zero(t, len(getExistingTranslations(t, s, ln)))
}

func TestServer_DeleteTranslationById(t *testing.T) {
	s := initTestServer()
	ln := "EN"

	jsonValue, _ := json.Marshal(translationRequest{Source: "test", Target: "test", Lang: ln})
	req, _ := http.NewRequest("POST", v1TranslationAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingTranslations(t, s, ln)[0].ID
	req, _ = http.NewRequest("DELETE", v1TranslationAPI+"/"+id, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Zero(t, len(getExistingTranslations(t, s, ln)))
}

func TestServer_DeleteTranslationByIdUnauthorised(t *testing.T) {
	s := initTestServer()
	ln := "EN"

	jsonValue, _ := json.Marshal(translationRequest{Source: "test", Target: "test", Lang: ln})
	req, _ := http.NewRequest("POST", v1TranslationAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingTranslations(t, s, ln)[0].ID
	req, _ = http.NewRequest("DELETE", v1TranslationAPI+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, 1, len(getExistingTranslations(t, s, ln)))
}

func TestServer_UpdateTranslation(t *testing.T) {
	s := initTestServer()
	ln := "EN"

	jsonValue, _ := json.Marshal(translationRequest{Source: "test", Target: "test", Lang: ln})
	req, _ := http.NewRequest("POST", v1TranslationAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingTranslations(t, s, ln)[0].ID

	transcription := "[updateTranscription]"
	tr := "UpdateTranslation"
	source := "UpdateText"
	example := "UpdateExample"

	request := translationRequest{
		Transcription: transcription,
		Target:        tr,
		Source:        source,
		Example:       example,
		Lang:          ln,
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
	assert.Equal(t, ln, record.Lang)
}

func TestServer_UpdateTranslationUnauthorised(t *testing.T) {
	s := initTestServer()
	originalTranslation := "originalTranslation"
	ln := "EN"

	request := translationRequest{Target: originalTranslation, Source: "test", Lang: ln}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TranslationAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingTranslations(t, s, ln)[0].ID

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

func getExistingTranslations(t *testing.T, s *HTTPServer, lang string) []translationResponse {
	req, _ := http.NewRequest("GET", v1TranslationAPI+"/last?pageSize=10&lang="+lang, http.NoBody)
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var response lastTranslationsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	return response.Translations
}
