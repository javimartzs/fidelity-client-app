package handlers

import (
	"encoding/json"
	"fidelity-client-app/services"
	"net/http"
	"strconv"
)

type PointsHandler struct {
	PointsService *services.PointsService
}

// AccumulatePoints maneja la solicitud para acumular puntos de un cliente
func (h *PointsHandler) AccumulatePoints(w http.ResponseWriter, r *http.Request) {

	// Configuramos el metodo permitido
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusBadRequest)
		return
	}

	// Obtener el UserID y purchaseAmount de los parametros de consulta
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "el user id es obligatorio", http.StatusBadRequest)
		return
	}

	purchaseAmountStr := r.URL.Query().Get("purchase_amount")
	if purchaseAmountStr == "" {
		http.Error(w, "El monto de la compra es obligatorio", http.StatusBadRequest)
		return
	}

	purchaseAmount, err := strconv.ParseFloat(purchaseAmountStr, 64)
	if err != nil {
		http.Error(w, "Monto de la compra no valido", http.StatusBadRequest)
		return
	}

	// Llamar al servicio para acumular puntos y mensaje de respuesta
	message, err := h.PointsService.AccumulatePoints(userID, purchaseAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver el mensaje en la respuesta JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}
