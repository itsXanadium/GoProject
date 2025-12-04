package services

import (
	"errors"
	"fmt"
	"sort"
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

func (s *CardServices) UpdateCard(card *models.Card, listPublicID string) error {
	oldCard, err := s.CardRepo.FetchCardPublicID(card.PublicID.String())
	if err != nil {
		return fmt.Errorf("card not found: %v", err)
	}
	newList, err := s.ListRepo.FetchByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %v", err)
	}
	trax := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			trax.Rollback()
			panic(r)
		}
	}()
	if oldCard.ListID != newList.InternalID {
		//Remove old list
		var oldPos models.CardPosition
		if err := trax.Where("list_internal_id = ?", oldCard.ListID).
			First(&oldPos).Error; err != nil {
			filtered := make(types.UUIDArray, 0, len(oldPos.CardOrder))
			for _, id := range oldPos.CardOrder {
				if id != oldCard.PublicID {
					filtered = append(filtered, id)
				}
			}
			if err := trax.Model(&models.CardPosition{}).Where("internal_id = ?", oldPos.InternalID).
				Update("card_order", types.UUIDArray(filtered)).Error; err != nil {
				trax.Rollback()
				return fmt.Errorf("failed to update old card position: %v", err)
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			trax.Rollback()
			return fmt.Errorf("failed to fetch old card position: %v", err)
		}
		var newPos models.CardPosition
		res := trax.Where("list_internal_id = ?", newList.InternalID).
			First(&newPos)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			newPos = models.CardPosition{
				PublicID:  uuid.New(),
				ListID:    newList.InternalID,
				CardOrder: types.UUIDArray{oldCard.PublicID},
			}
			if err := trax.Create(&newPos).Error; err != nil {
				trax.Rollback()
				return fmt.Errorf("failed to create card position for new list: %v", err)
			}
		} else if res.Error == nil {
			//append
			UpdateOrder := append(newPos.CardOrder, card.PublicID)
			if err := trax.Model(&models.CardPosition{}).
				Where("internal_id = ?", newPos.InternalID).
				Update("card_order", types.UUIDArray(UpdateOrder)).Error; err != nil {
				trax.Rollback()
				return fmt.Errorf("failed to update card position: %v", err)
			}
		} else {
			trax.Rollback()
			return fmt.Errorf("failed to get new card position: %v", res.Error)
		}
	}
	card.InternalID = oldCard.InternalID
	card.PublicID = oldCard.PublicID
	card.ListID = oldCard.ListID
	if err := trax.Save(card).Error; err != nil {
		trax.Rollback()
		return fmt.Errorf("failed to update card: %v", err)
	}
	return nil
}

func (s *CardServices) DeleteCard(id uint) error {
	return s.CardRepo.DeleteCard(id)
}

func (s *CardServices) FetchByListID(listPublicID string) ([]models.Card, error) {
	//Verify list
	list, err := s.ListRepo.FetchByPublicID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("list not found: %v", err)
	}
	//Fetch card pos
	pos, err := s.CardRepo.FetchCardPositionbyListID(list.InternalID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch card position: %v", err)
	}
	//Fetch all card in List
	cards, err := s.CardRepo.FindByListID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch the cards: %v", err)
	}
	//sorting
	if pos != nil && len(pos.CardOrder) > 0 {
		cards = sortCardByPosition(cards, pos.CardOrder)
	}
	return cards, nil
}

func sortCardByPosition(cards []models.Card, order []uuid.UUID) []models.Card {
	//Map for fast Search
	orderMap := make(map[uuid.UUID]int)
	for i, id := range order {
		orderMap[id] = i
	}
	defaultIndex := len(order)

	sort.SliceStable(cards, func(i, j int) bool {
		idxI, okI := orderMap[cards[i].PublicID]
		if !okI {
			idxI = defaultIndex
		}
		idxJ, okJ := orderMap[cards[j].PublicID]
		if !okJ {
			idxJ = defaultIndex
		}
		if idxI == idxJ {
			return cards[i].CreatedAt.Before(cards[j].CreatedAt)
		}
		return idxI < idxJ
	})
	return cards
}

func (s *CardServices) FetchByCardID(id uint) (*models.Card, error) {
	return s.CardRepo.FetchCardID(id)
}

func (s *CardServices) FetchByCardPublicID(publicID string) (*models.Card, error) {
	return s.CardRepo.FetchCardPublicID(publicID)
}
