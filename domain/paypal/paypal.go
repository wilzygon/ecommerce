package paypal

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/wilzygon/ecommerce/model"
)

type UseCase interface { //interface que va a utilizar el handler, simplemente va a procesar la petición
	ProcessRequest(header http.Header, body []byte) error
}

// UseCasePurchaseOrder va a permitir consultar la o0rden de compra por el ID
type UseCasePurchaseOrder interface {
	GetByID(ID uuid.UUID) (model.PurchaseOrder, error)
}

// UseCaseInvoice va a permitir crear la factura, porque cuando recibamos el pago éste caso de uso
// de Paypal va a crear una factura si todo salió bien, si el webhook es correcto
type UseCaseInvoice interface {
	Create(m *model.PurchaseOrder) error
}
