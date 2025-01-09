package handler

import (
	"app/internal"
	"app/platform/web/request"
	"app/platform/web/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func NewHandlerWarehouse(rp internal.RepositoryWarehouse) (h *HandlerWarehouse) {
	h = &HandlerWarehouse{
		rp: rp,
	}
	return
}

type HandlerWarehouse struct {
	rp internal.RepositoryWarehouse
}

type WarehouseJSON struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	Capacity  int    `json:"capacity"`
}

func (h *HandlerWarehouse) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// process
		// - find all warehouses
		warehouses, err := h.rp.FindAll()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, "internal server error: "+err.Error())
			return
		}

		// response
		// - serialize warehouses to JSON
		var warehouseResponses []WarehouseJSON
		for _, w := range warehouses {
			warehouseResponses = append(warehouseResponses, WarehouseJSON{
				Id:        w.Id,
				Name:      w.Name,
				Address:   w.Address,
				Telephone: w.Telephone,
				Capacity:  w.Capacity,
			})
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    warehouseResponses,
		})
	}
}

func (h *HandlerWarehouse) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - path parameter: id
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.JSON(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		// - find warehouse by id
		p, err := h.rp.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrRepositoryWarehouseNotFound):
				response.JSON(w, http.StatusNotFound, "warehouse not found")
			default:
				response.JSON(w, http.StatusInternalServerError, "internal server error"+err.Error())
			}
			return
		}

		// response
		// - serialize warehouse to JSON
		data := WarehouseJSON{
			Id:        p.Id,
			Name:      p.Name,
			Address:   p.Address,
			Telephone: p.Telephone,
			Capacity:  p.Capacity,
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

// RequestBodyWarehouseCreate is a request body for creating a warehouse.
type RequestBodyWarehouseCreate struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	Capacity  int    `json:"capacity"`
}

// Create creates a warehouse.
func (h *HandlerWarehouse) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - body
		var body RequestBodyWarehouseCreate
		err := request.JSON(r, &body)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, "invalid body")
			return
		}

		// process
		// - save warehouse
		wh := internal.Warehouse{
			Name:      body.Name,
			Address:   body.Address,
			Telephone: body.Telephone,
			Capacity:  body.Capacity,
		}
		err = h.rp.Save(&wh)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// response
		// - serialize warehouse to JSON
		data := WarehouseJSON{
			Id:        wh.Id,
			Name:      wh.Name,
			Address:   wh.Address,
			Telephone: wh.Telephone,
			Capacity:  wh.Capacity,
		}
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}
