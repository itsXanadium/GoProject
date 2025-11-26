package utils

import (
	"github.com/ADMex1/GoProject/models"
	"github.com/google/uuid"
)

func SortListByPos(lists []models.List, order []uuid.UUID) []models.List {
	orderedlist := make([]models.List, 0, len(order))

	listMap := make(map[uuid.UUID]models.List)
	for _, l := range lists {
		listMap[l.PublicID] = l
	}

	for _, id := range order {
		if list, ok := listMap[id]; ok {
			orderedlist = append(orderedlist, list)
		}
	}
	return orderedlist
}
