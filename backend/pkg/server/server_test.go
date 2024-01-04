package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/macyan13/webdict/backend/pkg/app"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/auth"
	"github.com/macyan13/webdict/backend/pkg/store/inmemory"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func initTestServer() *HTTPServer {
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
		SingleTranslation:  query.NewSingleTranslationHandler(translationRepo),
		SearchTranslations: query.NewSearchTranslationsHandler(translationRepo, validate),
		RandomTranslations: query.NewRandomTranslationsHandler(translationRepo),
		SingleTag:          query.NewSingleTagHandler(tagRepo),
		AllTags:            query.NewAllTagsHandler(tagRepo),
		SingleUser:         query.NewSingleUserHandler(userRepo),
		AllUsers:           query.NewAllUsersHandler(userRepo),
		SingleLang:         query.NewSingleLangHandler(langRepo),
		AllLangs:           query.NewAllLangsHandler(langRepo),
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
		userRepo:    userRepo,
	}

	s.buildRoutes()
	s.populateInitData()
	return &s
}

func setAdminAuthToken(t *testing.T, s *HTTPServer, r *http.Request) {
	token, err := s.authHandler.Authenticate(s.opts.Admin.AdminEmail, s.opts.Admin.AdminPasswd)
	assert.NoError(t, err)
	r.Header.Set("Authorization", token.Type+" "+token.Token)
}

func setAuthTokenWithCredentials(t *testing.T, s *HTTPServer, r *http.Request, email, passwd string) {
	token, err := s.authHandler.Authenticate(email, passwd)
	assert.NoError(t, err)
	r.Header.Set("Authorization", token.Type+" "+token.Token)
}
