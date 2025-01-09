package storage

import (
	"app/internal"
	"encoding/json"
	"os"
)

func NewCustomersStorage(customerJson string) *CustomersStorage {
	return &CustomersStorage{pathJson: customerJson}
}

type CustomersStorage struct {
	pathJson string
}

type CustomersJSON struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Condition int    `json:"condition"`
}

func (s *CustomersStorage) FindAll() (c []internal.Customer, err error) {
	// open file
	file, err := os.Open(s.pathJson)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var customersJSON []CustomersJSON
	err = json.NewDecoder(file).Decode(&customersJSON)
	if err != nil {
		return
	}

	// serialize
	for _, cust := range customersJSON {
		c = append(c, internal.Customer{
			Id: cust.Id,
			CustomerAttributes: internal.CustomerAttributes{
				FirstName: cust.FirstName,
				LastName:  cust.LastName,
				Condition: cust.Condition,
			},
		})
	}

	return
}
