package model

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	UserRoles []UserRoles `json:"user_roles"`
	UserSubscriptions []UserSubscriptions `json:"user_subscriptions"`
}

type UserRoles struct {
	RoleId  int    `json:"role_id"`
	RoleKey string `json:"role_key"`
}

type UserSubscriptions struct {
	SubscriptionID int    `json:"subscription_id"`
	Name           string `json:"name"`
	ProductID      int    `json:"product_id"`
}