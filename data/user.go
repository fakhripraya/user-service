package data

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fakhripraya/user-service/config"
	"github.com/fakhripraya/user-service/entities"
	"github.com/hashicorp/go-hclog"
)

// Claims determine the current user token holder
type Claims struct {
	Username string
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

// UpdateUser is a function to update the given user model
func (user *User) UpdateUser(targetUser *entities.User) error {

	// work with database
	// looking for an existing user to update
	var updateUser entities.User
	if err := config.DB.Where("username = ?", targetUser.Username).First(&updateUser).Error; err != nil {

		return fmt.Errorf("username does not exist")
	}

	updateUser.RoleID = targetUser.RoleID
	updateUser.DisplayName = targetUser.DisplayName
	updateUser.Email = targetUser.Email
	updateUser.Phone = targetUser.Phone
	updateUser.Modified = time.Now().Local()
	updateUser.ModifiedBy = "SYSTEM"

	config.DB.Save(updateUser)

	return nil
}
