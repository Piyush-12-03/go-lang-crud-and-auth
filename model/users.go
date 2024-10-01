package model

type Users struct {
	Id       int    `gorm:"primary_key;autoIncrement"`
	Name     string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Role     string `gorm:"type:varchar(255);not null"`
}
