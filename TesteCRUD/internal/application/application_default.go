package application

import (
	"app/internal/handler"
	"app/internal/repository"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

// NewApplicationDefault creates a new default application.
func NewApplicationDefault(addr, filePathStore string) (a *ApplicationDefault) {
	// default config
	defaultRouter := chi.NewRouter()
	defaultAddr := ":8080"
	if addr != "" {
		defaultAddr = addr
	}

	a = &ApplicationDefault{
		rt:            defaultRouter,
		addr:          defaultAddr,
		filePathStore: filePathStore,
	}
	return
}

// ApplicationDefault is the default application.
type ApplicationDefault struct {
	// rt is the router.
	rt *chi.Mux
	// addr is the address to listen.
	addr string
	// filePathStore is the file path to store.
	filePathStore string
	// db is the database connection.
	db *sql.DB
}

// TearDown tears down the application.
func (a *ApplicationDefault) TearDown() (err error) {
	return a.db.Close()
}

// SetUp sets up the application.
func (a *ApplicationDefault) SetUp() (err error) {
	dsn := "root:root@tcp(127.0.0.1:3308)/my_db3"
	a.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err = a.db.Ping(); err != nil {
		return err
	}

	rpWare := repository.NewRepositoryWarehouseDB(a.db)
	hdWare := handler.NewHandlerWarehouse(rpWare)

	rpProd := repository.NewRepositoryProductDB(a.db)
	hdProd := handler.NewHandlerProduct(rpProd, rpWare)

	// router
	// - middlewares
	a.rt.Use(middleware.Logger)
	a.rt.Use(middleware.Recoverer)
	// - endpoints
	a.rt.Route("/products", func(r chi.Router) {
		// GET /products/{id}
		r.Get("/", hdProd.GetAll())
		r.Get("/{id}", hdProd.GetById())
		r.Get("/warehouse/reportProducts", hdProd.GetReportProductsById())
		// POST /products
		r.Post("/", hdProd.Create())
		// PUT /products/{id}
		r.Put("/{id}", hdProd.UpdateOrCreate())
		// PATCH /products/{id}
		r.Patch("/{id}", hdProd.Update())
		// DELETE /products/{id}
		r.Delete("/{id}", hdProd.Delete())
	})

	a.rt.Route("/warehouse", func(r chi.Router) {
		r.Get("/", hdWare.GetAll())
		r.Get("/{id}", hdWare.GetById())
		r.Post("/", hdWare.Create())
	})

	return
}

// Run runs the application.
func (a *ApplicationDefault) Run() (err error) {
	log.Println("Server is running on", a.addr)
	err = http.ListenAndServe(a.addr, a.rt)
	return
}
