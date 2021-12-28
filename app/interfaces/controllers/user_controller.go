package controllers

import (
	"net/http"

	"kintai-bot/app/common"
	"kintai-bot/app/domain"
	"kintai-bot/app/interfaces/database"
	"kintai-bot/app/usecase"

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
		return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "Invalid Request"))
	}
	if err := controller.Interactor.Add(u); err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Could not create record"))
	}
	return c.JSON(http.StatusCreated, u)
}

func (controller *UserController) Show(c echo.Context) error {
	filter := c.Param("email")
	if filter == "" {
		return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "No Parameters"))
	}
	user, err := controller.Interactor.UserByEmail(filter)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewErrorResponse(404, "Not Found"))
	}
	return c.JSON(http.StatusOK, user)
}
