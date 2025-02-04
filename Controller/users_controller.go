package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	model "main/Model"
	"main/clients"
	"net/http"

	"github.com/gorilla/mux"
)

func GetUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        var u model.User
        err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).
            Scan(&u.ID, &u.Name, &u.Email)

        if err != nil {
            log.Fatal(err)
        }

        // Get user roles
        roles, err := clients.GetUserRoles(u.Email)
        if err != nil {
            log.Println("Error fetching user roles:", err)
            u.UserRoles = []model.UserRoles{} // Return an empty slice instead of nil
        } else {
            u.UserRoles = roles
        }

        // Get user subscriptions
        subscriptions, err := clients.GetUserSubscriptions(u.ID)
        if err != nil {
            log.Println("Error fetching user subscriptions:", err)
            u.UserSubscriptions = []model.UserSubscriptions{}
        } else {
            u.UserSubscriptions = subscriptions
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(u)
    }
}



// func GetUser(db *sql.DB) http.HandlerFunc {
//     // The returned function is the actual HTTP handler.
//     return func(w http.ResponseWriter, r *http.Request) {
//         vars := mux.Vars(r) // `mux.Vars` retrieves a map of route parameters.
//         id := vars["id"]
        
//         // Define a variable to hold the data for the user.
//         var u model.User

//         // Execute an SQL query to select the user with the given ID.
//         // `QueryRow` is used when the query is expected to return a single row.
//         err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).
//             Scan(&u.ID, &u.Name, &u.Email) // Map the result to the fields of the `User` struct.

//         if err != nil {
//             // If there's an error (e.g., no user found or database issue), log it and terminate the program.
//             log.Fatal(err)
//         }

//         // Get roles for the user using the GetUserRoles function from the clients package
//         roles := clients.GetUserRoles(u.Email)
//         u.UserRoles = roles

//         // Get subscriptions for the user using the GetUserSubscriptions function from the clients package
//         subscriptions, err := clients.GetUserSubscriptions(u.ID)
//         if err != nil {
//             log.Println("Error fetching subscriptions:", err)
//             // Optionally, handle the error (return empty list or specific error message)
//         }
//         u.UserSubscriptions = subscriptions

//         // Return the user with roles and subscriptions
//         w.Header().Set("Content-Type", "application/json")
//         json.NewEncoder(w).Encode(u)
//     }
// }



func GetUsers(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        includeRoles := r.URL.Query().Get("include-roles") == "true"
        includeSubscriptions := r.URL.Query().Get("include-subscriptions") == "true"

        rows, err := db.Query("SELECT * FROM users")
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        users := []map[string]interface{}{}

        for rows.Next() {
            var u model.User
            if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
                log.Fatal(err)
            }

            userData := map[string]interface{}{
                "id":    u.ID,
                "name":  u.Name,
                "email": u.Email,
            }

            // Fetch and add user roles if requested
            if includeRoles {
                roles, err := clients.GetUserRoles(u.Email)
                if err != nil {
                    log.Println("Error fetching user roles:", err)
                    userData["user_roles"] = []model.UserRoles{}
                } else {
                    userData["user_roles"] = roles
                }
            }

            // Fetch and add user subscriptions if requested
            if includeSubscriptions {
                subscriptions, err := clients.GetUserSubscriptions(u.ID)
                if err != nil {
                    log.Println("Error fetching user subscriptions:", err)
                    userData["user_subscriptions"] = []model.UserSubscriptions{}
                } else {
                    userData["user_subscriptions"] = subscriptions
                }
            }

            users = append(users, userData)
        }

        if err := rows.Err(); err != nil {
            log.Fatal(err)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    }
}

// func GetUsers(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Parse query parameters for roles and subscriptions inclusion
// 		includeRoles := r.URL.Query().Get("include-roles") == "true"
// 		includeSubscriptions := r.URL.Query().Get("include-subscriptions") == "true"

// 		// Query all users from the database
// 		rows, err := db.Query("SELECT * FROM users")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer rows.Close()

// 		users := []map[string]interface{}{}

// 		// Iterate over the result rows
// 		for rows.Next() {
// 			var u model.User
// 			if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
// 				log.Fatal(err)
// 			}

// 			userData := map[string]interface{}{
// 				"id":    u.ID,
// 				"name":  u.Name,
// 				"email": u.Email,
// 			}

// 			// Conditionally fetch and add user roles
// 			if includeRoles {
// 				userData["user_roles"] = clients.GetUserRoles(u.Email)
// 			}

// 			// Conditionally fetch and add user subscriptions
// 			if includeSubscriptions {
// 				subscriptions, err := clients.GetUserSubscriptions(u.ID)
// 				if err != nil {
// 					log.Println("Error fetching user subscriptions:", err)
// 					userData["user_subscriptions"] = []model.UserSubscriptions{}
// 				} else {
// 					userData["user_subscriptions"] = subscriptions
// 				}
// 			}

// 			// Add user data to the list
// 			users = append(users, userData)
// 		}

// 		// Error check for rows iteration
// 		if err := rows.Err(); err != nil {
// 			log.Fatal(err)
// 		}

// 		// Return the list of users with their roles and subscriptions
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(users)
// 	}
// }

// create user
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
			u.Name,  // Replace `$1` with the `Name` field from the User struct.
			u.Email, // Replace `$2` with the `Email` field from the User struct.
			id,      // Replace `$3` with the extracted user ID.
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
		_, err := db.Exec("DELETE FROM users WHERE id = $1", id)

		// Check if any error occurred during the deletion.
		if err != nil {
			// If an error occurred, log it and exit the program.
			log.Fatal(err)
		}

		// If the deletion was successful, send a success message to the client.
		json.NewEncoder(w).Encode("User deleted")
	}
}
