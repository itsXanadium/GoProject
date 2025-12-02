package services

import (
	"errors"
	"fmt"

	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/models/types"
	"github.com/ADMex1/GoProject/repositories"
	"github.com/ADMex1/GoProject/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	UpdateListPosition(boardPublicID string, positions []uuid.UUID) error
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
	fmt.Println(pos)
	fmt.Println(list)
	//Sort by pos
	orderedlist := utils.SortListByPos(list, pos)
	return &ListWithOrder{
		Position: pos,
		Lists:    orderedlist,
	}, nil
}

func (s *ListServices) FetchByID(id uint) (*models.List, error) {
	return s.listRepo.FetchByID(id)
}

func (s *ListServices) FetchByPublicID(publicID string) (*models.List, error) {
	return s.listRepo.FetchByPublicID(publicID)
}

func (s *ListServices) CreateList(list *models.List) error {
	//Validate Board
	board, err := s.boardRepo.FindByPublicID(string(list.BoardPublicId.String()))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("board not found")
		}
		return fmt.Errorf("failed to fetch board: %v", err)
	}
	list.BoardInternalId = board.InternalID

	if list.PublicID == uuid.Nil {
		list.PublicID = uuid.New()
	}

	trax := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			trax.Rollback()
		}
	}()
	//Store new list
	if err := trax.Create(list).Error; err != nil {
		trax.Rollback()
		return fmt.Errorf("failed to create list: %v", err)
	}
	//update pos
	var pos models.ListPosition
	res := trax.Where("board_internal_id = ?", board.InternalID).First(&pos)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		//Create new if not exist
		pos := models.ListPosition{
			PublicID:  uuid.New(),
			BoardID:   board.InternalID,
			ListOrder: types.UUIDArray{list.PublicID},
		}
		if err := trax.Create(&pos).Error; err != nil {
			trax.Rollback()
			return fmt.Errorf("failed to create list position: %v", err)
		}
	} else if res.Error != nil {
		trax.Rollback()
		return fmt.Errorf("failed to create list position: %v", res.Error)

	}
	//add new id
	pos.ListOrder = append(pos.ListOrder, list.PublicID)
	//update to Db
	if err := trax.Model(&pos).Where("internal_id = ?", pos.InternalID).Update("list_order", pos.ListOrder).Error; err != nil {
		trax.Rollback()
		return fmt.Errorf("failed to update list position: %v", err)
	}

	//Commit trax
	if err := trax.Commit().Error; err != nil {
		return fmt.Errorf("transation commit failed: %v", err)
	}
	return nil
}

func (s *ListServices) UpdateList(list *models.List) error {
	return s.listRepo.UpdateList(list)
}

func (s *ListServices) DeleteList(id uint) error {
	//
	return s.listRepo.DeleteList(id)
}

func (s *ListServices) UpdateListPosition(boardPublicID string, positions []uuid.UUID) error {
	//Verify
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}
	//Fetch List Pos
	position, err := s.listPosition.FetchByBoard(board.PublicID.String())
	if err != nil {
		return errors.New("position not found")
	}
	//update list order
	position.ListOrder = positions
	return s.listPosition.UpdateListOrder(position)
}
