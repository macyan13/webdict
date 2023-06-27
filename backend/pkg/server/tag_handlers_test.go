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
		Name: tg,
	}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
	setAdminAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	records := getExistingTags(t, s)
	created := records[0]

	assert.Equal(t, tg, created.Name)
}

func TestServer_CreateTagUnauthorised(t *testing.T) {
	s := initTestServer()
	tg := "CreateTag"

	request := tagRequest{
		Name: tg,
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

	request := tagRequest{Name: "test"}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
	setAdminAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingTags(t, s)[0].ID
	req, _ = http.NewRequest("DELETE", v1TagAPI+"/"+id, bytes.NewBuffer(jsonValue))
	setAdminAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Zero(t, len(getExistingTags(t, s)))
}

func TestServer_DeleteTagByIdUnauthorised(t *testing.T) {
	s := initTestServer()

	request := tagRequest{Name: "test"}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
	setAdminAuthToken(s, req)
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
	request := tagRequest{Name: "test"}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
	setAdminAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingTags(t, s)[0].ID
	tg := "UpdateTag"

	request = tagRequest{
		Name: tg,
	}

	jsonValue, _ = json.Marshal(request)
	req, _ = http.NewRequest("PUT", v1TagAPI+"/"+id, bytes.NewBuffer(jsonValue))
	setAdminAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", v1TagAPI+"/"+id, http.NoBody)
	setAdminAuthToken(s, req)
	w = httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record tagResponse
	err := json.Unmarshal(w.Body.Bytes(), &record)
	assert.Nil(t, err)

	assert.Equal(t, tg, record.Name)
}

func TestServer_UpdateTagUnauthorised(t *testing.T) {
	s := initTestServer()
	originalTag := "originalTag"
	request := tagRequest{
		Name: originalTag,
	}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
	setAdminAuthToken(s, req)
	s.engine.ServeHTTP(httptest.NewRecorder(), req)
	id := getExistingTags(t, s)[0].ID
	tg := "UpdateTag"

	request = tagRequest{
		Name: tg,
	}

	jsonValue, _ = json.Marshal(request)
	req, _ = http.NewRequest("PUT", v1TagAPI+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	req, _ = http.NewRequest("GET", v1TagAPI+"/"+id, http.NoBody)
	setAdminAuthToken(s, req)
	w = httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record tagResponse
	err := json.Unmarshal(w.Body.Bytes(), &record)
	assert.Nil(t, err)

	assert.Equal(t, originalTag, record.Name)
}

func TestServer_GetTags(t *testing.T) {
	s := initTestServer()

	tags := map[string]interface{}{"testTag1": nil, "testTag2": nil}

	for tag := range tags {
		request := tagRequest{Name: tag}
		jsonValue, _ := json.Marshal(request)
		req, _ := http.NewRequest("POST", v1TagAPI, bytes.NewBuffer(jsonValue))
		setAdminAuthToken(s, req)
		recorder := httptest.NewRecorder()
		s.engine.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusCreated, recorder.Code)
	}

	createdTags := getExistingTags(t, s)

	assert.Equal(t, len(tags), len(createdTags))

	for _, tag := range createdTags {
		_, exist := tags[tag.Name]
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
	setAdminAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var records []tagResponse
	err := json.Unmarshal(w.Body.Bytes(), &records)
	assert.Nil(t, err)
	return records
}
