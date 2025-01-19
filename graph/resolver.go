//go:generate go run generate.go

package graph

import (
	"github.com/biFebriansyah/goraphql/graph/service"
)

type Resolver struct {
	UserService    *service.UserService
	ProductService *service.ProductService
}
