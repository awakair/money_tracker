package handler

import (
	"encoding/json"
	"net/http"

	"money_tracker/internal/delivery/http/dto"
	httpinternal "money_tracker/internal/delivery/http/internal"
)

type UserHandler struct {
	userUseCase UserUseCase
}

func NewUserHandler(userUseCase UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User registration data"
// @Success 201 {object} dto.UserResponse "User created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 409 {object} dto.ErrorResponse "User already exists"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpinternal.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userUseCase.CreateUser(r.Context(), req.ToUsecaseDTO())
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.NewUserResponseFromUsecaseDTO(user)
	httpinternal.RespondWithJSON(w, http.StatusCreated, resp)
}
