package repositories

import (
	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/models"
)

type CardRepository interface {
	//Interface
	CreateCard(card *models.Card) error
	UpdateCard(card *models.Card) error
	DeleteCard(id uint) error
}

type CardRepositorys struct {
}

func NewCardRepository() CardRepository {
	return &CardRepositorys{}
}

func (r *CardRepositorys) CreateCard(card *models.Card) error {
	return config.DB.Create(card).Error
}
func (r *CardRepositorys) UpdateCard(card *models.Card) error {
	return config.DB.Save(card).Error
}

func (r *CardRepositorys) DeleteCard(id uint) error {
	return config.DB.Delete(&models.Card{}, id).Error
}
