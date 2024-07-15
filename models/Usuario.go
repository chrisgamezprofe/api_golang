package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Usuario struct{
	gorm.Model
	Nombre string `json:"nombre" gorm:"size:100;not null"`
	Correo string `json:"correo" gorm:"size:100;unique;not null"`
	Password string `json:"password" gorm:"default:true"`
	RolId   uint `json:"rol_id"`
	Rol    Rol `json:"rol"`
}

type UsuarioResponse struct{
	ID uint `json:"id"`
	Nombre string `json:"nombre"`
	Correo string `json:"correo"`
	RolId   uint `json:"rol_id"`
	Rol    Rol `json:"rol"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Usuario) TableName() string {
	return "usuarios"
}

func Hash(password string) ([]byte,error) {
	return bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
}


func (u *Usuario) BeforeSave(tx *gorm.DB) error{
	passwordHashed,err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(passwordHashed)
	return nil
}



