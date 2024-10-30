package main

import (
	"fidelity-client-app/config"
	"fidelity-client-app/database"
	"fidelity-client-app/handlers"
	"fidelity-client-app/services"
	"net/http"
)

func main() {

	// Cargar el fichero .env
	config.LoadEnv()

	// Conectarnos a base de datos
	DB := database.ConnectDB()

	// Inicializar servicios
	authService := services.AuthService{DB: DB}
	promotionService := services.PromotionService{DB: DB}

	// Inicializar handlers
	authHandler := handlers.AuthHandler{AuthService: authService}
	promotionHandler := handlers.PromotionHandler{PromotionService: &promotionService}

	// Iniciar enrutador
	mux := http.NewServeMux()

	// Rutas publicas de Auth
	mux.HandleFunc("/api/v1/register", authHandler.RegisterNewUser)
	mux.HandleFunc("/api/v1/login", authHandler.LoginUser)

	// Rutas para promociones
	mux.HandleFunc("/api/v1/promotions/create", promotionHandler.CreatePromotion)     // POST: Crear promoción
	mux.HandleFunc("/api/v1/promotions/active", promotionHandler.GetActivePromotions) // GET: Obtener promociones activas (con paginación)
	mux.HandleFunc("/api/v1/promotions", promotionHandler.GetPromotionByID)           // GET: Obtener promoción por ID
	mux.HandleFunc("/api/v1/promotions/update", promotionHandler.UpdatePromotion)     // PUT: Actualizar promoción
	mux.HandleFunc("/api/v1/promotions/delete", promotionHandler.DeletePromotion)     // DELETE: Eliminar promoción

	http.ListenAndServe(":8080", mux)
}
