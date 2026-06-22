package dto

type UpdateProfileReq struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type UpdateLocationReq struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type DeleteAccountReq struct {
	CurrentPassword string `json:"current_password" validate:"required"`
}
