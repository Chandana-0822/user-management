package handlers

import (
	"backend/models"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetUsersHandler(c echo.Context) error {
	users, err := services.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func CreateUserHandler(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}

	// Check if username already exists
	exists, err := services.SearchUsername(user.UserName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error while checking username"})
	}
	if exists {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Username already exists"})
	}

	err = services.CreateUser(&user)
	if err != nil {
		if err.Error() == "email already exists" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "email already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
}

func UpdateUserHandler(c echo.Context) error {
	// Bind request body to updatedUser (excluding email)
	type UpdateUserRequest struct {
		UserName   string `json:"user_name"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		UserStatus string `json:"user_status"` // Values: I, A, T
		Department string `json:"department,omitempty"`
	}
	var updateUserRequest UpdateUserRequest
	if err := c.Bind(&updateUserRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Map request data to the models.User struct
	updatedUser := models.User{
		UserName:   updateUserRequest.UserName,
		FirstName:  updateUserRequest.FirstName,
		LastName:   updateUserRequest.LastName,
		UserStatus: updateUserRequest.UserStatus,
		Department: updateUserRequest.Department,
	}

	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// Call the service to update the user
	err = services.UpdateUser(int64(userID), updatedUser)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, updateUserRequest)
}

func DeleteUserHandler(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	err = services.DeleteUser(int64(intID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func SearchUsernameHandler(c echo.Context) error {
	var req models.SearchUsernameRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Check if the username exists
	exists, err := services.SearchUsername(req.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if exists {
		// Generate smart suggestions
		suggestions, err := services.GenerateUsernameSuggestions(req.FirstName, req.LastName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating suggestions"})
		}
		response := models.SearchUsernameResponse{
			Message:     "Username already exists",
			Suggestions: suggestions,
		}
		return c.JSON(http.StatusOK, response)
	} else {
		response := models.SearchUsernameResponse{
			Message: "Username is Ready to use",
		}
		return c.JSON(http.StatusOK, response)
	}
}
