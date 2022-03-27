package sqlcstorage

// import (
// 	"context"
// 	"database/sql"
// 	"testing"
// 	"time"

// 	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/util"
// 	"github.com/stretchr/testify/require"
// )

// func createRandomEvent(t *testing.T) Event {
// 	t.Helper()
// 	expectedEvent := CreateEventParams{
// 		Title:       util.RandomTitle() + "_test",
// 		StartEvent:  time.Now().AddDate(0, 0, util.RandomInt(100)).Local().UTC(),
// 		EndEvent:    time.Now().AddDate(0, 0, util.RandomInt(100)).Local().UTC(),
// 		Description: sql.NullString{String: util.RandomDescription(), Valid: true},
// 		IDUser:      util.RandomUserID(),
// 	}

// 	actualEvent, err := testQueries.CreateEvent(context.Background(), expectedEvent)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, actualEvent)

// 	require.Equal(t, expectedEvent.Title, actualEvent.Title)
// 	require.Equal(t, expectedEvent.StartEvent.Local().UTC(), actualEvent.StartEvent.Local().UTC())
// 	require.Equal(t, expectedEvent.EndEvent.Local().UTC(), actualEvent.EndEvent.Local().UTC())
// 	require.Equal(t, expectedEvent.Description, actualEvent.Description)
// 	require.Equal(t, expectedEvent.IDUser, actualEvent.IDUser)

// 	return actualEvent
// }

// func TestCreateEvent(t *testing.T) {
// 	event := createRandomEvent(t)

// 	require.NotZero(t, event.ID)
// 	require.NotZero(t, event.Title)
// }

// func TestUpdateAccount(t *testing.T) {
// 	oldEvent := createRandomEvent(t)

// 	updateArg := UpdateEventParams{
// 		ID:          oldEvent.ID,
// 		Title:       util.RandomTitle() + "_test",
// 		StartEvent:  time.Now().Local().UTC(),
// 		EndEvent:    time.Now().AddDate(0, 0, util.RandomInt(100)).Local().UTC(),
// 		Description: sql.NullString{String: util.RandomDescription(), Valid: true},
// 		IDUser:      util.RandomUserID(),
// 	}

// 	updatedEvent, err := testQueries.UpdateEvent(context.Background(), updateArg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, updatedEvent)

// 	require.Equal(t, oldEvent.ID, updatedEvent.ID)
// 	require.Equal(t, updateArg.Title, updatedEvent.Title)
// 	require.Equal(t, updateArg.StartEvent.Local().UTC(), updatedEvent.StartEvent.Local().UTC())
// 	require.Equal(t, updateArg.EndEvent.Local().UTC(), updatedEvent.EndEvent.Local().UTC())
// 	require.Equal(t, updateArg.Description, updateArg.Description)
// 	require.Equal(t, updateArg.IDUser, updateArg.IDUser)
// }

// func TestDeleteEvent(t *testing.T) {
// 	newEvent := createRandomEvent(t)
// 	err := testQueries.DeleteEvent(context.Background(), newEvent.ID)
// 	require.NoError(t, err)

// 	deletedEvent, err := testQueries.GetEvent(context.Background(), newEvent.ID)
// 	require.Error(t, err)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// 	require.Empty(t, deletedEvent)
// }

// func TestTodayGetEvents(t *testing.T) {
// 	var lastEvent Event
// 	for i := 0; i < 20; i++ {
// 		lastEvent = createRandomEvent(t)
// 	}
// 	events, err := testQueries.GetDayEvents(context.Background(), lastEvent.StartEvent)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, events)

// 	day := time.Hour * 24
// 	var date time.Time
// 	for i := 0; i < len(events); i++ {
// 		date = events[i].StartEvent

// 		require.WithinDuration(t,
// 			lastEvent.StartEvent.Local().UTC(),
// 			date.Local().UTC(),
// 			day)
// 	}
// }

// func TestWeekEvents(t *testing.T) {
// 	var lastEvent Event
// 	for i := 0; i < 20; i++ {
// 		lastEvent = createRandomEvent(t)
// 	}

// 	events, err := testQueries.GetWeekEvents(context.Background(), lastEvent.StartEvent)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, events)

// 	week := time.Hour * 168
// 	var date time.Time
// 	for i := 0; i < len(events); i++ {
// 		date = events[i].StartEvent

// 		require.WithinDuration(t,
// 			lastEvent.StartEvent.Local().UTC(),
// 			date.Local().UTC(),
// 			week)
// 	}
// }

// func TestMonthEvents(t *testing.T) {
// 	var lastEvent Event
// 	for i := 0; i < 20; i++ {
// 		lastEvent = createRandomEvent(t)
// 	}

// 	events, err := testQueries.GetMonthEvents(context.Background(), lastEvent.StartEvent)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, events)

// 	month := time.Hour * 730
// 	var date time.Time
// 	for i := 0; i < len(events); i++ {
// 		date = events[i].StartEvent

// 		require.WithinDuration(t,
// 			lastEvent.StartEvent.Local().UTC(),
// 			date.Local().UTC(),
// 			month)
// 	}
// }

// func deleteAllTestEvents() {
// 	testQueries.DeleteTestEvents(context.Background())
// }
