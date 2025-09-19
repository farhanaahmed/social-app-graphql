package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/graphql-go/graphql"
)

// GraphQL schema
var socialAppSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    nil,
	Mutation: nil,
})

func main() {
	// Create a new router
	r := chi.NewRouter()

	// Define the GraphQL endpoint
	r.Post("/graphql", func(w http.ResponseWriter, req *http.Request) {
		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}

		if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Execute the GraphQL query
		result := graphql.Do(graphql.Params{
			Schema:         socialAppSchema,
			RequestString:  params.Query,
			VariableValues: params.Variables,
			OperationName:  params.OperationName,
			Context:        context.Background(),
		})

		// Return the result as a JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	// This is a GET endpoint to help you test with a browser
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GraphQL server is running at /graphql"))
	})

	log.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
