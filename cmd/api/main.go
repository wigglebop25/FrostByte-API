package main

import (
	"log"
	"net/http"
	"os"

	"frostbyte-api/internal/config"
	"frostbyte-api/internal/database"
	"frostbyte-api/internal/handlers"
	"frostbyte-api/internal/repository"
	"frostbyte-api/internal/service"
	"frostbyte-api/internal/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Database
	database.Connect(cfg.DSN())

	// Initialize WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize Repositories
	userRepo := repository.NewUserRepository(database.DB)
	productRepo := repository.NewProductRepository(database.DB)
	categoryRepo := repository.NewCategoryRepository(database.DB)
	orderRepo := repository.NewOrderRepository(database.DB)
	roleRepo := repository.NewRoleRepository(database.DB)

	// Initialize Services
	authService := service.NewAuthService(userRepo, roleRepo, cfg.JWTSecret)
	productService := service.NewProductService(productRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, userRepo, hub)
	roleService := service.NewRoleService(roleRepo, userRepo)
	userService := service.NewUserService(userRepo, roleRepo)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	orderHandler := handlers.NewOrderHandler(orderService)
	roleHandler := handlers.NewRoleHandler(roleService)
	userHandler := handlers.NewUserHandler(userService)

	// Initialize Router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("FrostByte API is running!"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Use JSON middleware for all API routes
		r.Use(handlers.JSONMiddleware)

		// Public Routes
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(handlers.AuthMiddleware(authService))

			r.Route("/products", func(r chi.Router) {
				r.Get("/", productHandler.GetAll)
				r.Post("/", productHandler.Create)
				r.Get("/{id}", productHandler.GetByID)
				r.Put("/{id}", productHandler.Update)
				r.Delete("/{id}", productHandler.Delete)
			})

			r.Route("/categories", func(r chi.Router) {
				r.Get("/", categoryHandler.GetAll)
				r.Post("/", categoryHandler.Create)
				r.Get("/{id}", categoryHandler.GetByID)
				r.Put("/{id}", categoryHandler.Update)
				r.Delete("/{id}", categoryHandler.Delete)
				r.Post("/product", categoryHandler.AddProduct)
				r.Post("/product/remove", categoryHandler.RemoveProduct)
			})

			r.Route("/orders", func(r chi.Router) {
				// Admin-only routes
				r.Group(func(r chi.Router) {
					r.Use(handlers.AdminMiddleware(userService))
					r.Get("/analytics", orderHandler.GetAnalytics)
				})

				// Standard authenticated routes
				r.Get("/", orderHandler.GetAll)
				r.Post("/", orderHandler.Create)
				r.Get("/{id}", orderHandler.GetByID)
				r.Put("/{id}/status", orderHandler.UpdateStatus)
				r.Get("/user/{username}", orderHandler.GetByUser)
				r.Get("/role/{role}", orderHandler.GetByRole)
			})

			r.Route("/roles", func(r chi.Router) {
				r.Get("/", roleHandler.GetAll)
				r.Post("/create", roleHandler.Create)
				r.Get("/{id}", roleHandler.GetByID)
				r.Put("/update/{id}", roleHandler.Update)
				r.Delete("/{id}", roleHandler.Delete)
				r.Post("/add_permission", roleHandler.AddPermission)
				r.Post("/assign", roleHandler.AssignRole)
			})

			r.Route("/users", func(r chi.Router) {
				r.Post("/", userHandler.Create)
				r.Get("/search", userHandler.Search)
				r.Get("/", userHandler.GetAll)
				r.Get("/{id}", userHandler.GetByID)
				r.Put("/{id}", userHandler.Update)
				r.Delete("/{id}", userHandler.Delete)
			})
		})

		// WebSocket
		r.Get("/ws/orders", hub.ServeWs)
	})

	log.Printf("Server starting on port %s (HTTPS)", cfg.ServerPort)
	// Check if certs exist
	if _, err := os.Stat("server.crt"); err == nil {
		if err := http.ListenAndServeTLS(":"+cfg.ServerPort, "server.crt", "server.key", r); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	} else {
		log.Printf("Certificates not found, falling back to HTTP")
		if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}
}
