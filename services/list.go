package services

import (
	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/repositories"
	"github.com/google/uuid"
)

type ListWithOrder struct {
	Position []uuid.UUID
	Lists    []models.List
}

type ListService interface {
	//interface
	FetchByBoardID(boardPublicID string) (*ListWithOrder, error)
	FetchByID(id uint) (*models.List, error)
	FetchByPublicID(publicID string) (*models.List, error)
	CreateList(list *models.List) error
	UpdateList(list *models.List) error
	DeleteList(id uint) error
	UpdateListPosition(boardPublicID string, pos []uuid.UUID) error
}
type ListServices struct {
	//struct
	listRepo     repositories.ListRepository
	boardRepo    repositories.BoardRepository
	listPosition repositories.ListPositioRepository
}

func NewListService(listRepo repositories.ListRepository, boardRepo repositories.BoardRepository, listPosition repositories.ListPositioRepository) ListService {
	// return &ListServices{listRepo, boardRepo, listPosition}
}

// func (s *ListServices) CreateList(list *models.List) error {
// 	//
// }
