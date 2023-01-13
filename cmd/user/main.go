package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	userserver "zpe/internal/user/server"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"

	"github.com/go-chi/chi/v5"
)

func main() {
	// load swagger for request validation
	swagger, err := userserver.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Setting up the router:
	r := chi.NewRouter()

	// Use validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))

	// register our userStore above as the handler for the interface
	userHandler := userserver.NewUserStore()
	userserver.HandlerFromMux(userHandler, r)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", 3000),
	}

	log.Fatal(s.ListenAndServe())
}
