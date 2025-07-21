package auth

import (
	"errors"
	"me-ai/internal/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWrongCredetials = "wrong credentials"
	ErrUserExists      = "user already exists"
)

type AuthService struct {
	UserRepository *models.UserRepository
}

func NewAuthService(userRepository *models.UserRepository) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (service *AuthService) Login(email, password string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser == nil {
		return "", errors.New(ErrWrongCredetials)
	}
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredetials)
	}
	return existedUser.Email, nil
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}
