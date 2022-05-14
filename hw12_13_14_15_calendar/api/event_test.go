package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	api "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/api/protoc"
	configs "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	models "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/models"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/require"
)

var (
	configFilePath = "../configs/calendar_config.toml"
	logFilePath    = "../log/logs.log"
)

type Period int

const (
	Day Period = iota
	Week
	Month
)

type EventModel struct {
	ID           int64     `json:"id,omitempty,string"`
	Title        string    `json:"title,omitempty"`
	StartEvent   time.Time `json:"startEvent,omitempty"`
	EndEvent     time.Time `json:"endEvent,omitempty"`
	Description  string    `json:"description,omitempty"`
	IDUser       int64     `json:"idUser,omitempty,string"`
	Notification time.Time `json:"notification,omitempty"`
}

type EventResponse struct {
	Event EventModel
}

func TestEventCRUD(t *testing.T) {
	ts, ctx, storage := createAndLaunchTestServer(t)
	defer ts.Close()

	client := &http.Client{}
	baseAppURL := ts.URL + "/v1/EventService"
	baseEvent := createRandomEvent()

	t.Run("Create event", func(t *testing.T) {
		defer storage.Close(ctx)

		jsonData, err := json.Marshal(baseEvent)
		require.Nil(t, err)

		req, err := http.NewRequestWithContext(ctx, "POST", baseAppURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}

		resp, err := client.Do(req)
		require.Nil(t, err)
		defer resp.Body.Close()

		var result EventResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		requireEqualMap(t, baseEvent, result)
	})

	t.Run("Get event", func(t *testing.T) {
		defer storage.Close(ctx)

		newEvent := createRandomDBEventModels(Day, 1)[0]
		eventID := addEventToStorage(ctx, storage, newEvent)
		getURL := baseAppURL + "/" + strconv.FormatInt(eventID, 10)
		req, err := http.NewRequestWithContext(ctx, "GET", getURL, nil)
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}

		resp, err := client.Do(req)
		require.Nil(t, err)
		defer resp.Body.Close()

		var result EventResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		requireEqual(t, newEvent, result)
	})
	t.Run("Update event", func(t *testing.T) {
		defer storage.Close(ctx)

		newEvent := createRandomDBEventModels(Day, 1)[0]
		newEventID := addEventToStorage(ctx, storage, newEvent)
		updateEvent := createRandomEvent()
		updateEvent["id"] = newEventID

		jsonData, err := json.Marshal(updateEvent)
		require.Nil(t, err)

		req, err := http.NewRequestWithContext(ctx, "PUT", baseAppURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}

		tempResp, err := client.Do(req)
		require.Nil(t, err)
		defer tempResp.Body.Close()

		updatedEventID := strconv.FormatInt(newEventID, 10)
		req, err = http.NewRequestWithContext(ctx, "GET", baseAppURL+"/"+updatedEventID, nil)
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}
		resp, err := client.Do(req)
		require.Nil(t, err)
		defer resp.Body.Close()

		var updatedEvent EventResponse
		json.NewDecoder(resp.Body).Decode(&updatedEvent)
		requireEqualMap(t, updateEvent, updatedEvent)
	})
	t.Run("Delete event", func(t *testing.T) {
		defer storage.Close(ctx)

		newEvent := createRandomDBEventModels(Day, 1)[0]
		newEventID := addEventToStorage(ctx, storage, newEvent)

		deleteEventID := strconv.FormatInt(newEventID, 10)
		req, err := http.NewRequestWithContext(ctx, "DELETE", baseAppURL+"/"+deleteEventID, nil)
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}
		resp, err := client.Do(req)
		require.Nil(t, err)
		defer resp.Body.Close()

		var result api.EmptyResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		require.True(t, result.Success)
	})
}

func TestEventPeriods(t *testing.T) {
	ts, ctx, storage := createAndLaunchTestServer(t)
	defer ts.Close()

	client := &http.Client{}
	baseAppURL := ts.URL + "/v1/EventService"

	t.Run("Get events by day", func(t *testing.T) {
		defer storage.Close(ctx)

		eventsCount := util.RandomInt(10)
		events := createRandomDBEventModels(Day, eventsCount)
		for _, ev := range events {
			addEventToStorage(ctx, storage, ev)
		}

		nowPlusDay := time.Now().UTC().Format("2006-01-02T15:04:05Z")
		req, err := http.NewRequestWithContext(ctx, "GET", baseAppURL+"/Day/"+nowPlusDay, nil)
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}

		resultEvents := executeGetEventsRequest(t, req, client)
		require.NotNil(t, resultEvents.Event)
		require.Equal(t, eventsCount, len(resultEvents.Event))
	})
	t.Run("Get events by week", func(t *testing.T) {
		defer storage.Close(ctx)

		eventsCount := util.RandomInt(10)
		events := createRandomDBEventModels(Week, eventsCount)
		for _, ev := range events {
			addEventToStorage(ctx, storage, ev)
		}

		nowPlusWeek := time.Now().UTC().AddDate(0, 0, 7).Format("2006-01-02T15:04:05Z")
		req, err := http.NewRequestWithContext(ctx, "GET", baseAppURL+"/Week/"+nowPlusWeek, nil)
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}

		resultEvents := executeGetEventsRequest(t, req, client)
		require.NotNil(t, resultEvents.Event)
		require.Equal(t, eventsCount, len(resultEvents.Event))
	})
	t.Run("Get events by month", func(t *testing.T) {
		defer storage.Close(ctx)

		eventsCount := util.RandomInt(10)
		events := createRandomDBEventModels(Month, eventsCount)
		for _, ev := range events {
			addEventToStorage(ctx, storage, ev)
		}

		nowPlusMonth := time.Now().UTC().AddDate(0, 0, 30).Format("2006-01-02T15:04:05Z")
		req, err := http.NewRequestWithContext(ctx, "GET", baseAppURL+"/Month/"+nowPlusMonth, nil)
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}

		resultEvents := executeGetEventsRequest(t, req, client)
		require.NotNil(t, resultEvents.Event)
		require.Equal(t, eventsCount, len(resultEvents.Event))
	})
}

func requireEqualMap(t *testing.T, expected map[string]interface{}, actual EventResponse) {
	t.Helper()
	require.NotNil(t, actual.Event)
	require.Equal(t, expected["title"], actual.Event.Title)
	require.Equal(t, expected["startEvent"], actual.Event.StartEvent)
	require.Equal(t, expected["endEvent"], actual.Event.EndEvent)
	require.Equal(t, expected["description"], actual.Event.Description)
	require.Equal(t, expected["idUser"], actual.Event.IDUser)
}

func requireEqual(t *testing.T, expected models.Event, actual EventResponse) {
	t.Helper()
	require.NotNil(t, actual.Event)
	require.Equal(t, expected.Title, actual.Event.Title)
	require.Equal(t, expected.StartEvent, actual.Event.StartEvent)
	require.Equal(t, expected.EndEvent, actual.Event.EndEvent)
	require.Equal(t, expected.Description, actual.Event.Description)
	require.Equal(t, expected.IDUser, actual.Event.IDUser)
}

func createAndLaunchTestServer(t *testing.T) (*httptest.Server, context.Context, *memorystorage.MemoryStorage) {
	t.Helper()
	config, err := configs.NewConfig(configFilePath)
	if err != nil {
		require.NoError(t, err, "can't read config file")
	}

	zapLogg := logger.New(logFilePath, config.Logger.Level)
	defer zapLogg.ZapLogger.Sync()

	storage := memorystorage.New(config)
	calendar := app.New(zapLogg, storage)
	mux := runtime.NewServeMux()

	ctx := context.Background()
	err = api.RegisterEventServiceHandlerServer(ctx, mux, calendar)
	if err != nil {
		require.NoError(t, err, "can't register event service server")
	}
	return httptest.NewServer(mux), ctx, storage
}

func createRandomEvent() map[string]interface{} {
	randEvent := map[string]interface{}{
		"title":        util.RandomTitle(),
		"startEvent":   time.Now().UTC(),
		"endEvent":     time.Now().AddDate(0, 0, util.RandomInt(100)).UTC(),
		"description":  util.RandomDescription(),
		"idUser":       util.RandomUserID(),
		"notification": time.Now().UTC().AddDate(0, 0, -1),
	}
	return randEvent
}

func createRandomDBEventModels(p Period, count int) []models.Event {
	events := make([]models.Event, count)
	for i := 0; i < count; i++ {
		startEvent := GetTimeStartPeriod(p)
		endEvent := startEvent.AddDate(0, 0, util.RandomInt(10)).UTC()

		events[i] = models.Event{
			Title:        util.RandomTitle(),
			StartEvent:   startEvent,
			EndEvent:     endEvent,
			Description:  util.RandomDescription(),
			IDUser:       util.RandomUserID(),
			Notification: time.Now().AddDate(0, 0, -1).UTC(),
		}
	}
	return events
}

func GetTimeStartPeriod(p Period) time.Time {
	switch p {
	case Day:
		return time.Now().UTC()
	case Week:
		return time.Now().AddDate(0, 0, 7).UTC()
	case Month:
		return time.Now().AddDate(0, 0, 30).UTC()
	}
	return time.Now().UTC()
}

func addEventToStorage(ctx context.Context, storage *memorystorage.MemoryStorage, ev models.Event) int64 {
	event, err := storage.CreateEvent(ctx, ev)
	if err != nil {
		log.Fatal(err)
	}
	return event.ID
}

func executeGetEventsRequest(t *testing.T, req *http.Request, client *http.Client) *api.GetEventsResponse {
	t.Helper()
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()

	var resultEvents api.GetEventsResponse
	json.NewDecoder(resp.Body).Decode(&resultEvents)
	return &resultEvents
}
