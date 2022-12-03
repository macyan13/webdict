package app

import (
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AddTranslation    command.AddTranslationHandler
	UpdateTranslation command.UpdateTranslationHandler
	DeleteTranslation command.DeleteTranslationHandler

	AddTag    command.AddTagHandler
	UpdateTag command.UpdateTagHandler
	DeleteTag command.DeleteTagHandler

	AddUser command.AddUserHandler
}

type Queries struct {
	SingleTranslation query.SingleTranslationHandler
	LastTranslations  query.LastTranslationsHandler

	SingleTag query.SingleTagHandler
	AllTags   query.AllTagsHandler
}