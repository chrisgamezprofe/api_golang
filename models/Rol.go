package models

import "gorm.io/gorm"

type Rol struct{
	gorm.Model
	Nombre string `json:"nombre" gorm:"unique;not null"`
	Activo bool `json:"activo" gorm:"default:true"`
	Usuarios []Usuario `json:"usuarios"`
}

func (Rol) TableName() string {
	return "roles"
}