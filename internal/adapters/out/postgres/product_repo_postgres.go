package postgres

import (
	"context"
	"errors"
	"go-ms/internal/core/domain"
	errorCore "go-ms/internal/core/errors"
	core "go-ms/internal/core/ports"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProductRepo es la implementación de ProductRepository usando PostgreSQL.
type ProductRepo struct {
	pool *pgxpool.Pool
}

// NewProductRepo constructor que recibe el pool de conexiones.
func NewProductRepo(pool *pgxpool.Pool) core.ProductRepository {
	return &ProductRepo{pool: pool}
}

// helper para contexto con timeout
func ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func (r *ProductRepo) Save(p domain.Product) error {
	const q = `
		INSERT INTO products (id, name, price, stock)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    price = EXCLUDED.price,
		    stock = EXCLUDED.stock`
	c, cancel := ctx()
	defer cancel()
	_, err := r.pool.Exec(c, q, p.Id, p.Name, p.Price, p.Stock)
	return err
}

func (r *ProductRepo) FindAll() ([]domain.Product, error) {
	const q = `SELECT id, name, price, stock FROM products`
	c, cancel := ctx()
	defer cancel()

	rows, err := r.pool.Query(c, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *ProductRepo) FindByIdProduct(id string) (domain.Product, error) {
	const q = `SELECT id, name, price, stock FROM products WHERE id = $1`
	c, cancel := ctx()
	defer cancel()

	var p domain.Product
	err := r.pool.QueryRow(c, q, id).Scan(&p.Id, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, errorCore.ErrInvalidProductID // definí ErrProductNotFound en core/errs
		}
		return domain.Product{}, err
	}
	return p, nil
}

func (r *ProductRepo) DeleteProductById(id string) error {
	const q = `DELETE FROM products WHERE id = $1`
	c, cancel := ctx()
	defer cancel()

	ct, err := r.pool.Exec(c, q, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errorCore.ErrProductNotFound
	}
	return nil
}
