package repositories

import (
	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/models"
	"github.com/google/uuid"
)

type ListRepository interface {
	//interface
	CreateList(list *models.List) error
	UpdateList(list *models.List) error
	DeleteList(id uint) error
	UpdatePosition(boardPublicID string, position []string) error
	FetchCardPosition(listPublicID string) ([]uuid.UUID, error)
	FetchByBoardID(boardID string) ([]models.List, error)
	FetchByID(id string) (*models.List, error)
}

type ListRepositorys struct {
	//struct

}

func NewListRepository() ListRepository {
	return &ListRepositorys{}
}

func (r *ListRepositorys) CreateList(list *models.List) error {
	return config.DB.Create(list).Error
}

func (r *ListRepositorys) UpdateList(list *models.List) error {
	return config.DB.Model(&models.List{}).Where("public_id = ?", list.PublicID).Updates(map[string]interface{}{
		"title": list.Title,
	}).Error
}

func (r *ListRepositorys) DeleteList(id uint) error {
	return config.DB.Delete(&models.List{}, id).Error
}

func (r *ListRepositorys) UpdatePosition(boardPublicID string, position []string) error {
	return config.DB.Model(&models.ListPosition{}).Where("board_internal_id = (Select internal_id FROM boards Where public_id = ?)", boardPublicID).Update("list_order", position).Error
}

func (r *ListRepositorys) FetchCardPosition(listPublicID string) ([]uuid.UUID, error) {
	var position models.CardPosition
	err := config.DB.Joins("JOIN lists on list.internal_id = card_positions.list_internal_id").
		Where("list.public_id = ?", listPublicID).Error
	return position.CardOrder, err
}

func (r *ListRepositorys) FetchByBoardID(boardID string) ([]models.List, error) {
	var list []models.List
	err := config.DB.Where("board_public_id = ?", boardID).Order("internal_id ASC").Find(&list).Error
	return list, err
}

func (r *ListRepositorys) FetchByID(publicID string) (*models.List, error) {
	var list models.List
	err := config.DB.Where("public_id = ?", publicID).First(&list).Error
	return &list, err
}
