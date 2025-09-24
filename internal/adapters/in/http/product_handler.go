package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	core "go-ms/internal/core/ports"
)

type createProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}
type productResponse struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}
type updateProductRequest struct {
	Id    string  `json:"idProduct"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type ProductHandler struct {
	uc core.ProductUseCase
}

func NewProductHandler(uc core.ProductUseCase) *ProductHandler { return &ProductHandler{uc: uc} }

func (h *ProductHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	// /products y /products/{id}
	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.listProducts)   // GET /products
		r.Post("/", h.createProduct) // POST /products

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.getProductByID)   // GET /products/{id}
			r.Put("/", h.updateProduct)    // PUT /products/{id}
			r.Delete("/", h.deleteProduct) // DELETE /products/{id}
		})
	})

	return r
}

// ===== Handlers concretos =====

func (h *ProductHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var in createProductRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	p, err := h.uc.CreateProduct(in.Name, in.Price, in.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("error al crear", err)
		return
	}
	writeJSON(w, http.StatusCreated, productResponse{Id: p.Id, Name: p.Name, Price: p.Price, Stock: p.Stock})
}

func (h *ProductHandler) listProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Leer query params: ?limit=10&offset=20
	limit := 10 // valor por defecto
	offset := 0 // valor por defecto

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	// Llamar al caso de uso con paginación
	list, err := h.uc.GetAllProducts(limit, offset)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		fmt.Println("error del getAll products:", err)
		return
	}

	// Mapear al response DTO
	out := make([]productResponse, 0, len(list))
	for _, p := range list {
		out = append(out, productResponse{
			Id: p.Id, Name: p.Name, Price: p.Price, Stock: p.Stock,
		})
	}

	writeJSON(w, http.StatusOK, out)
}

func (h *ProductHandler) getProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id") // 👈 acá te da el {id} ya parseado
	if id == "" {
		http.NotFound(w, r)
		return
	}
	p, err := h.uc.GetProductById(id)
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, productResponse{Id: p.Id, Name: p.Name, Price: p.Price, Stock: p.Stock})
}

func (h *ProductHandler) updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	var in updateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	updated, err := h.uc.UpdateProduct(id, in.Name, in.Price, in.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, productResponse{Id: updated.Id, Name: updated.Name, Price: updated.Price, Stock: updated.Stock})
}

func (h *ProductHandler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	if err := h.uc.DeleteProduct(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ===== util pequeño para JSON =====
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
