package models

type Promotion struct {
	ID            string `gorm:"primaryKey" json:"id"`
	Title         string `gorm:"size:50;not null" json:"title"`
	Description   string `gorm:"size:250" json:"description"`
	LevelRequired int    `gorm:"not null" json:"level_required"`
	StartDate     string `json:"start_date"`         // Formato esperado: YYYY-MM-DD
	EndDate       string `json:"end_date,omitempty"` // Formato esperado: YYYY-MM-DD
}
