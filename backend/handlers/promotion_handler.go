package handlers

import (
	"encoding/json"
	"fidelity-client-app/models"
	"fidelity-client-app/services"
	"net/http"
	"strconv"
	"time"
)

type PromotionHandler struct {
	PromotionService *services.PromotionService
}

// CreatePromotion maneja la solicitud para crear una nueva promocion
func (h *PromotionHandler) CreatePromotion(w http.ResponseWriter, r *http.Request) {

	// Establecemos metodo permitido
	if r.Method != http.MethodPost {
		http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	var promotion models.Promotion

	// Decodificamoso el JSON recibido en el cuerpo de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&promotion); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Llamamos al servicio para crear la promocion
	if err := h.PromotionService.CreatePromotion(&promotion); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Enviamos una respuesta 201 para indicar creacion exitosa
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(promotion)
}

// GetActivePromotions maneja la solicitud para obtener las promociones activas
func (h *PromotionHandler) GetActivePromotions(w http.ResponseWriter, r *http.Request) {

	// Establecemos metodo permitido
	if r.Method != http.MethodGet {
		http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Procesamos los parametros "page" y "pageSize" de la URL
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1 // valor predeterminado
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // valor predeterminado
	}

	// Usamos la fecha actual para filtrar las promociones activas
	currentDate := time.Now()

	// Obtenemos las promociones activas con el servicio
	promotions, total, err := h.PromotionService.GetActivePromotions(currentDate, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Estructuramos la respuesta incluyendo datos de paginacion
	response := map[string]interface{}{
		"promotions": promotions,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
	}

	// Devolvemos la respuesta como json
	json.NewEncoder(w).Encode(response)
}

// GetPromotionByID maneja la solicitud para obtener una promocion especifica
func (h *PromotionHandler) GetPromotionByID(w http.ResponseWriter, r *http.Request) {

	// Establecemos metodo permitido
	if r.Method != http.MethodGet {
		http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtenemos el ID de la promocion
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "el id de la promocion es obligatorio", http.StatusBadRequest)
		return
	}

	// Llamamos al servicio de GetPromotionByID
	promotion, err := h.PromotionService.GetPromotionByID(id)
	if err != nil {
		if err.Error() == "promotion not found" {
			http.Error(w, "Promoción no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error interno en el servidor", http.StatusInternalServerError)
		}
		return
	}

	// Enviamos los detalles de la promocion al servidor
	json.NewEncoder(w).Encode(promotion)
}

// UpdatePromotion maneja la solicitud para actualizar los datos de una promocion
func (h *PromotionHandler) UpdatePromotion(w http.ResponseWriter, r *http.Request) {

	// Establecemos metodo permitido
	if r.Method != http.MethodPut {
		http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtenemos el ID de los parametros de la URL
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "el id de la promocion es obligatorio", http.StatusBadRequest)
		return
	}

	// Decodificamos el JSON recibido en el cuerpo de la solicitud
	var updatedPromotion models.Promotion
	if err := json.NewDecoder(r.Body).Decode(&updatedPromotion); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Llamamos al servicio para actualizar la promocion
	if err := h.PromotionService.UpdatePromotion(id, &updatedPromotion); err != nil {
		if err.Error() == "promotion not found" {
			http.Error(w, "Promoción no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Enviamos una respuesta con el codigo 200
	w.WriteHeader(http.StatusOK)
}

// DeletePromotion maneja la solicitud para eliminar una promocion especifica
func (h *PromotionHandler) DeletePromotion(w http.ResponseWriter, r *http.Request) {

	// Establecemos metodo permitido
	if r.Method != http.MethodDelete {
		http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtenemos el ID de la promoción desde los parámetros de la URL
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "El ID de la promoción es obligatorio", http.StatusBadRequest)
		return
	}

	// Llamamos al servicio para eliminar la promoción
	if err := h.PromotionService.DeletePromotion(id); err != nil {
		if err.Error() == "promotion not found" {
			http.Error(w, "Promoción no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error interno en el servidor", http.StatusInternalServerError)
		}
		return
	}

	// Enviamos una respuesta con código 204 para indicar que la promoción fue eliminada
	w.WriteHeader(http.StatusNoContent)
}

// GetActivePromotionsForUser maneja la obtención de promociones activas no consumidas por el usuario
func (h *PromotionHandler) GetActivePromotionsForUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id es obligatorio", http.StatusBadRequest)
		return
	}

	promotions, err := h.PromotionService.GetActivePromotionsForUser(userID)
	if err != nil {
		http.Error(w, "Error al obtener promociones", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(promotions)
}

// ConsumePromotion maneja la solicitud de consumo de una promoción
func (h *PromotionHandler) ConsumePromotion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	promotionID := r.URL.Query().Get("promotion_id")
	if userID == "" || promotionID == "" {
		http.Error(w, "user_id y promotion_id son obligatorios", http.StatusBadRequest)
		return
	}

	err := h.PromotionService.ConsumePromotion(userID, promotionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Promoción consumida exitosamente",
	})
}

// CheckPromotionAvailability maneja la solicitud para verificar si una promoción ya fue consumida
func (h *PromotionHandler) CheckPromotionAvailability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	promotionID := r.URL.Query().Get("promotion_id")
	if userID == "" || promotionID == "" {
		http.Error(w, "user_id y promotion_id son obligatorios", http.StatusBadRequest)
		return
	}

	isConsumed, err := h.PromotionService.IsPromotionConsumed(userID, promotionID)
	if err != nil {
		http.Error(w, "Error al verificar la promoción", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{
		"is_consumed": isConsumed,
	})
}
