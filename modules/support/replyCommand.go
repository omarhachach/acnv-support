package support

import (
	"errors"
	"strings"

	"github.com/omarhachach/bear"
)

type ReplyCommand struct {
	Module *Module
}

func (r *ReplyCommand) GetCallers() []string {
	return []string{
		"-reply",
	}
}

func (r *ReplyCommand) GetHandler() func(*bear.Context) {
	return func(ctx *bear.Context) {
		cmdSplit := strings.Split(ctx.Message.Content, " ")
		if len(cmdSplit) < 3 {
			_, _ = ctx.SendErrorMessage("Error sending message")
			return
		}

		caseId := cmdSplit[1]
		replyMsg := strings.Join(cmdSplit[2:], " ")

		entry, err := r.Module.GetEntryFromSupportFile(caseId)
		if errors.Is(err, ErrNotFound) {
			_, _ = ctx.SendErrorMessage("Case ID %s was not found", caseId)
			return
		}

		replyCtx := &bear.Context{
			Log:       ctx.Log,
			Session:   ctx.Session,
			ChannelID: entry.ChannelID,
			Message:   ctx.Message,
		}

		_, _ = replyCtx.SendMessage(bear.InfoColor, "Reply", replyMsg)
	}
}
