package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/macyan13/webdict/backend/pkg/app"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/auth"
	"github.com/macyan13/webdict/backend/pkg/store/inmemory"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

type testHTTPServer struct {
	*HTTPServer
	userRepo user.Repository
}

func initTestServer() *testHTTPServer {
	authGroup := AuthGroup{}
	authGroup.TTL.Auth = time.Minute * 10
	authGroup.TTL.Refresh = time.Minute * 10
	authGroup.TTL.Cookie = time.Hour

	opts := Opts{
		Auth: authGroup,
		Admin: AdminGroup{
			AdminPasswd: "test_password",
			AdminEmail:  "test@email.com",
		},
		Port:       4000,
		WebdictURL: "",
		Dbg:        false,
	}

	tagRepo := inmemory.NewTagRepository()
	langRepo := inmemory.NewLangRepository()
	translationRepo := inmemory.NewTranslationRepository(*tagRepo, *langRepo)
	userRepo := inmemory.NewUserRepository(query.NewRoleMapper())

	cipher := auth.Cipher{}
	cmd := app.Commands{
		AddTranslation:    command.NewAddTranslationHandler(translationRepo, tagRepo, langRepo),
		UpdateTranslation: command.NewUpdateTranslationHandler(translationRepo, tagRepo, langRepo),
		DeleteTranslation: command.NewDeleteTranslationHandler(translationRepo),
		AddTag:            command.NewAddTagHandler(tagRepo),
		UpdateTag:         command.NewUpdateTagHandler(tagRepo),
		DeleteTag:         command.NewDeleteTagHandler(tagRepo, translationRepo),
		AddUser:           command.NewAddUserHandler(userRepo, cipher),
		UpdateUser:        command.NewUpdateUserHandler(userRepo, cipher),
		DeleteUser:        command.NewDeleteUserHandler(userRepo, langRepo, tagRepo, translationRepo),
		AddLang:           command.NewAddLangHandler(langRepo),
		UpdateLang:        command.NewUpdateLangHandler(langRepo),
		DeleteLang:        command.NewDeleteLangHandler(langRepo, translationRepo),
		UpdateProfile:     command.NewUpdateProfileHandler(userRepo, cipher, langRepo),
	}

	validate := validator.New()

	queries := app.Queries{
		SingleTranslation:  query.NewSingleTranslationHandler(translationRepo, validate),
		SearchTranslations: query.NewSearchTranslationsHandler(translationRepo, validate),
		RandomTranslations: query.NewRandomTranslationsHandler(translationRepo, validate),
		SingleTag:          query.NewSingleTagHandler(tagRepo, validate),
		AllTags:            query.NewAllTagsHandler(tagRepo, validate),
		SingleUser:         query.NewSingleUserHandler(userRepo, validate),
		AllUsers:           query.NewAllUsersHandler(userRepo),
		SingleLang:         query.NewSingleLangHandler(langRepo, validate),
		AllLangs:           query.NewAllLangsHandler(langRepo, validate),
		AllRoles:           query.NewAllRolesHandler(),
	}

	application := app.Application{
		Commands: cmd,
		Queries:  queries,
	}

	authHandler := auth.NewHandler(userRepo, cipher, auth.Params{
		AuthTTL:    opts.Auth.TTL.Auth,
		RefreshTTL: opts.Auth.TTL.Refresh,
		Secret:     opts.Auth.Secret,
	})

	router := gin.Default()

	s := HTTPServer{
		engine:      router,
		app:         &application,
		authHandler: authHandler,
		opts:        opts,
	}

	s.buildRoutes()
	s.populateInitData()
	return &testHTTPServer{HTTPServer: &s, userRepo: userRepo}
}

func setAdminAuthToken(t *testing.T, s *testHTTPServer, r *http.Request) {
	token, err := s.authHandler.Authenticate(s.opts.Admin.AdminEmail, s.opts.Admin.AdminPasswd)
	assert.NoError(t, err)
	r.Header.Set("Authorization", token.Type+" "+token.Token)
}

func setAuthTokenWithCredentials(t *testing.T, s *testHTTPServer, r *http.Request, email, passwd string) {
	token, err := s.authHandler.Authenticate(email, passwd)
	assert.NoError(t, err)
	r.Header.Set("Authorization", token.Type+" "+token.Token)
}
