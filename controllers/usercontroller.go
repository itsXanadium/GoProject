package controllers

import (
	"math"
	"strconv"

	"github.com/ADMex1/GoProject/models"
	"github.com/ADMex1/GoProject/services"
	"github.com/ADMex1/GoProject/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		return utils.BadReq(ctx, "Failed to Parse Data", err.Error())
	}

	if err := c.service.Register(user); err != nil {
		return utils.BadReq(ctx, "Registration Failed!", err.Error())
	}
	var UserRespons models.UserResponse
	_ = copier.Copy(&UserRespons, user)
	return utils.Success(ctx, "Registration success!", UserRespons)
}
func (l *UserController) Login(ctx *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request", err.Error())
	}
	user, err := l.service.Login(body.Email, body.Password)
	if err != nil {
		return utils.UnauthorizedAccess(ctx, "login failed", err.Error())
	}

	Token, _ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	RefreshToken, _ := utils.RefreshToken(user.InternalID)

	var UserRespons models.UserResponse
	_ = copier.Copy(&UserRespons, &user)

	return utils.Success(ctx, "Login Success", fiber.Map{
		"Token":         Token,
		"refresh_token": RefreshToken,
		"user":          UserRespons,
	})
	// return utils.Success(ctx, "Login Successful", fiber.Map{
	// 	"Data": user,
	// })
}
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.service.GetByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "User not Found", err.Error())
	}
	var UserRespons models.UserResponse
	err = copier.Copy(&UserRespons, &user)
	if err != nil {
		return utils.BadReq(ctx, "Internal Server Error", err.Error())
	}

	return utils.Success(ctx, "Data Found!", UserRespons)
}

func (c *UserController) FetchUserPaginated(ctx *fiber.Ctx) error {

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	users, total, err := c.service.FetchUsersPaginated(filter, sort, limit, offset)
	if err != nil {
		return utils.BadReq(ctx, "Failed to Fetch Data", err.Error())
	}
	var userResponse []models.User
	_ = copier.Copy(&userResponse, &users)
	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}

	if total == 0 {
		return utils.NotFoundPaginated(ctx, "Data Not Found", userResponse, meta)
	}

	return utils.SuccessPaginated(ctx, "Data Fetched", userResponse, meta)
}
