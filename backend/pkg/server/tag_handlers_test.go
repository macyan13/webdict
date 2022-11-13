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
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	records := getExistingTags(s)
	created := records[0]

	assert.Equal(t, tg, created.Tag)
}

func TestServer_DeleteTagById(t *testing.T) {
	s := initTestServer()

	request := tagRequest{}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagApi, bytes.NewBuffer(jsonValue))
	s.engine.ServeHTTP(httptest.NewRecorder(), req)

	id := getExistingTags(s)[0].Id
	req, _ = http.NewRequest("DELETE", v1TagApi+"/"+id, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Zero(t, len(getExistingTags(s)))
}

func TestServer_UpdateTag(t *testing.T) {
	s := initTestServer()
	request := tagRequest{}
	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", v1TagApi, bytes.NewBuffer(jsonValue))
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
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", v1TagApi+"/"+id, nil)
	w = httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record tagResponse
	json.Unmarshal(w.Body.Bytes(), &record)

	assert.Equal(t, tg, record.Tag)
}

func getExistingTags(s *HttpServer) []tagResponse {
	req, _ := http.NewRequest("GET", v1TagApi, nil)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var records []tagResponse
	json.Unmarshal(w.Body.Bytes(), &records)
	return records
}
