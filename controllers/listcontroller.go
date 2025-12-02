package controllers

import (
	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/services"
	"github.com/ADMex1/GoProject/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ListController struct {
	service services.ListService
}

func NewListController(s *services.ListService) *ListController {
	return &ListController{service: *s}
}

func (c *ListController) CreateList(ctx *fiber.Ctx) error {
	list := new(models.List)
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadReq(ctx, "Failed to read request", err.Error())
	}
	if err := c.service.CreateList(list); err != nil {
		return utils.BadReq(ctx, "unable to create list", err.Error())

	}
	return utils.Success(ctx, "List created", list)
}

func (c *ListController) UpdateList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new(models.List)

	if err := ctx.BodyParser(list); err != nil {
		return utils.BadReq(ctx, "unable to parse data", err.Error())
	}
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadReq(ctx, "Invalid ID", err.Error())
	}
	existingList, err := c.service.FetchByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}
	list.InternalID = existingList.InternalID
	list.PublicID = existingList.PublicID

	if err := c.service.UpdateList(list); err != nil {
		return utils.BadReq(ctx, "Update Failed!", err.Error())
	}
	updatedList, err := c.service.FetchByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}
	return utils.Success(ctx, "List updated", updatedList)
}

func (c *ListController) FetchListOnBoard(ctx *fiber.Ctx) error {
	boardPublicID := ctx.Params("board_id")

	if _, err := uuid.Parse(boardPublicID); err != nil {
		return utils.BadReq(ctx, "Invalid ID", err.Error())
	}

	list, err := c.service.FetchByBoardID(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "List Not Found", err.Error())
	}
	return utils.Success(ctx, "List Fetched", list)
}

func (c *ListController) DeleteList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.NotFound(ctx, "INVALID ID", err.Error())
	}

	lists, err := c.service.FetchByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List Not Found", err.Error())
	}

	if err := c.service.DeleteList(uint(lists.InternalID)); err != nil {
		return utils.InternalServerError(ctx, "Failed to delete list", err.Error())
	}
	return utils.Success(ctx, "List deleted", publicID)
}
