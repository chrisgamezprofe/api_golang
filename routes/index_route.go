package routes

import (
	"github.com/chrisgamezprofe/api_golang/controllers"
	"github.com/chrisgamezprofe/api_golang/middleware"
	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router  {
	rutas := mux.NewRouter()
	api := rutas.PathPrefix("/api").Subrouter()

	api.HandleFunc("",controllers.InitRoute).Methods("GET")


	apiRoles := api.PathPrefix("/roles").Subrouter()
	apiRoles.HandleFunc("",middleware.SetMiddlewareAuthentication(controllers.GetRoles)).Methods("GET")
	apiRoles.HandleFunc("/{id}",middleware.SetMiddlewareAuthentication(controllers.GetRol)).Methods("GET")
	apiRoles.HandleFunc("/{id}",middleware.SetMiddlewareAuthentication(controllers.UpdateRol)).Methods("PUT")
	apiRoles.HandleFunc("",middleware.SetMiddlewareAuthentication(controllers.NewRol)).Methods("POST")
	apiRoles.HandleFunc("/{id}",middleware.SetMiddlewareAuthentication(controllers.DeleteRol)).Methods("DELETE")

	apiUsuarios := api.PathPrefix("/usuarios").Subrouter()
	apiUsuarios.HandleFunc("",middleware.SetMiddlewareAuthentication(controllers.GetUsuarios)).Methods("GET")
	apiUsuarios.HandleFunc("/{id}",middleware.SetMiddlewareAuthentication(controllers.GetUsuario)).Methods("GET")
	apiUsuarios.HandleFunc("/{id}",middleware.SetMiddlewareAuthentication(controllers.UpdateUsuario)).Methods("PUT")
	apiUsuarios.HandleFunc("",middleware.SetMiddlewareAuthentication(controllers.NewUsuario)).Methods("POST")
	apiUsuarios.HandleFunc("/{id}",middleware.SetMiddlewareAuthentication(controllers.DeleteUsuario)).Methods("DELETE")


	apiAuth := api.PathPrefix("/auth").Subrouter()

	apiAuth.HandleFunc("/login",controllers.Login).Methods("POST")
	apiAuth.HandleFunc("/register",controllers.NewUsuario).Methods("POST")

	return rutas
}