package support

import (
	"errors"
	"strings"

	"github.com/omarhachach/bear"
	"gorm.io/gorm"

	"github.com/omarhachach/acnv-support/modules/support/model"
)

// ReplyCommand handles a reply to a support ticket.
type ReplyCommand struct {
	Module *Ticket
}

// GetCallers will associate the callers to this command.
func (r *ReplyCommand) GetCallers() []string {
	return []string{
		"-reply",
	}
}

// GetHandler returns the handler to this command.
func (r *ReplyCommand) GetHandler() func(*bear.Context) {
	return func(ctx *bear.Context) {
		cmdSplit := strings.Split(ctx.Message.Content, " ")
		if len(cmdSplit) < 3 {
			_, _ = ctx.SendErrorMessage("Error sending message")
			return
		}

		caseId := cmdSplit[1]
		replyMsg := strings.Join(cmdSplit[2:], " ")

		ticket := &model.Ticket{
			Model: model.Model{ID: caseId},
		}

		err := r.Module.DB.First(&ticket).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_, _ = ctx.SendErrorMessage("Case ID %s was not found", caseId)
			return
		}

		ctx.SendSuccessMessage("Successfully sent reply!")

		replyCtx := &bear.Context{
			Log:       ctx.Log,
			Session:   ctx.Session,
			ChannelID: ticket.ChannelID,
			Message:   ctx.Message,
		}

		_, _ = replyCtx.SendMessage(bear.InfoColor, "Reply From Staff", replyMsg)
	}
}
