package controllers

import (
	//"encoding/json"
	//"errors"
	"net/http"

	//"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/responses"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/services"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// TODO
func (c *UserController) Register(w http.ResponseWriter, r *http.Request)

func (c *UserController) Login(w http.ResponseWriter, r *http.Request)

func (c *UserController) Logout(w http.ResponseWriter, r *http.Request)

func (c *UserController) GetByID(w http.ResponseWriter, r *http.Request)

func (c *UserController) Update(w http.ResponseWriter, r *http.Request)

func (c *UserController) Delete(w http.ResponseWriter, r *http.Request)
