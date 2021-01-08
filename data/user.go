package data

import (
	"github.com/hashicorp/go-hclog"
)

// User defines a struct for user flow
type User struct {
	logger hclog.Logger
}

// NewUser is a function to create new User struct
func NewUser(newLogger hclog.Logger) *User {
	return &User{newLogger}
}
