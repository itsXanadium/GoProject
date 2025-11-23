package services

import (
	"errors"

	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/repositories"
	"github.com/ADMex1/GoProject/utils"
	"github.com/google/uuid"
)

type UserService interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByPublicID(PublicID string) (*models.User, error)
	FetchUsersPaginated(filter, sort string, limit, offset int) ([]models.User, int64, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(user *models.User) error {
	existingUser, _ := s.repo.FindByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return errors.New("email already registered")
	}

	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	user.Role = "user"
	user.PublicID = uuid.New()
	return s.repo.Create(user)
}

func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.repo.FindById(id)
}

func (s *userService) GetByPublicID(publicID string) (*models.User, error) {
	return s.repo.FindByPublicID(publicID)
}

func (s *userService) FetchUsersPaginated(filter, sort string, limit, offset int) ([]models.User, int64, error) {
	return s.repo.FetchAllWPagination(filter, sort, limit, offset)
}
