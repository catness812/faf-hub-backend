package service

import (
	"fmt"

	"github.com/catness812/faf-hub-backend/user_service/internal/models"
	"github.com/catness812/faf-hub-backend/user_service/internal/util"
	"github.com/gookit/slog"
)

type IUserRepository interface {
	SaveUser(user models.User) error
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(userID uint) (models.User, error)
	UpdateUser(user models.User) error
	CheckAdmin(userID uint) (bool, error)
}

type UserService struct {
	userRepository IUserRepository
}

func NewUserService(
	userRepo IUserRepository,
) *UserService {
	return &UserService{
		userRepository: userRepo,
	}
}

func (svc *UserService) CreateNew(user models.User) error {
	if _, err := svc.userRepository.GetUserByEmail(user.Email); err == nil {
		slog.Errorf("User already exists: %v", user.Email)
		return fmt.Errorf("user already exists")
	}

	if err := util.EncryptPass(&user); err != nil {
		slog.Errorf("Could not encrypt password: %v", err)
		return err
	}

	if err := svc.userRepository.SaveUser(user); err != nil {
		slog.Errorf("Could not create new user: %v", err)
		return err
	}

	slog.Infof("User successfully created: %s", user.Email)
	return nil
}

func (svc *UserService) LoginUser(email string, password string) (uint, error) {
	user, err := svc.userRepository.GetUserByEmail(email)
	if err != nil {
		slog.Errorf("Could not retrieve user: %v", err)
		return 0, err
	}

	if err := util.ValidatePassword(user.Password, password); err != nil {
		slog.Errorf("Could not validate password: %v", err)
		return 0, err
	}

	slog.Infof("User successfully logged in: %s", email)
	return user.ID, err
}

func (svc *UserService) GetUserByID(userID uint) (models.User, error) {
	user, err := svc.userRepository.GetUserByID(userID)
	if err != nil {
		slog.Errorf("Could not retrieve user: %v", err)
		return user, err
	}

	slog.Info("User successfully retrieved")
	return user, nil
}

func (svc *UserService) GoogleLogUser(email string, firstName string, lastName string) (uint, error) {
	user, err := svc.userRepository.GetUserByEmail(email)
	user.Email = email
	user.Password = "-"
	user.FirstName = firstName
	user.LastName = lastName

	if err != nil {
		if err.Error() == "user doesn't exist" {
			if err := svc.userRepository.SaveUser(user); err != nil {
				slog.Errorf("Could not create new user: %v", err)
				return 0, err
			}
			user, err := svc.userRepository.GetUserByEmail(email)
			if err != nil {
				slog.Errorf("Could not retrieve user: %v", err)
				return 0, err
			}
			slog.Infof("User successfully logged in: %s", email)
			return user.ID, err
		} else {
			slog.Errorf("Could not retrieve user: %v", err)
			return 0, err
		}
	} else {
		if err := svc.userRepository.UpdateUser(user); err != nil {
			slog.Errorf("Could not update user: %v", err)
			return 0, err
		}
		slog.Infof("User successfully logged in: %s", email)
		return user.ID, err
	}
}

func (svc *UserService) UpdateUser(userID uint, user models.UserInfo) error {
	newUser, err := svc.userRepository.GetUserByID(userID)
	if err != nil {
		slog.Errorf("Could not retrieve user: %v", err)
		return err
	}

	if user.Email != "" {
		newUser.Email = user.Email
	}
	if user.PhoneNumber != 0 {
		newUser.PhoneNumber = user.PhoneNumber
	}
	if user.FirstName != "" {
		newUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		newUser.LastName = user.LastName
	}
	if user.AcademicGroup != "" {
		newUser.AcademicGroup = user.AcademicGroup
	}

	if err := svc.userRepository.UpdateUser(newUser); err != nil {
		slog.Errorf("Could not update user: %v", err)
		return err
	}

	slog.Info("User successfully updated")
	return nil
}

func (svc *UserService) CheckAdmin(userID uint) (bool, error) {
	status, err := svc.userRepository.CheckAdmin(userID)
	if err != nil {
		slog.Errorf("Could not retrieve admin status: %v", err)
		return false, err
	}

	slog.Info("User admin status successfully retrieved")
	return status, nil
}
