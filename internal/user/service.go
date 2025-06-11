package user

import (
	"encoding/json"
	"errors"
	"mozho_chat/internal/models"
	"mozho_chat/internal/repository"
	"mozho_chat/internal/user/dto"
	"mozho_chat/pkg/auth"
	"gorm.io/datatypes"
)

type Service interface {
	Register(input dto.CreateUserRequest) (*dto.UserResponse, error)
	Login(input dto.LoginRequest) (string, error)
	GetProfile(userID string) (*dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) Service {
	return &userService{repo: repo}
}

func (s *userService) Register(input dto.CreateUserRequest) (*dto.UserResponse, error) {
	hashed, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	var profileJSON datatypes.JSON
	if input.Profile != nil {
		profileBytes, err := json.Marshal(input.Profile)
		if err != nil {
			return nil, err
		}
		profileJSON = profileBytes
	}

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashed,
		Profile:      profileJSON,
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Profile:  user.Profile,
	}, nil
}

func (s *userService) Login(input dto.LoginRequest) (string, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return "", err
	}

	if err := auth.CheckPasswordHash(input.Password, user.PasswordHash); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetProfile(userID string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Profile:  user.Profile,
	}, nil
}
