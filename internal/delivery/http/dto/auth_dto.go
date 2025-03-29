package dto

import usecaseDTO "money_tracker/internal/application/usecases/dto"

// LoginRequest represents the login credentials
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

func (r *LoginRequest) ToUsecaseDTO() usecaseDTO.LoginDTO {
	return usecaseDTO.LoginDTO{
		Email:    r.Email,
		Password: r.Password,
	}
}

// LoginResponse represents the authentication response
type LoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

func NewLoginResponseFromUsecaseDTO(authResp *usecaseDTO.AuthResponseDTO) LoginResponse {
	return LoginResponse{
		Token: authResp.Token,
		User:  NewUserResponseFromUsecaseDTO(authResp.User),
	}
}
