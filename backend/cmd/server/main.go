// @title StockApp API
// @version 1.0
// @description API para gestión de inventario y stock
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.stockapp.com/support
// @contact.email support@stockapp.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stockapp/backend/internal/config"
	"github.com/stockapp/backend/internal/handler"
	appMiddleware "github.com/stockapp/backend/internal/middleware"
	"github.com/stockapp/backend/internal/repository"
	"github.com/stockapp/backend/internal/service"

	_ "github.com/stockapp/backend/docs"
	"github.com/swaggo/http-swagger/v2"
)

func main() {
	// Load config
	cfg := config.Load()

	// Initialize database
	if err := repository.InitDB(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize services
	authService := service.NewAuthService(cfg.JWTSecret)
	categoriaService := service.NewCategoriaService()
	productoService := service.NewProductoService(cfg.UploadDir)
	movimientoService := service.NewMovimientoService()
	deviceService := service.NewDeviceService()
	// notificationService := service.NewNotificationService() // TODO: agregar NotificationHandler cuando esté listo

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	categoriaHandler := handler.NewCategoriaHandler(categoriaService)
	productoHandler := handler.NewProductoHandler(productoService)
	movimientoHandler := handler.NewMovimientoHandler(movimientoService)
	deviceHandler := handler.NewDeviceHandler(deviceService)
	reportHandler := handler.NewReportHandler()

	// Setup router
	r := chi.NewRouter()
	webDir := resolvePath("./web", "./backend/web")
	uploadDir := resolvePath(cfg.UploadDir, "./backend/uploads")

	// Custom CORS
	r.Use(appMiddleware.CORS)

	// Chi middleware
	r.Use(appMiddleware.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	// Static files (uploads)
	r.Handle("/uploads/*", http.StripPrefix("/uploads", http.FileServer(http.Dir(uploadDir))))

	// Web Admin Panel
	r.Handle("/web/*", http.StripPrefix("/web", http.FileServer(http.Dir(webDir))))
	// Serve index.html for /web
	r.Get("/web", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	})

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"),
	))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/v1/auth/register", authHandler.Register)
		r.Post("/api/v1/auth/login", authHandler.Login)
		
		// Health check endpoint - useful for Render health checks
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok"}`))
		})
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware(authService))

		// Categorías
		r.Route("/api/v1/categorias", func(r chi.Router) {
			r.Get("/", categoriaHandler.GetAll)
			r.Post("/", categoriaHandler.Create)
			r.Get("/{id}", categoriaHandler.GetByID)
			r.Put("/{id}", categoriaHandler.Update)
			r.Delete("/{id}", categoriaHandler.Delete)
		})

		// Productos
		r.Route("/api/v1/productos", func(r chi.Router) {
			r.Get("/", productoHandler.GetAll)
			r.Post("/", productoHandler.Create)
			r.Get("/stock-bajo", productoHandler.GetStockBajo)
			r.Get("/barcode/{codigo}", productoHandler.GetByBarcode)
			r.Get("/{id}", productoHandler.GetByID)
			r.Put("/{id}", productoHandler.Update)
			r.Post("/{id}/imagen", productoHandler.UploadImagen)
			r.Delete("/{id}", productoHandler.Delete)
		})

		// Movimientos
		r.Route("/api/v1/movimientos", func(r chi.Router) {
			r.Get("/", movimientoHandler.GetAll)
			r.Post("/", movimientoHandler.Create)
			r.Get("/producto/{producto_id}", movimientoHandler.GetByProducto)
		})

		// Devices (Push Notifications)
		r.Route("/api/v1/devices", func(r chi.Router) {
			r.Post("/", deviceHandler.Register)
			r.Get("/", deviceHandler.GetUserDevices)
			r.Delete("/{id}", deviceHandler.Delete)
		})

		// Dashboard
		r.Get("/api/v1/dashboard", movimientoHandler.GetDashboard)

		// Reportes
		r.Route("/api/v1/reportes", func(r chi.Router) {
			r.Get("/stock-categoria", reportHandler.GetStockPorCategoria)
			r.Get("/movimientos-fecha", reportHandler.GetMovimientosPorFecha)
			r.Get("/dashboard-avanzado", reportHandler.GetDashboardAvanzado)
		})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func resolvePath(paths ...string) string {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return paths[0]
}
