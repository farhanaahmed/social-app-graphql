package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

// Our GraphQL schema
var socialAppSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

const defaultPort = "8080"

func main() {
	// Database connection setup
	var err error
	connStr := "user=go_user password=gopass dbname=social_app_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close() // Ensure the connection is closed

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database!")

	// Router and GraphQL server setup
	r := chi.NewRouter()

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
		result := graphql.Do(graphql.Params{
			Schema:         socialAppSchema,
			RequestString:  params.Query,
			VariableValues: params.Variables,
			OperationName:  params.OperationName,
			//Context:      context.WithValue(req.Context(), "db", db), // Pass the DB connection via context
			// Simulating logged-in user with ID 1
			Context: context.WithValue(context.WithValue(req.Context(), "db", db), "userID", 1),
		})
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GraphQL server is running at /graphql"))
	})

	log.Println("Server is running on", defaultPort)
	if err := http.ListenAndServe(":"+defaultPort, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
