package memory

import (
	"fmt"
	"go-ms/internal/core/domain"
	core "go-ms/internal/core/ports"
	"sync"
)

type InMemoryProductRepo struct {
	mu sync.RWMutex
	db map[string]domain.Product
}

func NewInMemoryProductRepo() core.ProductRepository {
	return &InMemoryProductRepo{db: map[string]domain.Product{}}
}

func (r *InMemoryProductRepo) Save(p domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.db[p.Id] = p
	return nil
}

func (r *InMemoryProductRepo) FindAll(limit, offset int) ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Product, 0, len(r.db))
	for _, v := range r.db {
		out = append(out, v)
	}
	return out, nil
}

func (r *InMemoryProductRepo) FindByID(productId string) (domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if product, isExists := r.db[productId]; isExists {
		return product, nil
	}

	return domain.Product{}, fmt.Errorf("product with id %s not found", productId)
}

func (r *InMemoryProductRepo) DeleteByID(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.db[id]; !ok {
		return fmt.Errorf("product with id %s not found", id)
		// o mejor: return errs.ProductNotFound (si ya definiste un error tipado en tu capa común)
	}

	delete(r.db, id)
	return nil
}
