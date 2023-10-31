package handler

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/wilzygon/ecommerce/infrastructure/handler/invoice"
	"github.com/wilzygon/ecommerce/infrastructure/handler/login"
	"github.com/wilzygon/ecommerce/infrastructure/handler/paypal"
	"github.com/wilzygon/ecommerce/infrastructure/handler/product"
	"github.com/wilzygon/ecommerce/infrastructure/handler/proveedor"
	"github.com/wilzygon/ecommerce/infrastructure/handler/purchaseorder"
	"github.com/wilzygon/ecommerce/infrastructure/handler/user"
)

func InitRoutes(e *echo.Echo, dbPool *pgxpool.Pool) {
	health(e)

	//I
	invoice.NewRouter(e, dbPool)
	//L
	login.NewRouter(e, dbPool)

	//P
	paypal.NewRouter(e, dbPool)
	product.NewRouter(e, dbPool)
	proveedor.NewRouter(e, dbPool)
	purchaseorder.NewRouter(e, dbPool)

	//U
	user.NewRouter(e, dbPool)
}

func health(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(
			http.StatusOK,
			map[string]string{
				"time":         time.Now().String(),
				"message":      "Hola Mundo",
				"service_name": "",
			},
		)
	})
}
