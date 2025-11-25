package main

import (
	"log"

	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/controllers"
	"github.com/ADMex1/GoProject/database/seeder"
	"github.com/ADMex1/GoProject/repositories"
	"github.com/ADMex1/GoProject/routes"
	"github.com/ADMex1/GoProject/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	config.DBConnect()

	seeder.AdminSeeder()
	app := fiber.New()

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)
	boardMemberRepo := repositories.NewMemberRepository()

	boardRepo := repositories.NewBoardRepo()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boardController := controllers.NewBoardController(boardService)

	routes.Setup(app, userController, boardController)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"API Status": "Running",
		})
	})
	port := config.AppConfig.AppPort
	log.Println("Server is running on port: ", port)
	log.Fatal(app.Listen(":" + port))

}
