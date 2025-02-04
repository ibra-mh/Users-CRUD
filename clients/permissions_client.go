
package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"io"
	"net/http"
	"main/Model"
)
// func GetUserRoles(email string) model.UserRoles {
//     url := fmt.Sprintf("http://go-app:8000/user-roles?email=%s", email)
//     method := "GET"

//     client := &http.Client{}
//     req, err := http.NewRequest(method, url, nil)
//     if err != nil {
//         fmt.Println(err)
//         return model.UserRoles{}  // Return an empty UserRoles if an error occurs
//     }

//     res, err := client.Do(req)
//     if err != nil {
//         fmt.Println(err)
//         return model.UserRoles{}
//     }
//     defer res.Body.Close()

//     body, err := io.ReadAll(res.Body)
//     if err != nil {
//         fmt.Println(err)
//         return model.UserRoles{}
//     }

//     fmt.Println(string(body)) // Debugging the response

//     var userRoles []model.UserRoles  // Use the correct struct name, 'UserRoles' instead of 'UserRole'
//     err = json.Unmarshal(body, &userRoles)
//     if err != nil {
//         log.Println("Error unmarshalling JSON:", err)
//         return model.UserRoles{}  // Return empty UserRoles if unmarshalling fails
//     }

//     // Check if any user roles are found, and return the first one
//     if len(userRoles) > 0 {
//         return userRoles[0]  // Return the first role
//     }

//     return model.UserRoles{}  // Return an empty object if no roles found
// }

func GetUserRoles(email string) ([]model.UserRoles, error) {
    url := fmt.Sprintf("http://go-app:8000/user-roles?email=%s", email)
    method := "GET"

    client := &http.Client{}
    req, err := http.NewRequest(method, url, nil)
    if err != nil {
        fmt.Println(err)
        return nil, err  // Return nil slice and error
    }

    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }

    fmt.Println(string(body)) // Debugging the response

    var userRoles []model.UserRoles
    err = json.Unmarshal(body, &userRoles)
    if err != nil {
        log.Println("Error unmarshalling JSON:", err)
        return nil, err
    }

    return userRoles, nil  // Return the slice of roles
}



