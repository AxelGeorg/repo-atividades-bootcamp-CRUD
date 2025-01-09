package repository_test

import (
	"app/internal"
	"app/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestProductsMySQL_Save(t *testing.T) {
	t.Run("success - product saved", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
		}
		defer db.Close()

		repo := repository.NewProductsMySQL(db, nil)

		mock.ExpectExec("INSERT INTO products").
			WithArgs("New Product", 150.00).
			WillReturnResult(sqlmock.NewResult(1, 1))

		product := &internal.Product{
			ProductAttributes: internal.ProductAttributes{
				Description: "New Product",
				Price:       150.00,
			},
		}
		err = repo.Save(product)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

}

func TestProductsMySQL_FindAll(t *testing.T) {
	t.Run("success - products fetched", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
		}
		defer db.Close()

		repo := repository.NewProductsMySQL(db, nil)

		rows := sqlmock.NewRows([]string{"id", "description", "price"}).
			AddRow(1, "Product 1", 100.00).
			AddRow(2, "Product 2", 150.50)

		mock.ExpectQuery("SELECT `id`, `description`, `price` FROM products").
			WillReturnRows(rows)

		products, err := repo.FindAll()

		require.NoError(t, err)
		require.Len(t, products, 2)
		require.Equal(t, products[0].Description, "Product 1")
		require.Equal(t, products[1].Price, 150.50)
	})

}

func TestProductsMySQL_GetBestSelling(t *testing.T) {
	t.Run("success - best selling products fetched", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' occurred when opening a mock database connection: %s", err, err)
		}
		defer db.Close()

		repo := repository.NewProductsMySQL(db, nil)

		rows := sqlmock.NewRows([]string{"description", "total_sold"}).
			AddRow("Product 1", 100).
			AddRow("Product 2", 150)

		mock.ExpectQuery("SELECT p.description, SUM\\(s.quantity\\)").
			WillReturnRows(rows)

		bestSellingProducts, err := repo.GetBestSelling()

		require.NoError(t, err)
		require.Len(t, bestSellingProducts, 2)
		require.Equal(t, bestSellingProducts[0].Description, "Product 1")
		require.Equal(t, bestSellingProducts[1].Total, 150)
	})
}
