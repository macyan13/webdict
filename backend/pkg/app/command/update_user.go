package command

import (
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
)

// UpdateUser update existing user cmd
type UpdateUser struct {
	Name       string
	Email      string
	Password   string
	Role       user.Role
	IsAdminCMD bool
}

// UpdateUserHandler update existing user cmd handler
type UpdateUserHandler struct {
	tagRepo tag.Repository
}

//
// func NewUpdateTagHandler(tagRepo tag.Repository) UpdateTagHandler {
// 	return UpdateTagHandler{tagRepo: tagRepo}
// }
//
// // Handle applies cmd changes to tag and saves it to DB
// func (h UpdateTagHandler) Handle(cmd UpdateTag) error {
// 	tg, err := h.tagRepo.Get(cmd.TagID, cmd.AuthorID)
//
// 	if err != nil {
// 		return err
// 	}
//
// 	if err := tg.ApplyChanges(cmd.Tag); err != nil {
// 		return err
// 	}
//
// 	return h.tagRepo.Update(tg)
// }
