package controllers

import (
	"math"
	"strconv"

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

func (c *BoardController) AddBoardMember(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	var userIds []string
	if err := ctx.BodyParser(&userIds); err != nil {
		return utils.BadReq(ctx, "failed to parse data", err.Error())
	}
	if err := c.service.AddMemeber(publicID, userIds); err != nil {
		return utils.BadReq(ctx, "Failed to add new member!", err.Error())
	}
	return utils.Success(ctx, "New Member Added!", nil)
}

func (c *BoardController) RemoveBoardMember(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadReq(ctx, "Failed to Parse Data", err.Error())
	}
	if err := c.service.RemoveMember(publicID, userIDs); err != nil {
		return utils.BadReq(ctx, "Failed to Remove Members from the Board!", err.Error())
	}
	return utils.Success(ctx, "Members has Been Removed", nil)
}

func (c *BoardController) FetchMyBoardPaginated(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["pub_id"].(string)

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	boards, total, err := c.service.FetchAllPaginatedViaUser(userID, filter, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch Board Data", err.Error())
	}

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}
	return utils.SuccessPaginated(ctx, "Data Fetched", boards, meta)
}
