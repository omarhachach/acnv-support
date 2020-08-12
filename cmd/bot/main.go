package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/omarhachach/bear"

	"omarh.net/acnv-support/modules/support"
)

func main() {
	confFile, err := os.Open("config.json")
	if err != nil {
		panic("Couldn't open file: " + err.Error())
	}

	config := &Config{}
	confFileBytes, _ := ioutil.ReadAll(confFile)

	err = json.Unmarshal(confFileBytes, &config)
	if err != nil {
		panic("error reading config: " + err.Error())
	}

	b := bear.New(config.Config).RegisterModules(&support.Module{
		SupportChannelID: config.SupportChannelID,
	}).Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	b.Close()
}
