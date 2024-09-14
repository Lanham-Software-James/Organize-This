package routers

import (
	"chi-boilerplate/controllers"
	"chi-boilerplate/helpers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(router *chi.Mux) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		helpers.SuccessResponse(w, "alive ok")
	})

	//Entity Management Route
	router.Group(func(r chi.Router) {
		r.Post("/entity-management/building", controllers.CreateBuilding)
	})
}
