package support

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"

	"github.com/omarhachach/bear"
)

type Module struct {
	SupportChannelID string
	sync.Mutex
}

func (m *Module) GetName() string {
	return "Support"
}

func (m *Module) GetDesc() string {
	return "This allows users to send support tickets to the staff team."
}

func (m *Module) GetCommands() []bear.Command {
	return []bear.Command{
		&ReplyCommand{Module: m},
	}
}

func (m *Module) GetVersion() string {
	return "1.0.0"
}

func (m *Module) Init(b *bear.Bear) {
	_, err := os.Open("support-file.json")
	if errors.Is(err, os.ErrNotExist) {
		err := ioutil.WriteFile("support-file.json", []byte("{}"), 0755)
		if err != nil {
			b.Log.WithError(err).Fatal("Error starting module, failed while create file support-file.json")
			return
		}
	}

	b.AddHandler(onDirectMessageCreate(b.Log, m))
}

func (*Module) Close(*bear.Bear) {
}
