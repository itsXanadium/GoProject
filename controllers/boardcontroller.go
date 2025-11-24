package controllers

import (
	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/services"
	"github.com/ADMex1/GoProject/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{service: s}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	var userID uuid.UUID
	var errs error
	board := new(models.Board)
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadReq(ctx, "Failed to Read Request", err.Error())
	}

	userID, errs = uuid.Parse(claims["pub_id"].(string))

	if errs != nil {
		return utils.BadReq(ctx, "Failed to read Request", errs.Error())
	}
	board.OwnerPublicID = userID

	if err := c.service.CreateBoard(board); err != nil {
		return utils.BadReq(ctx, "Failed to create Board", err.Error())
	}
	return utils.Success(ctx, "Board Created!", board)
}
