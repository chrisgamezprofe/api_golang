package controllers

import (
	"encoding/json"
	"net/http"
)

type Saludo struct{
	Msg string `json:"message"`
	StatusCode int `json:"code"`
}


func InitRoute(w http.ResponseWriter, r *http.Request)  {
	saludo := Saludo{
		Msg: "Api Funcionando",
		StatusCode: 200,
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(saludo)
}