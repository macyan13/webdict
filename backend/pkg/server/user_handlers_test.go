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

const v1UserAPI = "/v1/api/users"

func TestHTTPServer_CreateUser(t *testing.T) {
	s := initTestServer()
	name := "John Do"
	email := "john@test.com"
	pwd := "testPassword"

	response := createUser(t, s, name, email, pwd)

	usr := getUserByID(t, s, response.ID)
	assert.Equal(t, name, usr.Name)
	assert.Equal(t, response.ID, usr.ID)
	assert.Equal(t, email, usr.Email)
	assert.Equal(t, int(user.Author), usr.Role.ID)
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
	setAuthTokenWithCredentials(t, s, req, email, pwd)
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
		setAdminAuthToken(t, s, req)
		w := httptest.NewRecorder()
		s.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

	}

	var records []userResponse
	req, _ := http.NewRequest("GET", v1UserAPI, http.NoBody)
	setAdminAuthToken(t, s, req)
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
		assert.Equal(t, int(user.Author), records[i].Role.ID)
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
	setAuthTokenWithCredentials(t, s, req, email, pwd)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHTTPServer_UpdateUser_Unauthorized(t *testing.T) {
	s := initTestServer()
	name := "John Do"
	email := "john@test.com"
	pwd := "testPassword"

	response := createUser(t, s, name, email, pwd)

	updRequest := userRequest{
		Name:     "test",
		Email:    "updated@test.com",
		Password: "newPasswd12345",
		Role:     int(user.Author),
	}

	jsonValue, _ := json.Marshal(updRequest)
	req, _ := http.NewRequest("PUT", v1UserAPI+"/"+response.ID, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	usr := getUserByID(t, s, response.ID)
	assert.Equal(t, name, usr.Name)
	assert.Equal(t, response.ID, usr.ID)
	assert.Equal(t, email, usr.Email)
	assert.Equal(t, int(user.Author), usr.Role.ID)
}

func TestHTTPServer_UpdateUser_NotAdmin(t *testing.T) {
	s := initTestServer()
	name := "John Do"
	email := "john@test.com"
	pwd := "testPassword"

	response := createUser(t, s, name, email, pwd)

	updRequest := userRequest{
		Name:     "test",
		Email:    "updated@test.com",
		Password: "newPasswd12345",
		Role:     int(user.Author),
	}

	jsonValue, _ := json.Marshal(updRequest)
	req, _ := http.NewRequest("PUT", v1UserAPI+"/"+response.ID, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	setAuthTokenWithCredentials(t, s, req, email, pwd)
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	usr := getUserByID(t, s, response.ID)
	assert.Equal(t, name, usr.Name)
	assert.Equal(t, response.ID, usr.ID)
	assert.Equal(t, email, usr.Email)
	assert.Equal(t, int(user.Author), usr.Role.ID)
}

func TestHTTPServer_UpdateUser(t *testing.T) {
	s := initTestServer()
	name := "John Do"
	email := "john@test.com"
	pwd := "testPassword"

	response := createUser(t, s, "test", "test@mail.com", "passwd")

	updRequest := userRequest{
		Name:     name,
		Email:    email,
		Password: pwd,
		Role:     int(user.Admin),
	}

	jsonValue, _ := json.Marshal(updRequest)
	req, _ := http.NewRequest("PUT", v1UserAPI+"/"+response.ID, bytes.NewBuffer(jsonValue))
	setAdminAuthToken(t, s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	usr := getUserByID(t, s, response.ID)
	assert.Equal(t, name, usr.Name)
	assert.Equal(t, response.ID, usr.ID)
	assert.Equal(t, email, usr.Email)
	assert.Equal(t, int(user.Admin), usr.Role.ID)
}

func TestHTTPServer_DeleteUser_Unauthorized(t *testing.T) {
	s := initTestServer()
	name := "John Do"
	email := "john@test.com"
	pwd := "testPassword"

	response := createUser(t, s, name, email, pwd)

	updRequest := userRequest{
		Name:     "test",
		Email:    "updated@test.com",
		Password: "newPasswd12345",
		Role:     int(user.Author),
	}

	jsonValue, _ := json.Marshal(updRequest)
	req, _ := http.NewRequest("DELETE", v1UserAPI+"/"+response.ID, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	usr := getUserByID(t, s, response.ID)
	assert.Equal(t, name, usr.Name)
}

func TestHTTPServer_DeleteUser_NotAdmin(t *testing.T) {
	s := initTestServer()
	name := "John Do"
	email := "john@test.com"
	pwd := "testPassword"

	response := createUser(t, s, name, email, pwd)

	updRequest := userRequest{
		Name:     "test",
		Email:    "updated@test.com",
		Password: "newPasswd12345",
		Role:     int(user.Author),
	}

	jsonValue, _ := json.Marshal(updRequest)
	req, _ := http.NewRequest("DELETE", v1UserAPI+"/"+response.ID, bytes.NewBuffer(jsonValue))
	setAuthTokenWithCredentials(t, s, req, email, pwd)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	usr := getUserByID(t, s, response.ID)
	assert.Equal(t, name, usr.Name)
}

func TestHTTPServer_DeleteUser(t *testing.T) {
	s := initTestServer()
	name := "John Do"
	email := "john@test.com"
	pwd := "testPassword"

	response := createUser(t, s, name, email, pwd)

	updRequest := userRequest{
		Name:     "test",
		Email:    "updated@test.com",
		Password: "newPasswd12345",
		Role:     int(user.Author),
	}

	jsonValue, _ := json.Marshal(updRequest)
	req, _ := http.NewRequest("DELETE", v1UserAPI+"/"+response.ID, bytes.NewBuffer(jsonValue))
	setAdminAuthToken(t, s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", v1UserAPI+"/"+response.ID, http.NoBody)
	setAdminAuthToken(t, s, req)
	w = httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func getUserByID(t *testing.T, s *testHTTPServer, id string) userResponse {
	req, _ := http.NewRequest("GET", v1UserAPI+"/"+id, http.NoBody)
	setAdminAuthToken(t, s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	var record userResponse
	err := json.Unmarshal(w.Body.Bytes(), &record)
	assert.Nil(t, err)

	return record
}

func createUser(t *testing.T, s *testHTTPServer, name, email, passwd string) idResponse {
	jsonValue, _ := json.Marshal(userRequest{
		Name:     name,
		Email:    email,
		Password: passwd,
	})
	req, _ := http.NewRequest("POST", v1UserAPI, bytes.NewBuffer(jsonValue))
	setAdminAuthToken(t, s, req)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response idResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	return response
}
