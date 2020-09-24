package support

import (
	"sync"

	"github.com/omarhachach/bear"
	"gorm.io/gorm"

	"omarh.net/acnv-support/modules/support/model"
)

// Ticket is the module for handling support tickets.
type Ticket struct {
	SupportChannelID string
	DB               *gorm.DB
	sync.Mutex
}

// GetName will return the name of the module for use in the help module.
func (*Ticket) GetName() string {
	return "Tickets"
}

// GetDesc will return the description of the module.
func (*Ticket) GetDesc() string {
	return "This allows users to send support tickets to the staff team."
}

// GetCommands will return the commands associated with this module.
func (t *Ticket) GetCommands() []bear.Command {
	return []bear.Command{
		&ReplyCommand{Module: t},
		&CloseCommand{Module: t},
	}
}

// GetVersion will return the version of the module.
func (*Ticket) GetVersion() string {
	return "1.1.0"
}

// Init will initialize the module.
func (t *Ticket) Init(b *bear.Bear) {
	err := t.DB.AutoMigrate(
		&model.Ticket{},
	)
	if err != nil {
		b.Log.WithError(err).Fatal("Error migrating database.")

		return
	}

	b.AddHandler(OnDirectMessageCreate(b.Log, t))
}

// Close handles the closing of the module.
func (*Ticket) Close(*bear.Bear) {
}
