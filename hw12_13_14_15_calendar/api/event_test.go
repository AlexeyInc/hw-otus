package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	api "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/api/protoc"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	models "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/models"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/require"
)

type Period string

const (
	_day   = "Day"
	_week  = "Week"
	_month = "Month"
)

var (
	configFilePath = "../configs/config.toml"
	logFilePath    = "../log/logs.log"
)

type EventModel struct {
	Id          int64     `json:"id,omitempty,string"`
	Title       string    `json:"title,omitempty"`
	StartEvent  time.Time `json:"startEvent,omitempty"`
	EndEvent    time.Time `json:"endEvent,omitempty"`
	Description string    `json:"description,omitempty"`
	IdUser      int64     `json:"idUser,omitempty,string"`
}

type EventResponse struct {
	Event EventModel
}

func TestEventAPI(t *testing.T) {
	ts, ctx, storage := createAndLaunchTestServer()
	defer ts.Close()

	fmt.Println(storage)

	client := &http.Client{}
	baseAppUrl := ts.URL + "/v1/EventService"

	baseEvent := createRandomEvent()

	t.Run("Create event", func(t *testing.T) {
		json_data, err := json.Marshal(baseEvent)
		require.Nil(t, err)

		req, err := http.NewRequestWithContext(ctx, "POST", baseAppUrl, bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}

		resp, err := client.Do(req)
		require.Nil(t, err)

		var result EventResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		requireEqual(t, baseEvent, result)
	})

	t.Run("Get event", func(t *testing.T) {
		getUrl := baseAppUrl + "/1"
		req, err := http.NewRequestWithContext(ctx, "GET", getUrl, nil)
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}
		resp, err := client.Do(req)
		require.Nil(t, err)

		var result EventResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		requireEqual(t, baseEvent, result)
	})
	t.Run("Update event", func(t *testing.T) {
		updateEventId := "1"
		eventId, err := strconv.ParseInt(updateEventId, 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		updateEvent := createRandomEvent()
		updateEvent["id"] = eventId

		json_data, err := json.Marshal(updateEvent)
		require.Nil(t, err)

		req, err := http.NewRequestWithContext(ctx, "PUT", baseAppUrl, bytes.NewBuffer(json_data))
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}

		_, err = client.Do(req)
		require.Nil(t, err)

		req, err = http.NewRequestWithContext(ctx, "GET", baseAppUrl+"/"+updateEventId, nil)
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}
		resp, err := client.Do(req)
		require.Nil(t, err)

		var updatedEvent EventResponse
		json.NewDecoder(resp.Body).Decode(&updatedEvent)

		requireEqual(t, updateEvent, updatedEvent)
	})
	t.Run("Delete event", func(t *testing.T) {
		deleteEventId := "/1"
		req, err := http.NewRequestWithContext(ctx, "DELETE", baseAppUrl+deleteEventId, nil)
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}
		resp, err := client.Do(req)
		require.Nil(t, err)

		var result api.EmptyResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		require.True(t, result.Success)
	})
	t.Run("Get events by day", func(t *testing.T) {
		defer storage.Close(ctx)

		eventsCount := util.RandomInt(10)

		events := createRandomDbEvents(_day, eventsCount)

		addEventsToStorage(ctx, storage, events)

		req, err := http.NewRequestWithContext(ctx, "GET", baseAppUrl+"/Day/2022-04-06T00:00:00Z", nil)
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

		events := createRandomDbEvents(_week, eventsCount)

		addEventsToStorage(ctx, storage, events)

		req, err := http.NewRequestWithContext(ctx, "GET", baseAppUrl+"/Week/2022-04-12T00:00:00Z", nil)
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

		events := createRandomDbEvents(_month, eventsCount)

		addEventsToStorage(ctx, storage, events)

		req, err := http.NewRequestWithContext(ctx, "GET", baseAppUrl+"/Month/2022-05-06T00:00:00Z", nil)
		if err != nil {
			log.Fatal("Request err: " + err.Error())
		}

		resultEvents := executeGetEventsRequest(t, req, client)

		require.NotNil(t, resultEvents.Event)
		require.Equal(t, eventsCount, len(resultEvents.Event))
	})
}

func requireEqual(t *testing.T, expected map[string]interface{}, actual EventResponse) {
	require.NotNil(t, actual.Event)
	require.Equal(t, expected["title"], actual.Event.Title)
	require.Equal(t, expected["startEvent"], actual.Event.StartEvent)
	require.Equal(t, expected["endEvent"], actual.Event.EndEvent)
	require.Equal(t, expected["description"], actual.Event.Description)
	require.Equal(t, expected["idUser"], actual.Event.IdUser)
}

func createAndLaunchTestServer() (*httptest.Server, context.Context, *memorystorage.MemoryStorage) {
	config, err := configs.NewConfig(configFilePath)
	if err != nil {
		log.Fatalln("can't read config file: " + err.Error())
	}

	zapLogg := logger.New(logFilePath, config.Logger.Level)
	defer zapLogg.ZapLogger.Sync()

	storage := memorystorage.New(config)

	calendar := app.New(zapLogg, storage)

	mux := runtime.NewServeMux()

	ctx := context.Background()

	err = api.RegisterEventServiceHandlerServer(ctx, mux, calendar)
	if err != nil {
		log.Fatal(err)
	}
	return httptest.NewServer(mux), ctx, storage
}

func createRandomEvent() map[string]interface{} {
	randEvent := map[string]interface{}{
		"title":       util.RandomTitle(),
		"startEvent":  time.Now().Local().UTC(),
		"endEvent":    time.Now().AddDate(0, 0, util.RandomInt(100)).Local().UTC(),
		"description": util.RandomDescription(),
		"idUser":      util.RandomUserID(),
	}
	return randEvent
}

func createRandomDbEvents(p Period, count int) []models.Event {
	events := make([]models.Event, count)
	for i := 0; i < count; i++ {
		startEvent := p.GetTimePeriod()
		endEvent := startEvent.AddDate(0, 0, util.RandomInt(10)).Local().UTC()

		events[i] = models.Event{
			Title:       util.RandomTitle(),
			StartEvent:  startEvent,
			EndEvent:    endEvent,
			Description: util.RandomDescription(),
			IDUser:      util.RandomUserID(),
		}
	}
	return events
}

func (s Period) GetTimePeriod() time.Time {
	switch s {
	case "Day":
		return time.Now().Local().UTC()
	case "Week":
		return time.Now().AddDate(0, 0, 7).Local().UTC()
	case "Month":
		return time.Now().AddDate(0, 0, 30).Local().UTC()
	}
	return time.Now().Local().UTC()
}

func addEventsToStorage(ctx context.Context, storage *memorystorage.MemoryStorage, events []models.Event) {
	for _, ev := range events {
		_, err := storage.CreateEvent(ctx, ev)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func executeGetEventsRequest(t *testing.T, req *http.Request, client *http.Client) api.GetEventsResponse {
	t.Helper()

	resp, err := client.Do(req)
	require.Nil(t, err)

	var resultEvents api.GetEventsResponse
	json.NewDecoder(resp.Body).Decode(&resultEvents)
	return resultEvents
}
