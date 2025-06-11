package user

import (
	"errors"
	"mozho_chat/internal/models"
	"mozho_chat/internal/repository"
	"mozho_chat/internal/user/dto"
	"mozho_chat/pkg/auth"
)

type Service interface {
	Register(input dto.CreateUserRequest) (*dto.UserResponse, error)
	Login(input dto.LoginRequest) (string, error)
	GetProfile(userID string) (*dto.UserResponse, error)

	UpdateProfile(userID string, input dto.UpdateUserRequest) (*dto.UserResponse, error)
	UpdatePassword(userID string, input dto.UpdatePasswordRequest) error
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

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashed,
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

func (s *userService) UpdateProfile(userID string, input dto.UpdateUserRequest) (*dto.UserResponse, error) {
    user, err := s.repo.FindByID(userID)
	if err != nil {
        return nil, err
    }

    if input.Username != nil {
        user.Username = *input.Username
    }
    if input.Email != nil {
        user.Email = *input.Email
    }

	if err := s.repo.Update(user); err != nil {
        return nil, err
    }

    return &dto.UserResponse{
        ID:       user.ID.String(),
        Username: user.Username,
        Email:    user.Email,
        Profile:  user.Profile,
    }, nil
}

func (s *userService) UpdatePassword(userID string, input dto.UpdatePasswordRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := auth.CheckPasswordHash(input.CurrentPassword, user.PasswordHash); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashed, err := auth.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashed

	return s.repo.Update(user)
}
