package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ExtraData struct {
	CallbackUrl string `json:"callbackUrl"`
}
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PacketPayload struct {
	ServiceCode         string    `json:"serviceCode"`
	MSISDN              string    `json:"MSISDN"`
	InvoiceNumber       string    `json:"invoiceNumber"`
	AccountNumber       string    `json:"accountNumber"`
	CountryCode         string    `json:"countryCode"`
	PayerTransactionID  string    `json:"payerTransactionID"`
	Amount              float64   `json:"amount"`
	HubID               string    `json:"hubID"`
	Narration           string    `json:"narration"`
	DatePaymentReceived time.Time `json:"datePaymentReceived"`
	ExtraData           ExtraData `json:"extraData"`
	CurrencyCode        string    `json:"currencyCode"`
	CustomerNames       string    `json:"customerNames"`
	PaymentMode         string    `json:"paymentMode"`
}
type Payload struct {
	Credentials Credential      `json:"credentials"`
	Packet      []PacketPayload `json:"packet"`
}
type PostpaymentPayload struct {
	CountryCode string  `json:"countryCode"`
	Function    string  `json:"function"`
	Payload     Payload `json:"payload"`
}

func main() {
	fmt.Println("we want to dissect cellulant global apis.so let's get going")

	//formulate credentials payload
	credentials := Credential{"sandboxuser_api", "ocTpFMlTTHZRhdi4IckBYvvHyPN5liEj"}
	extraData := ExtraData{CallbackUrl: "https://webhook.site/9d2b6845-2410-4cd4-80e2-d70403a84a17"}
	packetPayload := PacketPayload{ServiceCode: "UG-DSTV", MSISDN: "255744408022",
		InvoiceNumber: "", AccountNumber: "255744408022", CountryCode: "UG",
		PayerTransactionID: "21241512321323", Amount: 100, HubID: "",
		Narration: "Mobile Money Testing", DatePaymentReceived: time.Now(),
		ExtraData: extraData, CurrencyCode: "UGX", CustomerNames: "Irvin Kessler", PaymentMode: "STK_PUSH"}

	payload := Payload{Credentials: credentials, Packet: []PacketPayload{packetPayload}}

	postpaymentPayload := PostpaymentPayload{CountryCode: "UG", Function: "BEEP.postPayment", Payload: payload}
	jsonPayload, err := json.Marshal(postpaymentPayload)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Json Payload is ", string(jsonPayload))
	req, err := http.NewRequest(http.MethodPost, "https://globalapi.uat.tingg.africa/globalApi/v2/JSON/", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalf("error creating http request %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error initiating request %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading response body %v", err)

	}
	if resp.StatusCode == http.StatusOK {
		fmt.Println("response received successfully")
	} else {
		fmt.Printf("receive non-ok status: %s\n", resp.Status)
	}

	fmt.Println("body response is ", string(body))

}
