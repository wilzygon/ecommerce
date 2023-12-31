package paypal

import (
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wilzygon/ecommerce/domain/paypal"
	"github.com/wilzygon/ecommerce/infrastructure/handler/response"
)

type handler struct {
	useCasePayPal paypal.UseCase
	responser     response.API
}

func newHandler(ucp paypal.UseCase) handler {
	return handler{useCasePayPal: ucp}
}

func (h handler) Webhook(c echo.Context) error {
	//Recibe el body de la petición que nos están haciendo desde Paypal hacía nosotros, lo convertimos
	//en un slice de bytes
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return h.responser.BindFailed(err)
	}
	go func() {
		// Por medio de goroutines procesamos ese webhook
		err = h.useCasePayPal.ProcessRequest(c.Request().Header, body)
		if err != nil {
			log.Printf("useCasePayPal.ProcessRequest(): %v", err)
		}
	}()

	return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
}
