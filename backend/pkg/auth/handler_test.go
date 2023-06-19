package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_Authenticate(t *testing.T) {
	var cipher = Cipher{}
	var hashedPwd, _ = cipher.GenerateHash("password")

	type fields struct {
		userRepo user.Repository
		tokener  tokener
	}
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     AuthenticationToken
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"User does not exist",
			func() fields {
				repository := user.MockRepository{}
				repository.On("GetByEmail", "notExist").Return(nil, user.ErrNotFound)
				return fields{
					userRepo: &repository,
					tokener:  &mockTokener{},
				}
			},
			args{
				email:    "notExist",
				password: "test",
			},
			AuthenticationToken{},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, err == ErrInvalidCredentials, i)
				return true
			},
		},
		{
			"Password is not valid",
			func() fields {
				email := "test@email.com"
				existingUser, err := user.NewUser("test", email, hashedPwd, user.Admin)
				assert.NoError(t, err)

				repository := user.MockRepository{}
				repository.On("GetByEmail", email).Return(existingUser, nil)
				return fields{
					userRepo: &repository,
					tokener:  &mockTokener{},
				}
			},
			args{
				email:    "test@email.com",
				password: "test",
			},
			AuthenticationToken{},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, err == ErrInvalidCredentials, i)
				return true
			},
		},
		{
			"Error on token generation",
			func() fields {
				email := "test@email.com"
				existingUser, err := user.NewUser("test", email, hashedPwd, user.Admin)
				assert.NoError(t, err)

				repository := user.MockRepository{}
				repository.On("GetByEmail", email).Return(existingUser, nil)

				tokener := mockTokener{}
				tokener.On("generateToken", email, mock.IsType(time.Time{})).Return("", fmt.Errorf("noToken"))
				return fields{
					userRepo: &repository,
					tokener:  &tokener,
				}
			},
			args{
				email:    "test@email.com",
				password: "password",
			},
			AuthenticationToken{},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "noToken", err.Error(), i)
				return true
			},
		},
		{
			"Positive Case",
			func() fields {
				email := "test@email.com"
				existingUser, err := user.NewUser("test", email, hashedPwd, user.Admin)
				assert.NoError(t, err)

				repository := user.MockRepository{}
				repository.On("GetByEmail", email).Return(existingUser, nil)

				tokener := mockTokener{}
				tokener.On("generateToken", email, mock.IsType(time.Time{})).Return("validToken", nil)
				return fields{
					userRepo: &repository,
					tokener:  &tokener,
				}
			},
			args{
				email:    "test@email.com",
				password: "password",
			},
			AuthenticationToken{
				Token: "validToken",
				Type:  authType,
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsFn()
			h := Handler{
				userRepo: fields.userRepo,
				tokener:  fields.tokener,
				cipher:   cipher,
			}
			got, err := h.Authenticate(tt.args.email, tt.args.password)
			if !tt.wantErr(t, err, fmt.Sprintf("Authenticate(%v, %v)", tt.args.email, tt.args.password)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Authenticate(%v, %v)", tt.args.email, tt.args.password)
		})
	}
}

func TestHandler_GenerateRefreshToken(t *testing.T) {
	type fields struct {
		userRepo user.Repository
		tokener  tokener
	}
	type args struct {
		email string
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     RefreshToken
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Error on token generation",
			func() fields {
				repository := user.MockRepository{}
				tokener := mockTokener{}
				tokener.On("generateToken", "testEmail", mock.IsType(time.Time{})).Return("", fmt.Errorf("noToken"))

				return fields{
					userRepo: &repository,
					tokener:  &tokener,
				}
			},
			args{
				email: "testEmail",
			},
			RefreshToken{},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "noToken", err.Error(), i)
				return true
			},
		},
		{
			"Positive case",
			func() fields {
				repository := user.MockRepository{}
				tokener := mockTokener{}
				tokener.On("generateToken", "testEmail", mock.IsType(time.Time{})).Return("validToken", nil)

				return fields{
					userRepo: &repository,
					tokener:  &tokener,
				}
			},
			args{
				email: "testEmail",
			},
			RefreshToken{
				Token: "validToken",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.fieldsFn()
			h := Handler{
				userRepo: args.userRepo,
				tokener:  args.tokener,
			}
			got, err := h.GenerateRefreshToken(tt.args.email)
			if !tt.wantErr(t, err, fmt.Sprintf("GenerateRefreshToken(%v)", tt.args.email)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GenerateRefreshToken(%v)", tt.args.email)
		})
	}
}

func TestHandler_Refresh(t *testing.T) {
	type fields struct {
		userRepo user.Repository
		tokener  tokener
	}
	type args struct {
		token string
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     AuthenticationToken
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Error on token parse",
			func() fields {
				repository := user.MockRepository{}
				tokener := mockTokener{}
				tokener.On("parseToken", "testToken").Return(nil, fmt.Errorf("noValidToken"))

				return fields{
					userRepo: &repository,
					tokener:  &tokener,
				}
			},
			args{
				token: "testToken",
			},
			AuthenticationToken{},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "noValidToken", err.Error(), i)
				return true
			},
		},
		{
			"Error on token generation",
			func() fields {
				repository := user.MockRepository{}
				tokener := mockTokener{}
				claims := JWTClaim{
					Email: "testEmail",
				}
				tokener.On("parseToken", "testToken").Return(&claims, nil)
				tokener.On("generateToken", "testEmail", mock.IsType(time.Time{})).Return("", fmt.Errorf("noToken"))

				return fields{
					userRepo: &repository,
					tokener:  &tokener,
				}
			},
			args{
				token: "testToken",
			},
			AuthenticationToken{},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "noToken", err.Error(), i)
				return true
			},
		},
		{
			"Positive Case",
			func() fields {
				repository := user.MockRepository{}
				tokener := mockTokener{}
				claims := JWTClaim{
					Email: "testEmail",
				}
				tokener.On("parseToken", "testToken").Return(&claims, nil)
				tokener.On("generateToken", "testEmail", mock.IsType(time.Time{})).Return("validToken", nil)

				return fields{
					userRepo: &repository,
					tokener:  &tokener,
				}
			},
			args{
				token: "testToken",
			},
			AuthenticationToken{
				Token: "validToken",
				Type:  authType,
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.fieldsFn()
			h := Handler{
				userRepo: args.userRepo,
				tokener:  args.tokener,
			}
			got, err := h.Refresh(tt.args.token)
			if !tt.wantErr(t, err, fmt.Sprintf("Refresh(%v)", tt.args.token)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Refresh(%v)", tt.args.token)
		})
	}
}

func TestHandler_Middleware(t *testing.T) {
	type fields struct {
		userRepo user.Repository
		tokener  tokener
	}
	tests := []struct {
		name       string
		fieldsFn   func() fields
		contextFn  func(r *httptest.ResponseRecorder) *gin.Context
		validateFn func(t *testing.T, c *gin.Context, r *httptest.ResponseRecorder, tokener *mockTokener, repo *user.MockRepository)
	}{
		{
			"No Token Is Passed",
			func() fields {
				return fields{}
			},
			func(r *httptest.ResponseRecorder) *gin.Context {
				c, _ := gin.CreateTestContext(r)
				c.Request = &http.Request{}
				return c
			},
			func(t *testing.T, c *gin.Context, r *httptest.ResponseRecorder, tokener *mockTokener, repo *user.MockRepository) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
				assert.True(t, c.IsAborted())
			},
		},
		{
			"Token Is invalid",
			func() fields {
				tokener := mockTokener{}
				tokener.On("parseToken", "testToken").Return(nil, fmt.Errorf("invalidToken"))
				return fields{tokener: &tokener}
			},
			func(r *httptest.ResponseRecorder) *gin.Context {
				c, _ := gin.CreateTestContext(r)
				c.Request = &http.Request{Header: http.Header{"Authorization": {"Bearer testToken"}}}
				return c
			},
			func(t *testing.T, c *gin.Context, r *httptest.ResponseRecorder, tokener *mockTokener, repo *user.MockRepository) {
				tokener.AssertCalled(t, "parseToken", "testToken")
				assert.Equal(t, http.StatusUnauthorized, r.Code)
				assert.True(t, c.IsAborted())
			},
		},
		{
			"User from claims not exist",
			func() fields {
				tokener := mockTokener{}
				claims := &JWTClaim{Email: "testEmail"}
				tokener.On("parseToken", "testToken").Return(claims, nil)

				userRepo := user.MockRepository{}
				userRepo.On("GetByEmail", "testEmail").Return(nil, user.ErrNotFound)
				return fields{tokener: &tokener, userRepo: &userRepo}
			},
			func(r *httptest.ResponseRecorder) *gin.Context {
				c, _ := gin.CreateTestContext(r)
				c.Request = &http.Request{Header: http.Header{"Authorization": {"Bearer testToken"}}}
				return c
			},
			func(t *testing.T, c *gin.Context, r *httptest.ResponseRecorder, tokener *mockTokener, repo *user.MockRepository) {
				tokener.AssertCalled(t, "parseToken", "testToken")
				repo.AssertCalled(t, "GetByEmail", "testEmail")
				assert.Equal(t, http.StatusUnauthorized, r.Code)
				assert.True(t, c.IsAborted())
			},
		},
		{
			"Positive case",
			func() fields {
				tokener := mockTokener{}
				claims := &JWTClaim{Email: "test@email.com"}
				tokener.On("parseToken", "testToken").Return(claims, nil)

				userRepo := user.MockRepository{}
				userRepo.On("GetByEmail", "test@email.com").Return(user.NewUser("test", "test@email.com", "12345678", user.Admin))
				return fields{tokener: &tokener, userRepo: &userRepo}
			},
			func(r *httptest.ResponseRecorder) *gin.Context {
				c, _ := gin.CreateTestContext(r)
				c.Request = &http.Request{Header: http.Header{"Authorization": {"Bearer testToken"}}}
				return c
			},
			func(t *testing.T, c *gin.Context, r *httptest.ResponseRecorder, tokener *mockTokener, repo *user.MockRepository) {
				tokener.AssertCalled(t, "parseToken", "testToken")
				repo.AssertCalled(t, "GetByEmail", "test@email.com")

				usr, exist := c.Get(userContextKey)
				assert.True(t, exist)

				authUsr, _ := usr.(User)

				assert.Equal(t, "test@email.com", authUsr.Email)
				assert.Equal(t, user.Admin, authUsr.Role)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsFn()
			handler := Handler{
				userRepo: fields.userRepo,
				tokener:  fields.tokener,
			}
			responseRecorder := httptest.NewRecorder()
			context := tt.contextFn(responseRecorder)
			handler.Middleware()(context)

			mockTokener, _ := fields.tokener.(*mockTokener)
			mockRepo, _ := fields.userRepo.(*user.MockRepository)
			tt.validateFn(t, context, responseRecorder, mockTokener, mockRepo)
		})
	}
}

func TestHandler_tokenFromHeader(t *testing.T) {
	type fields struct {
		userRepo user.Repository
		tokener  tokener
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		argsFn func() args
		want   string
	}{
		{
			"No auth header set",
			fields{},
			func() args {
				return args{r: &http.Request{}}
			},
			"",
		},
		{
			"Too short auth header",
			fields{},
			func() args {
				r := http.Request{Header: http.Header{"Authorization": {"12345678"}}}
				return args{r: &r}
			},
			"",
		},
		{
			"Invalid auth type",
			fields{},
			func() args {
				r := http.Request{Header: http.Header{"Authorization": {"buarer 234234"}}}
				return args{r: &r}
			},
			"",
		},
		{
			"Positive Case",
			fields{},
			func() args {
				r := http.Request{Header: http.Header{"Authorization": {"Bearer authToken"}}}
				return args{r: &r}
			},
			"authToken",
		},
		{
			"Positive Case: check case sensitive",
			fields{},
			func() args {
				r := http.Request{Header: http.Header{"Authorization": {"BeAreR authToken"}}}
				return args{r: &r}
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{
				userRepo: tt.fields.userRepo,
				tokener:  tt.fields.tokener,
			}
			assert.Equalf(t, tt.want, h.tokenFromHeader(tt.argsFn().r), "tokenFromHeader(%v)", tt.argsFn().r)
		})
	}
}

func TestHandler_UserFromContext(t *testing.T) {
	type fields struct {
		userRepo user.Repository
		tokener  tokener
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name    string
		fields  fields
		argsFn  func() args
		want    User
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"User is not set",
			fields{},
			func() args {
				return args{c: &gin.Context{}}
			},
			User{},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not get authorized user", err.Error(), i)
				return true
			},
		},
		{
			"Under user key set different struct",
			fields{},
			func() args {
				context := &gin.Context{}
				context.Set(userContextKey, interface{}("test"))
				return args{c: context}
			},
			User{},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not cast authorized user", err.Error(), i)
				return true
			},
		},
		{
			"Positive case",
			fields{},
			func() args {
				context := &gin.Context{}
				context.Set(userContextKey, User{ID: "testId", Email: "testEmail", Role: user.Admin})
				return args{c: context}
			},
			User{ID: "testId", Email: "testEmail", Role: user.Admin},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{
				userRepo: tt.fields.userRepo,
				tokener:  tt.fields.tokener,
			}
			got, err := h.UserFromContext(tt.argsFn().c)
			if !tt.wantErr(t, err, fmt.Sprintf("UserFromContext(%v)", tt.argsFn().c)) {
				return
			}
			assert.Equalf(t, tt.want, got, "UserFromContext(%v)", tt.argsFn().c)
		})
	}
}

func TestHandler_AdminMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		contextFn  func(r *httptest.ResponseRecorder) *gin.Context
		validateFn func(t *testing.T, c *gin.Context, r *httptest.ResponseRecorder)
	}{
		{
			"User is not set",
			func(r *httptest.ResponseRecorder) *gin.Context {
				c, _ := gin.CreateTestContext(r)
				c.Request = &http.Request{}
				return c
			},
			func(t *testing.T, c *gin.Context, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
				assert.True(t, c.IsAborted())
			},
		},
		{
			"User is not admin",
			func(r *httptest.ResponseRecorder) *gin.Context {
				c, _ := gin.CreateTestContext(r)
				c.Set(userContextKey, User{
					ID:    "id",
					Email: "email",
					Role:  user.Author,
				})
				c.Request = &http.Request{}
				return c
			},
			func(t *testing.T, c *gin.Context, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
				assert.True(t, c.IsAborted())
			},
		},
		{
			"User is admin",
			func(r *httptest.ResponseRecorder) *gin.Context {
				c, _ := gin.CreateTestContext(r)
				c.Set(userContextKey, User{
					ID:    "id",
					Email: "email",
					Role:  user.Admin,
				})
				c.Request = &http.Request{}
				return c
			},
			func(t *testing.T, c *gin.Context, r *httptest.ResponseRecorder) {
				assert.False(t, r.Flushed)
				assert.False(t, c.IsAborted())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Handler{}
			responseRecorder := httptest.NewRecorder()
			context := tt.contextFn(responseRecorder)
			handler.AdminMiddleware()(context)
			tt.validateFn(t, context, responseRecorder)
		})
	}
}
