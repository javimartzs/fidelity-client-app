package main

import (
	"fidelity-client-app/config"
	"fidelity-client-app/database"
	"fidelity-client-app/handlers"
	"fidelity-client-app/services"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Cargar el fichero .env
	config.LoadEnv()
	// Conectarnos a base de datos
	DB := database.ConnectDB()

	authService := services.AuthService{DB: DB}
	authHandler := handlers.AuthHandler{
		AuthService: authService,
	}
	// Iniciar enrutador
	r := mux.NewRouter()

	// Rutas publicas
	r.HandleFunc("/api/v1/register", authHandler.RegisterNewUser).Methods("POST")
	r.HandleFunc("/api/v1/login", authHandler.LoginUser).Methods("POST")

	http.ListenAndServe(":8080", r)
}
