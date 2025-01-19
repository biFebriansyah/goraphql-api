package rest

type SigninInput struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type UserToken struct {
	Token string `json:"token" bson:"token"`
}
