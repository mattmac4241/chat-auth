package service

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// NewServer configures and returns a server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()
	db := &dataHandler{}
	initRoutes(mx, formatter, db)
	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, database Database) {
	mx.HandleFunc("/auth/register", registerUserHandler(formatter, database)).Methods("POST")
	mx.HandleFunc("/auth/login", loginUserHandler(formatter, database)).Methods("POST")
	mx.HandleFunc("/auth/token/{key}", tokenValidatorHandler(formatter, database)).Methods("GET")
}
