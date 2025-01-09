package storage

import (
	"app/internal"
	"encoding/json"
	"os"
)

func NewProductsStorage(productJson string) *ProductsStorage {
	return &ProductsStorage{pathJson: productJson}
}

type ProductsStorage struct {
	pathJson string
}

type ProductsJSON struct {
	Id          int     `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (s *ProductsStorage) FindAll() (p []internal.Product, err error) {
	// open file
	file, err := os.Open(s.pathJson)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var productsJSON []ProductsJSON
	err = json.NewDecoder(file).Decode(&productsJSON)
	if err != nil {
		return
	}

	// serialize
	for _, prod := range productsJSON {
		p = append(p, internal.Product{
			Id: prod.Id,
			ProductAttributes: internal.ProductAttributes{
				Description: prod.Description,
				Price:       prod.Price,
			},
		})
	}

	return
}
