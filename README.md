# Go GraphQL Social App

This is a demo social media backend API built with **Go** and **GraphQL**. It provides core functionality for post management, serving as a foundational learning project.

## 🚀 Features Implemented

* **GraphQL API**: A single, efficient `/graphql` endpoint for all data operations.
* **Post Management**: Users can create new posts and fetch all existing posts.
* **Database Integration**: Data is persisted in a **PostgreSQL** database.
* **Authentication (Mocked)**: A hardcoded user ID is injected into the request context to simulate a logged-in user and enable post creation.

## 🛠️ Technologies Used

* **Go**: The main backend language.
* **GraphQL**: The query language for our API, implemented using `graphql-go/graphql`.
* **Chi**: A lightweight router for Go.
* **PostgreSQL**: Our relational database.

## 📦 Getting Started

### Prerequisites

You need to have **Go** (1.18 or higher) and **PostgreSQL** installed.

### Installation

1.  **Clone the repository**:
    ```bash
    git clone [https://github.com/your-username/go-social-app.git](https://github.com/your-username/go-social-app.git)
    cd go-social-app
    ```
2.  **Set up the database**:
    Create a database named `social_app_db` and the `users` and `posts` tables.
    ```sql
    -- Connect to your PostgreSQL instance
    CREATE DATABASE social_app_db;
    \c social_app_db;
    CREATE TABLE users (id SERIAL PRIMARY KEY, username VARCHAR(50) UNIQUE NOT NULL, password_hash TEXT NOT NULL);
    CREATE TABLE posts (id SERIAL PRIMARY KEY, user_id INT NOT NULL REFERENCES users(id), content TEXT NOT NULL, created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP);
    -- Insert a dummy user to enable post creation
    INSERT INTO users (id, username, password_hash) VALUES (1, 'testuser', 'dummypass');
    ```
3.  **Install dependencies**:
    ```bash
    go mod tidy
    ```
4.  **Run the server**:
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8080`.

## 🕹️ API Usage

All requests should be **POST** requests sent to `http://localhost:8080/graphql`.

### Example Queries and Mutations

#### Create a Post (Mutation)

```json
{
  "query": "mutation { createPost(content: \"This is my first post via GraphQL!\") { id content user { username } } }"
}
```

#### Fetch All Posts (Query)

```json
{
  "query": "mutation { createPost(content: \"This is my first post via GraphQL!\") { id content user { username } } }"
}
```

## 📝 Learning Points

- Go's Concurrency: Go's lightweight threads, goroutines, make the API efficient and fast at handling many user requests simultaneously. This is crucial for a social app.

- GraphQL Schema Definition in Go: Defining types, queries, and mutations directly in Go code.

- Resolvers: Implementing the backend logic that connects a GraphQL operation to the database.

- Database Interactions: Connecting to and performing SQL queries on a PostgreSQL database.

## 🛣️ Future Work

- User Authentication: Implement a secure user login system.

- CRUD Operations: Add updatePost and deletePost mutations.

- Social Features: Implement logic for users to follow each other.
