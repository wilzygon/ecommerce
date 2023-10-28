package user

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wilzygon/ecommerce/domain/user"
	"github.com/wilzygon/ecommerce/infrastructure/handler/response"
	"github.com/wilzygon/ecommerce/model"
)

type handler struct {
	useCase   user.UseCase //UseCase es la interface que se conecta con el domino, y handler lo conoce
	responser response.API
}

func newHandler(uc user.UseCase) handler { //Cuando creamos un nuevo handler le enviamos ese UseCase, ya
	return handler{useCase: uc} //la implementacion como tal del domino
}

func (h handler) Create(c echo.Context) error {
	m := model.User{}

	if err := c.Bind(&m); err != nil {
		return h.responser.BindFailed(err)
		//return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.useCase.Create(&m); err != nil {
		return h.responser.Error(c, "useCase.Create()", err)
		//return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(h.responser.Created(m))
	//return c.JSON(http.StatusCreated, m)
}

// MySelf returns the data from my profile
func (h handler) MySelf(c echo.Context) error {
	ID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return h.responser.Error(c, "c.Get().(uuid.UUID)", errors.New("couldn´t parse the ID"))
	}

	u, err := h.useCase.GetByID(ID)
	if err != nil {
		return h.responser.Error(c, "useCase.GetWhere()", err)
	}

	return c.JSON(h.responser.Ok(u))
}

func (h handler) GetAll(c echo.Context) error {
	users, err := h.useCase.GetAll()
	if err != nil {
		return h.responser.Error(c, "useCase.GetAll()", err)
		//return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(h.responser.Ok(users))
	//return c.JSON(http.StatusOK, users)
}
