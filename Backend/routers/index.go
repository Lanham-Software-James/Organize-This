package routers

import (
	"chi-boilerplate/controllers"
	"chi-boilerplate/helpers"
	"chi-boilerplate/infra/database"
	"chi-boilerplate/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(r *chi.Mux) {
	handler := controllers.Handler{Repository: &repository.Repository{Database: database.DB}}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		helpers.SuccessResponse(w, "alive ok")
	})

	// v1 api routes
	r.Route("/v1", func(r chi.Router) {

		// enitity management routes
		r.Route("/entity-management", func(r chi.Router) {
			// Buildings
			r.Post("/building", handler.CreateBuilding)

			// Rooms
			r.Post("/room", handler.CreateRoom)
		})
	})
}
