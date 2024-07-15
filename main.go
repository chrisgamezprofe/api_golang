package main

import (
	"log"
	"net/http"

	"github.com/chrisgamezprofe/api_golang/data"
	"github.com/chrisgamezprofe/api_golang/models"
	"github.com/chrisgamezprofe/api_golang/routes"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando .env")
	}

	data.ConectarPostgres()

	data.DB.AutoMigrate(&models.Rol{})
	data.DB.AutoMigrate(&models.Usuario{})

	rutas := routes.InitRouter()
	log.Fatal(http.ListenAndServe(":8080", rutas))
}

