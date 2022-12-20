package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/auth"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"github.com/macyan13/webdict/backend/pkg/storage/mongo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

type HttpServer struct {
	engine      *gin.Engine
	app         *app.Application
	authHandler *auth.Handler
	opts        Opts
	userRepo    user.Repository // todo: need only for tests, remove after query for user
}

func InitServer(opts Opts) (*HttpServer, error) {
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

	tagRepo, err := mongo.NewTagRepo(ctx, dbConnect)
	if err != nil {
		return nil, err
	}

	translationRepo, err := mongo.NewTranslationRepo(ctx, dbConnect, tagRepo)
	if err != nil {
		return nil, err
	}

	userRepo, err := mongo.NewUserRepo(ctx, dbConnect)
	if err != nil {
		return nil, err
	}

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
	// 	"github.com/gin-contrib/cors"
	// router.Use(cors.Default()) - middleware for CORS support, maybe add later

	s := HttpServer{
		engine:      router,
		app:         &application,
		authHandler: authHandler,
		opts:        opts,
	}

	s.BuildRoutes()
	s.PopulateInitData()
	return &s, nil
}

func (s *HttpServer) Run() error {
	err := s.engine.Run(fmt.Sprintf(":%d", s.opts.Port))

	if err != nil {
		log.Printf("HttpServer - there was an error calling Run on engine: %v", err)
		return err
	}

	return nil
}

func (s *HttpServer) PopulateInitData() {
	if err := s.app.Commands.AddUser.Handle(command.AddUser{
		Name:     "admin",
		Email:    s.opts.Admin.AdminEmail,
		Password: s.opts.Admin.AdminPasswd,
	}); err != nil {
		log.Printf("[WARN] can not create admin user - %s", err)
	}
}

func (s *HttpServer) unauthorised(c *gin.Context, err error) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	log.Printf("[Error] Can not authorise action: %s:%s: %v", frame.File, frame.Function, err)
	c.JSON(http.StatusUnauthorized, nil)
}

func (s *HttpServer) badRequest(c *gin.Context, err error) {
	log.Printf("[Error] Can not handle request - %v", err)
	c.JSON(http.StatusBadRequest, nil)
}
