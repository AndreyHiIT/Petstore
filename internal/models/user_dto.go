package models

import (
	"pet-store/internal/infrastructure/db/types"
	"time"
)

type UserDTO struct {
	ID        int              `json:"id" db:"id" db_type:"BIGSERIAL primary key" db_default:"not null"`
	FirstName types.NullString `json:"firstname" db:"firstname" db_type:"varchar(255)" db_default:"default null" db_ops:"create,update"`
	UserName  types.NullString `json:"username" db:"username" db_type:"varchar(255)" db_default:"not null" db_ops:"create,update"`
	LastName  types.NullString `json:"lastname" db:"lastname" db_type:"varchar(255)" db_default:"default null" db_ops:"create,update"`
	Phone     types.NullString `json:"phone" db:"phone" db_type:"varchar(255)" db_default:"default null" db_index:"index,unique" db_ops:"create,update"`
	Email     types.NullString `json:"email" db:"email" db_type:"varchar(255)" db_default:"default null" db_index:"index,unique" db_ops:"create,update"`
	Password  types.NullString `json:"password" db:"password" db_type:"varchar(255)" db_default:"default null" db_ops:"create,update"`
	Status    int              `json:"status" db:"status" db_type:"int" db_default:"default 0" db_ops:"create,update"`
	DeletedAt types.NullTime   `json:"deleted_at" db:"deleted_at" db_type:"timestamp" db_default:"default null" db_index:"index"`
}

func (u *UserDTO) TableName() string {
	return "users"
}

func (u *UserDTO) OnCreate() []string {
	return []string{}
}

func (u *UserDTO) SetID(id int) *UserDTO {
	u.ID = id
	return u
}

func (u *UserDTO) GetID() int {
	return u.ID
}

func (u *UserDTO) SetUserName(username string) *UserDTO {
	u.UserName = types.NewNullString(username)
	return u
}

func (u *UserDTO) GetUserName() string {
	return u.UserName.String
}

func (u *UserDTO) SetPhone(phone string) *UserDTO {
	u.Phone = types.NewNullString(phone)
	return u
}

func (u *UserDTO) GetFirstName() string {
	return u.FirstName.String
}

func (u *UserDTO) GetLastName() string {
	return u.LastName.String
}

func (u *UserDTO) GetPhone() string {
	return u.Phone.String
}

func (u *UserDTO) SetEmail(email string) *UserDTO {
	u.Email = types.NewNullString(email)
	return u
}

func (u *UserDTO) GetEmail() string {
	return u.Email.String
}

func (u *UserDTO) SetPassword(password string) *UserDTO {
	u.Password = types.NewNullString(password)
	return u
}

func (u *UserDTO) SetFirstName(firstname string) *UserDTO {
	u.FirstName = types.NewNullString(firstname)
	return u
}

func (u *UserDTO) SetLastName(lastname string) *UserDTO {
	u.LastName = types.NewNullString(lastname)
	return u
}

func (u *UserDTO) GetPassword() string {
	return u.Password.String
}

func (u *UserDTO) SetStatus(status int) *UserDTO {
	u.Status = status
	return u
}

func (u *UserDTO) GetStatus() int {
	return u.Status
}

func (s *UserDTO) SetDeletedAt(deletedAt time.Time) *UserDTO {
	s.DeletedAt.Time.Time = deletedAt
	s.DeletedAt.Time.Valid = true
	return s
}

func (s *UserDTO) GetDeletedAt() time.Time {
	return s.DeletedAt.Time.Time
}
