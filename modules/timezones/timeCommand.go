package timezones

import (
	"errors"
	"strings"
	"time"

	"github.com/omarhachach/bear"
	"gorm.io/gorm"

	"omarh.net/acnv-support/modules/timezones/model"
)

// TimeCommand handles timezone convertion between users.
type TimeCommand struct {
	DB *gorm.DB
}

func (t *TimeCommand) GetCallers() []string {
	return []string{
		"-time",
		"-t",
	}
}

func (t *TimeCommand) GetHandler() func(*bear.Context) {
	return func(ctx *bear.Context) {
		cmdSplit := strings.Split(ctx.Message.Content, " ")
		if len(cmdSplit) == 2 {
			if len(ctx.Message.Mentions) == 0 {
				ctx.SendErrorMessage("You have to mention the user!")
				return
			}

			mentioned := ctx.Message.Mentions[0]
			user := &model.UserTimezone{ID: mentioned.ID}

			err := t.DB.First(&user).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.SendErrorMessage("User has not set a timezone.")
				return
			}

			loc, err := time.LoadLocation(user.TimeZone)
			if err != nil {
				ctx.SendErrorMessage("%s is an invalid timezone.", cmdSplit[2])
				return
			}

			member, err := ctx.Session.GuildMember(ctx.Message.GuildID, mentioned.ID)
			if err != nil {
				ctx.SendErrorMessage("Internal error occured.")
				ctx.Log.WithError(err).Error("Error retrieving guild member.")

				return
			}

			name := member.Nick
			if name == "" {
				name = mentioned.Username
			}
			ctx.SendSuccessMessage("%s's time is %s", name, time.Now().In(loc).Format(time.Kitchen))
		} else if len(cmdSplit) >= 3 || len(cmdSplit) <= 4 {
			if len(ctx.Message.Mentions) == 0 {
				ctx.SendErrorMessage("You have to mention the user!")
				return
			}

			mentioned := ctx.Message.Mentions[0]
			user := &model.UserTimezone{ID: mentioned.ID}

			err := t.DB.First(&user).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.SendErrorMessage("User has not set a timezone.")
				return
			}

			userTz, _ := time.LoadLocation(user.TimeZone)
			authorTz := &time.Location{}

			if len(cmdSplit) == 4 {
				authorTz, err = time.LoadLocation(cmdSplit[3])
				if err != nil {
					ctx.SendErrorMessage("%s is not a valid timezone.", cmdSplit[3])
					return
				}
			} else {
				author := &model.UserTimezone{ID: ctx.Message.Author.ID}

				err = t.DB.First(&author).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					ctx.SendErrorMessage("You have not set a timezone.")
					return
				}

				authorTz, _ = time.LoadLocation(author.TimeZone)
			}

			localTime := strings.ToUpper(cmdSplit[2])
			if len(localTime) == 3 {
				localTime = "0" + localTime
			}

			parsedTime, err := time.ParseInLocation("03PM", localTime, authorTz)
			if err != nil {
				ctx.SendErrorMessage("Make sure to specify time with AM/PM, fx. 3PM")
				return
			}

			member, err := ctx.Session.GuildMember(ctx.Message.GuildID, mentioned.ID)
			if err != nil {
				ctx.SendErrorMessage("Internal error occured.")
				ctx.Log.WithError(err).Error("Error retrieving guild member.")

				return
			}

			name := member.Nick
			if name == "" {
				name = mentioned.Username
			}
			ctx.SendSuccessMessage("%s is %s for %s", localTime, parsedTime.In(userTz).Format("03PM"), name)
		} else {
			ctx.SendErrorMessage("Correct usage is -time <user> or -time <user> <time> [timezone]")
		}
	}
}
