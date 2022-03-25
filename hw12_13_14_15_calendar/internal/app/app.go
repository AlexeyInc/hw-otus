package app

import (
	"context"
	"time"

	models "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/models"
)

type App struct { // TODO
	storage Storage
}

type Logger interface { // TODO
}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	CreateEvent(ctx context.Context, event models.Event) (models.Event, error)
	UpdateEvent(ctx context.Context, event models.Event) (models.Event, error)
	DeleteEvent(ctx context.Context, id int64) error
	GetEvent(ctx context.Context, id int64) (models.Event, error)
	GetDayEvents(ctx context.Context, day time.Time) ([]models.Event, error)
	GetWeekEvents(ctx context.Context, weekStart time.Time) ([]models.Event, error)
	GetMonthEvents(ctx context.Context, monthStart time.Time) ([]models.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context,
	title string, start time.Time, end time.Time, idUser int64) (models.Event, error) {
	newEvent := models.Event{
		Title:      title,
		StartEvent: start,
		EndEvent:   end,
		IDUser:     idUser,
	}

	return a.storage.CreateEvent(ctx, newEvent)
}

func (a *App) GetEvent(ctx context.Context, eventID int64) (models.Event, error) {
	return a.storage.GetEvent(context.Background(), eventID)
}

func (a *App) UpdateEvent(ctx context.Context,
	eventID int64, title string, start time.Time, end time.Time, idUser int64) (models.Event, error) {
	updateEvent := models.Event{
		ID:         eventID,
		Title:      title,
		StartEvent: start,
		EndEvent:   end,
		IDUser:     idUser,
	}
	return a.storage.UpdateEvent(context.Background(), updateEvent)
}

func (a *App) DeleteEvent(ctx context.Context, eventID int64) error {
	return a.storage.DeleteEvent(ctx, eventID)
}

func (a *App) DayEvents(ctx context.Context, day time.Time) ([]models.Event, error) {
	return a.storage.GetDayEvents(ctx, day)
}

func (a *App) DayWeek(ctx context.Context, weekStart time.Time) ([]models.Event, error) {
	return a.storage.GetWeekEvents(ctx, weekStart)
}

func (a *App) DayMonth(ctx context.Context, monthStart time.Time) ([]models.Event, error) {
	return a.storage.GetMonthEvents(ctx, monthStart)
}
