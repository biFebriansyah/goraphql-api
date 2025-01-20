package rest

import "time"

type SigninInput struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type SignupInput struct {
	Name      string     `json:"name" bson:"name"`
	Email     string     `json:"email" bson:"email"`
	Password  string     `json:"password" bson:"password"`
	Admin     *bool      `json:"admin,omitempty" bson:"admin"`
	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at"`
}

type UserToken struct {
	Token string `json:"token" bson:"token"`
}
