package models

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	FirstName  string `jsnon:"firstname"`
	LastName   string `json:"lastname"`
	UserStatus int    `json:"status"`
}
