package main

import (
	"localiza-compra/backend/internal/api/middleware"
	"localiza-compra/backend/internal/api/product"
	"localiza-compra/backend/internal/api/shoppinglist"
	"localiza-compra/backend/internal/api/stock"
	"localiza-compra/backend/internal/api/store"
	"localiza-compra/backend/internal/api/user"
	"localiza-compra/backend/internal/database"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	db := database.Connect()
	defer db.Close()

	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	storeRepo := store.NewRepository(db)
	storeService := store.NewService(storeRepo)
	storeHandler := store.NewHandler(storeService)

	stockItemRepo := stock.NewRepository(db)
	stockItemService := stock.NewService(stockItemRepo)
	stockItemHandler := stock.NewHandler(stockItemService)

	shoppinglistRepo := shoppinglist.NewRepository(db)
	shoppinglistService := shoppinglist.NewService(shoppinglistRepo)
	shoppinglistHandler := shoppinglist.NewHandler(shoppinglistService)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3005"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "database": "connected"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {

		// --- Sub-grupo de Rotas Públicas ---
		r.Group(func(r chi.Router) {
			r.Post("/login", userHandler.Login)
			r.Post("/users", userHandler.Create)
			r.Route("/products", func(r chi.Router) {
				r.Get("/", productHandler.GetAll)
				// Nossa nova rota de busca!
				r.Get("/search", productHandler.SearchByName)
			})
		})

		// --- Sub-grupo de Rotas Protegidas ---
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth) // Segurança geral para este grupo

			r.Get("/users/me", userHandler.GetMe)
			r.Get("/stores", storeHandler.GetAll)

			// Rotas de Listas de Compras do utilizador
			r.Route("/shopping-lists", func(r chi.Router) {
				r.Post("/", shoppinglistHandler.CreateList)
				r.Get("/", shoppinglistHandler.GetAllByUserID)
				r.Post("/{listID}/items", shoppinglistHandler.CreateItem)
				r.Get("/{listID}/items", shoppinglistHandler.GetAllItemsByListID)
				r.Patch("/{listID}/items/{itemID}", shoppinglistHandler.UpdateItemStatus)
			})

			// --- Sub-grupo de Rotas SÓ PARA ADMINS ---
			r.Group(func(r chi.Router) {
				r.Use(middleware.AdminOnly) // Segurança extra

				r.Post("/stores", storeHandler.Create)
				r.Get("/stores/{storeID}/products", stockItemHandler.GetAllByStoreId)
				r.Post("/stores/{storeID}/products/{productID}", stockItemHandler.Create)

				// Rotas de gestão de produtos para admins
				r.Post("/products", productHandler.Create)
				r.Put("/products/{id}", productHandler.Update)
				r.Delete("/products/{id}", productHandler.Delete)
			})
		})
	})

	log.Println("Iniciando API na porta 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
