package infrastructure

import (
	"os"

	"kintai-bot/app/interfaces/controllers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Init() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userController := controllers.NewUserController(NewDBHandler())
	tokenController := controllers.NewTokenController(NewTokenDBHandler())

	e.GET("/auth", tokenController.Auth)
	e.GET("/oauth2/callback", tokenController.AuthCallback)
	e.GET("/refresh", tokenController.Refresh)

	e.GET("/user/:email", userController.Show)
	//e.GET("/user/:email/clock_in", userController.Show, auth.Auth())
	e.POST("/user", userController.Create)

	e.GET("/token/:company_id", tokenController.Show)
	e.PUT("/token", tokenController.Update)

	serverPort := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(serverPort))
}
