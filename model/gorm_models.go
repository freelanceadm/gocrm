package model

import "gorm.io/gorm"

// gorm.Model definition
// type Model struct {
// 	ID        uint           `gorm:"primaryKey"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
//   }

type User struct {
	gorm.Model
	Email        string
	PasswordHash string
}
