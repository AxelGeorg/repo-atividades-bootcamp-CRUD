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
	dsn := "user:user@tcp(127.0.0.1:3306)/my_db"
	a.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err = a.db.Ping(); err != nil {
		return err
	}

	// - repository
	rp := repository.NewRepositoryProductDB(a.db)

	// - handler
	hd := handler.NewHandlerProduct(rp)

	// router
	// - middlewares
	a.rt.Use(middleware.Logger)
	a.rt.Use(middleware.Recoverer)
	// - endpoints
	a.rt.Route("/products", func(r chi.Router) {
		// GET /products/{id}
		r.Get("/{id}", hd.GetById())
		// POST /products
		r.Post("/", hd.Create())
		// PUT /products/{id}
		r.Put("/{id}", hd.UpdateOrCreate())
		// PATCH /products/{id}
		r.Patch("/{id}", hd.Update())
		// DELETE /products/{id}
		r.Delete("/{id}", hd.Delete())
	})

	return
}

// Run runs the application.
func (a *ApplicationDefault) Run() (err error) {
	log.Println("Server is running on", a.addr)
	err = http.ListenAndServe(a.addr, a.rt)
	return
}
