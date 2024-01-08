package main

import (
	"github.com/Adedunmol/mycart/internal/app"
)

func main() {
	app.Run()
	// r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello World!"))
	// })
	// http.ListenAndServe(":3000", r)
}
