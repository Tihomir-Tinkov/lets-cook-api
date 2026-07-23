package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/responses"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/ports"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/repositories"
)

type UserController struct {
	userService ports.UserService
}

func NewUserController(
	userService ports.UserService,
) *UserController {
	return &UserController{
		userService: userService,
	}
}

// @Summary Get user
// @Description Get user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User UUID"
// @Success 200 {object} models.User
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /users/{id} [get]
func (c *UserController) GetUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		responses.JSONError(
			w,
			r,
			errors.New("invalid user id"),
			http.StatusBadRequest,
		)
		return
	}

	user, err := c.userService.GetUser(
		r.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			responses.JSONError(
				w,
				r,
				err,
				http.StatusNotFound,
			)
			return
		}

		responses.JSONError(
			w,
			r,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	responses.JSONResponse(
		w,
		http.StatusOK,
		models.NewUserResponse(user),
	)
}
