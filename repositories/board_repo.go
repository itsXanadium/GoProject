package repositories

import (
	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/models"
)

type BoardRepository interface {
	//Board interface
	CreateBoard(board *models.Board) error
}

type boardRepository struct {
	//Board Struct

}

func NewBoardRepo() BoardRepository {
	return &boardRepository{}
}

func (r *boardRepository) CreateBoard(board *models.Board) error {
	return config.DB.Create(board).Error
}
