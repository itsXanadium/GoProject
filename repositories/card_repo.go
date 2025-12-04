package repositories

import (
	"fmt"
	"path/filepath"

	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/models"
	"gorm.io/gorm"
)

type CardRepository interface {
	//Interface
	CreateCard(card *models.Card) error
	UpdateCard(card *models.Card) error
	DeleteCard(id uint) error
	FetchCardID(id uint) (*models.Card, error)
	FetchCardPublicID(PublicID string) (*models.Card, error)
	FindByListID(listID string) ([]models.Card, error)
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

func (r *CardRepositorys) FetchCardID(id uint) (*models.Card, error) {
	var Card models.Card
	err := config.DB.Preload("Labels").Preload("Assigments").First(&Card, id).Error

	return &Card, err
}

func (r *CardRepositorys) FetchCardPublicID(publicID string) (*models.Card, error) {
	var Card models.Card
	if err := config.DB.Preload("Assignees.user", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("internal_id", "public_id", "name", "email")
	}).Preload("Attachment").Where("public_id = ?", publicID).First(&Card).Error; err != nil {
		return nil, err
	}

	baseUrl := config.AppConfig.APPURL

	for i := range Card.Attatchment {
		Card.Attatchment[i].FileURL = fmt.Sprintf("%s/files/%s",
			baseUrl,
			filepath.Base(Card.Attatchment[i].File),
		)
	}
	return &Card, nil
}

func (r *CardRepositorys) FindByListID(listID string) ([]models.Card, error) {
	var cards []models.Card

	err := config.DB.Joins(`JOIN lists on lists.internal_id = card.list_internal_id`).
		Where("lists.public_id = ?", listID).
		Order("position ASC").
		Find(&cards).Error
	return cards, err
}
