package timezones

import (
	"github.com/omarhachach/bear"
	"gorm.io/gorm"

	"omarh.net/acnv-support/modules/timezones/model"
)

// TimeZone handles interactions for timezones.
type TimeZone struct {
	DB *gorm.DB
}

func (t *TimeZone) GetName() string {
	return "Time Zones"
}

func (t *TimeZone) GetDesc() string {
	return "This module handles interactions for timezones for different users."
}

func (t *TimeZone) GetCommands() []bear.Command {
	return []bear.Command{
		&TimezoneCommand{},
		&SetTimeCommand{DB: t.DB},
		&TimeCommand{DB: t.DB},
	}
}

func (t *TimeZone) GetVersion() string {
	return "1.0.0"
}

func (t *TimeZone) Init(b *bear.Bear) {
	err := t.DB.AutoMigrate(
		&model.UserTimezone{},
	)
	if err != nil {
		b.Log.WithError(err).Fatal("Error migrating database.")

		return
	}
}

func (t *TimeZone) Close(*bear.Bear) {
}
