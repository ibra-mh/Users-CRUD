package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	model "main/Model"
	"net/http"

	"github.com/gorilla/mux"
)

func GetUser(db *sql.DB) http.HandlerFunc {
	// The returned function is the actual HTTP handler.
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r) // `mux.Vars` retrieves a map of route parameters.
		id := vars["id"]
		// Define a variable to hold the data for the user.
		var u model.User

		// Execute an SQL query to select the user with the given ID.
		// `QueryRow` is used when the query is expected to return a single row.
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).
			Scan(&u.ID, &u.Name, &u.Email) // Map the result to the fields of the `User` struct.

		if err != nil {
			// If there's an error (e.g., no user found or database issue), log it and terminate the program.
			log.Fatal(err)
		}

		
		json.NewEncoder(w).Encode(u)
	}
}

func GetUsers(db *sql.DB) http.HandlerFunc {
    // The returned function is the actual HTTP handler.
    return func(w http.ResponseWriter, r *http.Request) {
        // Execute an SQL query to select all rows from the "users" table.
        rows, err := db.Query("SELECT * FROM users")
        if err != nil {
            // If there's an error while querying, log it and terminate the program.
            log.Fatal(err)
        }
        // Ensure the database rows are closed when the function finishes.
        defer rows.Close()

        // Create a slice to store the list of users. 
        // `[]User{}` initializes an empty slice of type `User`.
        users := []model.User{}

        // Iterate through each row returned by the query.
        for rows.Next() {
            // Create a variable to hold the data for one user.
            var u model.User

            // Map the current row's columns to the fields of the `User` struct.
            // `Scan` assigns values from the row to the fields of `u`.
            if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
                // If there's an error during scanning, log it and terminate the program.
                log.Fatal(err)
            }
            // Append the scanned user to the `users` slice.
            users = append(users, u)
        }

        // Check for any errors that might have occurred during the row iteration.
        if err := rows.Err(); err != nil {
            // If an error occurred, log it and terminate the program.
            log.Fatal(err)
        }

        // Encode the `users` slice into JSON format and write it to the HTTP response.
        // `w` is the ResponseWriter, which sends data back to the client.
        json.NewEncoder(w).Encode(users)
    }
}

//create user
func CreateUser(db *sql.DB) http.HandlerFunc {

    // This inner function is the actual HTTP handler
    return func(w http.ResponseWriter, r *http.Request) { 

        // Create a new User struct to hold the data
        var u model.User 

        // Decode the JSON data from the request body into the User struct
        // If there's an error, send a bad request response
        if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Insert a new row into the "users" table with the provided name and email
        // The RETURNING clause retrieves the generated ID for the new user
        err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)

        // If there's an error during the database operation, send an internal server error
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Encode the created User struct as JSON and send it as the response
        // If there's an error during encoding, send an internal server error
        if err := json.NewEncoder(w).Encode(u); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}

// updateUser is a function that returns an HTTP handler function to update a user in the database.
func UpdateUser(db *sql.DB) http.HandlerFunc {
    // The function returns another function, which is an HTTP handler.
    return func(w http.ResponseWriter, r *http.Request) {
        // Declare a variable of type User to hold the incoming request data.
        var u model.User
        
        // Decode the JSON request body into the User struct.
        // If the body doesn't match the structure, an error will occur.
        json.NewDecoder(r.Body).Decode(&u)

        // Extract the variables from the request using the mux router.
        vars := mux.Vars(r) // `vars` is a map containing path parameters from the request.
        
        // Retrieve the "id" parameter from the `vars` map.
        id := vars["id"] 

        // Execute the SQL query to update the user in the database.
        // The query uses placeholders `$1`, `$2`, and `$3` to prevent SQL injection.
        _, err := db.Exec(
            "UPDATE users SET name = $1, email = $2 WHERE id = $3", // SQL query to update the user.
            u.Name, // Replace `$1` with the `Name` field from the User struct.
            u.Email, // Replace `$2` with the `Email` field from the User struct.
            id,     // Replace `$3` with the extracted user ID.
        )
        
        // If there is an error during the database operation, log the error and stop the application.
        if err != nil {
            log.Fatal(err) // Note: `log.Fatal` will terminate the program. Consider using `log.Println` to avoid crashing the server.
        }

        // Encode the updated user data into JSON and write it to the response.
        json.NewEncoder(w).Encode(u)
    }
}

// deleteUser function takes a database connection (db) as input
// and returns an http.HandlerFunc that handles deleting a user.
func DeleteUser(db *sql.DB) http.HandlerFunc {

    // This inner function is the actual handler for HTTP requests.
    // It receives the response writer (w) and the request object (r) as arguments.
    return func(w http.ResponseWriter, r *http.Request) {

        // Extract the user ID from the URL parameters.
        vars := mux.Vars(r) 
        id := vars["id"] 

        // Execute the SQL DELETE statement.
        // The $1 placeholder is used to prevent SQL injection.
        _,err := db.Exec("DELETE FROM users WHERE id = $1", id) 

        // Check if any error occurred during the deletion.
        if err != nil { 
            // If an error occurred, log it and exit the program.
            log.Fatal(err) 
        }

        // If the deletion was successful, send a success message to the client.
        json.NewEncoder(w).Encode("User deleted") 
    }
}

