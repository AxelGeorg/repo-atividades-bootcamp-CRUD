package repository_test

import (
	"testing"

	"app/internal"
	"app/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCustomersMySQL_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
	}
	defer db.Close()

	repo := repository.NewCustomersMySQL(db, nil)
	mock.ExpectExec("INSERT INTO customers").
		WithArgs("John", "Doe", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	customer := &internal.Customer{
		CustomerAttributes: internal.CustomerAttributes{
			FirstName: "John",
			LastName:  "Doe",
			Condition: 1,
		},
	}
	err = repo.Save(customer)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCustomersMySQL_GetTotalValues(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
	}
	defer db.Close()

	repo := repository.NewCustomersMySQL(db, nil)

	rows := sqlmock.NewRows([]string{"condition", "total_value"}).
		AddRow(1, 100.00).
		AddRow(0, 50.50)

	mock.ExpectQuery(`(?i)SELECT c.condition,\s*ROUND\(SUM\(s.quantity \* p.price\), 2\)`).
		WillReturnRows(rows)

	totalValues, err := repo.GetTotalValues()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(totalValues) != 2 {
		t.Errorf("expected 2 total values, got %d", len(totalValues))
	}
	if totalValues[0].TotalValue != 100.00 {
		t.Errorf("expected 100.00, got %v", totalValues[0].TotalValue)
	}
	if totalValues[1].TotalValue != 50.50 {
		t.Errorf("expected 50.50, got %v", totalValues[1].TotalValue)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCustomersMySQL_GetSpentMoreMoney(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not open db: %v", err)
	}
	defer db.Close()

	repo := repository.NewCustomersMySQL(db, nil)

	rows := sqlmock.NewRows([]string{"first_name", "last_name", "total_spent"}).
		AddRow("John", "Doe", 200.00).
		AddRow("Jane", "Doe", 150.00)
	mock.ExpectQuery("SELECT c.first_name, c.last_name, ROUND").
		WillReturnRows(rows)

	spentMoreMoney, err := repo.GetSpentMoreMoney()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(spentMoreMoney) != 2 {
		t.Errorf("expected 2 customers, got %d", len(spentMoreMoney))
	}
	if spentMoreMoney[0].FirstName != "John" || spentMoreMoney[0].LastName != "Doe" || spentMoreMoney[0].Amount != 200.00 {
		t.Errorf("unexpected data for first customer: got %+v", spentMoreMoney[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
