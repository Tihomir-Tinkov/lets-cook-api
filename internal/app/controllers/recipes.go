package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/controllers/responses"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/dto"
	"github.com/Tihomir-Tinkov/lets-cook-api/internal/app/services"
)

var ErrInvalidRecipeID = errors.New("invalid recipe id")

type RecipeController struct {
	recipeService *services.RecipeService
}

func NewRecipeController(recipeService *services.RecipeService) *RecipeController {
	return &RecipeController{
		recipeService: recipeService,
	}
}

func (c *RecipeController) Create(w http.ResponseWriter, r *http.Request) {
	auth, err := requireAuth(r)
	if err != nil {
		responses.JSONError(w, r, err, http.StatusUnauthorized)
		return
	}

	var req dto.RecipeCreateRequest

	if err := decodeJSON(r, &req); err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	recipe, err := c.recipeService.Create(r.Context(), auth.UserID, req)

	if err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	responses.JSONResponse(w, http.StatusCreated, recipe)
}

func (c *RecipeController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	recipe, err := c.recipeService.GetByID(
		r.Context(),
		id,
	)

	if err != nil {
		responses.JSONError(w, r, err, http.StatusNotFound)
		return
	}

	responses.JSONResponse(w, http.StatusOK, recipe)
}

func (c *RecipeController) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	recipes, err := c.recipeService.List(
		r.Context(),
		limit,
		offset,
	)

	if err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	responses.JSONResponse(w, http.StatusOK, recipes)
}

func (c *RecipeController) Update(w http.ResponseWriter, r *http.Request) {
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

	var req dto.RecipeUpdateRequest

	if err := decodeJSON(r, &req); err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	recipe, err := c.recipeService.Update(r.Context(), id, auth.UserID, req)

	if err != nil {
		if errors.Is(err, services.ErrRecipeForbidden) {
			responses.JSONError(w, r, err, http.StatusForbidden)
			return
		}

		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	responses.JSONResponse(w, http.StatusOK, recipe)
}

func (c *RecipeController) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := c.recipeService.Delete(
		r.Context(),
		auth.UserID,
		id,
	); err != nil {
		if errors.Is(err, services.ErrRecipeForbidden) {
			responses.JSONError(w, r, err, http.StatusForbidden)
			return
		}

		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
