package paypal

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"github.com/wilzygon/ecommerce/domain/invoice"
	"github.com/wilzygon/ecommerce/domain/paypal"
	"github.com/wilzygon/ecommerce/domain/purchaseorder"

	storageInvoice "github.com/wilzygon/ecommerce/infrastructure/postgres/invoice"
	storagePurchaseOrder "github.com/wilzygon/ecommerce/infrastructure/postgres/purchaseorder"
)

func NewRouter(e *echo.Echo, dbPool *pgxpool.Pool) {
	h := buildHandler(dbPool)

	publicRoutes(e, h)
}

func buildHandler(dbPool *pgxpool.Pool) handler {
	//Crea las dependencias
	//El adaptador de la órden de compra
	purchaseOrderUseCase := purchaseorder.New(storagePurchaseOrder.New(dbPool))
	//El adaptador de facturación
	invoiceUseCase := invoice.New(storageInvoice.New(dbPool), nil)
	//Crea un nuevo caso de uso de Paypal
	useCase := paypal.New(purchaseOrderUseCase, invoiceUseCase)

	return newHandler(useCase)
}

// publicRoutes handle the routes that not requires a validation of any kind to be use
func publicRoutes(e *echo.Echo, h handler) {
	//Recomendación
	//Cambiar la palabra paypal por otra palabra en lo ruta, para evitar ataques
	route := e.Group("/api/v1/public/paypal")

	route.POST("", h.Webhook)
}
