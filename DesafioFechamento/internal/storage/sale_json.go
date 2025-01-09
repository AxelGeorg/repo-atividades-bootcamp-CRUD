package storage

import (
	"app/internal"
	"encoding/json"
	"os"
)

func NewSalesStorage(saleJson string) *SalesStorage {
	return &SalesStorage{pathJson: saleJson}
}

type SalesStorage struct {
	pathJson string
}

type SalesJSON struct {
	Id        int `json:"id"`
	Quantity  int `json:"quantity"`
	ProductId int `json:"product_id"`
	InvoiceId int `json:"invoice_id"`
}

func (s *SalesStorage) FindAll() (l []internal.Sale, err error) {
	// open file
	file, err := os.Open(s.pathJson)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var salesJSON []SalesJSON
	err = json.NewDecoder(file).Decode(&salesJSON)
	if err != nil {
		return
	}

	// serialize
	for _, sale := range salesJSON {
		l = append(l, internal.Sale{
			Id: sale.Id,
			SaleAttributes: internal.SaleAttributes{
				Quantity:  sale.Quantity,
				ProductId: sale.ProductId,
				InvoiceId: sale.InvoiceId,
			},
		})
	}

	return
}
