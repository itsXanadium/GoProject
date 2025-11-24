package services

import (
	"errors"

	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/repositories"
	"github.com/google/uuid"
)

type BoardService interface {
	//Interface
	CreateBoard(board *models.Board) error
	UpdateBoard(board *models.Board) error
	GetBoardPublicID(publicID string) (*models.Board, error)
}

type boardService struct {
	//Struct
	BoardRepo repositories.BoardRepository
	userRepo  repositories.UserRepository
}

func NewBoardService(boardRepo repositories.BoardRepository, userRepo repositories.UserRepository,
) BoardService {
	return &boardService{boardRepo, userRepo}
}

func (s *boardService) CreateBoard(board *models.Board) error {
	users, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("user not exist")
	}
	board.PublicID = uuid.New()
	board.OwnerID = users.InternalID
	return s.BoardRepo.CreateBoard(board)
}
func (s *boardService) UpdateBoard(board *models.Board) error {
	return s.BoardRepo.UpdateBoard(board)
}

func (s *boardService) GetBoardPublicID(publicID string) (*models.Board, error) {
	return s.BoardRepo.FindByPublicID(publicID)
}
