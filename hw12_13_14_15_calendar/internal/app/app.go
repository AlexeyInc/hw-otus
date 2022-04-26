package app

import (
	"context"
	"time"

	api "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/api/protoc"
	models "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type App struct {
	api.UnimplementedEventServiceServer

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

func (a *App) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
	eventDto := models.Event{
		Title:        req.Title,
		StartEvent:   req.StartEvent.AsTime(),
		EndEvent:     req.EndEvent.AsTime(),
		Description:  req.Description,
		IDUser:       req.IdUser,
		Notification: req.Notification.AsTime(),
	}

	createdEvent, err := a.storage.CreateEvent(ctx, eventDto)

	response := &api.CreateEventResponse{
		Event: &api.Event{
			Id:           createdEvent.ID,
			Title:        createdEvent.Title,
			StartEvent:   timestamppb.New(createdEvent.StartEvent),
			EndEvent:     timestamppb.New(createdEvent.EndEvent),
			Description:  createdEvent.Description,
			IdUser:       createdEvent.IDUser,
			Notification: timestamppb.New(createdEvent.Notification),
		},
	}

	return response, err
}

func (a *App) GetEvent(ctx context.Context, req *api.GetEventRequest) (*api.GetEventResponse, error) {
	eventDto, err := a.storage.GetEvent(context.Background(), req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	response := &api.GetEventResponse{
		Event: &api.Event{
			Id:           eventDto.ID,
			Title:        eventDto.Title,
			StartEvent:   timestamppb.New(eventDto.StartEvent),
			EndEvent:     timestamppb.New(eventDto.EndEvent),
			Description:  eventDto.Description,
			IdUser:       eventDto.IDUser,
			Notification: timestamppb.New(eventDto.Notification),
		},
	}

	return response, err
}

func (a *App) UpdateEvent(ctx context.Context, req *api.UpdateEventRequest) (*api.UpdateEventResponse, error) {
	eventDto := models.Event{
		ID:           req.Id,
		Title:        req.Title,
		StartEvent:   req.StartEvent.AsTime(),
		EndEvent:     req.EndEvent.AsTime(),
		Description:  req.Description,
		IDUser:       req.IdUser,
		Notification: req.Notification.AsTime(),
	}

	updatedEvent, err := a.storage.UpdateEvent(context.Background(), eventDto)

	response := &api.UpdateEventResponse{
		Event: &api.Event{
			Id:           updatedEvent.ID,
			Title:        updatedEvent.Title,
			StartEvent:   timestamppb.New(updatedEvent.StartEvent),
			EndEvent:     timestamppb.New(updatedEvent.EndEvent),
			Description:  updatedEvent.Description,
			IdUser:       updatedEvent.IDUser,
			Notification: timestamppb.New(updatedEvent.Notification),
		},
	}
	return response, err
}

func (a *App) DeleteEvent(ctx context.Context,
	req *api.DeleteEventRequest) (response *api.EmptyResponse, err error) {
	err = a.storage.DeleteEvent(ctx, req.Id)
	response = &api.EmptyResponse{
		Success: err == nil,
	}
	return
}

func (a *App) GetDayEvents(ctx context.Context,
	day *api.GetEventsByDayRequest) (*api.GetEventsResponse, error) {
	eventsDto, err := a.storage.GetDayEvents(ctx, day.Day.AsTime())

	return toResposeModels(eventsDto), err
}

func (a *App) GetWeekEvents(ctx context.Context,
	weekStart *api.GetEventsByWeekRequest) (*api.GetEventsResponse, error) {
	eventsDto, err := a.storage.GetWeekEvents(ctx, weekStart.WeekStart.AsTime())

	return toResposeModels(eventsDto), err
}

func (a *App) GetMonthEvents(ctx context.Context,
	monthStart *api.GetEventsByMonthRequest) (*api.GetEventsResponse, error) {
	eventsDto, err := a.storage.GetMonthEvents(ctx, monthStart.MonthStart.AsTime())

	return toResposeModels(eventsDto), err
}

func toResposeModels(events []models.Event) *api.GetEventsResponse {
	results := &api.GetEventsResponse{
		Event: make([]*api.Event, len(events)),
	}
	for i, ev := range events {
		results.Event[i] = &api.Event{
			Id:           ev.ID,
			Title:        ev.Title,
			StartEvent:   timestamppb.New(ev.StartEvent),
			EndEvent:     timestamppb.New(ev.EndEvent),
			Description:  ev.Description,
			IdUser:       ev.IDUser,
			Notification: timestamppb.New(ev.Notification),
		}
	}
	return results
}
