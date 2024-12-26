package tests

import (
	"backend/database"
	"backend/models"
	"backend/services"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Replace the database instance
	database.DB = &database.SQLDB{DB: db}

	// Mock query and result
	rows := sqlmock.NewRows([]string{"user_id", "user_name", "first_name", "last_name", "email", "user_status", "department"}).
		AddRow(1, "testuser1", "Test", "User1", "test.user1@example.com", "A", "Engineering").
		AddRow(2, "testuser2", "Test", "User2", "test.user2@example.com", "I", "HR")
	mock.ExpectQuery("SELECT user_id, user_name, first_name, last_name, email, user_status, department FROM users").WillReturnRows(rows)

	// Call the function
	users, err := services.GetAllUsers()

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "testuser1", users[0].UserName)
	assert.Equal(t, "Engineering", users[0].Department)
	assert.Equal(t, "testuser2", users[1].UserName)
	assert.Equal(t, "HR", users[1].Department)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Replace the database instance
	database.DB = &database.SQLDB{DB: db}

	// Mock existing email check
	mock.ExpectQuery("SELECT user_id FROM users WHERE email = \\$1").
		WithArgs("test.user1@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"user_id"})) // No existing rows

	// Mock exec
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser1", "Test", "User1", "test.user1@example.com", "A", "Engineering").
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

	// Call the function
	user := &models.User{
		UserName:   "testuser1",
		FirstName:  "Test",
		LastName:   "User1",
		Email:      "test.user1@example.com",
		UserStatus: "A",
		Department: "Engineering",
	}
	err = services.CreateUser(user)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Replace the database instance
	database.DB = &database.SQLDB{DB: db}

	// Mock user existence check
	mock.ExpectQuery("SELECT user_id FROM users WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

	// Mock exec
	mock.ExpectExec("UPDATE users").
		WithArgs("testuser1_updated", "Test", "User1 Updated", "A", "Engineering Updated", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the function
	user := models.User{
		UserName:   "testuser1_updated",
		FirstName:  "Test",
		LastName:   "User1 Updated",
		UserStatus: "A",
		Department: "Engineering Updated",
	}
	err = services.UpdateUser(1, user)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Replace the database instance
	database.DB = &database.SQLDB{DB: db}

	// Mock exec
	mock.ExpectExec("DELETE FROM users WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the function
	err = services.DeleteUser(1)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchUsername(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Replace the database instance
	database.DB = &database.SQLDB{DB: db}

	// Case 1: Username exists
	mock.ExpectQuery("SELECT user_name FROM users WHERE user_name = \\$1").
		WithArgs("existing_user").
		WillReturnRows(sqlmock.NewRows([]string{"user_name"}).AddRow("existing_user"))

	exists, err := services.SearchUsername("existing_user")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Case 2: Username does not exist
	mock.ExpectQuery("SELECT user_name FROM users WHERE user_name = \\$1").
		WithArgs("non_existing_user").
		WillReturnError(sql.ErrNoRows)

	exists, err = services.SearchUsername("non_existing_user")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Case 3: Database error
	mock.ExpectQuery("SELECT user_name FROM users WHERE user_name = \\$1").
		WithArgs("error_user").
		WillReturnError(sql.ErrConnDone) // Simulate a database connection error

	exists, err = services.SearchUsername("error_user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.False(t, exists)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
