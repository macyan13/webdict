package server

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const v1TagApi = "/v1/api/tags"

func TestServer_CreateTag(t *testing.T) {
	s := initTestServer()
	tg := "CreateTag"

	request := tagRequest{
		Tag: tg,
	}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagApi, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	records := getExistingTags(s)
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
	req, _ := http.NewRequest("POST", v1TagApi, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Zero(t, len(getExistingTags(s)))
}

func TestServer_DeleteTagById(t *testing.T) {
	s := initTestServer()

	request := tagRequest{}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagApi, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingTags(s)[0].Id
	req, _ = http.NewRequest("DELETE", v1TagApi+"/"+id, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Zero(t, len(getExistingTags(s)))
}

func TestServer_DeleteTagByIdUnauthorised(t *testing.T) {
	s := initTestServer()

	request := tagRequest{}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagApi, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingTags(s)[0].Id
	req, _ = http.NewRequest("DELETE", v1TagApi+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, 1, len(getExistingTags(s)))
}

func TestServer_UpdateTag(t *testing.T) {
	s := initTestServer()
	request := tagRequest{}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagApi, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingTags(s)[0].Id
	tg := "UpdateTag"

	request = tagRequest{
		Tag: tg,
	}

	jsonValue, _ = json.Marshal(request)
	req, _ = http.NewRequest("PUT", v1TagApi+"/"+id, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", v1TagApi+"/"+id, nil)
	setAuthToken(s, req)
	w = httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record tagResponse
	json.Unmarshal(w.Body.Bytes(), &record)

	assert.Equal(t, tg, record.Tag)
}

func TestServer_UpdateTagUnauthorised(t *testing.T) {
	s := initTestServer()
	originalTag := "originalTag"
	request := tagRequest{
		Tag: originalTag,
	}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagApi, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingTags(s)[0].Id
	tg := "UpdateTag"

	request = tagRequest{
		Tag: tg,
	}

	jsonValue, _ = json.Marshal(request)
	req, _ = http.NewRequest("PUT", v1TagApi+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	req, _ = http.NewRequest("GET", v1TagApi+"/"+id, nil)
	setAuthToken(s, req)
	w = httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record tagResponse
	json.Unmarshal(w.Body.Bytes(), &record)

	assert.Equal(t, originalTag, record.Tag)
}

func TestServer_GetTags(t *testing.T) {
	s := initTestServer()

	tags := map[string]interface{}{"testTag1": nil, "testTag2": nil}

	for tag, _ := range tags {
		request := tagRequest{Tag: tag}
		jsonValue, _ := json.Marshal(request)
		req, _ := http.NewRequest("POST", v1TagApi, bytes.NewBuffer(jsonValue))
		setAuthToken(s, req)
		recorder := httptest.NewRecorder()
		s.engine.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code)
	}

	createdTags := getExistingTags(s)

	assert.Equal(t, len(tags), len(createdTags))

	for _, tag := range createdTags {
		_, exist := tags[tag.Tag]
		assert.True(t, exist)
	}
}

func TestServer_GetTagsUnauthorised(t *testing.T) {
	s := initTestServer()

	req, _ := http.NewRequest("GET", v1TagApi, nil)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func getExistingTags(s *HttpServer) []tagResponse {
	req, _ := http.NewRequest("GET", v1TagApi, nil)
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var records []tagResponse
	json.Unmarshal(w.Body.Bytes(), &records)
	return records
}