package controllers

import (
	"errors"
	"net/http"

	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/controllers/responses"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/models"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/repositories"
	"github.com/Tihomir-Tinkov/cooking-site-project/internal/app/services"
)

var ErrInvalidCommentID = errors.New("invalid comment id")

type CommentController struct {
	commentService *services.CommentService
}

func NewCommentController(commentService *services.CommentService) *CommentController {
	return &CommentController{
		commentService: commentService,
	}
}

type createCommentRequest struct {
	Body   string `json:"body"`
	Rating int    `json:"rating"`
}

type updateCommentRequest struct {
	Body   string `json:"body"`
	Rating int    `json:"rating"`
}

func (c *CommentController) Create(w http.ResponseWriter, r *http.Request) {
	auth, err := requireAuth(r)
	if err != nil {
		responses.JSONError(w, r, err, http.StatusUnauthorized)
		return
	}

	var req createCommentRequest

	if err := decodeJSON(r, &req); err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	recipeID, err := parseUUIDParam(r, "recipeId")
	if err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	comment := models.Comment{
		RecipeID: recipeID,
		AuthorID: auth.UserID,
		Body:     req.Body,
		Rating:   req.Rating,
	}

	if err := c.commentService.Create(
		r.Context(),
		&comment,
	); err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	responses.JSONResponse(w, http.StatusCreated, comment)
}

func (c *CommentController) GetByID(w http.ResponseWriter, r *http.Request) {
	recipeID, err := parseUUIDParam(r, "recipeId")
	if err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	commentID, err := parseUUIDParam(r, "id")
	if err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	comment, err := c.commentService.GetByID(r.Context(), commentID)
	if err != nil {
		responses.JSONError(w, r, err, http.StatusNotFound)
		return
	}

	if comment.RecipeID != recipeID {
		responses.JSONError(
			w,
			r,
			repositories.ErrCommentNotFound,
			http.StatusNotFound,
		)
		return
	}

	responses.JSONResponse(w, http.StatusOK, comment)
}

func (c *CommentController) GetByRecipeID(w http.ResponseWriter, r *http.Request) {
	recipeID, err := parseUUIDParam(r, "recipeId")
	if err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	comments, err := c.commentService.GetByRecipeID(r.Context(), recipeID)
	if err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	responses.JSONResponse(w, http.StatusOK, comments)
}

func (c *CommentController) Update(w http.ResponseWriter, r *http.Request) {
	recipeID, err := parseUUIDParam(r, "recipeId")
	if err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

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

	var req updateCommentRequest

	if err := decodeJSON(r, &req); err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	comment := models.Comment{
		ID:     id,
		Body:   req.Body,
		Rating: req.Rating,
	}

	if err := c.commentService.Update(
		r.Context(),
		recipeID,
		auth.UserID,
		&comment,
	); err != nil {
		if errors.Is(err, services.ErrCommentForbidden) {
			responses.JSONError(w, r, err, http.StatusForbidden)
			return
		}

		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	recipeID, err := parseUUIDParam(r, "recipeId")
	if err != nil {
		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

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

	if err := c.commentService.Delete(
		r.Context(),
		recipeID,
		auth.UserID,
		id,
	); err != nil {
		if errors.Is(err, services.ErrCommentForbidden) {
			responses.JSONError(w, r, err, http.StatusForbidden)
			return
		}

		responses.JSONError(w, r, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
