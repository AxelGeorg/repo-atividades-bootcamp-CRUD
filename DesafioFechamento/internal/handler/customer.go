package handler

import (
	"log"
	"net/http"

	"app/internal"

	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
)

// NewCustomersDefault returns a new CustomersDefault
func NewCustomersDefault(sv internal.ServiceCustomer) *CustomersDefault {
	return &CustomersDefault{sv: sv}
}

// CustomersDefault is a struct that returns the customer handlers
type CustomersDefault struct {
	// sv is the customer's service
	sv internal.ServiceCustomer
}

// CustomerJSON is a struct that represents a customer in JSON format
type CustomerJSON struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Condition int    `json:"condition"`
}

type TotalValueJSON struct {
	Condition  int     `json:"condition"`
	TotalValue float64 `json:"total_value"`
}

type SpentMoreMoneyJSON struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Amount    float64 `json:"amount"`
}

// GetAll returns all customers
func (h *CustomersDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...

		// process
		c, err := h.sv.FindAll()
		if err != nil {
			log.Println(err)
			response.Error(w, http.StatusInternalServerError, "error getting customers")
			return
		}

		// response
		// - serialize
		csJSON := make([]CustomerJSON, len(c))
		for ix, v := range c {
			csJSON[ix] = CustomerJSON{
				Id:        v.Id,
				FirstName: v.FirstName,
				LastName:  v.LastName,
				Condition: v.Condition,
			}
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "customers found",
			"data":    csJSON,
		})
	}
}

func (h *CustomersDefault) GetTotalValues() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		totalValues, err := h.sv.GetTotalValues()
		if err != nil {
			log.Println(err)
			response.Error(w, http.StatusInternalServerError, "error getting total values")
			return
		}

		tvJSON := make([]TotalValueJSON, len(totalValues))
		for ix, v := range totalValues {
			tvJSON[ix] = TotalValueJSON{
				Condition:  v.Condition,
				TotalValue: v.TotalValue,
			}
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "total values found",
			"data":    tvJSON,
		})
	}
}

func (h *CustomersDefault) GetSpentMoreMoney() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spentMoreMoney, err := h.sv.GetSpentMoreMoney()
		if err != nil {
			log.Println(err)
			response.Error(w, http.StatusInternalServerError, "error getting spent more money")
			return
		}

		smmJSON := make([]SpentMoreMoneyJSON, len(spentMoreMoney))
		for ix, s := range spentMoreMoney {
			smmJSON[ix] = SpentMoreMoneyJSON{
				FirstName: s.FirstName,
				LastName:  s.LastName,
				Amount:    s.Amount,
			}
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "spent more money found",
			"data":    smmJSON,
		})
	}
}

// RequestBodyCustomer is a struct that represents the request body for a customer
type RequestBodyCustomer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Condition int    `json:"condition"`
}

// Create creates a new customer
func (h *CustomersDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - body
		var reqBody RequestBodyCustomer
		err := request.JSON(r, &reqBody)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "error deserializing request body")
			return
		}

		// process
		// - deserialize
		c := internal.Customer{
			CustomerAttributes: internal.CustomerAttributes{
				FirstName: reqBody.FirstName,
				LastName:  reqBody.LastName,
				Condition: reqBody.Condition,
			},
		}
		// - save
		err = h.sv.Save(&c)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error saving customer")
			return
		}

		// response
		// - serialize
		cs := CustomerJSON{
			Id:        c.Id,
			FirstName: c.FirstName,
			LastName:  c.LastName,
			Condition: c.Condition,
		}
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "customer created",
			"data":    cs,
		})
	}
}
