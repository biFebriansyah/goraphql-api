package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/biFebriansyah/goraphql/graph"
	"github.com/biFebriansyah/goraphql/graph/service"
	"github.com/biFebriansyah/goraphql/rest"
	"github.com/biFebriansyah/goraphql/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	server := fiber.New()
	server.Use(cors.New())

	mongoDB := utils.NewMongo()
	userCollection := mongoDB.Collection("users")
	productCollection := mongoDB.Collection("product")
	userService := service.NewUserService(userCollection)
	productService := service.NewProductService(productCollection)

	resolver := graph.Resolver{
		UserService:    userService,
		ProductService: productService,
	}

	restHandler := rest.RestHandler{UserService: userService}
	server.Post("/signin", restHandler.SignIn)
	server.Post("/signup", restHandler.SignUp)

	graphServer := graph.GraphServer(resolver)
	server.All("/query", restHandler.AuthMiddleware, adaptor.HTTPHandler(graphServer))
	server.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/query")))

	if err := server.Listen(os.Getenv("APP_PORT")); err != nil {
		log.Fatal(err)
	}
}
