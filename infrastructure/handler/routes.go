package handler

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/wilzygon/ecommerce/infrastructure/handler/product"
	"github.com/wilzygon/ecommerce/infrastructure/handler/purchaseorder"
	"github.com/wilzygon/ecommerce/infrastructure/handler/user"
)

func InitRoutes(e *echo.Echo, dbPool *pgxpool.Pool) {
	health(e)

	//P
	product.NewRouter(e, dbPool)
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