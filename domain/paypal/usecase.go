package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/wilzygon/ecommerce/model"
)

const (
	ExpectedVerification = "SUCCESS"
	ExpectedStatus       = "completed"
)

const (
	EventTypeProduct = "PAYMENT.CAPTURE.COMPLETED"
)

type PayPal struct {
	useCasePurchaseOrder UseCasePurchaseOrder
	useCaseInvoice       UseCaseInvoice
}

// New recibe las dependencias
func New(ucpo UseCasePurchaseOrder, uci UseCaseInvoice) PayPal {
	return PayPal{
		useCasePurchaseOrder: ucpo,
		useCaseInvoice:       uci,
	}
}

// ProcessRequest recibe un header o los headers de la petición de Paypal, porque ahí vienen los
// datos de con qué se cifró, el certificado, etc, y el body que es el webhook
func (pp PayPal) ProcessRequest(header http.Header, body []byte) error {
	//Lo primero que hacemos es convertir el body y los headers en la estructura que espera Paypal
	//para validar el webhook
	payPalRequestValidator, payPalRequestData, err := pp.parsePayPalRequestAndData(header, body)
	if err != nil {
		errMsg := fmt.Errorf("%s %w", "pp.parsePayPalRequest()", err)
		log.Println(errMsg)
		return errMsg
	}

	err = pp.validate(&payPalRequestValidator)
	if err != nil {
		log.Println(err)
		return err
	}

	err = pp.processPayment(&payPalRequestData)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (pp PayPal) parsePayPalRequestAndData(headers http.Header, body []byte) (model.PayPalRequestValidator, model.PayPalRequestData, error) {
	data := model.PayPalRequestData{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return model.PayPalRequestValidator{}, model.PayPalRequestData{}, fmt.Errorf("%s %w", "json.Unmarshal()", err)
	}

	if data.EventType != EventTypeProduct {
		return model.PayPalRequestValidator{}, model.PayPalRequestData{}, fmt.Errorf("the event_type %q is not allowed", data.EventType)
	}

	return model.PayPalRequestValidator{
		AuthAlgo:         headers.Get("Paypal-Auth-Algo"),
		CertURL:          headers.Get("Paypal-Cert-Url"),
		TransmissionID:   headers.Get("Paypal-Transmission-Id"),
		TransmissionSig:  headers.Get("Paypal-Transmission-Sig"),
		TransmissionTime: headers.Get("Paypal-Transmission-Time"),
		WebhookEvent:     body,
		WebhookID:        os.Getenv("WEBHOOK_ID"),
	}, data, nil
}

func (pp PayPal) validate(p *model.PayPalRequestValidator) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, os.Getenv("VALIDATION_URL"), bytes.NewReader(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(os.Getenv("CLIENT_ID"), os.Getenv("SECRET_ID"))

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(response)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("PayPal response with status code %d, body: %s", response.StatusCode, string(body))
	}

	bodyMap := make(map[string]string)
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return err
	}

	if bodyMap["verification_status"] != ExpectedVerification {
		return fmt.Errorf("verification status is %s", bodyMap["verification_status"])
	}

	return nil
}

func (pp PayPal) processPayment(data *model.PayPalRequestData) error {
	if !strings.EqualFold(data.Resource.Status, ExpectedStatus) {
		return fmt.Errorf("el estado de la transacción: %q no es el estado esperado: %q", data.Resource.Status, ExpectedStatus)
	}

	ID, err := uuid.Parse(data.Resource.CustomID)
	if err != nil {
		return fmt.Errorf("%s %w", "uuid.Parse()", err)
	}

	order, err := pp.useCasePurchaseOrder.GetByID(ID)
	if err != nil {
		return fmt.Errorf("%s %w", "useCasePurchaseOrder.GetWhere()", err)
	}

	value, err := strconv.ParseFloat(data.Resource.Amount.Value, 64)
	if err != nil {
		return err
	}

	totalAmount := order.TotalAmount()
	if totalAmount != value {
		return fmt.Errorf("el valor recibido: %0.2f, es diferente al valor esperado %0.2f", value, totalAmount)
	}

	return pp.useCaseInvoice.Create(&order)
}
