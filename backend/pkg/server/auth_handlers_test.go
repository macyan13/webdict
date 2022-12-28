package server

import (
	"bytes"
	"encoding/json"
	"github.com/macyan13/webdict/backend/pkg/auth"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const authApi = "/v1/api/auth"

func TestServer_SighIn(t *testing.T) {
	s := initTestServer()

	request := SignInRequest{
		Email:    s.opts.Admin.AdminEmail,
		Password: s.opts.Admin.AdminPasswd,
	}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", authApi+"/signin", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response AuthTokenResponse

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, "Bearer", response.Type, "Route:SignIn - Auth type must be Bearer")
	assert.NotEmpty(t, response.AccessToken, "Route:SignIn -AccessToken must present")

	cookies := w.Header().Get("Set-Cookie") // Todo: add proper tests for cookie elements
	assert.Contains(t, cookies, "refreshToken=", "Route:SignIn - RefreshToken must present in cookies")
}

func TestServer_Refresh(t *testing.T) {
	s := initTestServer()

	req, _ := http.NewRequest("POST", authApi+"/refresh", &bytes.Buffer{})
	req.AddCookie(&http.Cookie{
		Name:  "refreshToken",
		Value: getRefreshToken(s),
	})
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response AuthTokenResponse

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, "Bearer", response.Type, "Route:SignIn - Auth type must be Bearer")
	assert.NotEmpty(t, response.AccessToken, "Route:SignIn -AccessToken must present")
}

func getRefreshToken(s *HttpServer) string {
	request := SignInRequest{
		Email:    s.opts.Admin.AdminEmail,
		Password: s.opts.Admin.AdminPasswd,
	}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", authApi+"/signin", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	token := []rune(strings.Split(w.Header().Get("Set-Cookie"), ";")[0])
	return string(token[13:]) // `refreshToken=` - 13
}

var authenticationToken *auth.AuthenticationToken

func setAuthToken(s *HttpServer, r *http.Request) {
	if authenticationToken == nil {
		token, _ := s.authHandler.Authenticate(s.opts.Admin.AdminEmail, s.opts.Admin.AdminPasswd)
		authenticationToken = &token
	}
	r.Header.Set("Authorization", authenticationToken.Type+" "+authenticationToken.Token)
}
