package models

import "gorm.io/gorm"

type Distribusi struct {
	gorm.Model
	UserID        uint `gorm:"uniqueIndex;not null" json:"user_id"`
	User          User `gorm:"foreignKey:UserID" json:"user"`
	JumlahKantong int  `gorm:"default:0;not null" json:"jumlah_kantong"`
}