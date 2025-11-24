package repositories

import (
	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/models"
)

type BoardRepository interface {
	//Board interface
	CreateBoard(board *models.Board) error
	UpdateBoard(board *models.Board) error
	FindByPublicID(publicID string) (*models.Board, error)
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

func (r *boardRepository) UpdateBoard(board *models.Board) error {
	return config.DB.Model(&models.Board{}).Where("public_id = ?", board.PublicID).Updates(map[string]interface{}{
		"title":       board.Title,
		"description": board.Description,
		"due_date":    board.Duedate,
	}).Error
}

func (r *boardRepository) FindByPublicID(publicID string) (*models.Board, error) {
	var board models.Board
	err := config.DB.Where("public_id = ?", publicID).First(&board).Error
	return &board, err
}
