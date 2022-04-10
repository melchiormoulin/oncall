package gcal

import (
	"errors"
	"fmt"
	"time"
)

type SupportPlan struct {
	start    time.Time
	end      time.Time
	calendar ICalendar
}

func InitSupportPlan(start time.Time, end time.Time, calendar ICalendar) (SupportPlan, error) {
	duration := end.Sub(start)
	if duration.Hours() < 24 {
		return SupportPlan{}, errors.New(fmt.Sprintf("Duration between start (%s) and end (%s) date should be more than 1 day to do a support plan", start.GoString(), end.GoString()))
	}
	return SupportPlan{
		start:    start,
		end:      end,
		calendar: calendar,
	}, nil
}

func (supportPlan *SupportPlan) CreateSchedule(members []string) error {
	if members == nil || len(members) < 1 {
		return errors.New("we should have one member at minimum")
	}
	currentDay := supportPlan.start
	for currentDay.Unix() < supportPlan.end.Unix() {
		for _, member := range members {
			if currentDay.Weekday() == time.Saturday {
				for i := 1; i <= 2; i++ {
					if currentDay.Unix() >= supportPlan.end.Unix() {
						break
					}
					err := supportPlan.calendar.CreateDailyEvent(currentDay, member)
					if err != nil {
						return err
					}
					currentDay = currentDay.Add(24 * time.Hour)
				}
			}
			if currentDay.Unix() >= supportPlan.end.Unix() {
				break
			}
			err := supportPlan.calendar.CreateDailyEvent(currentDay, member)
			if err != nil {
				return err
			}
			currentDay = currentDay.Add(24 * time.Hour)
		}
	}
	return nil
}
