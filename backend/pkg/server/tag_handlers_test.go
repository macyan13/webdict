package server

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const v1TagAPI = "/v1/api/tags"

func TestServer_CreateTag(t *testing.T) {
	s := initTestServer()
	tg := "CreateTag"

	request := tagRequest{
		Tag: tg,
	}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	records := getExistingTags(t, s)
	created := records[0]

	assert.Equal(t, tg, created.Tag)
}

func TestServer_CreateTagUnauthorised(t *testing.T) {
	s := initTestServer()
	tg := "CreateTag"

	request := tagRequest{
		Tag: tg,
	}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Zero(t, len(getExistingTags(t, s)))
}

func TestServer_DeleteTagById(t *testing.T) {
	s := initTestServer()

	request := tagRequest{Tag: "test"}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingTags(t, s)[0].ID
	req, _ = http.NewRequest("DELETE", v1TagAPI+"/"+id, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Zero(t, len(getExistingTags(t, s)))
}

func TestServer_DeleteTagByIdUnauthorised(t *testing.T) {
	s := initTestServer()

	request := tagRequest{Tag: "test"}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingTags(t, s)[0].ID
	req, _ = http.NewRequest("DELETE", v1TagAPI+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, 1, len(getExistingTags(t, s)))
}

func TestServer_UpdateTag(t *testing.T) {
	s := initTestServer()
	request := tagRequest{Tag: "test"}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingTags(t, s)[0].ID
	tg := "UpdateTag"

	request = tagRequest{
		Tag: tg,
	}

	jsonValue, _ = json.Marshal(request)
	req, _ = http.NewRequest("PUT", v1TagAPI+"/"+id, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", v1TagAPI+"/"+id, http.NoBody)
	setAuthToken(s, req)
	w = httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record tagResponse
	err := json.Unmarshal(w.Body.Bytes(), &record)
	assert.Nil(t, err)

	assert.Equal(t, tg, record.Tag)
}

func TestServer_UpdateTagUnauthorised(t *testing.T) {
	s := initTestServer()
	originalTag := "originalTag"
	request := tagRequest{
		Tag: originalTag,
	}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingTags(t, s)[0].ID
	tg := "UpdateTag"

	request = tagRequest{
		Tag: tg,
	}

	jsonValue, _ = json.Marshal(request)
	req, _ = http.NewRequest("PUT", v1TagAPI+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	req, _ = http.NewRequest("GET", v1TagAPI+"/"+id, http.NoBody)
	setAuthToken(s, req)
	w = httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record tagResponse
	err := json.Unmarshal(w.Body.Bytes(), &record)
	assert.Nil(t, err)

	assert.Equal(t, originalTag, record.Tag)
}

func TestServer_GetTags(t *testing.T) {
	s := initTestServer()

	tags := map[string]interface{}{"testTag1": nil, "testTag2": nil}

	for tag := range tags {
		request := tagRequest{Tag: tag}
		jsonValue, _ := json.Marshal(request)
		req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
		setAuthToken(s, req)
		recorder := httptest.NewRecorder()
		s.engine.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code)
	}

	createdTags := getExistingTags(t, s)

	assert.Equal(t, len(tags), len(createdTags))

	for _, tag := range createdTags {
		_, exist := tags[tag.Tag]
		assert.True(t, exist)
	}
}

func TestServer_GetTagsUnauthorised(t *testing.T) {
	s := initTestServer()

	req, _ := http.NewRequest("GET", v1TagAPI, http.NoBody)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func getExistingTags(t *testing.T, s *HTTPServer) []tagResponse {
	req, _ := http.NewRequest("GET", v1TagAPI, http.NoBody)
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var records []tagResponse
	err := json.Unmarshal(w.Body.Bytes(), &records)
	assert.Nil(t, err)
	return records
}
