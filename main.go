package main

import (
	"flag"
	"log"
	"oncall/gcal"
	"strings"
	"time"
)

func main() {
	start, end, onCallPersons := flagsParams()
	googleCalendar, err := gcal.InitGoogleCalendar("Support plan")
	if err != nil {
		log.Fatalln(err)
	}
	supportPlan, err := gcal.InitSupportPlan(start, end, &googleCalendar)
	if err != nil {
		log.Fatalln(err)
	}
	err = supportPlan.CreateSchedule(onCallPersons)
	if err != nil {
		log.Fatalln(err)
	}

}

func flagsParams() (time.Time, time.Time, []string) {
	var now bool
	var nbDays time.Duration
	var onCalls string
	var startStr string
	flag.BoolVar(&now, "now", true, "start should start now. Should not be used with start param")
	flag.DurationVar(&nbDays, "nbDays", 24*time.Hour*365, "number days for the support plan from ")
	flag.StringVar(&onCalls, "onCalls", "alice,bob,carol,dave", "onCalls persons separated by comma")
	flag.StringVar(&startStr, "start", "", "start date with this layout 2006-01-02 . Should not be used with now param")
	if startStr != "" && now {
		log.Fatalln("now and start could be used at the same time")
	}
	flag.Parse()

	var start time.Time
	if now {
		start = time.Now()
	} else {
		var err error
		start, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			log.Fatalln("could not  parse correctly start parameter. should be with this layout: 2006-01-02")
		}
	}
	end := start.Add(nbDays)

	onCallPersons := strings.Split(onCalls, ",")
	return start, end, onCallPersons
}
