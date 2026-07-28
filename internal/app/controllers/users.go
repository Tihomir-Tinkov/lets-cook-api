package controllers

import (
	"errors"
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/responses"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/services"
)

var (
	ErrInvalidUserID = errors.New("invalid user id")
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

type registerRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := decodeJSON(r, &req); err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	user, err := c.userService.Register(
		r.Context(),
		req.DisplayName,
		req.Email,
		req.Password,
	)

	if err != nil {
		// Example:
		// switch {
		// case errors.Is(err, services.ErrEmailTaken):
		//     status = http.StatusConflict
		// default:
		//     status = http.StatusBadRequest
		// }

		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	responses.JSONResponse(w, http.StatusCreated, user)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := decodeJSON(r, &req); err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	token, err := c.userService.Login(
		r.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		responses.JSONError(w, r, err, http.StatusUnauthorized)
		return
	}

	responses.JSONResponse(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func (c *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	// Stateless JWT logout:
	// client removes the token.
	//
	// When refresh tokens are introduced,
	// this endpoint can revoke them.

	responses.JSONResponse(
		w,
		http.StatusOK,
		map[string]string{
			"message": "logged out",
		},
	)
}

func (c *UserController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")

	if err != nil {
		responses.JSONError(
			w,
			r,
			err,
			http.StatusBadRequest,
		)
		return
	}

	user, err := c.userService.GetByID(
		r.Context(),
		id,
	)

	if err != nil {
		responses.JSONError(
			w,
			r,
			err,
			http.StatusNotFound,
		)
		return
	}

	responses.JSONResponse(
		w,
		http.StatusOK,
		user,
	)
}

func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	auth, err := requireAuth(r)
	if err != nil {
		responses.JSONError(w, r, err, http.StatusUnauthorized)
		return
	}

	if id != auth.UserID {
		responses.JSONError(
			w,
			r,
			ErrForbidden,
			http.StatusForbidden,
		)
		return
	}

	var req updateUserRequest

	if err := decodeJSON(r, &req); err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:          auth.UserID,
		DisplayName: req.DisplayName,
		Email:       req.Email,
	}

	if err := c.userService.Update(
		r.Context(),
		auth.UserID,
		&user,
		req.Password,
	); err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	auth, err := requireAuth(r)
	if err != nil {
		responses.JSONError(w, r, err, http.StatusUnauthorized)
		return
	}

	if id != auth.UserID {
		responses.JSONError(
			w,
			r,
			ErrForbidden,
			http.StatusForbidden,
		)
		return
	}

	if err := c.userService.Delete(
		r.Context(),
		auth.UserID,
		id,
	); err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
