package services

import (
	"errors"

	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/repositories"
	"github.com/google/uuid"
)

type BoardService interface {
	//Interface
	CreateBoard(board *models.Board) error
	UpdateBoard(board *models.Board) error
	GetBoardPublicID(publicID string) (*models.Board, error)
	AddMemeber(boardPublicID string, userPublicIDs []string) error
}

type boardService struct {
	//Struct
	BoardRepo   repositories.BoardRepository
	userRepo    repositories.UserRepository
	BoardMember repositories.BoardMemberRepository
}

func NewBoardService(boardRepo repositories.BoardRepository, userRepo repositories.UserRepository, boardMember repositories.BoardMemberRepository,
) BoardService {
	return &boardService{boardRepo, userRepo, boardMember}
}

func (s *boardService) CreateBoard(board *models.Board) error {
	users, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("user not exist")
	}
	board.PublicID = uuid.New()
	board.OwnerID = users.InternalID
	return s.BoardRepo.CreateBoard(board)
}
func (s *boardService) UpdateBoard(board *models.Board) error {
	return s.BoardRepo.UpdateBoard(board)
}

func (s *boardService) GetBoardPublicID(publicID string) (*models.Board, error) {
	return s.BoardRepo.FindByPublicID(publicID)
}

func (s *boardService) AddMemeber(boardPublicID string, userPublicIDs []string) error {
	board, err := s.BoardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	var UserInternalIDs []uint
	for _, userPublicID := range userPublicIDs {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			return errors.New("user not found")
		}
		UserInternalIDs = append(UserInternalIDs, uint(user.InternalID))
	}
	//Members
	existingMember, err := s.BoardMember.GetMembers(string(board.PublicID.String()))
	if err != nil {
		return err
	}
	//quick check
	memberMap := make(map[uint]bool)
	for _, member := range existingMember {
		memberMap[uint(member.InternalID)] = true //memberMap[1] == True
	}

	var NewMemberIDs []uint
	for _, userID := range UserInternalIDs {
		if !memberMap[userID] {
			NewMemberIDs = append(NewMemberIDs, userID)
		}
	}
	if len(NewMemberIDs) == 0 {
		return nil
	}
	return s.BoardRepo.AddMember(uint(board.InternalID), NewMemberIDs)
}
