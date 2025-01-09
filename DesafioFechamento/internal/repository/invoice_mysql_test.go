package repository_test

import (
	"app/internal"
	"app/internal/repository"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestInvoicesMySQL_Save(t *testing.T) {
	t.Run("success - invoice saved", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
		}
		defer db.Close()

		repo := repository.NewInvoicesMySQL(db, nil)

		mock.ExpectExec("INSERT INTO invoices").
			WithArgs(sqlmock.AnyArg(), 100.00, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		invoice := &internal.Invoice{
			InvoiceAttributes: internal.InvoiceAttributes{
				Datetime:   "2023-01-01 12:00:00",
				Total:      100.00,
				CustomerId: 1,
			},
		}
		err = repo.Save(invoice)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("error - failed to save invoice", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
		}
		defer db.Close()

		repo := repository.NewInvoicesMySQL(db, nil)

		mock.ExpectExec("INSERT INTO invoices").
			WithArgs(sqlmock.AnyArg(), 100.00, 1).
			WillReturnError(sql.ErrNoRows)

		invoice := &internal.Invoice{
			InvoiceAttributes: internal.InvoiceAttributes{
				Datetime:   "2023-01-01 12:00:00",
				Total:      100.00,
				CustomerId: 1,
			},
		}
		err = repo.Save(invoice)

		require.Error(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestInvoicesMySQL_FindAll(t *testing.T) {
	t.Run("success - invoices fetched", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
		}
		defer db.Close()

		repo := repository.NewInvoicesMySQL(db, nil)

		rows := sqlmock.NewRows([]string{"id", "datetime", "total", "customer_id"}).
			AddRow(1, "2023-01-01 12:00:00", 100.00, 1).
			AddRow(2, "2023-01-02 12:00:00", 200.00, 2)

		mock.ExpectQuery("SELECT `id`, `datetime`, `total`, `customer_id` FROM invoices").
			WillReturnRows(rows)

		invoices, err := repo.FindAll()

		require.NoError(t, err)
		require.Len(t, invoices, 2)
		require.Equal(t, invoices[0].Total, 100.00)
	})

	t.Run("error - failed to fetch invoices", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
		}
		defer db.Close()

		repo := repository.NewInvoicesMySQL(db, nil)

		mock.ExpectQuery("SELECT `id`, `datetime`, `total`, `customer_id` FROM invoices").
			WillReturnError(sql.ErrConnDone)

		invoices, err := repo.FindAll()

		require.Error(t, err)
		require.Empty(t, invoices)
	})
}
