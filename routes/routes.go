package routes

import (
	"log"

	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/controllers"
	"github.com/ADMex1/GoProject/utils"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App, uc *controllers.UserController, bc *controllers.BoardController, lc *controllers.ListController, cc *controllers.CardController) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error unable to load .env")
	}

	app.Post("/v1/auth/register", uc.Register)
	app.Post("/v1/auth/login", uc.Login)

	//JWT ROUTE
	api := app.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey: []byte(config.AppConfig.JWTSecret),
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.UnauthorizedAccess(c, "Unauthorized Access!", err.Error())
		},
	}))

	userGroup := api.Group("/users")
	userGroup.Get("/page", uc.FetchUserPaginated) //Paginated
	userGroup.Get("/:id", uc.GetUser)             //Using Public_id to call
	userGroup.Put("/:id", uc.UserUpdate)
	userGroup.Delete("/:id", uc.DeleteUser)

	boardGroup := api.Group("/boards")
	boardGroup.Post("/create", bc.CreateBoard)
	boardGroup.Put("/:id", bc.UpdateBoard)
	boardGroup.Post("/:id/add/members", bc.AddBoardMember)
	boardGroup.Delete("/:id/remove/members", bc.RemoveBoardMember)
	boardGroup.Get("/:board_id/list", lc.FetchListOnBoard)
	boardGroup.Get("/my", bc.FetchMyBoardPaginated)

	listGroup := api.Group("/list")
	listGroup.Post("/create", lc.CreateList)
	listGroup.Put("/update/:id", lc.UpdateList)
	listGroup.Delete("/delete/:id", lc.DeleteList)

	cardGroup := api.Group("/card")
	cardGroup.Post("/create", cc.CreateCard)
	cardGroup.Put("/update/:id", cc.UpdateCard)
}
