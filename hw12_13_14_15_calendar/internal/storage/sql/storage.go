package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	calendarconfig "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	sqlc "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/internal/storage/sql/sqlc"
	domainModels "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/models"

	_ "github.com/lib/pq"
)

type Storage struct {
	db        *sql.DB
	DbQueries *sqlc.Queries

	Driver string
	Source string
}

func New(c calendarconfig.Config) *Storage {
	return &Storage{
		Driver: c.Storage.Driver,
		Source: c.Storage.Source,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sql.Open(s.Driver, s.Source)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	s.db = db
	connErr := s.db.PingContext(ctx)
	if connErr != nil {
		return connErr
	}

	s.DbQueries = sqlc.New(db)

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event domainModels.Event) (domainModels.Event, error) {
	createEvent := sqlc.CreateEventParams{
		Title:        event.Title,
		StartEvent:   event.StartEvent,
		EndEvent:     event.EndEvent,
		Description:  sql.NullString{String: event.Description, Valid: true},
		IDUser:       event.IDUser,
		Notification: sql.NullTime{Time: event.Notification, Valid: true},
	}

	createdModel, err := s.DbQueries.CreateEvent(ctx, createEvent)

	return toViewModel(createdModel), err
}

func (s *Storage) UpdateEvent(ctx context.Context, event domainModels.Event) (domainModels.Event, error) {
	updateEvent := sqlc.UpdateEventParams{
		ID:           event.ID,
		Title:        event.Title,
		StartEvent:   event.StartEvent,
		EndEvent:     event.EndEvent,
		Description:  sql.NullString{String: event.Description, Valid: true},
		IDUser:       event.IDUser,
		Notification: sql.NullTime{Time: event.Notification, Valid: true},
	}

	updatedEvent, err := s.DbQueries.UpdateEvent(ctx, updateEvent)

	return toViewModel(updatedEvent), err
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	return s.DbQueries.DeleteEvent(ctx, id)
}

func (s *Storage) GetEvent(ctx context.Context, id int64) (eventModel domainModels.Event, err error) {
	event, err := s.DbQueries.GetEvent(ctx, id)
	if err != nil {
		return eventModel, err
	}
	return toViewModel(event), err
}

func (s *Storage) GetDayEvents(ctx context.Context, day time.Time) (eventModels []domainModels.Event, err error) {
	events, err := s.DbQueries.GetDayEvents(ctx, day)
	if err != nil {
		return eventModels, err
	}
	return toViewModels(events), err
}

func (s Storage) GetWeekEvents(ctx context.Context, weekStart time.Time) (eventModels []domainModels.Event, err error) {
	events, err := s.DbQueries.GetWeekEvents(ctx, weekStart)
	if err != nil {
		return eventModels, err
	}
	return toViewModels(events), err
}

func (s Storage) GetMonthEvents(ctx context.Context, monthStart time.Time) (evModels []domainModels.Event, err error) {
	events, err := s.DbQueries.GetMonthEvents(ctx, monthStart)
	if err != nil {
		return evModels, err
	}
	return toViewModels(events), err
}

func toViewModel(ev sqlc.Event) domainModels.Event {
	return domainModels.Event{
		ID:           ev.ID,
		Title:        ev.Title,
		StartEvent:   ev.StartEvent,
		EndEvent:     ev.EndEvent,
		Description:  ev.Description.String,
		IDUser:       ev.IDUser,
		Notification: ev.Notification.Time,
	}
}

func toViewModels(events []sqlc.Event) []domainModels.Event {
	result := make([]domainModels.Event, len(events))
	for i := 0; i < len(events); i++ {
		result[i] = toViewModel(events[i])
	}
	return result
}
