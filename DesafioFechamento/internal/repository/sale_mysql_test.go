package repository_test

import (
	"app/internal"
	"app/internal/repository"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestSalesMySQL_Save(t *testing.T) {
	t.Run("success - sale saved", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
		}
		defer db.Close()

		repo := repository.NewSalesMySQL(db, nil)

		mock.ExpectExec("INSERT INTO sales").
			WithArgs(10, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		sale := &internal.Sale{
			SaleAttributes: internal.SaleAttributes{
				Quantity:  10,
				ProductId: 1,
				InvoiceId: 1,
			},
		}
		err = repo.Save(sale)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("error - failed to save sale", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
		}
		defer db.Close()

		repo := repository.NewSalesMySQL(db, nil)

		mock.ExpectExec("INSERT INTO sales").
			WithArgs(10, 1, 1).
			WillReturnError(sql.ErrConnDone)

		sale := &internal.Sale{
			SaleAttributes: internal.SaleAttributes{
				Quantity:  10,
				ProductId: 1,
				InvoiceId: 1,
			},
		}
		err = repo.Save(sale)

		require.Error(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestSalesMySQL_FindAll(t *testing.T) {
	t.Run("success - sales fetched", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
		}
		defer db.Close()

		repo := repository.NewSalesMySQL(db, nil)

		rows := sqlmock.NewRows([]string{"id", "quantity", "product_id", "invoice_id"}).
			AddRow(1, 10, 1, 1).
			AddRow(2, 5, 2, 1)

		mock.ExpectQuery("SELECT `id`, `quantity`, `product_id`, `invoice_id` FROM sales").
			WillReturnRows(rows)

		sales, err := repo.FindAll()

		require.NoError(t, err)
		require.Len(t, sales, 2)
		require.Equal(t, sales[0].Quantity, 10)
		require.Equal(t, sales[1].ProductId, 2)
	})

	t.Run("error - failed to fetch sales", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
		}
		defer db.Close()

		repo := repository.NewSalesMySQL(db, nil)

		mock.ExpectQuery("SELECT `id`, `quantity`, `product_id`, `invoice_id` FROM sales").
			WillReturnError(sql.ErrConnDone)

		sales, err := repo.FindAll()

		require.Error(t, err)
		require.Empty(t, sales)
	})
}
