package controller

type GetUserResponseSuccess struct {
	Success   bool     `json:"success"`
	ErrorCode int      `json:"error_code,omitempty"`
	DataUser  DataUser `json:"datauser"`
}

type DataUser struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	UserStatus int    `json:"userstatus"`
}

type UserResponseErr struct {
	Success   bool   `json:"success"`
	ErrorCode int    `json:"error_code,omitempty"`
	Data      string `json:"data"`
}

type UserUpdateResponse struct{
	Success bool `json:"success"`
	Message string `json:"message"`
}