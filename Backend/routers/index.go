// Package routers provides all the details of our chi router.
package routers

import (
	"net/http"
	"organize-this/controllers"
	"organize-this/helpers"
	"organize-this/infra/cache"
	"organize-this/infra/cognito"
	"organize-this/infra/database"
	"organize-this/repository"
	"organize-this/routers/middlewares"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(r *chi.Mux) {
	handler := controllers.Handler{
		Repository: &repository.Repository{
			Database: database.GetDB(),
			Cache:    cache.GetClient(),
		},
		CognitoClient: cognito.GetClient(),
		TokenHelper:   &helpers.DefaultTokenHelper{},
	}

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		helpers.SuccessResponse(w, "alive ok")
	})

	// v1 api routes
	r.Route("/v1", func(r chi.Router) {

		// Users
		r.Post("/user", handler.SignUp)
		r.Put("/user", handler.ConfirmSignUp)
		r.Post("/token", handler.SignIn)
		r.Put("/token", handler.Refresh)
		r.Delete("/token", handler.LogOut)

		// Protected endpoints
		r.Group(func(r chi.Router) {
			r.Use(middlewares.JWTAuth(handler))

			// Entities
			r.Post("/entity", handler.CreateEntity)
			r.Get("/entity/{category}/{id}", handler.GetEntity)
			r.Get("/entities", handler.GetEntities)
		})

	})
}
