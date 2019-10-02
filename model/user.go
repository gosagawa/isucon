package model

import (
	"time"

	"github.com/wcl48/valval"
)

// User ユーザ
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Tel       *string   `json:"tel"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserValidate ユーザバリデート
func UserValidate(user User) error {
	Validator := valval.Object(valval.M{
		"Name": valval.String(
			valval.MaxLength(20),
		),
	})

	return Validator.Validate(user)
}
