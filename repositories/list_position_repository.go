package repositories

import (
	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/models"
	"github.com/google/uuid"
)

type ListPositioRepository interface {
	FetchByBoard(boardPublicID string) (*models.ListPosition, error)
	CreateOrUpdate(boardPublicID string, listOrder []uuid.UUID) error
	FetchListOrder(boardPublicID string) ([]uuid.UUID, error)
	UpdateListOrder(position *models.ListPosition) error
}

type ListPositionRepositorys struct {
}

func NewListPositionRepository() ListPositioRepository {
	return &ListPositionRepositorys{}
}

func (r *ListPositionRepositorys) FetchByBoard(boardPublicID string) (*models.ListPosition, error) {
	var position models.ListPosition
	err := config.DB.Joins("Join boards ON boards.internal_id = list_positions.board_internal_id").
		Where("boards.public_id = ?", boardPublicID).Error

	return &position, err
}

func (r *ListPositionRepositorys) CreateOrUpdate(boardPublicID string, listOrder []uuid.UUID) error {
	return config.DB.Exec(`
INSERT INTO list_position (board_internal_id, list_order)
SELECT internal_id, ? FROM boards Where public_id = ?  
ON CONFLICT (board_internal_id)
DO UPDATE SET list_order = EXCLUDE.list_order)`, listOrder, boardPublicID).Error
}

func (r *ListPositionRepositorys) FetchListOrder(boardPublicID string) ([]uuid.UUID, error) {
	pos, err := r.FetchByBoard(boardPublicID)
	if err != nil {
		return nil, err
	}
	return pos.ListOrder, err
}

func (r *ListPositionRepositorys) UpdateListOrder(position *models.ListPosition) error {
	return config.DB.Model(position).
		Where("internal_id = ?", position.InternalID).
		Update("list_order", position.ListOrder).Error
}
