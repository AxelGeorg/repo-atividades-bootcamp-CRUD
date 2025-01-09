package repository_test

import (
	"app/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestProductRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "quantity", "code_value", "is_published", "expiration", "price", "id_warehouse"}).
		AddRow(1, "Corn Shoots", 244, "0009-1111", "0", "2022-01-08", 23.27, 1).
		AddRow(2, "Shrimp - Baby, Cold Water", 174, "49288-0877", "0", "2022-08-04", 52.12, 1)

	mock.ExpectQuery("SELECT id, name, quantity, code_value, is_published, expiration, price, id_warehouse FROM products").
		WillReturnRows(rows)

	repo := repository.NewRepositoryProductDB(db)
	products, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, products, 2, "Expected 2 products")

	assert.Equal(t, 1, products[0].Id)
	assert.Equal(t, "Corn Shoots", products[0].Name)
	assert.Equal(t, 244, products[0].Quantity)
	assert.Equal(t, "0009-1111", products[0].CodeValue)
	assert.False(t, products[0].IsPublished)
	assert.Equal(t, 23.27, products[0].Price)
	assert.Equal(t, 1, products[0].IdWarehouse)

	assert.Equal(t, 2, products[1].Id)
	assert.Equal(t, "Shrimp - Baby, Cold Water", products[1].Name)
	assert.Equal(t, 174, products[1].Quantity)
	assert.Equal(t, "49288-0877", products[1].CodeValue)
	assert.False(t, products[1].IsPublished)
	assert.Equal(t, 52.12, products[1].Price)
	assert.Equal(t, 1, products[1].IdWarehouse)

	assert.NoError(t, mock.ExpectationsWereMet())
}
