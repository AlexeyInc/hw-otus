package app

import (
	"context"
	"time"

	pb "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/api/protoc"
	models "github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type App struct {
	pb.UnimplementedEventServiceServer
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

func (a *App) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	eventDto := models.Event{
		Title:       req.Title,
		StartEvent:  req.StartEvent.AsTime(),
		EndEvent:    req.EndEvent.AsTime(),
		Description: req.Description,
		IDUser:      req.IdUser,
	}

	createdEvent, err := a.storage.CreateEvent(ctx, eventDto)

	response := &pb.CreateEventResponse{
		Event: &pb.Event{
			Id:          createdEvent.ID,
			Title:       createdEvent.Title,
			StartEvent:  timestamppb.New(createdEvent.StartEvent),
			EndEvent:    timestamppb.New(createdEvent.EndEvent),
			Description: createdEvent.Description,
			IdUser:      createdEvent.IDUser,
		},
	}

	return response, err
}

func (a *App) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	eventDto, err := a.storage.GetEvent(context.Background(), req.Id)

	response := &pb.GetEventResponse{
		Event: &pb.Event{
			Id:          eventDto.ID,
			Title:       eventDto.Title,
			StartEvent:  timestamppb.New(eventDto.StartEvent),
			EndEvent:    timestamppb.New(eventDto.EndEvent),
			Description: eventDto.Description,
			IdUser:      eventDto.IDUser,
		},
	}

	return response, err
}

func (a *App) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	eventDto := models.Event{
		ID:          req.Id,
		Title:       req.Title,
		StartEvent:  req.StartEvent.AsTime(),
		EndEvent:    req.EndEvent.AsTime(),
		Description: req.Description,
		IDUser:      req.IdUser,
	}

	updatedEvent, err := a.storage.UpdateEvent(context.Background(), eventDto)

	response := &pb.UpdateEventResponse{
		Event: &pb.Event{
			Id:          updatedEvent.ID,
			Title:       updatedEvent.Title,
			StartEvent:  timestamppb.New(updatedEvent.StartEvent),
			EndEvent:    timestamppb.New(updatedEvent.EndEvent),
			Description: updatedEvent.Description,
			IdUser:      updatedEvent.IDUser,
		},
	}

	return response, err
}

func (a *App) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (response *pb.EmptyResponse, err error) {
	err = a.storage.DeleteEvent(ctx, req.Id)
	response = &pb.EmptyResponse{
		Success: err == nil,
	}
	return
}

func (a *App) GetDayEvents(ctx context.Context, day *pb.GetEventsByDayRequest) (*pb.GetEventsResponse, error) {
	eventsDto, err := a.storage.GetDayEvents(ctx, day.Day.AsTime())

	return toResposeModels(eventsDto), err
}

func (a *App) GetWeekEvents(ctx context.Context, weekStart *pb.GetEventsByWeekRequest) (*pb.GetEventsResponse, error) {
	eventsDto, err := a.storage.GetWeekEvents(ctx, weekStart.WeekStart.AsTime())

	return toResposeModels(eventsDto), err
}

func (a *App) GetMonthEvents(ctx context.Context, monthStart *pb.GetEventsByMonthRequest) (*pb.GetEventsResponse, error) {
	eventsDto, err := a.storage.GetMonthEvents(ctx, monthStart.MonthStart.AsTime())

	return toResposeModels(eventsDto), err
}

func toResposeModels(events []models.Event) *pb.GetEventsResponse {
	results := &pb.GetEventsResponse{
		Event: make([]*pb.Event, len(events)),
	}

	for i, ev := range events {
		results.Event[i] = &pb.Event{
			Id:          ev.ID,
			Title:       ev.Title,
			StartEvent:  timestamppb.New(ev.StartEvent),
			EndEvent:    timestamppb.New(ev.EndEvent),
			Description: ev.Description,
			IdUser:      ev.IDUser,
		}
	}
	return results
}
