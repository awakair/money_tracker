package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"money_tracker/internal/delivery/http/dto"
	httpinternal "money_tracker/internal/delivery/http/internal"

	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	categoryUseCase CategoryUseCase
}

func NewCategoryHandler(categoryUseCase CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
	}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category for the authenticated user
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param request body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} dto.CategoryResponse "Category created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/categories [post]
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	var req dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpinternal.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	category, err := h.categoryUseCase.CreateCategory(r.Context(), userID, req.ToUsecaseDTO())
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.NewCategoryResponseFromUsecaseDTO(category)
	httpinternal.RespondWithJSON(w, http.StatusCreated, resp)
}

// ListCategories godoc
// @Summary List user categories
// @Description Get a paginated list of categories for the authenticated user
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param limit query int false "Maximum number of categories to return"
// @Param offset query int false "Number of categories to skip"
// @Success 200 {object} dto.CategoryListResponse "List of categories"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/categories [get]
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	var limit, offset *int

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 0 {
			httpinternal.RespondWithError(w, http.StatusBadRequest, "invalid limit parameter")
			return
		}
		limit = &parsedLimit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			httpinternal.RespondWithError(w, http.StatusBadRequest, "invalid offset parameter")
			return
		}
		offset = &parsedOffset
	}

	categories, total, err := h.categoryUseCase.ListCategories(r.Context(), userID, limit, offset)
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.CategoryListResponse{
		Categories: make([]dto.CategoryResponse, len(categories)),
		Total:      total,
		Limit:      limit,
		Offset:     offset,
	}

	for i, cat := range categories {
		resp.Categories[i] = dto.NewCategoryResponseFromUsecaseDTO(cat)
	}

	httpinternal.RespondWithJSON(w, http.StatusOK, resp)
}

// GetCategory godoc
// @Summary Get a category by ID
// @Description Get detailed information about a specific category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param id path string true "Category ID"
// @Success 200 {object} dto.CategoryResponse "Category details"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Category not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/categories/{id} [get]
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	categoryID := chi.URLParam(r, "id")

	category, err := h.categoryUseCase.GetCategory(r.Context(), userID, categoryID)
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.NewCategoryResponseFromUsecaseDTO(category)
	httpinternal.RespondWithJSON(w, http.StatusOK, resp)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param id path string true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Updated category data"
// @Success 200 {object} dto.CategoryResponse "Category updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Forbidden"
// @Failure 404 {object} dto.ErrorResponse "Category not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/{userId}/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	categoryID := chi.URLParam(r, "id")

	var req dto.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpinternal.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	category, err := h.categoryUseCase.UpdateCategory(r.Context(), userID, categoryID, req.ToUsecaseDTO())
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.NewCategoryResponseFromUsecaseDTO(category)
	httpinternal.RespondWithJSON(w, http.StatusOK, resp)
}
