package main


import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
type application struct{
	config config

}


func(app *application) mount() http.Handler{
r:= chi.NewRouter()
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)

r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte (".root"))
})

http.ListenAndServe(":8080", r)
return r



}
type config struct {
	address string 
	db dbConfig
}
type dbConfig struct {
	dsn string 
}