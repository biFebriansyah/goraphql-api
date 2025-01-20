// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Mutation struct {
}

type NewProduct struct {
	Name      string     `json:"name" bson:"name"`
	Price     int32      `json:"price" bson:"price"`
	Stock     int32      `json:"stock" bson:"stock"`
	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at"`
}

type ProductDetail struct {
	Data []*Products  `json:"data" bson:"data"`
	Meta *ProductMeta `json:"meta" bson:"meta"`
}

type ProductMeta struct {
	Total int32 `json:"total" bson:"total"`
	Prev  int32 `json:"prev" bson:"prev"`
	Next  int32 `json:"next" bson:"next"`
}

type Products struct {
	ID        string     `json:"_id" bson:"_id"`
	Name      string     `json:"name" bson:"name"`
	Price     int32      `json:"price" bson:"price"`
	Stock     int32      `json:"stock" bson:"stock"`
	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at"`
}

type Query struct {
}

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

type UpdateInput struct {
	ID       string  `json:"_id" bson:"_id"`
	Name     *string `json:"name,omitempty" bson:"name"`
	Email    *string `json:"email,omitempty" bson:"email"`
	Password *string `json:"password,omitempty" bson:"password"`
}

type UpdateProduct struct {
	ID    string  `json:"_id" bson:"_id"`
	Name  *string `json:"name,omitempty" bson:"name"`
	Price *int32  `json:"price,omitempty" bson:"price"`
	Stock *int32  `json:"stock,omitempty" bson:"stock"`
}

type UserMeta struct {
	Total int32 `json:"total" bson:"total"`
	Prev  int32 `json:"prev" bson:"prev"`
	Next  int32 `json:"next" bson:"next"`
}

type UserToken struct {
	Token string `json:"token" bson:"token"`
}

type Users struct {
	ID        string     `json:"_id" bson:"_id"`
	Name      string     `json:"name" bson:"name"`
	Email     string     `json:"email" bson:"email"`
	Password  string     `json:"password" bson:"password"`
	Admin     *bool      `json:"admin,omitempty" bson:"admin"`
	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at"`
}

type UsersDetail struct {
	Data []*Users  `json:"data" bson:"data"`
	Meta *UserMeta `json:"meta" bson:"meta"`
}
