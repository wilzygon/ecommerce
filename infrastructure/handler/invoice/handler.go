package invoice

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wilzygon/ecommerce/domain/invoice"
	"github.com/wilzygon/ecommerce/infrastructure/handler/response"
	"github.com/wilzygon/ecommerce/model"
)

type handler struct {
	useCase  invoice.UseCase
	response response.API
}

func newHandler(useCase invoice.UseCase) handler {
	return handler{useCase: useCase}
}

// MyShops returns the shops from a logged user
// MyShops consulta las facturas de la persona loggeada
func (h handler) MyShops(c echo.Context) error {
	userID, ok := c.Get("userID").(uuid.UUID) //Del token obtiene el ID
	if !ok {
		err := model.NewError()
		err.APIMessage = "couldn't parse ID user from token"
		err.StatusHTTP = http.StatusInternalServerError

		return h.response.Error(c, "c.Get().(uuid.UUID)", &err)
	}

	data, err := h.useCase.GetByUserID(userID)
	if err != nil {
		errResp := model.NewError()
		errResp.Err = err
		errResp.APIMessage = err.Error()

		return h.response.Error(c, "useCase.GetWhere()", &errResp)
	}

	return c.JSON(h.response.Ok(data))
}

func (h handler) GetAll(c echo.Context) error {
	data, err := h.useCase.GetAll()
	if err != nil {
		errResp := model.NewError()
		errResp.Err = err
		errResp.APIMessage = err.Error()

		return h.response.Error(c, "useCase.GetWhere()", &errResp)
	}

	return c.JSON(h.response.Ok(data))
}
