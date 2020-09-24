package support

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/omarhachach/bear"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"omarh.net/acnv-support/modules/support/model"
)

// OnDirectMessageCreate is an event handler for *discordgo.MessageCreate but will check whether it is a direct message.
// It will create a new support ticket if there is no current support ticket associated with the user.
func OnDirectMessageCreate(log *logrus.Logger, module *Ticket) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(session *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == session.State.User.ID {
			return
		}

		if m.GuildID != "" {
			return
		}

		ticket := &model.Ticket{}

		err := module.DB.Where("sender_id = ?", m.Author.ID).First(&ticket).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ticket.SenderID = m.Author.ID
				ticket.ChannelID = m.ChannelID

				err = module.DB.Create(&ticket).Error
				if err != nil {
					log.WithError(err).Error("Error while saving ticket")
					return
				}

				handleSupportMessage(log, session, m, module, ticket)

				return
			}

			log.WithError(err).Error("Error while retrieving ticket.")
			return
		}

		handleSupportMessage(log, session, m, module, ticket)
	}
}

func handleSupportMessage(log *logrus.Logger, session *discordgo.Session, m *discordgo.MessageCreate, module *Ticket, ticket *model.Ticket) {
	ctx := bear.Context{
		Log:       log,
		Session:   session,
		ChannelID: m.ChannelID,
		Message:   m,
	}

	ctx.SendSuccessMessage("Thanks for your support message!")

	ctx.ChannelID = module.SupportChannelID

	content := m.ContentWithMentionsReplaced()
	ctx.SendMessage(bear.InfoColor, fmt.Sprintf("Support Message [%s]", ticket.ID), "<@%s> %s", m.Author.ID, content)

	for _, attachment := range m.Attachments {
		_, _ = session.ChannelMessageSendEmbed(module.SupportChannelID, &discordgo.MessageEmbed{
			Type: "image",
			Title: fmt.Sprintf("Support Message [%s]", ticket.ID),
			Image: &discordgo.MessageEmbedImage{
				URL:      attachment.URL,
				ProxyURL: attachment.ProxyURL,
				Width:    attachment.Width,
				Height:   attachment.Height,
			},
		})
	}

	return
}
