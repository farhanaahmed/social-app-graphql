# Social App GraphQL

This is a demo social media backend API built with Go and GraphQL.
It provides core functionality for user authentication and content management, serving as a learning project to demonstrate the principles of building a robust and efficient API.

## 🚀 Features

* **User Authentication**: Secure user registration and login with password hashing and JWT-based authentication.
* **Post Management**: Users can create, update, and delete their own text-based posts.
* **Data Fetching**: Efficient data retrieval using GraphQL queries, allowing clients to request exactly the data they need.
* **Scalable Architecture**: Built with a clear separation of concerns, making it easy to extend for future features.

## 🛠️ Technologies Used

* **Go**: The main programming language for the backend.
* **GraphQL**: The query language and server-side runtime for our API.
* **gqlgen**: A powerful code-generation library for Go that simplifies GraphQL server development.
* **Chi**: A lightweight, idiomatic router for Go.
* **PostgreSQL**: A robust and scalable relational database for data persistence.
* **bcrypt**: A library for secure password hashing.
* **jwt**: A library for implementing JSON Web Tokens for authentication.

## 📦 Getting Started

### Prerequisites

You need to have the following installed on your machine:

* **Go**: Version 1.18 or higher.
* **PostgreSQL**: A running instance of a PostgreSQL database.

### Installation

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/your-username/social-app.git](https://github.com/your-username/social-app.git)
    cd social-app
    ```
2.  **Set up the database:**
    Create a new PostgreSQL database and two tables for `users` and `posts`.
    ```sql
    -- Connect to your PostgreSQL instance
    CREATE DATABASE social_app_db;

    \c social_app_db;

    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        password_hash TEXT NOT NULL
    );

    CREATE TABLE posts (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL REFERENCES users(id),
        content TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
    ```
3.  **Configure the database connection:**
    Update the connection string in your `main.go` file with your database credentials.

4.  **Install dependencies:**
    ```bash
    go mod tidy
    ```
5.  **Run the GraphQL server:**
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8080`.

## 🕹️ Usage

You can interact with the API using the GraphQL Playground, accessible at `http://localhost:8080/` in your web browser.

### Example Queries and Mutations

#### User Registration (Mutation)

```graphql
mutation RegisterUser {
  register(username: "testuser", password: "password123") {
    token
  }
}
```
#### User Login (Mutation)

```graphql
mutation LoginUser {
  login(username: "testuser", password: "password123") {
    token
  }
}
```
#### Create a Post (Mutation - requires a JWT in the Authorization header)

```graphql
mutation CreatePost {
  createPost(content: "This is my first post on this app!") {
    id
    content
    user {
      username
    }
  }
}
```
#### Fetch All Posts (Query)
```graphql
query AllPosts {
  posts {
    id
    content
    user {
      username
    }
  }
}
```
## 📝 Learning Points

- GraphQL Schema Definition: Understanding how to define types, queries, and mutations, which form the "contract" between the client and server.

- Resolvers: Implementing the backend logic that fetches and manipulates data. This is where you'll use Go to connect to the database and fulfill the requests defined in the schema.

- Code Generation: Leveraging tools like gqlgen to automate boilerplate code. Go's strong typing makes this process especially seamless and reliable.

- Concurrency with Goroutines: Go's native support for lightweight threads (goroutines) allows the API to handle many concurrent user requests efficiently, a critical feature for a social media app.

- Go's Performance: Go is a compiled language that translates directly to machine code, providing excellent speed and low latency, which is essential for a fast and responsive API.

- Authentication Flow: Implementing secure JWT-based authentication in a GraphQL context. You'll learn how to validate tokens and secure specific API operations (mutations) that require a logged-in user.

- Database Interactions: Connecting to and performing CRUD (Create, Read, Update, Delete) operations on a PostgreSQL database using Go's database/sql package and external drivers