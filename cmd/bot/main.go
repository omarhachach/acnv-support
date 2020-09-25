package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/omarhachach/bear"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"omarh.net/acnv-support/modules/support"
	"omarh.net/acnv-support/modules/timezones"
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

	db, err := gorm.Open(sqlite.Open(config.DB), &gorm.Config{})
	if err != nil {
		panic("error opening db: " + err.Error())
	}

	b := bear.New(config.Config).RegisterModules(
		&support.Ticket{
			SupportChannelID: config.SupportChannelID,
			DB:               db,
		},
		&timezones.TimeZone{DB: db},
	).Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	b.Close()
}
