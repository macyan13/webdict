package server

import (
	"bytes"
	"encoding/json"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const v1UserAPI = "/v1/api/users"

func TestHTTPServer_CreateUser(t *testing.T) {
	s := initTestServer()
	name := "John Do"
	email := "john@test.com"
	pwd := "testPassword"

	response := createUser(t, s, name, email, pwd)

	usr := getUserById(t, s, response.ID)
	assert.Equal(t, name, usr.Name)
	assert.Equal(t, response.ID, usr.ID)
	assert.Equal(t, email, usr.Email)
	assert.Equal(t, int(user.Author), usr.Role)
}

func TestHTTPServer_CreateUser_Unauthorized(t *testing.T) {
	s := initTestServer()

	jsonValue, _ := json.Marshal(userRequest{
		Name:     "test",
		Email:    "test",
		Password: "test",
	})
	req, _ := http.NewRequest("POST", v1UserAPI, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHTTPServer_CreateUser_NotAdmin(t *testing.T) {
	s := initTestServer()
	name := "John Do"
	email := "john@test.com"
	pwd := "testPassword"
	createUser(t, s, name, email, pwd)

	jsonValue, _ := json.Marshal(userRequest{})
	req, _ := http.NewRequest("POST", v1UserAPI, bytes.NewBuffer(jsonValue))
	setAuthTokenWithCredentials(s, req, email, pwd)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHTTPServer_GetUsers(t *testing.T) {
	s := initTestServer()

	users := map[string]userRequest{
		"john@test.com": {
			Name:     "John Do",
			Email:    "john@test.com",
			Password: "testPassword",
		},
		"bob@test.com": {
			Name:     "Bob",
			Email:    "bob@test.com",
			Password: "testPassword",
		},
	}

	for usr := range users {
		jsonValue, _ := json.Marshal(users[usr])
		req, _ := http.NewRequest("POST", v1UserAPI, bytes.NewBuffer(jsonValue))
		setAuthToken(s, req)
		w := httptest.NewRecorder()
		s.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

	}

	var records []userResponse
	req, _ := http.NewRequest("GET", v1UserAPI, http.NoBody)
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &records)
	assert.Nil(t, err)

	for i := range records {
		usr, ok := users[records[i].Email]

		if !ok {
			continue
		}
		assert.Equal(t, usr.Email, records[i].Email)
		assert.Equal(t, usr.Name, records[i].Name)
		assert.Equal(t, int(user.Author), records[i].Role)
	}
}

func TestHTTPServer_GetUsers_Unauthorized(t *testing.T) {
	s := initTestServer()
	jsonValue, _ := json.Marshal(userRequest{})
	req, _ := http.NewRequest("GET", v1UserAPI, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHTTPServer_GetUsers_NotAdmin(t *testing.T) {
	s := initTestServer()
	name := "John Do"
	email := "john@test.com"
	pwd := "testPassword"
	createUser(t, s, name, email, pwd)

	req, _ := http.NewRequest("GET", v1UserAPI, http.NoBody)
	setAuthTokenWithCredentials(s, req, email, pwd)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func getUserById(t *testing.T, s *HTTPServer, id string) userResponse {
	req, _ := http.NewRequest("GET", v1UserAPI+"/"+id, http.NoBody)
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record userResponse
	err := json.Unmarshal(w.Body.Bytes(), &record)
	assert.Nil(t, err)

	return record
}

func createUser(t *testing.T, s *HTTPServer, name string, email string, passwd string) idResponse {
	jsonValue, _ := json.Marshal(userRequest{
		Name:     name,
		Email:    email,
		Password: passwd,
	})
	req, _ := http.NewRequest("POST", v1UserAPI, bytes.NewBuffer(jsonValue))
	setAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response idResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	return response
}
