package service

import (
	"github.com/focuscw0w/microservices/internal/user/errors"
	"github.com/focuscw0w/microservices/internal/user/repository"
	"github.com/focuscw0w/microservices/internal/user/security"
)

// Service dependency
type Service struct {
	userRepo repository.UserRepository
}

func NewService(repo repository.UserRepository) *Service {
	return &Service{userRepo: repo}
}

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
}

type UserDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (s *Service) SignUp(req *SignUpRequest) (*UserDTO, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.ErrEmptyCredentials
	}

	_, err := s.userRepo.GetUserByUsername(req.Username)
	if err == nil {
		return nil, errors.ErrUserAlreadyExist
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &repository.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	userDTO := &UserDTO{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
	}

	return userDTO, nil
}

func (s *Service) SignIn(req *SignInRequest) (*UserDTO, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.ErrEmptyCredentials
	}

	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	err = security.VerifyPassword(user.Password, req.Password)
	if err != nil {
		return nil, errors.ErrInvalidPassword
	}

	userDTO := &UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return userDTO, nil
}

func (s *Service) GetUsers() ([]*UserDTO, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	usersDTO := make([]*UserDTO, len(users))
	for i, u := range users {
		usersDTO[i] = &UserDTO{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
		}
	}

	return usersDTO, nil
}

func (s *Service) GetUser(id int) (*UserDTO, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	userDTO := &UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return userDTO, nil
}

func (s *Service) DeleteUser(id int) error {
	err := s.userRepo.DeleteUser(id)
	if err != nil {
		return errors.ErrDeleteUserFailed
	}

	return nil
}

func (s *Service) UpdateUser(id int, req UpdateUserRequest) (*UserDTO, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	err = s.userRepo.UpdateUser(user.ID, req.Username)
	if err != nil {
		return nil, errors.ErrUpdateUserFailed
	}

	userDTO := &UserDTO{
		ID:       user.ID,
		Username: req.Username,
		Email:    user.Email,
	}

	return userDTO, nil
}
