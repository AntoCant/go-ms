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
type ProductRepository struct {
	pool *pgxpool.Pool
}

// NewProductRepo constructor que recibe el pool de conexiones.
func NewProductRepo(pool *pgxpool.Pool) core.ProductRepository {
	return &ProductRepository{pool: pool}
}

// helper para contexto con timeout
func contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func (productRepository *ProductRepository) Save(product domain.Product) error {
	const saveQueryProduct = `
		INSERT INTO products (id, name, price, stock)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    price = EXCLUDED.price,
		    stock = EXCLUDED.stock
		WHERE products.name IS DISTINCT FROM EXCLUDED.name
		   OR products.price IS DISTINCT FROM EXCLUDED.price
		   OR products.stock IS DISTINCT FROM EXCLUDED.stock`

	ctx, cancel := contextWithTimeout()
	defer cancel()
	_, err := productRepository.pool.Exec(
		ctx,
		saveQueryProduct,
		product.Id,
		product.Name,
		product.Price,
		product.Stock)
	return err
}

func (productRepository *ProductRepository) FindAll(limit, offset int) ([]domain.Product, error) {
	const selectAllProductsQuery = `
		SELECT id, name, price, stock FROM products
		ORDER BY id
		LIMIT $1 OFFSET $2`
	ctx, cancel := contextWithTimeout()
	defer cancel()

	rows, err := productRepository.pool.Query(ctx, selectAllProductsQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.
			Scan(
				&product.Id,
				&product.Name,
				&product.Price,
				&product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (productRepository *ProductRepository) FindByID(id string) (domain.Product, error) {
	const selectByIdProductQuery = `SELECT id, name, price, stock FROM products WHERE id = $1`
	ctx, cancel := contextWithTimeout()
	defer cancel()

	var product domain.Product
	err := productRepository.pool.QueryRow(ctx, selectByIdProductQuery, id).
		Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.Stock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, errorCore.ErrInvalidProductID
		}
		return domain.Product{}, err
	}
	return product, nil
}

func (productRepository *ProductRepository) DeleteByID(id string) error {
	const deleteProductQuery = `DELETE FROM products WHERE id = $1`
	ctx, cancel := contextWithTimeout()
	defer cancel()

	execResult, err := productRepository.pool.Exec(ctx, deleteProductQuery, id)
	if err != nil {
		return err
	}
	if execResult.RowsAffected() == 0 {
		return errorCore.ErrProductNotFound
	}
	return nil
}
