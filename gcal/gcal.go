package gcal

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"oncall/oauth"
	"time"
)

//go:generate mockgen -destination=mocks/mock_calendar.go -package=mocks . ICalendar
type ICalendar interface {
	CreateDailyEvent(day time.Time, member string) error
}

type GoogleCalendar struct {
	calendarService *calendar.Service
	calendar        *calendar.Calendar
}

func CreateGoogleCalendar(srv *calendar.Service, name string) (*calendar.Calendar, error) {
	oncallCalendar, err := srv.Calendars.Insert(&calendar.Calendar{Summary: name}).Do()
	if err != nil {
		return nil, err
	}
	log.Printf("calendar %s with id %s created", oncallCalendar.Id, oncallCalendar.Summary)
	return oncallCalendar, nil
}

func CreateCalendarService() (*calendar.Service, error) {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to read client secret file: %v", err))
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to parse client secret file to config: %v", err))
	}
	client := oauth.GetClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to retrieve Calendar client: %v", err))
	}
	return srv, nil
}

func InitGoogleCalendar(calendarName string) (GoogleCalendar, error) {
	calendarService, err := CreateCalendarService()
	if err != nil {
		return GoogleCalendar{}, err
	}
	calendar, err := CreateGoogleCalendar(calendarService, calendarName)
	if err != nil {
		return GoogleCalendar{}, err
	}
	return GoogleCalendar{
		calendarService: calendarService,
		calendar:        calendar,
	}, nil
}

func (googleCalendar *GoogleCalendar) CreateDailyEvent(day time.Time, member string) error {
	eventDateTime := calendar.EventDateTime{
		Date: day.Format("2006-01-02"),
	}
	event := calendar.Event{Summary: member, Start: &eventDateTime, End: &eventDateTime}
	eventResp, err := googleCalendar.calendarService.Events.Insert(googleCalendar.calendar.Id, &event).Do()
	if err != nil {
		return fmt.Errorf("event creation failed on calendar %s with member %s on %s : %w", googleCalendar.calendar.Summary, member, day.Format("2006-01-02"), err)
	}
	log.Printf("Event created on calendar %s with member %s on %s : %s\n", googleCalendar.calendar.Summary, member, eventDateTime.Date, eventResp.HtmlLink)
	//No batch requests in the sdk https://developers.googleblog.com/2018/03/discontinuing-support-for-json-rpc-and.html
	//No time to do it manually with http or using backoff algorithm
	time.Sleep(time.Millisecond * 200)
	return nil
}
