package domain

import "time"

type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	CreatedAt  time.Time `json:"member_since"`
	ModifiedAt time.Time `json:"last_login"`
	// TODO:
	// - Activate / Validate accouot
	// - Enable MFA
	// - Login with Google account
	// - Force reset password
}
