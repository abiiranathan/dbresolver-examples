package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	adaptor "github.com/gwatts/gin-adapter"

	"github.com/abiiranathan/dbresolver/dbresolver"
	"github.com/gin-gonic/gin"
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

	app := gin.Default()
	// Convert the middleware to a Fiber/v2 middleware using gin-adaptor
	app.Use(adaptor.Wrap(resolver.Middleware))

	app.GET("/", func(ctx *gin.Context) {
		// Do not use ctx.Get or ctx.MustGet
		db := ctx.Request.Context().Value(dbresolver.ConnectionContextKey).(*gorm.DB)

		fmt.Println(db)
		// Perform your queries with the db
		ctx.String(http.StatusOK, "Hello from Gin")
	})

	app.Run(":8080")
}
