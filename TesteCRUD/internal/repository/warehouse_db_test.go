package repository_test

import (
	"app/internal"
	"app/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestWarehouseRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "address", "telephone", "capacity"}).
		AddRow(1, "Main Warehouse", "221 Baker Street", "4555666", 100)

	mock.ExpectQuery("SELECT id, name, address, telephone, capacity FROM warehouses").
		WillReturnRows(rows)

	repo := repository.NewRepositoryWarehouseDB(db)
	warehouses, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, warehouses, 1)
	assert.Equal(t, 1, warehouses[0].Id)
	assert.Equal(t, "Main Warehouse", warehouses[0].Name)
	assert.Equal(t, "221 Baker Street", warehouses[0].Address)
	assert.Equal(t, "4555666", warehouses[0].Telephone)
	assert.Equal(t, 100, warehouses[0].Capacity)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWarehouseRepository_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "name", "address", "telephone", "capacity"}).
		AddRow(1, "Main Warehouse", "221 Baker Street", "4555666", 100)

	mock.ExpectQuery("SELECT id, name, address, telephone, capacity FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnRows(mockRows)

	repo := repository.NewRepositoryWarehouseDB(db)
	warehouse, err := repo.FindById(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, warehouse.Id)
	assert.Equal(t, "Main Warehouse", warehouse.Name)
	assert.Equal(t, "221 Baker Street", warehouse.Address)
	assert.Equal(t, "4555666", warehouse.Telephone)
	assert.Equal(t, 100, warehouse.Capacity)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWarehouseRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT COALESCE\\(MAX\\(id\\), 0\\) FROM warehouses").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))

	mock.ExpectExec("INSERT INTO warehouses").
		WithArgs(1, "New Warehouse", "123 Test St", "1234567", 150).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewRepositoryWarehouseDB(db)

	wh := internal.Warehouse{
		Name:      "New Warehouse",
		Address:   "123 Test St",
		Telephone: "1234567",
		Capacity:  150,
	}

	err = repo.Save(&wh)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
