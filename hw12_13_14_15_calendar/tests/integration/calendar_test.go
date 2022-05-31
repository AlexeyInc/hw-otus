//go:build integration
// +build integration

package integration_test

import (
	"database/sql"
	"os"
	"strconv"
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	calendar "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/api/protoc"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/util"

	sqlcstorage "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/sql/sqlc"

	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

const (
	_day                            = "Day"
	_week                           = "Week"
	_month                          = "Month"
	_testEventTitleSufix            = "_test"
	_minEventCount                  = 5
	_maxEventCount                  = 10
	_dbCallTime                     = 1
	_notificationSendedStatus int32 = 2
)

var checkEventsNotifFreq int

type CalendarSuite struct {
	suite.Suite
	ctx             context.Context
	serverConn      *grpc.ClientConn
	sqlDB           *sql.DB
	calendarStorage *sqlcstorage.Queries
	calendarClient  calendar.EventServiceClient
}

func (s *CalendarSuite) SetupSuite() {
	calendarHost := os.Getenv("CALENDAR_SERVER_HOST")
	if calendarHost == "" {
		calendarHost = "localhost:8081"
	}
	var err error
	s.serverConn, err = grpc.Dial(calendarHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.ctx = context.Background()

	s.calendarClient = calendar.NewEventServiceClient(s.serverConn)

	storageSource := os.Getenv("CALENDAR_DB_SOURCE")
	if storageSource == "" {
		storageSource = "postgres://alex:secret@localhost:5432/calendar?sslmode=disable"
	}
	s.sqlDB, err = sql.Open("postgres", storageSource)
	s.Require().NoError(err)

	s.calendarStorage = sqlcstorage.New(s.sqlDB)

	timeSleep := os.Getenv("SCHEDULER_CHECKNOTIFICATIONFREQSECONDS")
	if timeSleep == "" {
		checkEventsNotifFreq = 3
		return
	}

	numSec, err := strconv.Atoi(timeSleep)
	s.Require().NoError(err)
	checkEventsNotifFreq = numSec
}

func (s *CalendarSuite) TearDownSuite() {
	s.sqlDB.Close()
	s.serverConn.Close()
}

func (s *CalendarSuite) TearDownTest() {
	s.calendarStorage.DeleteTestEvents(s.ctx)
}

func (s *CalendarSuite) TestGetDefaultEvent() {
	resp, err := s.calendarClient.GetEvent(s.ctx, &calendar.GetEventRequest{
		Id: 1,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp.GetEvent())
}

func (s *CalendarSuite) TestEventCreation() {
	expectedEvent := s.generateRandomEvents(_day, 1)[0]

	resp := s.createEvent(expectedEvent)

	actualEvent := s.getEventById(resp.Event.Id)

	s.Require().Equal(expectedEvent.Title, actualEvent.Title)
	s.Require().Equal(expectedEvent.StartEvent.Seconds, actualEvent.StartEvent.Seconds)
	s.Require().Equal(expectedEvent.EndEvent.Seconds, actualEvent.EndEvent.Seconds)
	s.Require().Equal(expectedEvent.Description, actualEvent.Description)
	s.Require().Equal(expectedEvent.IdUser, actualEvent.IdUser)
	s.Require().Equal(expectedEvent.Notification.Seconds, actualEvent.Notification.Seconds)
}

func (s *CalendarSuite) TestGetDayEvents() {
	dayEventsCount := util.RandomIntRange(_minEventCount, _maxEventCount)
	events := s.generateRandomEvents(_day, dayEventsCount)

	for _, event := range events {
		s.createEvent(event)
	}

	resp, err := s.calendarClient.GetDayEvents(s.ctx, &calendar.GetEventsByDayRequest{
		Day: timestamppb.New(util.Period(_day).GetTimePeriod()),
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp.GetEvent())
	s.Require().Equal(dayEventsCount, len(resp.Event))
}

func (s *CalendarSuite) TestGetWeekEvents() {
	weekEventsCount := util.RandomIntRange(_minEventCount, _maxEventCount)
	events := s.generateRandomEvents(_week, weekEventsCount)

	for _, event := range events {
		s.createEvent(event)
	}

	resp, err := s.calendarClient.GetWeekEvents(s.ctx, &calendar.GetEventsByWeekRequest{
		WeekStart: timestamppb.New(util.Period(_week).GetTimePeriod()),
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp.GetEvent())
	s.Require().Equal(weekEventsCount, len(resp.Event))
}

func (s *CalendarSuite) TestGetMonthEvents() {
	monthEventsCount := util.RandomIntRange(_minEventCount, _maxEventCount)
	events := s.generateRandomEvents(_month, monthEventsCount)

	for _, event := range events {
		s.createEvent(event)
	}

	resp, err := s.calendarClient.GetMonthEvents(s.ctx, &calendar.GetEventsByMonthRequest{
		MonthStart: timestamppb.New(util.Period(_month).GetTimePeriod()),
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp.GetEvent())
	s.Require().Equal(monthEventsCount, len(resp.Event))
}

func (s *CalendarSuite) TestNotificationSending() {
	event := s.generateRandomEvents(_week, 1)[0]

	event.Notification = timestamppb.New(time.Now().UTC())

	resp := s.createEvent(event)

	updateEventNotifStatusTime := checkEventsNotifFreq + 1

	time.Sleep(time.Duration(updateEventNotifStatusTime) * time.Second)

	resultEvent, err := s.calendarStorage.GetEvent(s.ctx, resp.Event.Id)

	s.Require().NoError(err)
	s.Require().Equal(_notificationSendedStatus, resultEvent.Notificationstatus.Int32)
}

func (s *CalendarSuite) getEventById(eventId int64) *calendar.Event {
	resp, err := s.calendarClient.GetEvent(s.ctx, &calendar.GetEventRequest{
		Id: eventId,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp.GetEvent())
	return resp.Event
}

func (s *CalendarSuite) generateRandomEvents(p util.Period, eventsCount int) []*calendar.CreateEventRequest {
	events := make([]*calendar.CreateEventRequest, eventsCount)
	for i := 0; i < eventsCount; i++ {
		startEvent := p.GetTimePeriod()
		notification := startEvent.AddDate(0, 0, util.RandomInt(2)).UTC()
		endEvent := startEvent.AddDate(0, 0, util.RandomInt(10)).UTC()

		events[i] = &calendar.CreateEventRequest{
			Title:        util.RandomTitle() + _testEventTitleSufix,
			StartEvent:   timestamppb.New(startEvent),
			EndEvent:     timestamppb.New(endEvent),
			Notification: timestamppb.New(notification),
			Description:  util.RandomDescription(),
			IdUser:       util.RandomUserID(),
		}
	}
	return events
}

func (s CalendarSuite) createEvent(event *calendar.CreateEventRequest) *calendar.CreateEventResponse {
	resp, err := s.calendarClient.CreateEvent(s.ctx, event)
	s.Require().NoError(err)
	s.Require().NotNil(resp.GetEvent())
	return resp
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}
