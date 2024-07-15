package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/chrisgamezprofe/api_golang/data"
	"github.com/chrisgamezprofe/api_golang/models"
	"github.com/chrisgamezprofe/api_golang/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUsuarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var usuarios []models.UsuarioResponse
	data.DB.Preload("Rol").Find(&usuarios)
	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Lista de usuarios",
		StatusCode: http.StatusOK,
		Data:       usuarios,
	})
}

func GetUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var usuario models.UsuarioResponse
	data.DB.Preload("Rol").First(&usuario, params["id"])
	if usuario.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data.DB.Model(&usuario).Association("Roles").Find(&usuario.RolId)

	json.NewEncoder(w).Encode(&usuario)
}

func NewUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var usuario models.Usuario
	json.NewDecoder(r.Body).Decode(&usuario)
	createdUsuario := data.DB.Create(&usuario)
	err := createdUsuario.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	if err := data.DB.Preload("Rol").First(&usuario, usuario.ID).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(utils.Respuesta{
            Msg:        "Error al cargar Rol",
            StatusCode: http.StatusInternalServerError,
            Data:       err.Error(),
        })
        return
	}

	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Usuario registrado con éxito",
		StatusCode: http.StatusOK,
		Data:       usuario,
	})
}

func UpdateUsuario(w http.ResponseWriter, r *http.Request) {
    var usuario models.Usuario
	w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    id := params["id"]

    if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        
        json.NewEncoder(w).Encode(utils.Respuesta{
            Msg:        "Error al decodificar la solicitud",
            StatusCode: http.StatusBadRequest,
            Data:       err.Error(),
        })
        return
    }

    var usuarioExistente models.Usuario
    if err := data.DB.Preload("Rol").First(&usuarioExistente, id).Error; err != nil {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(utils.Respuesta{
            Msg:        "Usuario no encontrado",
            StatusCode: http.StatusNotFound,
            Data:       err.Error(),
        })
        return
    }

    usuarioExistente.Nombre = usuario.Nombre
    usuarioExistente.Correo = usuario.Correo
    usuarioExistente.RolId = usuario.RolId
    usuarioExistente.Password = usuario.Password
    if err := data.DB.Save(&usuarioExistente).Error; err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(utils.Respuesta{
            Msg:        "Error al actualizar",
            StatusCode: http.StatusInternalServerError,
            Data:       err.Error(),
        })
        return
    }

    
    respuesta := utils.Respuesta{
        Msg:        "Actualizado con éxito",
        StatusCode: http.StatusOK,
        Data:       usuarioExistente,
    }

    json.NewEncoder(w).Encode(&respuesta)
}

func DeleteUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var usuario models.Usuario
	data.DB.Preload("Rol").First(&usuario, params["id"])

	if usuario.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data.DB.Delete(&usuario)
	w.WriteHeader(http.StatusOK)
}

type Credenciales struct{
	Correo string `json:"correo"`
	Password string `json:"password"`
}

type Claims struct{
	Correo string `json:"correo"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var credenciales Credenciales

	if err := json.NewDecoder(r.Body).Decode(&credenciales); err !=nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := data.DB.Where("correo = ?",credenciales.Correo).First(&usuario).Error; err !=nil {
		if err == gorm.ErrRecordNotFound{
			http.Error(w,"Datos de acceso incorrectos",http.StatusUnauthorized)
		}else{
			http.Error(w,"Error al consultar usuario",http.StatusInternalServerError)
		}
		return
	}

	if err := VerificarPassword(usuario.Password,credenciales.Password); err != nil{
		http.Error(w,"Datos de acceso incorrectos pwrd",http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Correo: usuario.Correo,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(os.Getenv("API_SECRET"))) // Convert to byte slice
	
    if err != nil {
        http.Error(w, "Error al crear el token", http.StatusInternalServerError)
        return
    }

	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Autentización exitosa",
		StatusCode: http.StatusOK,
		Data:       tokenString,
	})
	
}

func VerificarPassword(passwordHashed string,password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHashed),[]byte(password))
}