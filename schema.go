package main

import (
	"database/sql"
	"fmt"

	"github.com/graphql-go/graphql"
)

// Struct to represent a Post from the database
type Post struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

// Struct to represent a User from the database
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// Define the GraphQL User object type
var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.Int,
				Description: "The ID of the user.",
			},
			"username": &graphql.Field{
				Type:        graphql.String,
				Description: "The username of the user.",
			},
		},
	},
)

// Define the GraphQL Post object type
var postType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.Int,
				Description: "The ID of the post.",
			},
			"content": &graphql.Field{
				Type:        graphql.String,
				Description: "The content of the post.",
			},
			"user": &graphql.Field{
				Type:        userType,
				Description: "The user who created the post.",
			},
		},
	},
)

// Define the GraphQL root query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"posts": &graphql.Field{
			Type:        graphql.NewList(postType), // This query returns a list of Post objects
			Description: "Get a list of all posts",
			//Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			//	// Get the database connection from the request's context
			//	db := p.Context.Value("db").(*sql.DB)
			//
			//	// Run a SQL query to get all posts from the table
			//	rows, err := db.Query("SELECT id, user_id, content FROM posts")
			//	if err != nil {
			//		// If there's a problem with the query, we'll return an error
			//		return nil, err
			//	}
			//	defer rows.Close()
			//
			//	var posts []Post
			//	for rows.Next() {
			//		var post Post
			//		// Copy the data from the database row into our Go struct
			//		if err := rows.Scan(&post.ID, &post.UserID, &post.Content); err != nil {
			//			return nil, err
			//		}
			//		posts = append(posts, post)
			//	}
			//	// Return our list of posts
			//	return posts, nil
			//},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				db, ok := p.Context.Value("db").(*sql.DB)
				if !ok {
					return nil, fmt.Errorf("database connection not found in context")
				}

				// Retrieve the simulated user ID from the context
				userID, ok := p.Context.Value("userID").(int)
				if !ok {
					return nil, fmt.Errorf("user not authenticated")
				}

				content, _ := p.Args["content"].(string)

				var newPost Post
				err := db.QueryRow(
					"INSERT INTO posts (user_id, content) VALUES ($1, $2) RETURNING id, user_id, content",
					userID, content,
				).Scan(&newPost.ID, &newPost.UserID, &newPost.Content)

				if err != nil {
					return nil, fmt.Errorf("failed to create post: %w", err)
				}

				return newPost, nil
			},
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createPost": &graphql.Field{
			Type:        postType,
			Description: "Creates a new post",
			Args: graphql.FieldConfigArgument{
				"content": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Get the database connection and simulated user ID from context
				db, ok := p.Context.Value("db").(*sql.DB)
				if !ok {
					return nil, fmt.Errorf("database connection not found in context")
				}
				userID, ok := p.Context.Value("userID").(int)
				if !ok {
					return nil, fmt.Errorf("user not authenticated")
				}

				content, _ := p.Args["content"].(string)

				var newPost Post
				err := db.QueryRow(
					"INSERT INTO posts (user_id, content) VALUES ($1, $2) RETURNING id, user_id, content",
					userID, content,
				).Scan(&newPost.ID, &newPost.UserID, &newPost.Content)

				if err != nil {
					return nil, fmt.Errorf("failed to create post: %w", err)
				}

				return newPost, nil
			},
		},
	},
})
