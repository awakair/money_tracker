package dto

type LoginDTO struct {
	Email    string
	Password string
}

type AuthResponseDTO struct {
	Token string
	User  *UserDTO
}
