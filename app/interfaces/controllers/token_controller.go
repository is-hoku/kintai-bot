package controllers

import (
	"net/http"
	"strconv"

	"kintai-bot/app/common"
	"kintai-bot/app/domain"
	"kintai-bot/app/interfaces/database"
	"kintai-bot/app/usecase"

	"github.com/labstack/echo"
)

type TokenController struct {
	Interactor usecase.TokenInteractor
}

func NewTokenController(dbHandler database.TokenDBHandler) *TokenController {
	return &TokenController{
		Interactor: usecase.TokenInteractor{
			TokenRepository: &database.TokenRepository{
				TokenDBHandler: dbHandler,
			},
		},
	}
}

func (controller *TokenController) Update(c echo.Context) error {
	t := domain.Token{}
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "Invalid Request"))
	}
	if err := controller.Interactor.Update(t); err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorResponse(500, "Could not update record"))
	}
	return c.JSON(http.StatusCreated, t)
}

func (controller *TokenController) Show(c echo.Context) error {
	filter, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewErrorResponse(400, "No Parameters"))
	}
	token, err := controller.Interactor.TokenByCompanyID(filter)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewErrorResponse(404, "Not Found"))
	}
	return c.JSON(http.StatusOK, token)
}
