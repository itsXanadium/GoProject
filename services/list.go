package services

import (
	"errors"

	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/repositories"
	"github.com/ADMex1/GoProject/utils"
	"github.com/google/uuid"
)

type ListWithOrder struct {
	Position []uuid.UUID
	Lists    []models.List
}

type ListService interface {
	//interface
	FetchByBoardID(boardPublicID string) (*ListWithOrder, error)
	// FetchByID(id string) (*models.List, error)
	// FetchByPublicID(publicID string) (*models.List, error)
	// CreateList(list *models.List) error
	// UpdateList(list *models.List) error
	// DeleteList(id uint) error
	// UpdateListPosition(boardPublicID string, pos []uuid.UUID) error
}

type ListServices struct {
	//struct
	listRepo     repositories.ListRepository
	boardRepo    repositories.BoardRepository
	listPosition repositories.ListPositioRepository
}

func NewListService(listRepo repositories.ListRepository, boardRepo repositories.BoardRepository, listPosition repositories.ListPositioRepository) ListService {
	return &ListServices{listRepo, boardRepo, listPosition}
}

func (s *ListServices) FetchByBoardID(boardPublicID string) (*ListWithOrder, error) {
	//verify the board
	_, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return nil, errors.New("board not found")
	}
	pos, err := s.listPosition.FetchListOrder(boardPublicID)
	if err != nil {
		return nil, errors.New("unable to fetch list order: " + err.Error())
	}
	list, err := s.listRepo.FetchByBoardID(boardPublicID)
	if err != nil {
		return nil, errors.New("failed to fetch board: " + err.Error())
	}
	//Sort by pos
	orderedlist := utils.SortListByPos(list, pos)
	return &ListWithOrder{
		Position: pos,
		Lists:    orderedlist,
	}, nil
}

// func (s *ListServices) CreateList(list *models.List) error {
// }
