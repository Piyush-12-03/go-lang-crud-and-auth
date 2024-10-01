package model

type Neche struct {
	Id        int    `gorm:"primary_key;autoIncrement"`
	NecheType string `gorm:"type:varchar(255);not null"`
	TagID     int    `gorm:"not null"`
	Tag       Tags   `gorm:"foreignKey:TagID"`
}