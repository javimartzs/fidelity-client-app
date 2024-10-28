package models

type User struct {
	ID        string `gorm:"size:36;unique;not null;primaryKey"`
	FirstName string `gorm:"size:100;not null"`
	LastName  string `gorm:"size:100;not null"`
	BirthDate string `gorm:"size:25;not null"`
	Gender    string `gorm:"size:25;not null"`
	Email     string `gorm:"size:100;not null;unique"`
	Password  string `gorm:"size:100;not null"`
	Role      string `gorm:"default:customer-client"`
	Points    int    `gorm:"default:1"`
	Level     int    `gorm:"default:1"`
}
