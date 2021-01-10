package models

import "github.com/gin-gonic/gin"

// UserDTO for dto
type UserDTO struct {
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	Organization Organization `json:"organization"`
}

// User user object
type User struct {
	Email        string       `firestore:"email,omitempty"`
	Password     string       `firestore:"password,omitempty"`
	Organization Organization `firestore:"organization,omitempty"`
}

// UserToJSON to json
func UserToJSON(user User) interface{} {
	return gin.H{
		"email":        user.Email,
		"password":     user.Password,
		"organization": user.Organization,
	}
}
