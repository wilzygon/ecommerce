package model

import "encoding/json"

type PayPalRequestValidator struct {
	// Headers
	AuthAlgo string `json:"auth_algo"` //algoritmo de autenticación
	CertURL  string `json:"cert_url"`  //url donde se encuentra el certificado que Paypal utilizó
	//para firmar esa petición
	TransmissionID   string `json:"transmission_id"`   //id de la transmisión
	TransmissionSig  string `json:"transmission_sig"`  //signature o firma de la transmisión
	TransmissionTime string `json:"transmission_time"` //hora de la transmisión

	// Body
	WebhookID    string          `json:"webhook_id"`
	WebhookEvent json.RawMessage `json:"webhook_event"`
}

type PayPalRequestData struct {
	EventType string `json:"event_type"`
	ID        string `json:"id"`
	Resource  struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		CustomID string `json:"custom_id"`
		Amount   struct {
			Value string `json:"value"`
		} `json:"amount"`
	} `json:"resource"`
}
