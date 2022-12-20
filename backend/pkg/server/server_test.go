package server

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/auth"
	"github.com/macyan13/webdict/backend/pkg/storage/inmemory"
	"time"
)

func initTestServer() *HttpServer {
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
	translationRepo := inmemory.NewTranslationRepository(*tagRepo)
	userRepo := inmemory.NewUserRepository()

	cipher := auth.Cipher{}

	cmd := app.Commands{
		AddTranslation:    command.NewAddTranslationHandler(translationRepo, tagRepo),
		UpdateTranslation: command.NewUpdateTranslationHandler(translationRepo, tagRepo),
		DeleteTranslation: command.NewDeleteTranslationHandler(translationRepo, tagRepo),
		AddTag:            command.NewAddTagHandler(tagRepo),
		UpdateTag:         command.NewUpdateTagHandler(tagRepo),
		DeleteTag:         command.NewDeleteTagHandler(tagRepo),
		AddUser:           command.NewAddUserHandler(userRepo, cipher),
	}

	queries := app.Queries{
		SingleTranslation: query.NewSingleTranslationHandler(translationRepo),
		LastTranslations:  query.NewLastTranslationsHandler(translationRepo),
		SingleTag:         query.NewSingleTagHandler(tagRepo),
		AllTags:           query.NewAllTagsHandler(tagRepo),
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

	s := HttpServer{
		engine:      router,
		app:         &application,
		authHandler: authHandler,
		opts:        opts,
		userRepo:    userRepo,
	}

	s.BuildRoutes()
	s.PopulateInitData()
	return &s
}
