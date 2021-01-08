package entities

import "time"

// User is an entity to communicate with MasterUser table in database
type User struct {
	ID             uint      `json:"id"`
	RoleID         uint      `json:"role_id"`
	Username       string    `json:"username"`
	DisplayName    string    `json:"displayname"`
	Password       []byte    `json:"password"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	LoginFailCount uint      `json:"login_fail_count"`
	IsVerified     bool      `json:"is_verified"`
	IsActive       bool      `json:"is_active"`
	Created        time.Time `json:"created"`
	CreatedBy      string    `json:"created_by"`
	Modified       time.Time `json:"modified"`
	ModifiedBy     string    `json:"modified_by"`
}
