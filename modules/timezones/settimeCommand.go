package timezones

import (
	"errors"
	"strings"
	"time"

	"github.com/omarhachach/bear"
	"gorm.io/gorm"

	"omarh.net/acnv-support/modules/timezones/model"
)

// SetTimeCommand will set the time for a user.
type SetTimeCommand struct {
	DB *gorm.DB
}

func (t *SetTimeCommand) GetCallers() []string {
	return []string{
		"-set-time",
		"-settime",
		"-st",
	}
}

func (t *SetTimeCommand) GetHandler() func(*bear.Context) {
	return func(ctx *bear.Context) {
		cmdSplit := strings.Split(ctx.Message.Content, " ")
		loc, err := time.LoadLocation(cmdSplit[1])
		if err != nil {
			ctx.SendErrorMessage("%s is an invalid timezone.", cmdSplit[2])
			return
		}

		user := &model.UserTimezone{ID: ctx.Message.Author.ID, TimeZone: loc.String()}

		err = t.DB.First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = t.DB.Create(&user).Error
			if err != nil {
				ctx.SendErrorMessage("An error occured.")
				ctx.Log.WithError(err).Error("Error creating new UserTimeZone.")

				return
			}

			ctx.SendSuccessMessage("Sucessfully saved your timezone.")
			return
		}

		err = t.DB.Save(user).Error
		if err != nil {
			ctx.SendErrorMessage("An error occured.")
			ctx.Log.WithError(err).Error("Error update UserTimeZone.")

			return
		}

		ctx.SendSuccessMessage("Successfully updated your timezone")
	}
}

