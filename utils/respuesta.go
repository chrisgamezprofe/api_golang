package utils

type Respuesta struct{
	Msg string `json:"message"`
	Data interface{} `json:"data"`
	StatusCode int `json:"code"`
}