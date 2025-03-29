package handler

import (
	"encoding/json"
	"net/http"

	"money_tracker/internal/delivery/http/dto"
	httpinternal "money_tracker/internal/delivery/http/internal"
)

type AuthHandler struct {
	authUseCase AuthUseCase
}

func NewAuthHandler(authUseCase AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse "Successfully authenticated"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Invalid credentials"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpinternal.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	authResp, err := h.authUseCase.Login(r.Context(), req.ToUsecaseDTO())
	if err != nil {
		processError(w, err)
		return
	}

	resp := dto.NewLoginResponseFromUsecaseDTO(authResp)
	httpinternal.RespondWithJSON(w, http.StatusOK, resp)
}
