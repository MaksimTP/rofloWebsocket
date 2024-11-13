package server

import (
	"main/server/endpoints"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/ws", endpoints.HandleWS).Methods("GET")
	r.HandleFunc("/change", endpoints.ChangeSomething).Methods("POST")

	r.PathPrefix("/").Handler(
		negroni.Classic(),
		// negroni.HandlerFunc(middlewares.SomeMiddleware),
	)

	return r
}
