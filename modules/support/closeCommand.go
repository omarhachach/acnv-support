package support

import (
	"errors"
	"strings"

	"github.com/omarhachach/bear"
	"gorm.io/gorm"

	"github.com/omarhachach/acnv-support/modules/support/model"
)

// CloseCommand will close a support ticket.
type CloseCommand struct {
	Module *Ticket
}

// GetCallers will associate the callers to this command.
func (r *CloseCommand) GetCallers() []string {
	return []string{
		"-close",
	}
}

// GetHandler returns the handler to this command.
func (r *CloseCommand) GetHandler() func(*bear.Context) {
	return func(ctx *bear.Context) {
		cmdSplit := strings.Split(ctx.Message.Content, " ")
		if len(cmdSplit) < 2 || len(cmdSplit) > 3 {
			_, _ = ctx.SendErrorMessage("Error sending message.")
			return
		}

		caseId := cmdSplit[1]

		ticket := &model.Ticket{
			Model: model.Model{ID: caseId},
		}

		err := r.Module.DB.Delete(&ticket).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_, _ = ctx.SendErrorMessage("Case ID %s was not found", caseId)

			return
		}

		ctx.SendSuccessMessage("Successfully closed Case %s.", caseId)

		ctx.ChannelID = ticket.ChannelID

		ctx.SendInfoMessage("Your support ticket has been resolved.")
	}
}
