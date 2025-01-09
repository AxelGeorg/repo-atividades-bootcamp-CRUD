package storage

import (
	"app/internal"
	"encoding/json"
	"os"
)

func NewInvoicesStorage(invoiceJson string) *InvoicesStorage {
	return &InvoicesStorage{pathJson: invoiceJson}
}

type InvoicesStorage struct {
	pathJson string
}

type InvoicesJSON struct {
	Id         int     `json:"id"`
	Datetime   string  `json:"datetime"`
	Total      float64 `json:"total"`
	CustomerId int     `json:"customer_id"`
}

func (s *InvoicesStorage) FindAll() (i []internal.Invoice, err error) {
	// open file
	file, err := os.Open(s.pathJson)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var invoicesJSON []InvoicesJSON
	err = json.NewDecoder(file).Decode(&invoicesJSON)
	if err != nil {
		return
	}

	// serialize
	for _, invo := range invoicesJSON {
		i = append(i, internal.Invoice{
			Id: invo.Id,
			InvoiceAttributes: internal.InvoiceAttributes{
				Datetime:   invo.Datetime,
				Total:      invo.Total,
				CustomerId: invo.CustomerId,
			},
		})
	}

	return
}
