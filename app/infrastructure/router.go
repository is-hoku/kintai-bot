package infrastructure

import (
	"os"

	"github.com/is-hoku/kintai-bot/interfaces/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Init() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userController := controllers.NewUserController(NewDBHandler())

	e.GET("/user/:email", userController.Show)
	e.POST("/user", userController.Create)

	serverPort := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(serverPort))
}
