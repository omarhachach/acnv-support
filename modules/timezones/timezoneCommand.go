package timezones

import (
	"strings"
	"time"

	"github.com/omarhachach/bear"
)

// TimezoneCommand will handle convertion of time from TimeZone to another.
type TimezoneCommand struct {
}

func (t *TimezoneCommand) GetCallers() []string {
	return []string{
		"-timezone",
		"-tz",
		"-time-zone",
	}
}

// -time 1am MST Europe/Copenhagen
func (t *TimezoneCommand) GetHandler() func(*bear.Context) {
	return func(ctx *bear.Context) {
		cmdSplit := strings.Split(ctx.Message.Content, " ")
		if len(cmdSplit) == 2 {
			loc, err := time.LoadLocation(cmdSplit[1])
			if err != nil {
				ctx.SendErrorMessage("%s is an invalid timezone.", cmdSplit[2])
				return
			}

			ctx.SendSuccessMessage("Time in %s is %s", cmdSplit[1], time.Now().In(loc).Format(time.Kitchen))
			return
		}

		if len(cmdSplit) == 4 {
			loc, err := time.LoadLocation(cmdSplit[2])
			if err != nil {
				ctx.SendErrorMessage("%s is an invalid timezone.", cmdSplit[2])
				return
			}

			loc2, err := time.LoadLocation(cmdSplit[3])
			if err != nil {
				ctx.SendErrorMessage("%s is an invalid timezone.", cmdSplit[3])
				return
			}

			localTime := strings.ToUpper(cmdSplit[1])
			if len(localTime) == 3 {
				localTime = "0" + localTime
			}

			parsedTime, err := time.ParseInLocation("03PM", localTime, loc)
			if err != nil {
				ctx.SendErrorMessage("Make sure to specify time with AM/PM, fx. 3PM")
				return
			}

			ctx.SendSuccessMessage("%s in %s is %s in %s", localTime, cmdSplit[2], parsedTime.In(loc2).Format("03PM"), cmdSplit[3])
			return
		}

		ctx.SendErrorMessage("Invalid usage: -time <timezone> or -time <time> <timezone> <timezone>")
	}
}

