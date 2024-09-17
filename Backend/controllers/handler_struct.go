package controllers

import (
	"chi-boilerplate/repository"
)

type Handler struct {
	Repository *repository.Repository
}
