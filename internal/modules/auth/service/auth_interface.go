package service

import "context"

type Auther interface {
	CreateUser(ctx context.Context, in CreateUserIn) CreateUserOut
	AuthorizeEmail(ctx context.Context, in AuthorizeEmailIn) AuthorizeOut
}

type CreateUserIn struct {
	Username  string
	Password  string
	Email     string
	FirstName string
	LastName  string
	Phone     string
}

type CreateUserOut struct {
	Status    int
	ErrorCode int
}

type AuthorizeEmailIn struct {
	Email          string
	Password       string
	RetypePassword string
}

type AuthorizeOut struct {
	UserID       int
	AccessToken  string
	RefreshToken string
	ErrorCode    int
}