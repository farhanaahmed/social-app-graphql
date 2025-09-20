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
				// Add this Resolve function
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db, ok := p.Context.Value("db").(*sql.DB)
					if !ok {
						return nil, fmt.Errorf("database connection not found in context")
					}

					// The 'Source' field holds the Post object from the parent resolver
					post, ok := p.Source.(Post)
					if !ok {
						return nil, fmt.Errorf("source is not a Post object")
					}

					var user User
					err := db.QueryRow("SELECT id, username FROM users WHERE id = $1", post.UserID).
						Scan(&user.ID, &user.Username)

					if err != nil {
						return nil, fmt.Errorf("failed to fetch user: %w", err)
					}

					return user, nil
				},
			},
		},
	},
)

// Define the GraphQL root query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"posts": &graphql.Field{
			Type:        graphql.NewList(postType),
			Description: "Get a list of all posts",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// This is the correct logic for a query
				db, ok := p.Context.Value("db").(*sql.DB)
				if !ok {
					return nil, fmt.Errorf("database connection not found in context")
				}

				rows, err := db.Query("SELECT id, user_id, content FROM posts")
				if err != nil {
					// It's crucial to return a valid slice even if the query fails
					return nil, err
				}
				defer rows.Close()

				var posts []Post
				for rows.Next() {
					var post Post
					if err := rows.Scan(&post.ID, &post.UserID, &post.Content); err != nil {
						return nil, err
					}
					posts = append(posts, post)
				}
				return posts, nil
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
		"updatePost": &graphql.Field{
			Type:        postType,
			Description: "Updates an existing post",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"content": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				db, ok := p.Context.Value("db").(*sql.DB)
				if !ok {
					return nil, fmt.Errorf("database connection not found in context")
				}
				userID, ok := p.Context.Value("userID").(int)
				if !ok {
					return nil, fmt.Errorf("user not authenticated")
				}
				postID, _ := p.Args["id"].(int)
				content, _ := p.Args["content"].(string)

				// **Authorization Check:** Find the post's owner ID
				var ownerID int
				err := db.QueryRow("SELECT user_id FROM posts WHERE id = $1", postID).Scan(&ownerID)
				if err != nil {
					return nil, fmt.Errorf("post not found or failed to check ownership: %w", err)
				}
				if ownerID != userID {
					return nil, fmt.Errorf("unauthorized: you can only update your own posts")
				}
				// Execute the update
				var updatedPost Post
				err = db.QueryRow(
					"UPDATE posts SET content = $1 WHERE id = $2 RETURNING id, user_id, content",
					content, postID,
				).Scan(&updatedPost.ID, &updatedPost.UserID, &updatedPost.Content)
				if err != nil {
					return nil, fmt.Errorf("failed to update post: %w", err)
				}
				return updatedPost, nil

				return nil, nil
			},
		},
		"deletePost": &graphql.Field{
			Type:        graphql.Boolean, // Returns true/false for success
			Description: "Deletes an existing post",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				db, ok := p.Context.Value("db").(*sql.DB)
				if !ok {
					return nil, fmt.Errorf("database connection not found in context")
				}
				userID, ok := p.Context.Value("userID").(int)
				if !ok {
					return nil, fmt.Errorf("user not authenticated")
				}
				postID, _ := p.Args["id"].(int)
				// **Authorization Check:** Find the post's owner ID
				var ownerID int
				err := db.QueryRow("SELECT user_id FROM posts WHERE id = $1", postID).Scan(&ownerID)
				if err != nil {
					return nil, fmt.Errorf("post not found or failed to check ownership: %w", err)
				}
				if ownerID != userID {
					return nil, fmt.Errorf("unauthorized: you can only delete your own posts")
				}
				// Execute the delete
				res, err := db.Exec("DELETE FROM posts WHERE id = $1", postID)
				if err != nil {
					return nil, fmt.Errorf("failed to delete post: %w", err)
				}
				rowsAffected, _ := res.RowsAffected()
				return rowsAffected > 0, nil
			},
		},
	},
})
