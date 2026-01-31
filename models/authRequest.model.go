package models

type AuthRequest struct {
	Uid      string `json:"uid"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type AuthResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
	User  User   `json:"user"`
}
