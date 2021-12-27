package controllers

import (
	"net/http"

	"github.com/is-hoku/kintai-bot/domain"
	"github.com/is-hoku/kintai-bot/interfaces/database"
	"github.com/is-hoku/kintai-bot/usecase"
	"github.com/labstack/echo"
)

type UserController struct {
	Interactor usecase.UserInteractor
}

func NewUserController(dbHandler database.DBHandler) *UserController {
	return &UserController{
		Interactor: usecase.UserInteractor{
			UserRepository: &database.UserRepository{
				DBHandler: dbHandler,
			},
		},
	}
}

func (controller *UserController) Create(c echo.Context) error {
	u := domain.User{}
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := controller.Interactor.Add(u); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, u)
}

func (controller *UserController) Show(c echo.Context) error {
	filter := c.Param("email")
	if filter == "" {
		return c.JSON(http.StatusBadRequest, "No Parameters")
	}
	user, err := controller.Interactor.UserByEmail(filter)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, user)
}
