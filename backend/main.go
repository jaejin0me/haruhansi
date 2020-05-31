package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var renderer *render.Render

func init() {
	renderer = render.New()
}

func main() {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "Simple chat!"})
	})

	n := negroni.Classic()
	n.UseHandler(router)

	n.Run(":3000")
}
