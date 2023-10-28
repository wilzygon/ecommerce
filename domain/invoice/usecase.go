package invoice

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/wilzygon/ecommerce/model"
)

// Invoice implements UseCase
type Invoice struct {
	storage                    Storage
	storageInvoiceDetailReport StorageInvoiceDetailReport
}

// New returns a new Invoice
func New(s Storage, sidr StorageInvoiceDetailReport) Invoice {
	return Invoice{storage: s, storageInvoiceDetailReport: sidr}
}

// Create creates a model.Invoice
func (i Invoice) Create(po *model.PurchaseOrder) error {
	if err := po.Validate(); err != nil {
		return fmt.Errorf("invoice: %w", err)
	}
	//invoiceFromPurchaseOrder, Llenamos la estructura de factura y detalle factura de acuerdo
	//a la orden de compra que recibimos
	invoice, invoiceDetails, err := invoiceFromPurchaseOrder(po)
	if err != nil {
		return fmt.Errorf("%s %w", "invoiceFromPurchaseOrder()", err)
	}

	err = i.storage.Create(&invoice, invoiceDetails) //Le decimos al storage que crée el invoice y el invoiceDetails
	if err != nil {
		return fmt.Errorf("%s %w", "storage.Create()", err)
	}

	return nil
}

func (i Invoice) GetByUserID(userID uuid.UUID) (model.InvoicesReport, error) {
	invoicesHead, err := i.storageInvoiceDetailReport.HeadsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invoice: %w", err)
	}

	var invoicesReport model.InvoicesReport
	for _, invoiceHead := range invoicesHead {
		invoiceDetails, err := i.storageInvoiceDetailReport.AllDetailsByInvoiceID(invoiceHead.Invoice.ID)
		if err != nil {
			return nil, fmt.Errorf("%s %w", "storageInvoiceDetail.AllDetailsByInvoiceID()", err)
		}

		invoiceHead.InvoiceDetailsReport = invoiceDetails
		invoicesReport = append(invoicesReport, invoiceHead)
	}

	return invoicesReport, nil
}

// GetAll returns a model.Invoices according to filters and sorts
func (i Invoice) GetAll() (model.InvoicesReport, error) {
	invoices, err := i.storageInvoiceDetailReport.AllHead()
	if err != nil {
		return nil, fmt.Errorf("invoice: %w", err)
	}

	var invoicesReport model.InvoicesReport
	for _, v := range invoices {
		invoiceDetails, err := i.storageInvoiceDetailReport.AllDetailsByInvoiceID(v.Invoice.ID)
		if err != nil {
			return nil, fmt.Errorf("%s %w", "storageInvoiceDetailReport.AllDetailsByInvoiceID()", err)
		}

		v.InvoiceDetailsReport = invoiceDetails
		invoicesReport = append(invoicesReport, v)
	}

	return invoicesReport, nil
}

// invoiceFromPurchaseOrder recibe una orden de compra, nos devuelve un ecabezado, un detalle, y un error
func invoiceFromPurchaseOrder(po *model.PurchaseOrder) (model.Invoice, model.InvoiceDetails, error) {
	ID, err := uuid.NewUUID()
	if err != nil {
		return model.Invoice{}, nil, fmt.Errorf("%s %w", "uuid.NewUUID()", err)
	}

	invoice := model.Invoice{ //Creamos el encabezado
		ID:              ID,
		UserID:          po.UserID, //UserID que viene con la orden de compra
		PurchaseOrderID: po.ID,
		CreatedAt:       time.Now().Unix(),
	}

	var products model.ProductToPurchases        //Creamos un array de productos a comprar
	err = json.Unmarshal(po.Products, &products) //Hacemos Unmarshal y guardamos en products
	if err != nil {
		return model.Invoice{}, nil, fmt.Errorf("%s %w", "json.Unmarshal()", err)
	}

	var invoiceDetails model.InvoiceDetails //Creamos un array de detalle de productos a comprar
	for _, v := range products {
		detailID, err := uuid.NewUUID()
		if err != nil {
			return model.Invoice{}, nil, fmt.Errorf("%s %w", "uuid.NewUUID()", err)
		}

		detail := model.InvoiceDetail{ //Creamos el detalle
			ID:        detailID,
			InvoiceID: invoice.ID,
			ProductID: v.ProductID, //ProductID que viene con la orden de compra
			Amount:    v.Amount,
			UnitPrice: v.UnitPrice,
			CreatedAt: time.Now().Unix(),
		}

		invoiceDetails = append(invoiceDetails, detail)
	}

	return invoice, invoiceDetails, nil
}
