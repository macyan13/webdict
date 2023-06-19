package server

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const v1RoleAPI = "/v1/api/roles"

func TestHTTPServer_GetRoles_Unauthorized(t *testing.T) {
	s := initTestServer()

	jsonValue, _ := json.Marshal("")
	req, _ := http.NewRequest("GET", v1RoleAPI, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHTTPServer_GetRoles_NotAdmin(t *testing.T) {
	s := initTestServer()
	name := "John Do"
	email := "john@test.com"
	pwd := "testPassword"
	createUser(t, s, name, email, pwd)

	jsonValue, _ := json.Marshal("")
	req, _ := http.NewRequest("GET", v1RoleAPI, bytes.NewBuffer(jsonValue))
	setAuthTokenWithCredentials(s, req, email, pwd)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHTTPServer_GetRoles_Authorized(t *testing.T) {
	s := initTestServer()

	jsonValue, _ := json.Marshal("")
	req, _ := http.NewRequest("GET", v1RoleAPI, bytes.NewBuffer(jsonValue))
	setAdminAuthToken(s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response rolesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(response.Roles))
}
