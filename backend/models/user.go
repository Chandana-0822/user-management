package models

type User struct {
	UserID     int64  `json:"user_id"`
	UserName   string `json:"user_name"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	UserStatus string `json:"user_status"` // Values: I, A, T
	Department string `json:"department,omitempty"`
}

type SearchUsernameRequest struct {
	Username  string `json:"username" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type SearchUsernameResponse struct {
	Message     string   `json:"message,omitempty"`     // message if username already exists
	Suggestions []string `json:"suggestions,omitempty"` // Suggested usernames if the username is available
}
