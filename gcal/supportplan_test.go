package gcal

import (
	"github.com/golang/mock/gomock"
	"oncall/gcal/mocks"
	"testing"
	"time"
)

func TestInitSupportPlan(t *testing.T) {
	start, _ := time.Parse("2006-01-02", "2020-04-10")
	end, _ := time.Parse("2006-01-02", "2021-04-10")
	iCalendar := mocks.MockICalendar{}
	_, err := InitSupportPlan(start, end, &iCalendar)
	if err != nil {
		t.Fatalf("Init support plan should be created")
	}
}

func TestFailInitSupportPlanBadStartEnd(t *testing.T) {
	start, _ := time.Parse("2006-01-02", "2020-04-10")
	end, _ := time.Parse("2006-01-02", "2019-04-10")
	iCalendar := mocks.MockICalendar{}
	_, err := InitSupportPlan(start, end, &iCalendar)
	if err == nil {
		t.Fatalf("start should be before end time")
	}
}
func TestFailInitSupportPlanShouldBeMoreThan24Hours(t *testing.T) {
	start, _ := time.Parse("2006-01-02", "2020-04-10")
	end, _ := time.Parse("2006-01-02", "2020-04-10")
	iCalendar := mocks.MockICalendar{}
	_, err := InitSupportPlan(start, end, &iCalendar)
	if err == nil {
		t.Fatalf("duration (end - start ) time should be greater than 24 hours")
	}
}

func TestCreateScheduleShouldCall365timesCreateDailyEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	iCalendar := mocks.NewMockICalendar(ctrl)
	iCalendar.EXPECT().CreateDailyEvent(gomock.Any(), gomock.Any()).Times(365).Return(nil) //365 for 1 year duration see start and end variables

	start, _ := time.Parse("2006-01-02", "2020-04-10")
	end, _ := time.Parse("2006-01-02", "2021-04-10")
	supportPlan, _ := InitSupportPlan(start, end, iCalendar)
	onCallPersons := []string{"alice", "bob", "carol", "dave"}

	supportPlan.CreateSchedule(onCallPersons)

}
func TestFailCreateScheduleMemberEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	iCalendar := mocks.NewMockICalendar(ctrl)

	start, _ := time.Parse("2006-01-02", "2020-04-10")
	end, _ := time.Parse("2006-01-02", "2021-04-10")
	supportPlan, _ := InitSupportPlan(start, end, iCalendar)
	onCallPersons := []string{}

	err := supportPlan.CreateSchedule(onCallPersons)
	if err == nil {
		t.Fatalf("start should be before end time")
	}

}
