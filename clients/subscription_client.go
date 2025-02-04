package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"main/Model"
)

func GetUserSubscriptions(userID int) ([]model.UserSubscriptions, error) {
    if userID == 0 {
        return nil, fmt.Errorf("invalid user ID")
    }
    // Define the URL to fetch user subscriptions
    url := fmt.Sprintf("http://go-app:8002/user_subscriptions?user_id=%d", userID)
    
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    res, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error making request: %v", err)
    }
    defer res.Body.Close()

    // Read the response body
    body, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %v", err)
    }

    var subscriptions []model.UserSubscriptions
    err = json.Unmarshal(body, &subscriptions)
    if err != nil {
        return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
    }

    return subscriptions, nil
}



// package clients

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"main/Model"
// )

// func GetUserSubscriptions(ID int) *model.UserSubscriptions {
// 	if ID == 0 {
// 		return nil // Prevent invalid request
// 	}
// 	fmt.Println("---------------------------------------------------1")
// 	url := fmt.Sprintf("http://go-app:8002/user_subscriptions/%d", ID)
// 	log.Println(url)
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		log.Println("Error creating request:", err)
// 		return nil
// 	}
// 	fmt.Println("---------------------------------------------------2")
// 	res, err := client.Do(req)
// 	if err != nil {
// 		log.Println("Error making request:", err)
// 		return nil
// 	}
// 	defer res.Body.Close()

// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		log.Println("Error reading response body:", err)
// 		return nil
// 	}
// 	fmt.Println("---------------------------------------------------3")
// 	fmt.Println("Response Body:", string(body)) // Debugging

// 	var userSubscriptions model.UserSubscriptions
// 	err = json.Unmarshal(body, &userSubscriptions)
// 	if err != nil {
// 		fmt.Println("---------------------------------------------------5")
// 		log.Println("Error unmarshalling JSON:", err)
// 		return nil
// 	}
// 	fmt.Println("---------------------------------------------------4")
// 		return &userSubscriptions // Return the first subscription

	
// }
