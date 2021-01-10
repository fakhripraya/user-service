package data

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/go-hclog"
)

// Claims determine the current user token holder
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// User defines a struct for user flow
type User struct {
	logger hclog.Logger
}

// NewUser is a function to create new User struct
func NewUser(newLogger hclog.Logger) *User {
	return &User{newLogger}
}
