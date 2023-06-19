package server

import (
	"bytes"
	"encoding/json"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const v1ProfileAPI = "/v1/api/profile"

func TestHTTPServer_GetProfile(t *testing.T) {
	s := initTestServer()

	name := "John Do"
	email := "john@test.com"
	passwd := "testPassword"
	createUser(t, s, name, email, passwd)

	profile := getProfile(t, s, email, passwd)
	assert.Equal(t, name, profile.Name)
	assert.Equal(t, email, profile.Email)
	assert.Equal(t, int(user.Author), profile.Role.ID)
}

func TestHTTPServer_GetProfile_Unauthorized(t *testing.T) {
	s := initTestServer()
	jsonValue, _ := json.Marshal(userRequest{})
	req, _ := http.NewRequest("GET", v1ProfileAPI, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHTTPServer_TestHTTPServer_UpdateProfile(t *testing.T) {
	s := initTestServer()
	currentPasswd := "testPassword"
	email := "john@test.com"

	createUser(t, s, "John Do", email, currentPasswd)

	updatedName := "test"
	updatedEmail := "updated@test.com"
	newPasswd := "newPasswd12345"

	updRequest := updateProfileRequest{
		Name:            updatedName,
		Email:           updatedEmail,
		CurrentPassword: currentPasswd,
		NewPassword:     newPasswd,
	}

	jsonValue, _ := json.Marshal(updRequest)
	req, _ := http.NewRequest("PUT", v1ProfileAPI, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	setAuthTokenWithCredentials(s, req, email, currentPasswd)
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	profile := getProfile(t, s, updatedEmail, newPasswd)
	assert.Equal(t, updatedName, profile.Name)
	assert.Equal(t, updatedEmail, profile.Email)
}

func TestHTTPServer_TestHTTPServerUnauthorized(t *testing.T) {
	s := initTestServer()
	currentPasswd := "testPassword"
	email := "john@test.com"

	createUser(t, s, "John Do", email, currentPasswd)

	updatedName := "test"
	updatedEmail := "updated@test.com"
	newPasswd := "newPasswd12345"

	updRequest := updateProfileRequest{
		Name:            updatedName,
		Email:           updatedEmail,
		CurrentPassword: currentPasswd,
		NewPassword:     newPasswd,
	}

	jsonValue, _ := json.Marshal(updRequest)
	req, _ := http.NewRequest("PUT", v1ProfileAPI, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func getProfile(t *testing.T, s *HTTPServer, email, passwd string) userResponse {
	var profile userResponse
	req, _ := http.NewRequest("GET", v1ProfileAPI, http.NoBody)
	setAuthTokenWithCredentials(s, req, email, passwd)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &profile)
	assert.Nil(t, err)
	return profile
}
