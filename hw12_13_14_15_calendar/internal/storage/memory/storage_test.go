package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage"
	models "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/models"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/util"
	"github.com/stretchr/testify/require"
)

func createRandomEvent(t *testing.T) models.Event {
	t.Helper()
	event := models.Event{
		Title:       util.RandomTitle(),
		StartEvent:  time.Now().Local().UTC(),
		EndEvent:    time.Now().AddDate(0, 0, util.RandomInt(100)).Local().UTC(),
		Description: util.RandomDescription(),
		IDUser:      util.RandomUserID(),
	}

	ev, err := memoryStorage.CreateEvent(context.Background(), event)
	require.NoError(t, err)
	require.NotEmpty(t, ev)

	return ev
}

func TestCreateEvent(t *testing.T) {
	createRandomEvent(t)
}

func TestGetEvent(t *testing.T) {
	newEvent := createRandomEvent(t)
	eventID := newEvent.ID

	event, err := memoryStorage.GetEvent(context.Background(), eventID)

	require.NoError(t, err)
	require.NotEmpty(t, event)
	require.Equal(t, eventID, event.ID)
	require.Equal(t, newEvent.Title, event.Title)
	require.Equal(t, newEvent.StartEvent, event.StartEvent)
	require.Equal(t, newEvent.EndEvent, event.EndEvent)
	require.Equal(t, newEvent.Description, event.Description)
	require.Equal(t, newEvent.IDUser, event.IDUser)
}

func TestDeleteEvent(t *testing.T) {
	newEvent := createRandomEvent(t)
	eventID := newEvent.ID

	err := memoryStorage.DeleteEvent(context.Background(), eventID)

	require.NoError(t, err)

	notExistEvent, err := memoryStorage.GetEvent(context.Background(), eventID)

	require.Error(t, err)
	require.EqualError(t, err, storage.ErrEventNotFound.Error())
	require.Empty(t, notExistEvent)
}

func TestUpdateEvent(t *testing.T) {
	event := models.Event{
		ID:          1,
		Title:       util.RandomTitle() + "_test",
		StartEvent:  time.Now().Local().UTC(),
		EndEvent:    time.Now().AddDate(0, 0, util.RandomInt(100)).Local().UTC(),
		Description: util.RandomDescription(),
		IDUser:      util.RandomUserID(),
	}

	updatedEvent, err := memoryStorage.UpdateEvent(context.Background(), event)
	require.NoError(t, err)

	require.Equal(t, event.Title, updatedEvent.Title)
	require.Equal(t, event.StartEvent, updatedEvent.StartEvent)
	require.Equal(t, event.EndEvent, updatedEvent.EndEvent)
	require.Equal(t, event.Description, updatedEvent.Description)
	require.Equal(t, event.IDUser, updatedEvent.IDUser)
}

func TestGetWeekEvents(t *testing.T) {
	var lastEvent models.Event
	for i := 0; i < 20; i++ {
		lastEvent = createRandomEvent(t)
	}

	events, err := memoryStorage.GetWeekEvents(context.Background(), lastEvent.StartEvent)
	require.NoError(t, err)
	require.NotEmpty(t, events)

	week := time.Hour * 168
	var date time.Time
	for i := 0; i < len(events); i++ {
		date = events[i].StartEvent

		require.WithinDuration(t,
			lastEvent.StartEvent.Local().UTC(),
			date.Local().UTC(),
			week)
	}
}

func TestGetMonthEvents(t *testing.T) {
	var lastEvent models.Event
	for i := 0; i < 20; i++ {
		lastEvent = createRandomEvent(t)
	}

	events, err := memoryStorage.GetMonthEvents(context.Background(), lastEvent.StartEvent)
	require.NoError(t, err)
	require.NotEmpty(t, events)

	month := time.Hour * 730
	var date time.Time
	for i := 0; i < len(events); i++ {
		date = events[i].StartEvent

		require.WithinDuration(t,
			lastEvent.StartEvent.Local().UTC(),
			date.Local().UTC(),
			month)
	}
}

func TestDataRaceOnCRUD(t *testing.T) {
	for i := 0; i < 30; i++ {
		event := models.Event{
			Title:       util.RandomTitle(),
			StartEvent:  time.Now().Local().UTC(),
			EndEvent:    time.Now().AddDate(0, 0, util.RandomInt(100)).Local().UTC(),
			Description: util.RandomDescription(),
			IDUser:      util.RandomUserID(),
		}
		go memoryStorage.CreateEvent(context.Background(), event)
		go memoryStorage.UpdateEvent(context.Background(), event)
		go memoryStorage.DeleteEvent(context.Background(), event.ID)
		go memoryStorage.GetEvent(context.Background(), event.ID)
	}
}
