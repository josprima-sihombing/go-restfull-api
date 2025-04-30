package model

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserCredential struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type Profile struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type ProfileDetail struct {
	ID        string `json:"id"`
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type UpdateProfile struct {
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}
