package purchaseorder

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/wilzygon/ecommerce/domain/purchaseorder"
	"github.com/wilzygon/ecommerce/infrastructure/handler/middle"
	purchaseorderStorage "github.com/wilzygon/ecommerce/infrastructure/postgres/purchaseorder"
)

// NewRouter returns a router to handle model.PurchaseOrder requests
func NewRouter(e *echo.Echo, dbPool *pgxpool.Pool) {
	h := buildHandler(dbPool)

	authMiddleWare := middle.New()

	privateRoutes(e, h, authMiddleWare.IsValid)
}

func buildHandler(dbPool *pgxpool.Pool) handler {
	useCase := purchaseorder.New(purchaseorderStorage.New(dbPool))
	return newHandler(useCase)
}

// privateRoutes handle the routes that requires a token
func privateRoutes(e *echo.Echo, h handler, middlewares ...echo.MiddlewareFunc) {
	route := e.Group("/api/v1/private/purchase-orders", middlewares...)

	route.POST("", h.Create)
}
