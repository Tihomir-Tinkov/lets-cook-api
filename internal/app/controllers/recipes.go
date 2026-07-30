package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/responses"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/services"
	"github.com/google/uuid"
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

type recipeImageRequest struct {
	ImageID      uuid.UUID `json:"image_id"`
	DisplayOrder int       `json:"display_order"`
}

type createRecipeRequest struct {
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Ingredients  string               `json:"ingredients"`
	Instructions string               `json:"instructions"`
	PrepTimeMin  int                  `json:"prep_time_min"`
	Servings     int                  `json:"servings"`
	Difficulty   int                  `json:"difficulty"`
	Images       []recipeImageRequest `json:"images"`
}

type updateRecipeRequest struct {
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Ingredients  string               `json:"ingredients"`
	Instructions string               `json:"instructions"`
	PrepTimeMin  int                  `json:"prep_time_min"`
	Servings     int                  `json:"servings"`
	Difficulty   int                  `json:"difficulty"`
	Images       []recipeImageRequest `json:"images"`
}

func (c *RecipeController) Create(w http.ResponseWriter, r *http.Request) {
	auth, err := requireAuth(r)
	if err != nil {
		responses.JSONError(w, r, err, http.StatusUnauthorized)
		return
	}

	var req createRecipeRequest

	if err := decodeJSON(r, &req); err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	recipe := models.Recipe{
		AuthorID:     auth.UserID,
		Title:        req.Title,
		Description:  req.Description,
		Ingredients:  req.Ingredients,
		Instructions: req.Instructions,
		PrepTimeMin:  req.PrepTimeMin,
		Servings:     req.Servings,
		Difficulty:   req.Difficulty,
	}

	for _, img := range req.Images {
		recipe.Images = append(recipe.Images, models.RecipeImage{
			ImageID:      img.ImageID,
			DisplayOrder: img.DisplayOrder,
		})
	}

	if err := c.recipeService.Create(
		r.Context(),
		&recipe,
	); err != nil {
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

	var req updateRecipeRequest

	if err := decodeJSON(r, &req); err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	recipe := models.Recipe{
		ID:           id,
		Title:        req.Title,
		Description:  req.Description,
		Ingredients:  req.Ingredients,
		Instructions: req.Instructions,
		PrepTimeMin:  req.PrepTimeMin,
		Servings:     req.Servings,
		Difficulty:   req.Difficulty,
	}

	for _, img := range req.Images {
		recipe.Images = append(recipe.Images, models.RecipeImage{
			ImageID:      img.ImageID,
			DisplayOrder: img.DisplayOrder,
		})
	}

	if err := c.recipeService.Update(
		r.Context(),
		auth.UserID,
		&recipe,
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
