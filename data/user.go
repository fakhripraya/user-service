package data

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fakhripraya/user-service/config"
	"github.com/fakhripraya/user-service/database"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"github.com/srinathgs/mysqlstore"
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

// GetCurrentUser will get the current user login info
func (user *User) GetCurrentUser(rw http.ResponseWriter, r *http.Request, store *mysqlstore.MySQLStore) (*database.MasterUser, error) {

	// Get a session (existing/new)
	session, err := store.Get(r, "session-name")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		return nil, err
	}

	// check the logged in user from the session
	// if user available, get the user info from the session
	if session.Values["userLoggedin"] == nil {
		rw.WriteHeader(http.StatusUnauthorized)

		return nil, fmt.Errorf("Error 401")
	}

	// work with database
	// look for the current user logged in in the db
	var currentUser database.MasterUser
	if err := config.DB.Where("username = ?", session.Values["userLoggedin"].(string)).First(&currentUser).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)

		return nil, err
	}

	return &currentUser, nil

}

// UpdateUser is a function to update the given user model
func (user *User) UpdateUser(targetUser *database.MasterUser) error {

	// proceed to update the user with transaction scope
	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// work with database
		// looking for an existing user to update
		var updateUser database.MasterUser
		var dbErr error

		if dbErr = config.DB.Where("username = ?", targetUser.Username).First(&updateUser).Error; dbErr != nil {

			return fmt.Errorf("username does not exist")
		}

		updateUser.RoleID = targetUser.RoleID
		updateUser.DisplayName = targetUser.DisplayName
		updateUser.Email = targetUser.Email
		updateUser.Phone = targetUser.Phone
		updateUser.Modified = time.Now().Local()
		updateUser.ModifiedBy = targetUser.Username

		// update the user
		dbErr = config.DB.Save(updateUser).Error

		// if transaction error, return the error
		if dbErr != nil {
			return dbErr
		}

		return nil

	})

	// if transaction error, return the error
	if err != nil {
		return err
	}

	return nil
}
