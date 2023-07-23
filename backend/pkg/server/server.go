package server

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/auth"
	"github.com/macyan13/webdict/backend/pkg/store/cache"
	"github.com/macyan13/webdict/backend/pkg/store/mongo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

type HTTPServer struct {
	engine      *gin.Engine
	app         *app.Application
	authHandler *auth.Handler
	opts        Opts
	userRepo    user.Repository // todo: need only for tests, remove after query for user
}

func InitServer(opts Opts) (*HTTPServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Printf("[WARN] interrupt signal")
		cancel()
	}()

	dbConnect, err := mongo.InitDatabase(ctx, mongo.Opts{
		Database: opts.Mongo.Database,
		Host:     opts.Mongo.Host,
		Port:     opts.Mongo.Port,
		Username: opts.Mongo.Username,
		Passwd:   opts.Mongo.Passwd,
	})
	if err != nil {
		return nil, err
	}

	cacheOpts := cache.Opts{TagCacheTTL: opts.Cache.TagCacheTTL, TranslationCacheTTL: opts.Cache.TranslationCacheTTL, LangCacheTTL: opts.Cache.LangCacheTTL}

	tagRepo, err := mongo.NewTagRepo(dbConnect)
	if err != nil {
		return nil, err
	}

	cachedTagRepo := cache.NewTagRepo(ctx, tagRepo, tagRepo, cacheOpts)

	langRepo, err := mongo.NewLangRepo(dbConnect)
	if err != nil {
		return nil, err
	}

	cachedLangRepo := cache.NewLangRepo(ctx, langRepo, langRepo, cacheOpts)

	translationRepo, err := mongo.NewTranslationRepo(dbConnect, cachedTagRepo, cachedLangRepo)
	if err != nil {
		return nil, err
	}

	userRepo, err := mongo.NewUserRepo(dbConnect, cachedLangRepo, query.NewRoleMapper())
	if err != nil {
		return nil, err
	}

	cipher := auth.Cipher{}

	cachedTranslationRepo := cache.NewTranslationRepo(ctx, translationRepo, translationRepo, cacheOpts.TranslationCacheTTL)

	cmd := app.Commands{
		AddTranslation:    command.NewAddTranslationHandler(cachedTranslationRepo, cachedTagRepo, cachedLangRepo),
		UpdateTranslation: command.NewUpdateTranslationHandler(cachedTranslationRepo, cachedTagRepo, cachedLangRepo),
		DeleteTranslation: command.NewDeleteTranslationHandler(cachedTranslationRepo),
		AddTag:            command.NewAddTagHandler(cachedTagRepo),
		UpdateTag:         command.NewUpdateTagHandler(cachedTagRepo),
		DeleteTag:         command.NewDeleteTagHandler(cachedTagRepo, cachedTranslationRepo),
		AddUser:           command.NewAddUserHandler(userRepo, cipher),
		UpdateUser:        command.NewUpdateUserHandler(userRepo, cipher),
		DeleteUser:        command.NewDeleteUserHandler(userRepo, cachedLangRepo, cachedTagRepo, cachedTranslationRepo),
		AddLang:           command.NewAddLangHandler(cachedLangRepo),
		UpdateLang:        command.NewUpdateLangHandler(cachedLangRepo),
		DeleteLang:        command.NewDeleteLangHandler(cachedLangRepo, cachedTranslationRepo),
		UpdateProfile:     command.NewUpdateProfileHandler(userRepo, cipher, cachedLangRepo),
	}

	queries := app.Queries{
		SingleTranslation:  query.NewSingleTranslationHandler(cachedTranslationRepo),
		LastTranslations:   query.NewLastTranslationsHandler(cachedTranslationRepo),
		RandomTranslations: query.NewRandomTranslationsHandler(cachedTranslationRepo),
		SingleTag:          query.NewSingleTagHandler(cachedTagRepo),
		AllTags:            query.NewAllTagsHandler(cachedTagRepo),
		SingleUser:         query.NewSingleUserHandler(userRepo),
		AllUsers:           query.NewAllUsersHandler(userRepo),
		SingleLang:         query.NewSingleLangHandler(cachedLangRepo),
		AllLangs:           query.NewAllLangsHandler(cachedLangRepo),
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
	router.Use(cors.Default())

	s := HTTPServer{
		engine:      router,
		app:         &application,
		authHandler: authHandler,
		opts:        opts,
	}

	s.buildRoutes()
	s.loadStaticData()
	s.populateInitData()
	return &s, nil
}

func (s *HTTPServer) Run() error {
	err := s.engine.Run(fmt.Sprintf(":%d", s.opts.Port))

	if err != nil {
		log.Printf("HTTPServer - there was an error calling Run on engine: %v", err)
		return err
	}

	return nil
}

func (s *HTTPServer) loadStaticData() {
	s.engine.LoadHTMLGlob("./public/*.html")
	s.engine.StaticFile("favicon.ico", "./public/favicon.ico")
	s.engine.Static("/static", "./public/static")
}

func (s *HTTPServer) populateInitData() {
	if _, err := s.app.Commands.AddUser.Handle(command.AddUser{
		Name:     "admin",
		Email:    s.opts.Admin.AdminEmail,
		Password: s.opts.Admin.AdminPasswd,
		Role:     user.Admin,
	}); err != nil {
		log.Printf("[WARN] can not create admin user - %s", err)
	}
}

func (s *HTTPServer) unauthorized(c *gin.Context, err error) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	log.Printf("[ERROR] Can not authorize action: %s:%s: %v", frame.File, frame.Function, err)
	c.JSON(http.StatusUnauthorized, nil)
}

func (s *HTTPServer) badRequest(c *gin.Context, err error) {
	log.Printf("[ERROR] Can not handle request - %v", err)
	c.JSON(http.StatusBadRequest, err.Error())
}
