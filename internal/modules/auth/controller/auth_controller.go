package controller

import (
	"net/http"
	"pet-store/internal/infrastructure/component"
	"pet-store/internal/infrastructure/errors"
	"pet-store/internal/infrastructure/responder"
	"pet-store/internal/modules/auth/service"

	"github.com/go-playground/validator"
	"github.com/ptflp/godecoder"
)

type Auther interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type Auth struct {
	auth service.Auther
	responder.Responder
	godecoder.Decoder
}

func NewAuth(service service.Auther, components *component.Components) Auther {
	return &Auth{auth: service, Responder: components.Responder, Decoder: components.Decoder}
}

// @Summary CreateUser
// @Tags user
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body RegisterRequest true "account info"
// @Success 200 {integer} integer 1
// @Router /user [post]
func (a *Auth) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := a.Decode(r.Body, &req)
	if err != nil {
		a.ErrorBadRequest(w, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		a.ErrorBadRequest(w, err)
		return
	}

	out := a.auth.CreateUser(r.Context(), service.CreateUserIn{
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	})

	if out.ErrorCode != errors.NoError {
		msg := "register error"
		if out.ErrorCode == errors.UserServiceUserAlreadyExists {
			msg = "User already exists, please check your username"
		}
		a.OutputJSON(w, RegisterResponse{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data: Data{
				Message: msg,
			},
		})
		return
	}

	a.OutputJSON(w, RegisterResponse{
		Success: true,
		Data: Data{
			Message: "you have registered a user named: " + req.Username,
		},
	})
}

// @Summary Login
// @Tags user
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body LoginRequest true "credentials"
// @Success 200 {string} string "token"
// @Router /user/login [get]
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := a.Decode(r.Body, &req)
	if err != nil {
		a.ErrorBadRequest(w, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		a.ErrorBadRequest(w, err)
		return
	}

	out := a.auth.AuthorizeEmail(r.Context(), service.AuthorizeEmailIn{
		Email:    req.Email,
		Password: req.Password,
	})

	if out.ErrorCode == errors.AuthServiceUserNotVerified {
		a.OutputJSON(w, AuthResponse{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data: LoginData{
				Message: "user email is not verified",
			},
		})
		return
	}

	if out.ErrorCode != errors.NoError {
		a.OutputJSON(w, AuthResponse{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data: LoginData{
				Message: "login or password mismatch",
			},
		})
		return
	}

	a.OutputJSON(w, AuthResponse{
		Success: true,
		Data: LoginData{
			Message:      "success login",
			AccessToken:  out.AccessToken,
			RefreshToken: out.RefreshToken,
		},
	})
}
