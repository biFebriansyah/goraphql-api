package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/biFebriansyah/goraphql/graph"
	"github.com/biFebriansyah/goraphql/graph/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	_ "github.com/joho/godotenv/autoload"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	server := fiber.New()

	server.All("/query", adaptor.HTTPHandler(graphServer()))
	server.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/query")))

	if err := server.Listen(":8081"); err != nil {
		log.Fatal(err)
	}
}

func graphServer() *handler.Server {
	mongoDB := NewMongo()

	userCollection := mongoDB.Collection("users")
	productCollection := mongoDB.Collection("product")
	userService := service.NewUserService(userCollection)
	productService := service.NewProductService(productCollection)

	resolver := graph.Resolver{
		UserService:    userService,
		ProductService: productService,
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}
