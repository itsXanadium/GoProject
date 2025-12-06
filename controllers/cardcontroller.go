package controllers

import (
	"time"

	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/services"
	"github.com/ADMex1/GoProject/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CardController struct {
	service services.CardService
}

func NewCardController(s *services.CardService) *CardController {
	return &CardController{service: *s}
}

func (c *CardController) CreateCard(ctx *fiber.Ctx) error {
	type CreateCardRequest struct {
		ListPubliCD string    `json:"list_id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		DueDate     time.Time `json:"DueDate"`
		Position    int       `json:"position"`
	}
	var req CreateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadReq(ctx, "Unable to Fetch Data", err.Error())
	}
	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		Duedate:     &req.DueDate,
		Position:    req.Position,
	}
	if err := c.service.CreateCard(card, req.ListPubliCD); err != nil {
		return utils.InternalServerError(ctx, "Unable to create card", err.Error())
	}
	return utils.Success(ctx, "card created", card)
}

func (c *CardController) UpdateCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	type UpdateCardReq struct {
		ListPubliCD string    `json:"list_id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		DueDate     time.Time `json:"DueDate"`
		Position    int       `json:"position"`
	}
	var req UpdateCardReq
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadReq(ctx, "Unable to Fetch Data", err.Error())
	}
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadReq(ctx, "INVALID ID", err.Error())
	}
	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		Duedate:     &req.DueDate,
		Position:    req.Position,
		PublicID:    uuid.MustParse(publicID),
	}
	if err := c.service.UpdateCard(card, req.ListPubliCD); err != nil {
		return utils.InternalServerError(ctx, "Unable to Update card", err.Error())
	}
	return utils.Success(ctx, "Card Updated!", card)
}
