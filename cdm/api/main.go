package main

import (
	"fmt"
	httpin "go-ms/internal/adapters/in/http"
	"go-ms/internal/adapters/out/postgres"
	core "go-ms/internal/core/ports"
	pgconn "go-ms/internal/infra/DBs/postgres"
	"os"

	// "go-ms/internal/adapters/out/memory"
	"go-ms/internal/application"
	"net/http"
)

func main() {

	// ==============================
	// 1. Configuración de infraestructura
	// ==============================
	// Leer el DSN desde variable de entorno (PG_DSN),
	// o usar un valor por defecto (útil en dev/local).
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		dsn = "postgres://admin:admin@localhost:5433/productsdb?sslmode=disable"
	}

	// Crear el pool de conexiones con pgxpool
	pool, err := pgconn.NewPool(dsn)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	// ==============================
	// 2. Crear los adaptadores de salida
	// ==============================
	// Usamos Postgres como implementación concreta de ProductRepository
	var repo core.ProductRepository = postgres.NewProductRepo(pool)

	// 👇 Repo alternativo en memoria (comentado, pero disponible para tests rápidos)
	// repo := memory.NewInMemoryProductRepo()

	// ==============================
	// 3. Crear el caso de uso (application service)
	// ==============================
	productService := application.NewProductService(repo)

	// ==============================
	// 4. Crear el adaptador de entrada (HTTP handler)
	// ==============================
	productHandler := httpin.NewProductHandler(productService)

	// ==============================
	// 5. Levantar el servidor HTTP
	// ==============================
	serverAddress := ":8080"
	fmt.Println("listening on", serverAddress)
	_ = http.ListenAndServe(serverAddress, productHandler.Routes())
}
