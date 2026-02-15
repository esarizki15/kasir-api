package main

import (
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/middlewares"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port    string `mapstructure:"PORT"`
	DBConn  string `mapstructure:"DB_CONN"`
	API_KEY string `mapstructure:"API_KEY"`
}

func main() {

	// ===============================
	// Load ENV
	// ===============================
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:    viper.GetString("PORT"),
		DBConn:  viper.GetString("DB_CONN"),
		API_KEY: viper.GetString("API_KEY"),
	}

	// ===============================
	// Setup Database
	// ===============================
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// ===============================
	// Repository → Service → Handler
	// ===============================

	// Product
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Category
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Report
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	// ===============================
	// Middleware Setup
	// ===============================

	apiKeyMiddleware := middlewares.APIKey(config.API_KEY)

	// Helper untuk chaining middleware
	chain := func(h http.HandlerFunc, m ...middlewares.Middleware) http.HandlerFunc {
		for i := len(m) - 1; i >= 0; i-- {
			h = m[i](h)
		}
		return h
	}

	// Public middleware stack
	public := func(h http.HandlerFunc) http.HandlerFunc {
		return chain(h,
			middlewares.Logger,
			middlewares.CORS,
		)
	}

	// Protected middleware stack
	protected := func(h http.HandlerFunc) http.HandlerFunc {
		return chain(h,
			apiKeyMiddleware,
			middlewares.Logger,
			middlewares.CORS,
		)
	}

	// ===============================
	// Routes
	// ===============================

	// Public
	http.HandleFunc("/api/produk", public(productHandler.HandleProducts))
	http.HandleFunc("/api/categories", public(categoryHandler.HandleCategories))

	// Protected
	http.HandleFunc("/api/produk/", protected(productHandler.HandleProductByID))
	http.HandleFunc("/api/categories/", protected(categoryHandler.HandleCategoryByID))
	http.HandleFunc("/api/checkout", protected(transactionHandler.HandleCheckout))
	http.HandleFunc("/api/report/hari-ini", protected(reportHandler.HandleTodayReport))
	http.HandleFunc("/api/report", protected(reportHandler.HandleReportByDate))

	// ===============================
	// Start Server
	// ===============================
	addr := "0.0.0.0:" + config.Port
	fmt.Println("🚀 Server running di", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Gagal running server:", err)
	}
}
