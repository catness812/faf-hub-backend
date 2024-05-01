package rpc

import (
	"context"

	"github.com/catness812/faf-hub-backend/user_service/internal/models"
	"github.com/catness812/faf-hub-backend/user_service/internal/pb"
)

type IUserService interface {
	CreateNew(user models.User) error
	LoginUser(email string, password string) (uint, error)
	GetUserByID(userID uint) (models.User, error)
	GoogleLogUser(email string, firstName string, lastName string) (uint, error)
	UpdateUser(userID uint, user models.UserInfo) error
	CheckAdmin(userID uint) (bool, error)
}

type Server struct {
	pb.UserServiceServer
	UserService IUserService
}

func (s *Server) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	newUser := models.User{
		Email:         req.User.Email,
		Password:      req.User.Password,
		PhoneNumber:   int(req.User.PhoneNumber),
		FirstName:     req.User.FirstName,
		LastName:      req.User.LastName,
		AcademicGroup: req.User.AcademicGroup,
	}

	if req.User.Email == "board.faf@gmail.com" {
		newUser.Admin = true
	} else {
		newUser.Admin = false
	}

	if err := s.UserService.CreateNew(newUser); err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Message: "user created successfully",
	}, nil
}

func (s *Server) Login(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	userID, err := s.UserService.LoginUser(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Message: "user logged in successfully",
		UserId:  int32(userID),
	}, nil
}

func (s *Server) GetUser(_ context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.UserService.GetUserByID(uint(req.UserId))
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Email:         user.Email,
			PhoneNumber:   int32(user.PhoneNumber),
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			AcademicGroup: user.AcademicGroup,
		},
	}, nil
}

func (s *Server) GoogleAuth(_ context.Context, req *pb.GoogleAuthRequest) (*pb.LoginResponse, error) {
	userID, err := s.UserService.GoogleLogUser(req.Email, req.FirstName, req.LastName)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Message: "user logged in successfully",
		UserId:  int32(userID),
	}, nil
}

func (s *Server) UpdateUser(_ context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := models.UserInfo{
		Email:         req.User.Email,
		PhoneNumber:   int(req.User.PhoneNumber),
		FirstName:     req.User.FirstName,
		LastName:      req.User.LastName,
		AcademicGroup: req.User.AcademicGroup,
	}

	if err := s.UserService.UpdateUser(uint(req.UserId), user); err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{
		Message: "user updated successfully",
	}, nil
}

func (s *Server) CheckAdmin(_ context.Context, req *pb.CheckAdminRequest) (*pb.CheckAdminResponse, error) {
	status, err := s.UserService.CheckAdmin(uint(req.UserId))
	if err != nil {
		return nil, err
	}

	return &pb.CheckAdminResponse{
		Admin: status,
	}, nil
}
