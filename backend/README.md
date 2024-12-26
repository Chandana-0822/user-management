# User Management Backend Service

This README provides instructions for setting up and running the User Management backend service locally, including database setup, migration script execution, and API testing.

## Prerequisites

1. **Install Go:**
   - Download and install Go from [https://go.dev/dl/](https://go.dev/dl/).
   - Verify installation:
     ```bash
     go version
     ```

2. **Install PostgreSQL:**
   - Download and install PostgreSQL from [https://www.postgresql.org/download/](https://www.postgresql.org/download/).
   - Verify installation:
     ```bash
     psql --version
     ```

3. **Install `psql`:**
   - `psql` is included with PostgreSQL. After installing PostgreSQL, you can use `psql` from the command line.

4. **Set Up psql Path (Windows):**

   - Locate the PostgreSQL installation directory (e.g., C:\Program Files\PostgreSQL\<version>\bin).
   - Add it to the Path environment variable:
   - Open System Properties > Environment Variables.
   - Under System Variables, select Path and click Edit.
   - Add the PostgreSQL bin directory path.
   - Verify installation:
     ```bash
     psql --version
     ```

## Setting Up the Database

1. **Start PostgreSQL:**
   - Ensure the PostgreSQL service is running.

2. **Create the Database:**
   - Open a terminal and run the following commands:
     ```bash
     psql -U postgres
     ```
   - In the `psql` prompt:
     ```sql
     CREATE DATABASE user_management;
     ```

3. **Run the Migration Script:**
   - Exit the `psql` prompt by typing `\q`.
   - Navigate to the `database/migrations/` directory.
   - Run the migration script:
     ```bash
     psql -U postgres -d user_management -f migration.sql
     ```

   **Note:** Replace `postgres` with your PostgreSQL username if it's different.

4. **Verify the Table:**
   - Log back into the database:
     ```bash
     psql -U postgres -d user_management
     ```
   - Check if the table exists:
     ```sql
     \dt
     ```

## Running the Application Locally

1. **Clone the Repository:**
   ```bash
   git clone <repository_url>
   cd backend
   ```

2. **Install Dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the Application:**
   ```bash
   go run main.go
   ```

4. **Access the Application:**
   - The application runs on `http://localhost:8080`.

## Testing the APIs

You can use `curl` commands or import them into Postman.

### Sample `curl` Commands

1. **Create a User:**
   ```bash
   curl -X POST http://localhost:8080/v1/users \
   -H "Content-Type: application/json" \
   -d '{"user_name": "johndoe","first_name": "John","last_name": "Doe","user_status": "A","department": "Engineering"}'
   ```

2. **Get All Users:**
   ```bash
   curl -X GET http://localhost:8080/v1/users
   ```

3. **Update a User:**
   ```bash
   curl -X PUT http://localhost:8080/v1/users/1 \
   -H "Content-Type: application/json" \
   -d '{"user_name": "updated_username","first_name": "UpdatedFirstName","last_name": "UpdatedLastName","user_status": "I","department": "UpdatedDepartment"}'
   ```

4. **Delete a User:**
   ```bash
   curl -X DELETE http://localhost:8080/v1/users/1
   ```

### Importing into Postman

1. Open Postman.
2. Create a new request collection.
3. Add each `curl` command as a request by copying the URL, headers, and body.
4. Save the requests to your collection.
5. Test the APIs.

## License
This project is licensed under the MIT License.
