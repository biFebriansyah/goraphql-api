//go:generate go run generate.go

package graph

import (
	"github.com/biFebriansyah/goraphql/graph/service"
)

type Resolver struct {
	UserService    *service.UserService
	ProductService *service.ProductService
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }
