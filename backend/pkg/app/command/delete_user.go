package command

import (
	"errors"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
)

type DeleteUser struct {
	AuthorID string
}

type DeleteUserHandler struct {
	userRepo        user.Repository
	langRepo        lang.Repository
	tagRepo         tag.Repository
	translationRepo translation.Repository
}

func NewDeleteUserHandler(userRepo user.Repository, langRepo lang.Repository, tagRepo tag.Repository, translationRepo translation.Repository) DeleteUserHandler {
	return DeleteUserHandler{userRepo: userRepo, langRepo: langRepo, tagRepo: tagRepo, translationRepo: translationRepo}
}

// Handle removes user and all related content, no transaction support so far
func (h *DeleteUserHandler) Handle(cmd DeleteUser) (int, error) {
	var err error
	userCount, err1 := h.userRepo.Delete(cmd.AuthorID)
	err = errors.Join(err, err1)

	tagCount, err2 := h.tagRepo.DeleteByAuthorID(cmd.AuthorID)
	err = errors.Join(err, err2)

	LangCount, err3 := h.langRepo.DeleteByAuthorID(cmd.AuthorID)
	err = errors.Join(err, err3)

	translationCount, err4 := h.translationRepo.DeleteByAuthorID(cmd.AuthorID)
	err = errors.Join(err, err4)

	return userCount + tagCount + LangCount + translationCount, err
}
