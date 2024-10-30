package services

import (
	"errors"
	"fidelity-client-app/models"
	"fmt"

	"gorm.io/gorm"
)

type PointsService struct {
	DB *gorm.DB
}

// CalculateLevel determina el nivel en funcion de los puntos totales del usuario
func CalculateLevel(points int) int {

	// Calcular el nivel diviendo los puntos entre 100 y sumando 1
	level := (points / 100) + 1
	return level
}

// AccumulatePoints incrementa los puntos y actualiza el nivel del usuario
func (s *PointsService) AccumulatePoints(userID string, purchaseAmount float64) (string, error) {

	// Calcular los puntos acumulados por la compra (20% de la compra)
	pointsEarned := int(purchaseAmount * 1)

	// Obtener el usuario desde la base de datos
	var user models.User
	if err := s.DB.First(&user, "id = ?", userID).Error; err != nil {
		return "", errors.New("usuario no encontrado")
	}

	// Incrementar los puntos del usuario y determinar si hay un cambio de nivel
	user.Points += pointsEarned
	newLevel := CalculateLevel(user.Points)
	levelUp := false

	if user.Level != newLevel {
		user.Level = newLevel
		levelUp = true
	}

	// Guardar los cambios del usuario en la base de datos
	if err := s.DB.Save(&user).Error; err != nil {
		return "", errors.New("error al actualizar el nivel del usuario")
	}

	// Generar mensaje de confirmacion
	message := fmt.Sprintf("Puntos acumulados: %d puntos añadidos.", pointsEarned)
	if levelUp {
		message += fmt.Sprintf(" Felicidades! Has alcanzado el nivel %d", newLevel)
	}

	return message, nil
}
