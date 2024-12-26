package services

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func GetAllUsers() ([]models.User, error) {
	// Build the query to fetch all users
	query := squirrel.Select("user_id", "user_name", "first_name", "last_name", "email", "user_status", "department").
		From("users").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	// Execute the query
	rows, err := database.DB.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map rows to User models
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.UserStatus, &user.Department); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func CreateUser(user *models.User) error {
	// Check if the email already exists
	selectQuery := squirrel.Select("user_id").
		From("users").
		Where(squirrel.Eq{"email": user.Email}).
		PlaceholderFormat(squirrel.Dollar)

	selectSQL, selectArgs, err := selectQuery.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	var existingID int64
	err = database.DB.QueryRow(selectSQL, selectArgs...).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing email: %w", err)
	}

	if existingID != 0 {
		return fmt.Errorf("email already exists")
	}

	// Insert the user without specifying user_id
	insertQuery := squirrel.Insert("users").
		Columns("user_name", "first_name", "last_name", "email", "user_status", "department").
		Values(user.UserName, user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department).
		Suffix("RETURNING user_id").
		PlaceholderFormat(squirrel.Dollar)

	insertSQL, insertArgs, err := insertQuery.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = database.DB.QueryRow(insertSQL, insertArgs...).Scan(&user.UserID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func UpdateUser(userID int64, updatedUser models.User) error {
	// Check if the user exists
	selectQuery := squirrel.Select("user_id").
		From("users").
		Where(squirrel.Eq{"user_id": userID}).
		PlaceholderFormat(squirrel.Dollar)

	selectSQL, selectArgs, err := selectQuery.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	var existingID int64
	err = database.DB.QueryRow(selectSQL, selectArgs...).Scan(&existingID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Build the UPDATE query (excluding email)
	updateQuery := squirrel.Update("users").
		Set("user_name", updatedUser.UserName).
		Set("first_name", updatedUser.FirstName).
		Set("last_name", updatedUser.LastName).
		Set("user_status", updatedUser.UserStatus).
		Set("department", updatedUser.Department).
		Where(squirrel.Eq{"user_id": userID}).
		PlaceholderFormat(squirrel.Dollar)

	updateSQL, updateArgs, err := updateQuery.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Execute the UPDATE query
	_, err = database.DB.Exec(updateSQL, updateArgs...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func DeleteUser(userID int64) error {
	// Build the DELETE query
	query := squirrel.Delete("users").
		Where(squirrel.Eq{"user_id": userID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	// Execute the DELETE query
	_, err = database.DB.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func SearchUsername(username string) (bool, error) {
	query := squirrel.Select("user_name").
		From("users").
		Where(squirrel.Eq{"user_name": username}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var existingUsername string
	err = database.DB.QueryRow(sql, args...).Scan(&existingUsername)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil // Username not found
		}
		return false, fmt.Errorf("database error: %w", err)
	}

	return true, nil // Username exists
}
