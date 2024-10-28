package handlers

import (
	"encoding/json"
	"fidelity-client-app/models"
	"fidelity-client-app/services"

	"net/http"
)

type AuthHandler struct {
	AuthService services.AuthService
}

// RegisterNewUser (controlador para registrar un nuevo usuario)
func (h *AuthHandler) RegisterNewUser(w http.ResponseWriter, r *http.Request) {

	// Verificamos que el metodo sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Creamos una estructura para decodificar los campos
	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		BirthDate   string `json:"birth_date"`
		Gender      string `json:"gender"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PassConfirm string `json:"pass_confirm"`
	}

	// Decodificamos los campos de la request
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Creamos un nuevo usuario con los datos proporcionados
	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		BirthDate: input.BirthDate,
		Gender:    input.Gender,
		Email:     input.Email,
		Password:  input.Password,
	}

	// Ejecutamos el servicio RegisterNewUser
	if err := h.AuthService.RegisterNewUser(&user, input.PassConfirm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Devolvemos un mensaje de exito si el registro funcionó
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "usuario registrado correctamente"})
}

// LoginUser (controlador para el login de un usuario existente)
func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {

	// Verificamos que el metodo sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Creamos una estructura para decodificar
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decodificamos el json del formulario
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ejecutamos el servicio de LoginUser
	user, err := h.AuthService.LoginUser(input.Email, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Generamos un Json Web Token de autenticacion para el usuario
	token, err := h.AuthService.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Devolvemos el token y el ID del usuario en el cuerpo de la respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"uuid":  user.ID,
	})

}
