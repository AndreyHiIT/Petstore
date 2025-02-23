package controller

type RegisterRequest struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Phone     string `json:"phone"`
}

type RegisterResponse struct {
	Success   bool `json:"success"`
	ErrorCode int  `json:"error_code,omitempty"`
	Data      Data `json:"data"`
}

type Data struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Success   bool      `json:"success"`
	ErrorCode int       `json:"error_code,omitempty"`
	Data      LoginData `json:"data"`
}

type LoginData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
}
