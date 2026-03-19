package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"leave-management-system/internal/config"
	"leave-management-system/internal/database"
	"leave-management-system/internal/handler"
	"leave-management-system/internal/middleware"
	"leave-management-system/internal/repository"
	"leave-management-system/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Println("Configuration loaded successfully")

	// Connect to database
	db := database.Connect(cfg)
	defer database.Close(db)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	leaveRepo := repository.NewLeaveRepository(db)
	balanceRepo := repository.NewBalanceRepository(db)
	leaveTypeRepo := repository.NewLeaveTypeRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, balanceRepo, leaveTypeRepo)
	userService := service.NewUserService(userRepo, balanceRepo, leaveTypeRepo)
	leaveService := service.NewLeaveService(leaveRepo, balanceRepo, leaveTypeRepo, userRepo)
	balanceService := service.NewBalanceService(balanceRepo, leaveTypeRepo, userRepo)
	dashboardService := service.NewDashboardService(leaveRepo, userRepo, leaveTypeRepo, db)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	leaveHandler := handler.NewLeaveHandler(leaveService)
	adminHandler := handler.NewAdminHandler(leaveService, userService, balanceService, dashboardService)

	// Setup router
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.LoggerMiddleware)
	router.Use(middleware.CORSMiddleware)

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

		// Public routes (no authentication required)
		api.HandleFunc("/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
		api.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")
		api.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST", "OPTIONS")
		api.HandleFunc("/departments", authHandler.GetDepartments).Methods("GET", "OPTIONS")

	// Protected routes (authentication required)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Auth routes
	protected.HandleFunc("/auth/me", authHandler.GetMe).Methods("GET", "OPTIONS")

	// Employee routes (all authenticated users)
	protected.HandleFunc("/leaves", leaveHandler.CreateLeave).Methods("POST", "OPTIONS")
	protected.HandleFunc("/leaves", leaveHandler.GetUserLeaves).Methods("GET", "OPTIONS")
	protected.HandleFunc("/leaves/{id}", leaveHandler.GetLeaveByID).Methods("GET", "OPTIONS")
	protected.HandleFunc("/leaves/{id}", leaveHandler.CancelLeave).Methods("DELETE", "OPTIONS")
	protected.HandleFunc("/leave-types", leaveHandler.GetLeaveTypes).Methods("GET", "OPTIONS")
	protected.HandleFunc("/profile/balance", leaveHandler.GetUserBalance).Methods("GET", "OPTIONS")

	// Admin routes (admin and manager only)
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.RequireAdminOrManager)

	admin.HandleFunc("/leaves", adminHandler.GetAllLeaves).Methods("GET", "OPTIONS")
	admin.HandleFunc("/leaves/{id}/approve", adminHandler.ApproveLeave).Methods("PUT", "OPTIONS")
	admin.HandleFunc("/leaves/{id}/reject", adminHandler.RejectLeave).Methods("PUT", "OPTIONS")
	admin.HandleFunc("/users", adminHandler.GetAllUsers).Methods("GET", "OPTIONS")
	admin.HandleFunc("/users", adminHandler.CreateUser).Methods("POST", "OPTIONS")
	admin.HandleFunc("/users/{id}", adminHandler.UpdateUser).Methods("PUT", "OPTIONS")
	admin.HandleFunc("/dashboard", adminHandler.GetDashboardStats).Methods("GET", "OPTIONS")
	admin.HandleFunc("/balances", adminHandler.AllocateBalance).Methods("POST", "OPTIONS")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		log.Printf("Environment: %s", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Server shutting down...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
