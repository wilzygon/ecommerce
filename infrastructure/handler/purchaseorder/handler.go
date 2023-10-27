package purchaseorder

import (
	"errors"

	"github.com/google/uuid"
	"github.com/wilzygon/ecommerce/domain/purchaseorder"
	"github.com/wilzygon/ecommerce/infrastructure/handler/response"
	"github.com/wilzygon/ecommerce/model"

	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase  purchaseorder.UseCase
	response response.API
}

func newHandler(useCase purchaseorder.UseCase) handler {
	return handler{useCase: useCase}
}

// Create handles the creation of a model.PurchaseOrder
func (h handler) Create(c echo.Context) error {
	m := model.PurchaseOrder{}
	if err := c.Bind(&m); err != nil {
		return h.response.BindFailed(err)
	}

	userID, ok := c.Get("userID").(uuid.UUID) //Obtiene el userID que hemos pasado por el contexto
	if !ok {                                  //del middleware que ha validado el token
		return h.response.Error(c, "c.Get().(uuid.UUID)", errors.New("can´t parse uuid"))
	}

	m.UserID = userID
	if err := h.useCase.Create(&m); err != nil {
		return h.response.Error(c, "useCase.Create()", err)
	}

	return c.JSON(h.response.Created(m))
}
