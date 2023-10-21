package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wilzygon/ecommerce/domain/user"
	"github.com/wilzygon/ecommerce/model"
)

type handler struct {
	useCase user.UseCase //UseCase es la interface que se conecta con el domino, y handler lo conoce
}

func newHandler(uc user.UseCase) handler { //Cuando creamos un nuevo handler le enviamos ese UseCase, ya
	return handler{useCase: uc} //la implementacion como tal del domino
}

func (h handler) Create(c echo.Context) error {
	m := model.User{}

	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.useCase.Create(&m); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, m)
}

func (h handler) GetAll(c echo.Context) error {
	users, err := h.useCase.GetAll()
	if err != nil {

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}
