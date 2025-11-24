package controllers

import (
	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/services"
	"github.com/ADMex1/GoProject/utils"
	"github.com/gofiber/fiber/v2"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{service: s}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	board := new(models.Board)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadReq(ctx, "Failed to Read Request", err.Error())
	}
	if err := c.service.CreateBoard(board); err != nil {
		return utils.BadReq(ctx, "Failed to create Board", err.Error())
	}
	return utils.Success(ctx, "Board Created!", board)
}
