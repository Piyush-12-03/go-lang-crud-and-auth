package model

type Tags struct {
	Id     int     `gorm:"primary_key;autoIncrement"`
	Name   string  `gorm:"type:varchar(255);not null"`
	Neches []Neche `gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE;"`
}
