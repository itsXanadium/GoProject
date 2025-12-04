package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/models/types"
	"github.com/ADMex1/GoProject/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardService interface {
	CreateCard(card *models.Card, listPublicID string) error
	UpdateCard(card *models.Card, listPublicID string) error
	DeleteCard(id uint) error
	FetchByListID(listPublicID string) ([]models.Card, error)
	FetchByCardID(id uint) (*models.Card, error)
	FetchByCardPublicID(publicID string) (*models.Card, error)
}

type CardServices struct {
	CardRepo repositories.CardRepository
	ListRepo repositories.ListRepository
	UserRepo repositories.UserRepository
}

func NewCardService(CardRepo repositories.CardRepository, ListRepo repositories.ListRepository, UserRepo repositories.UserRepository) CardService {
	return &CardServices{CardRepo, ListRepo, UserRepo}
}

func (s *CardServices) CreateCard(card *models.Card, listPublicID string) error {
	list, err := s.ListRepo.FetchByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %v", err)
	}

	card.ListID = list.InternalID
	//Creates UUID for the card if not exist yet
	if card.PublicID == uuid.Nil {
		card.PublicID = uuid.New()
	}
	card.CreatedAt = time.Now()

	//Transaction
	trax := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			trax.Rollback()
			panic(r)
		}
	}()

	//Store Card
	if err := trax.Create(card).Error; err != nil {
		trax.Rollback()
		return fmt.Errorf("failed to create card: %v", err)
	}

	//Update/Create Card position
	var position models.CardPosition
	if err := trax.Model(&models.CardPosition{}).
		Where("list_internal_id = ?", list.InternalID).
		First(&position).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Insert new
			position = models.CardPosition{
				PublicID:  uuid.New(),
				ListID:    list.InternalID,
				CardOrder: types.UUIDArray{card.PublicID},
			}
			if err := trax.Create(&position).Error; err != nil {
				trax.Rollback()
				return fmt.Errorf("failed to create card position: %v", err)
			}
		} else {
			//Add New card to position
			position.CardOrder = append(position.CardOrder, card.PublicID)
			if err := trax.Model(&models.CardPosition{}).
				Where("internal_id = ?", position.InternalID).
				Update("card_order", position.CardOrder).Error; err != nil {
				trax.Rollback()
				return fmt.Errorf("failed to update card position: %v", err)
			}
		}
		if err := trax.Commit().Error; err != nil {
			return fmt.Errorf("transaction commit failed: %v", err)
		}
	}
	return nil
}
