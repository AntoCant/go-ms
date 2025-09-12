package main

import (
	"fmt"
	httpin "go-ms/internal/adapters/in/http"
	"go-ms/internal/adapters/out/memory"
	"go-ms/internal/application"
	"net/http"
)

func main() {

	repo := memory.NewInMemoryProductRepo()
	uc := application.NewProductService(repo)
	h := httpin.NewProductHandler(uc)

	addr := ":8080"
	fmt.Println("listening on", addr)
	_ = http.ListenAndServe(addr, h.Routes())
}
