package support

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/omarhachach/bear"
	"github.com/sirupsen/logrus"
)

func onDirectMessageCreate(log *logrus.Logger, m *Module) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(session *discordgo.Session, messageCreate *discordgo.MessageCreate) {
		if messageCreate.Author.ID == session.State.User.ID {
			return
		}

		if messageCreate.GuildID != "" {
			return
		}

		content := messageCreate.ContentWithMentionsReplaced()

		entry, err := m.AddEntryToSupportFile(Entry{
			MessageContent: content,
			ChannelID:      messageCreate.ChannelID,
			SenderID:       messageCreate.Author.ID,
		})
		if err != nil {
			log.WithError(err).Error("Error adding entry to support file.")
			return
		}

		ctx := bear.Context{
			Log:       log,
			Session:   session,
			ChannelID: messageCreate.ChannelID,
			Message:   messageCreate,
		}

		_, _ = ctx.SendSuccessMessage("Thanks for your support message!")

		ctx.ChannelID = m.SupportChannelID

		_, _ = ctx.SendMessage(bear.InfoColor, fmt.Sprintf("Support Message [%s]", entry.ID), "<@%s> %s", entry.SenderID, entry.MessageContent)

		for _, attachment := range messageCreate.Attachments {
			_, _ = session.ChannelMessageSendEmbed(m.SupportChannelID, &discordgo.MessageEmbed{
				URL:         "",
				Type:        "image",
				Title:       "",
				Description: "",
				Timestamp:   "",
				Color:       0,
				Footer:      nil,
				Image: &discordgo.MessageEmbedImage{
					URL:      attachment.URL,
					ProxyURL: attachment.ProxyURL,
					Width:    attachment.Width,
					Height:   attachment.Height,
				},
				Thumbnail: nil,
				Video:     nil,
				Provider:  nil,
				Author:    nil,
				Fields:    nil,
			})
		}
	}
}
