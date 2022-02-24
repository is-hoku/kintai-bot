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

	e.GET("/auth", tokenController.Auth)                    // 認可
	e.GET("/oauth2/callback", tokenController.AuthCallback) // 認可コードが渡されアクセストークンなどを返す

	e.GET("/user/:email", userController.Show)
	e.POST("/user", userController.Create)

	e.GET("/token/:company_id", tokenController.Show)
	e.PUT("/token", tokenController.Update)

	e.POST("/dakoku/:freee_id", tokenController.Dakoku) // freee API

	serverPort := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(serverPort))
}
