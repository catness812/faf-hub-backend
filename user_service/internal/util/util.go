package util

import (
	"fmt"

	"github.com/catness812/faf-hub-backend/user_service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPass(user *models.User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	return nil
}

func ValidatePassword(hash string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("wrong password")
	}
	return nil
}
