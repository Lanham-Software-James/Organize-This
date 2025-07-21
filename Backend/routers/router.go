// Package routers provides all the details of our chi router.
package routers

import (
	"willowsuite-vault/routers/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SetupRoute configures our Chi Router.
func SetupRoute() *chi.Mux {

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.Cors())
	RegisterRoutes(router) //routes register

	return router
}
