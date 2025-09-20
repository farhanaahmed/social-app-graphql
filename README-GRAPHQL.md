# The GraphQL Schema Architecture

This project follows a specific architectural logic centered around the GraphQL schema as the single source of truth for the API.
The schema defines a tree-like structure of data and operations. All API capabilities are explicitly defined here.

## Schema Structure

The schema is organized into three main parts:

- GraphQL Schema
    - RootQuery (Read-only operations)
        - Query Field (e.g., "posts")
            - Resolve function (fetches data)
            - Return Type (e.g., `[Post!]`)
    - RootMutation (Write operations)
        - Mutation Field (e.g., "createPost")
            - Resolve function (modifies data)
            - Return Type (e.g., `Post`)
    - Object Types (Reusable data structures)
        - Post Type
            - id: ID!
            - content: String!
            - user: User!
        - User Type
            - id: ID!
            - username: String!

---

##  Logic & Components

1.  **Object Types**: We define Go structs (e.g., `Post`, `User`) that mirror our database tables. These structs are then used to create corresponding **GraphQL object types** (`postType`, `userType`). This ensures strong typing from the database to the API.

2.  **`RootQuery`**: This object holds all the **read-only operations**. For every field inside `RootQuery` (e.g., the `posts` field), we write a `Resolve` function. This function contains the logic to fetch data from our PostgreSQL database and return it to the client.

3.  **`RootMutation`**: This object contains all the **write operations** that change the state of our data. Similar to queries, each mutation field (e.g., `createPost`) has a `Resolve` function that handles the logic for creating, updating, or deleting data in the database.

4.  **`Resolve` Function**: The **resolve function** is the core of our API's logic. It's a function that tells GraphQL how to fulfill a request for a specific field. It acts as the bridge between our API's schema and the underlying data source (the database).

By following this logical structure, our API becomes predictable, self-documenting, and easy to scale.