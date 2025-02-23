package controller

import (
	"fmt"
	"net/http"
	"pet-store/internal/infrastructure/component"
	"pet-store/internal/infrastructure/errors"
	"pet-store/internal/infrastructure/responder"
	"pet-store/internal/modules/user/service"

	"github.com/go-chi/chi"
	"github.com/ptflp/godecoder"
)

type Userer interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
}

type User struct {
	service service.Userer
	responder.Responder
	godecoder.Decoder
}

func NewUser(service service.Userer, components *component.Components) Userer {
	return &User{service: service, Responder: components.Responder, Decoder: components.Decoder}
}


// @Summary Get a user by username
// @Tags user
// @Description Fetches a user by their username.
// @ID GetUser
// @Accept  json
// @Produce  json
// @Param username path string true "Username of the user"
// @Success 200 {object} GetUserResponseSuccess "Successfully retrieved user"
// @Router /user/{username} [get]
func (u *User) GetUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	out := u.service.GetByUsername(r.Context(), username)
	if out.ErrorCode != errors.NoError {
		msg := "getuser error"
		u.OutputJSON(w, UserResponseErr{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data:      msg,
		})
		return
	}

	u.OutputJSON(w, GetUserResponseSuccess{
		Success: true,
		DataUser: DataUser{
			ID:         out.User.ID,
			Username:   out.User.Username,
			FirstName:  out.User.FirstName,
			LastName:   out.User.LastName,
			Email:      out.User.Email,
			Password:   out.User.Password,
			Phone:      out.User.Phone,
			UserStatus: out.User.UserStatus,
		},
	})
}

// @Summary Update user information
// @Security ApiKeyAuth
// @Tags user
// @Description Updates the information of an existing user identified by username.
// @ID UpdateUser
// @Accept  json
// @Produce  json
// @Param username path string true "Username of the user to update"
// @Param user body service.UpdateUserRequest true "User data to update"
// @Success 200 {object} UserUpdateResponse "Successfully updated user data"
// @Router /user/{username} [put]
func (u *User) UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	var userdata service.UpdateUserRequest
	err := u.Decode(r.Body, &userdata)
	if err != nil {
		u.ErrorBadRequest(w, err)
		return
	}

	userdata.UserName = username

	out := u.service.UpdateUser(r.Context(), userdata)
	if out.ErrorCode != errors.NoError {
		msg := "updateuser error"
		u.OutputJSON(w, UserResponseErr{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data:      msg,
		})
	}
	message := fmt.Sprintf("%s data has been updated", username)
	u.OutputJSON(w, UserUpdateResponse{
		Success: true,
		Message: message,
	})
}
