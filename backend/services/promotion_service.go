package services

import (
	"errors"
	"fidelity-client-app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PromotionService struct {
	DB *gorm.DB
}

const dateFormat = "2006-01-02"

// CreatePromotion (logica de negocio para crear una nueva promocion)
func (s *PromotionService) CreatePromotion(promotion *models.Promotion) error {

	// Validar que el titulo esta completo
	if promotion.Title == "" {
		return errors.New("el titulo es obligatorio")
	}

	// Validamos el nivel requerido
	if promotion.LevelRequired < 1 {
		return errors.New("el nivel de la promocion tiene que ser mayor que cero")
	}

	// Validar que StartDate siga el formato `YYYY-MM-DD`
	_, err := time.Parse(dateFormat, promotion.StartDate)
	if err != nil {
		return errors.New("el formato de start_date debe ser YYYY-MM-DD")
	}

	// Validar que EndDate siga el formato `YYYY-MM-DD` y no sea anterior a StartDate
	if promotion.EndDate != "" {
		endDate, err := time.Parse(dateFormat, promotion.EndDate)
		if err != nil {
			return errors.New("el formato de end_date debe ser YYYY-MM-DD")
		}

		startDate, _ := time.Parse(dateFormat, promotion.StartDate)
		if endDate.Before(startDate) {
			return errors.New("la fecha de fin no puede ser anterior a la fecha de inicio")
		}
	}
	// Generamos el ID de la promocion
	promotion.ID = uuid.NewString()

	// Guardar la promocion en la base de datos
	if err := s.DB.Save(promotion).Error; err != nil {
		return errors.New("error al guardar la promocion")
	}

	return nil
}

// GetActivePromotions (maneja la logica de negocio a la hora de ver las promociones)
func (s *PromotionService) GetActivePromotions(currentDate time.Time, page int, pageSize int) ([]models.Promotion, int64, error) {

	var promotions []models.Promotion
	var total int64

	// Convertir la fecha actual a `YYYY-MM-DD` para hacer comparaciones de strings
	currentDateString := currentDate.Format("2006-01-02")

	// Contar el número total de promociones activas
	if err := s.DB.Model(&models.Promotion{}).
		Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", currentDateString, currentDateString).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calcular el offset para la paginación
	offset := (page - 1) * pageSize

	// Consultar las promociones activas con paginación
	err := s.DB.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", currentDateString, currentDateString).
		Offset(offset).
		Limit(pageSize).
		Find(&promotions).Error

	if err != nil {
		return nil, 0, err
	}

	// Devolvemos el listado de promociones y el total de resultados
	return promotions, total, nil
}

// GetPromotionByID (maneja la logica de negocio para obtener unua promocion especifica)
func (s *PromotionService) GetPromotionByID(promotionID string) (*models.Promotion, error) {

	var promotion models.Promotion
	// Buscamos la promocion en la base de datos y la almacenamos en promotion
	if err := s.DB.First(&promotion, "id = ?", promotionID).Error; err != nil {
		return nil, errors.New("promotion not found")
	}

	// Devolvemos los datos de la promocion listada por ID
	return &promotion, nil

}

// UpdatePromotion (maneja la logica de negocio para actualizar los datos de una promocion existente)
func (s *PromotionService) UpdatePromotion(id string, updatedPromotion *models.Promotion) error {

	// Validar que el nivel introducido no sea menor que 1
	// Validar que el titulo esta completo
	if updatedPromotion.Title == "" {
		return errors.New("el titulo es obligatorio")
	}

	// Validamos el nivel requerido
	if updatedPromotion.LevelRequired < 1 {
		return errors.New("el nivel de la promocion tiene que ser mayor que cero")
	}

	// Validar formato de StartDate
	_, err := time.Parse(dateFormat, updatedPromotion.StartDate)
	if err != nil {
		return errors.New("el formato de start_date debe ser YYYY-MM-DD")
	}

	// Validar EndDate (si existe) y comprobar que no sea anterior a StartDate
	if updatedPromotion.EndDate != "" {
		endDate, err := time.Parse(dateFormat, updatedPromotion.EndDate)
		if err != nil {
			return errors.New("el formato de end_date debe ser YYYY-MM-DD")
		}

		startDate, _ := time.Parse(dateFormat, updatedPromotion.StartDate)
		if endDate.Before(startDate) {
			return errors.New("la fecha de fin no puede ser anterior a la fecha de inicio")
		}
	}
	// Buscar y actualizar la promocino
	var promotion models.Promotion
	if err := s.DB.First(&promotion, "id = ?", id).Error; err != nil {
		return errors.New("promotion not found")
	}

	// Modificar datos en la base de datos
	if err := s.DB.Model(&promotion).Updates(updatedPromotion).Error; err != nil {
		return errors.New("error updating promotion")
	}

	return nil
}

// DeletePromotion (maneja la logica de negocio para eliminar promociones existentes)
func (s *PromotionService) DeletePromotion(id string) error {

	if err := s.DB.Delete(&models.Promotion{}, "id = ?", id).Error; err != nil {
		return errors.New("promotion not found")
	}
	return nil
}

// GetActivePromotionsForUser obtiene las promociones activas que el usuario no ha consumido
func (s *PromotionService) GetActivePromotionsForUser(userID string) ([]models.Promotion, error) {
	var promotions []models.Promotion
	currentDate := time.Now().Format(dateFormat)

	var consumedPromotions []string
	s.DB.Model(&models.PromotionUsage{}).Where("user_id = ?", userID).Pluck("promotion_id", &consumedPromotions)

	if err := s.DB.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", currentDate, currentDate).
		Where("id NOT IN ?", consumedPromotions).
		Find(&promotions).Error; err != nil {
		return nil, err
	}
	return promotions, nil
}

// ConsumePromotion permite a un usuario consumir una promoción si cumple con los requisitos
func (s *PromotionService) ConsumePromotion(userID, promotionID string) error {
	var promotion models.Promotion
	if err := s.DB.First(&promotion, "id = ?", promotionID).Error; err != nil {
		return errors.New("promoción no encontrada")
	}

	currentDate := time.Now().Format(dateFormat)
	if promotion.StartDate > currentDate || (promotion.EndDate != "" && promotion.EndDate < currentDate) {
		return errors.New("la promoción no está activa en este momento")
	}

	var user models.User
	if err := s.DB.First(&user, "id = ?", userID).Error; err != nil {
		return errors.New("usuario no encontrado")
	}

	if user.Level < promotion.LevelRequired {
		return errors.New("el usuario no tiene el nivel necesario para consumir esta promoción")
	}

	var usage models.PromotionUsage
	if err := s.DB.Where("user_id = ? AND promotion_id = ?", userID, promotionID).First(&usage).Error; err == nil {
		return errors.New("esta promoción ya ha sido consumida por el usuario")
	}

	usage = models.PromotionUsage{
		ID:          uuid.NewString(),
		UserID:      userID,
		PromotionID: promotionID,
		ConsumedAt:  time.Now(),
	}
	if err := s.DB.Create(&usage).Error; err != nil {
		return errors.New("error al registrar el consumo de la promoción")
	}
	return nil
}

// IsPromotionConsumed verifica si una promoción ha sido consumida por el usuario
func (s *PromotionService) IsPromotionConsumed(userID, promotionID string) (bool, error) {
	var usage models.PromotionUsage
	if err := s.DB.Where("user_id = ? AND promotion_id = ?", userID, promotionID).First(&usage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
