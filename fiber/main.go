package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abiiranathan/dbresolver/dbresolver"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s config.yaml", os.Args[0])
	}

	dbresolver.SetHeaderName("apikey") // default: x-api-key
	configFile := os.Args[1]
	cfg, err := dbresolver.ConfigFromYAMLFile(configFile)
	if err != nil {
		panic(err)
	}

	resolver, err := dbresolver.New(cfg)
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	// Convert the middleware to a Fiber/v2 middleware.
	fiberMiddleware := adaptor.HTTPMiddleware(resolver.Middleware)
	app.Use(fiberMiddleware)

	app.Get("/", func(c *fiber.Ctx) error {
		// Get the database connection from the request
		// You can also use c.Context().Value()
		if db, ok := c.Locals(dbresolver.ConnectionContextKey).(*gorm.DB); ok {
			fmt.Println(db)
			// Run your queries or call your controllers here
			return c.SendString("Got a valid database connection")
		}
		return c.SendString("No database connection in context")

	})
	app.Listen(":8080")
}
