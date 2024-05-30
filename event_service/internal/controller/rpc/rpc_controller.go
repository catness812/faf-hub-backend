package rpc

import (
	"context"

	"github.com/catness812/faf-hub-backend/event_service/internal/models"
	"github.com/catness812/faf-hub-backend/event_service/internal/pb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IEventService interface {
	CreateEvent(event models.Event) (uint, error)
	UpdateEvent(eventID uint, event models.Event) error
	DeleteEvent(eventID uint) error
	GetEventByID(eventID uint) (models.Event, error)
	GetAllEvents() ([]models.Event, error)
}

type IRegistrationService interface {
	RegisterForEvent(registration models.Registration) error
	GetEventRegistrations(eventID uint) ([]models.Registration, error)
	GetUserEventRegistration(eventID uint, userID uint) (models.Registration, error)
	GetUserEvents(userID uint) ([]models.Event, error)
	UpdateRegistration(registration models.Registration) error
}

type Server struct {
	pb.EventServiceServer
	EventService        IEventService
	RegistrationService IRegistrationService
}

func (s *Server) CreateEvent(_ context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	newEvent := models.Event{
		Name:                req.Event.Name,
		StartDateTime:       req.Event.Start.AsTime(),
		EndDateTime:         req.Event.End.AsTime(),
		Location:            req.Event.Location,
		ApplicationDeadline: req.Event.Deadline.AsTime(),
		Cover:               req.Event.Cover,
		Description:         req.Event.Desc,
	}

	eventID, err := s.EventService.CreateEvent(newEvent)
	if err != nil {
		return nil, err
	}

	return &pb.CreateEventResponse{
		Message: "event created successfully",
		EventId: int32(eventID),
	}, nil
}

func (s *Server) EditEvent(_ context.Context, req *pb.EditEventRequest) (*pb.EditEventResponse, error) {
	event := models.Event{
		Name:                req.Event.Name,
		StartDateTime:       req.Event.Start.AsTime(),
		EndDateTime:         req.Event.End.AsTime(),
		Location:            req.Event.Location,
		ApplicationDeadline: req.Event.Deadline.AsTime(),
		Cover:               req.Event.Cover,
		Description:         req.Event.Desc,
	}

	if err := s.EventService.UpdateEvent(uint(req.Event.EventId), event); err != nil {
		return nil, err
	}

	return &pb.EditEventResponse{
		Message: "event updated successfully",
	}, nil
}

func (s *Server) DeleteEvent(_ context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	if err := s.EventService.DeleteEvent(uint(req.EventId)); err != nil {
		return nil, err
	}

	return &pb.DeleteEventResponse{
		Message: "event deleted successfully",
	}, nil
}

func (s *Server) GetEvent(_ context.Context, req *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	event, err := s.EventService.GetEventByID(uint(req.EventId))
	if err != nil {
		return nil, err
	}

	return &pb.GetEventResponse{
		Event: &pb.Event{
			EventId:  int32(event.ID),
			Name:     event.Name,
			Start:    timestamppb.New(event.StartDateTime),
			End:      timestamppb.New(event.EndDateTime),
			Location: event.Location,
			Deadline: timestamppb.New(event.ApplicationDeadline),
			Cover:    event.Cover,
			Desc:     event.Description,
		},
	}, nil
}

func (s *Server) GetEvents(_ context.Context, req *empty.Empty) (*pb.GetEventsResponse, error) {
	events, err := s.EventService.GetAllEvents()
	if err != nil {
		return nil, err
	}

	getEventsRes := make([]*pb.Event, len(events))

	for i := range getEventsRes {
		event := events[i]
		getEventsRes[i] = &pb.Event{
			EventId:  int32(event.ID),
			Name:     event.Name,
			Start:    timestamppb.New(event.StartDateTime),
			End:      timestamppb.New(event.EndDateTime),
			Location: event.Location,
			Deadline: timestamppb.New(event.ApplicationDeadline),
			Cover:    event.Cover,
			Desc:     event.Description,
		}
	}

	return &pb.GetEventsResponse{
		Events: getEventsRes,
	}, nil
}

func (s *Server) RegisterForEvent(_ context.Context, req *pb.RegisterForEventRequest) (*pb.RegisterForEventResponse, error) {
	newReg := models.Registration{
		EventID:         uint(req.Registration.EventId),
		UserID:          uint(req.Registration.UserId),
		FirstName:       req.Registration.FirstName,
		LastName:        req.Registration.LastName,
		Email:           req.Registration.Email,
		PhoneNumber:     int(req.Registration.PhoneNumber),
		AcademicGroup:   req.Registration.AcademicGroup,
		TeamMembers:     req.Registration.TeamMembers,
		ShirtSize:       req.Registration.ShirtSize,
		FoodPreferences: req.Registration.FoodPref,
		Motivation:      req.Registration.Motivation,
		Questions:       req.Registration.Questions,
	}

	if err := s.RegistrationService.RegisterForEvent(newReg); err != nil {
		return nil, err
	}

	return &pb.RegisterForEventResponse{
		Message: "event registration completed successfully",
	}, nil
}

func (s *Server) GetEventRegistrations(_ context.Context, req *pb.GetEventRegistrationsRequest) (*pb.GetEventRegistrationsResponse, error) {
	regs, err := s.RegistrationService.GetEventRegistrations(uint(req.EventId))
	if err != nil {
		return nil, err
	}

	getEventRegsRes := make([]*pb.EventRegistration, len(regs))

	for i := range getEventRegsRes {
		reg := regs[i]
		getEventRegsRes[i] = &pb.EventRegistration{
			EventId:       int32(reg.EventID),
			UserId:        int32(reg.UserID),
			FirstName:     reg.FirstName,
			LastName:      reg.LastName,
			Email:         reg.Email,
			PhoneNumber:   int32(reg.PhoneNumber),
			AcademicGroup: reg.AcademicGroup,
			TeamMembers:   reg.TeamMembers,
			ShirtSize:     reg.ShirtSize,
			FoodPref:      reg.FoodPreferences,
			Motivation:    reg.Motivation,
			Questions:     reg.Questions,
		}
	}

	return &pb.GetEventRegistrationsResponse{
		Registrations: getEventRegsRes,
	}, nil
}

func (s *Server) GetEventUserRegistration(_ context.Context, req *pb.GetEventUserRegistrationRequest) (*pb.GetEventUserRegistrationResponse, error) {
	reg, err := s.RegistrationService.GetUserEventRegistration(uint(req.EventId), uint(req.UserId))
	if err != nil {
		return nil, err
	}

	registration := &pb.EventRegistration{
		EventId:       int32(reg.EventID),
		UserId:        int32(reg.UserID),
		FirstName:     reg.FirstName,
		LastName:      reg.LastName,
		Email:         reg.Email,
		PhoneNumber:   int32(reg.PhoneNumber),
		AcademicGroup: reg.AcademicGroup,
		TeamMembers:   reg.TeamMembers,
		ShirtSize:     reg.ShirtSize,
		FoodPref:      reg.FoodPreferences,
		Motivation:    reg.Motivation,
		Questions:     reg.Questions,
	}

	return &pb.GetEventUserRegistrationResponse{
		Registration: registration,
	}, nil
}

func (s *Server) GetUserEvents(_ context.Context, req *pb.GetUserEventsRequest) (*pb.GetUserEventsResponse, error) {
	events, err := s.RegistrationService.GetUserEvents(uint(req.UserId))
	if err != nil {
		return nil, err
	}

	getUserEvents := make([]*pb.Event, len(events))

	for i := range getUserEvents {
		event := events[i]
		getUserEvents[i] = &pb.Event{
			EventId:  int32(event.ID),
			Name:     event.Name,
			Start:    timestamppb.New(event.StartDateTime),
			End:      timestamppb.New(event.EndDateTime),
			Location: event.Location,
			Deadline: timestamppb.New(event.ApplicationDeadline),
			Cover:    event.Cover,
			Desc:     event.Description,
		}
	}

	return &pb.GetUserEventsResponse{
		Events: getUserEvents,
	}, nil
}

func (s *Server) EditRegistration(_ context.Context, req *pb.EditRegistrationRequest) (*pb.EditRegistrationResponse, error) {
	registration := models.Registration{
		EventID:         uint(req.Registration.EventId),
		UserID:          uint(req.Registration.UserId),
		FirstName:       req.Registration.FirstName,
		LastName:        req.Registration.LastName,
		Email:           req.Registration.Email,
		PhoneNumber:     int(req.Registration.PhoneNumber),
		AcademicGroup:   req.Registration.AcademicGroup,
		TeamMembers:     req.Registration.TeamMembers,
		ShirtSize:       req.Registration.ShirtSize,
		FoodPreferences: req.Registration.FoodPref,
		Motivation:      req.Registration.Motivation,
		Questions:       req.Registration.Questions,
		Feedback:        req.Registration.Feedback,
	}

	if err := s.RegistrationService.UpdateRegistration(registration); err != nil {
		return nil, err
	}

	return &pb.EditRegistrationResponse{
		Message: "user registration updated successfully",
	}, nil
}
