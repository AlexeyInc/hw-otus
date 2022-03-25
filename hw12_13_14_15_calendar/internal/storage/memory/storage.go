package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage"
	models "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/models"
)

var globalNewEventID int64 = 1

type MemoryStorage struct {
	events map[int64]models.Event
	mutex  *sync.RWMutex
}

func New(c configs.Config) *MemoryStorage {
	return &MemoryStorage{
		events: make(map[int64]models.Event),
		mutex:  new(sync.RWMutex),
	}
}

func (s *MemoryStorage) CreateEvent(ctx context.Context, ev models.Event) (models.Event, error) {
	var newEvent models.Event
	s.mutex.Lock()
	newEvent.ID = globalNewEventID
	newEvent.Title = ev.Title
	newEvent.StartEvent = ev.StartEvent
	newEvent.EndEvent = ev.EndEvent
	newEvent.Description = ev.Description
	newEvent.IDUser = ev.IDUser
	s.events[globalNewEventID] = newEvent
	globalNewEventID++
	s.mutex.Unlock()

	return newEvent, nil
}

func (s *MemoryStorage) UpdateEvent(ctx context.Context, event models.Event) (models.Event, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ev, ex := s.events[event.ID]
	if !ex {
		return ev, storage.ErrEventNotFound
	}
	ev.Title = event.Title
	ev.StartEvent = event.StartEvent
	ev.EndEvent = event.EndEvent
	ev.Description = event.Description
	ev.IDUser = event.IDUser

	s.events[event.ID] = ev

	return ev, nil
}

func (s *MemoryStorage) DeleteEvent(ctx context.Context, eventID int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ev, ex := s.events[eventID]
	if ex && ev.ID == eventID {
		delete(s.events, eventID)
		return nil
	}

	return storage.ErrEventNotFound
}

func (s *MemoryStorage) GetEvent(ctx context.Context, eventID int64) (models.Event, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ev, ex := s.events[eventID]
	if ex && ev.ID == eventID {
		return ev, nil
	}

	return ev, storage.ErrEventNotFound
}

func (s *MemoryStorage) GetDayEvents(ctx context.Context, day time.Time) ([]models.Event, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	events := make([]models.Event, 0)
	for _, ev := range s.events {
		if ev.StartEvent == day || day.UTC().After(ev.StartEvent.UTC()) && day.UTC().Before(ev.StartEvent.UTC()) {
			events = append(events, ev)
		}
	}
	return events, nil
}

func (s *MemoryStorage) GetWeekEvents(ctx context.Context, weekStart time.Time) ([]models.Event, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	events := make([]models.Event, 0)
	for _, ev := range s.events {
		if ev.StartEvent == weekStart ||
			weekStart.UTC().After(ev.StartEvent.UTC()) && weekStart.UTC().Before(ev.StartEvent.AddDate(0, 0, 7)) {
			events = append(events, ev)
		}
	}

	return events, nil
}

func (s *MemoryStorage) GetMonthEvents(ctx context.Context, monthStart time.Time) ([]models.Event, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	events := make([]models.Event, 0)
	for _, ev := range s.events {
		if ev.StartEvent == monthStart ||
			monthStart.UTC().After(ev.StartEvent.UTC()) && monthStart.UTC().Before(ev.StartEvent.AddDate(0, 1, 0)) {
			events = append(events, ev)
		}
	}

	return events, nil
}

func (s *MemoryStorage) Connect(ctx context.Context) error {
	return nil
}

func (s *MemoryStorage) Close(ctx context.Context) error {
	return nil
}
