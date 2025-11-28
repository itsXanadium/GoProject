package controllers

import (
	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/services"
	"github.com/ADMex1/GoProject/utils"
	"github.com/gofiber/fiber/v2"
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
