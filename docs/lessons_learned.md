# Lessons Learned

## Project Overview
- **Purpose:** Develop a blogging platform API with functionalities such as creating, reading, updating, and deleting posts.
- **Technologies Used:** Golang, Docker, PostgreSQL, REST API principles, logging and web development frameworks.

## Lessons Learned
### API Development with Golang
- **Routing:** Learned how make it easier the web development using gorilla/mux to define routes and handle HTTP requests.
- **Restrict Operations:** Define which HTTP Methods each endpoint supports.  
  Example:  
  ```r.HandleFunc("/posts", handlers.HandlePosts).Methods("GET", "POST")```

- **Handling HTTP Requests:** I have understood how to handle requests when the URL contains a search term:  
  Example:  
  For requests such as http://localhost:8081/posts?term=Working, I can handle the term like this:  
  `searchTerm := r.URL.Query().Get("term")`

- **Handlers:** Built RESTFull endpoints adhering to clean architecture principles.
  I opted for the REST protocol because of its simplicity, performance, and wide compatibility, especially for projects that don't require the flexibility or complexity of GraphQL.

- **Interacting with the API:** I used CURL to send HTTP requests as it supports various protocols like HTTP, HTTPS, FTP, and more, making it ideal for testing and debugging.

- **Data management:** I have developed familiarity the following:
    - `Decode` the request body (json format) into a Post struct using a pointer. Example: `json.NewDecoder(r.Body).Decode(&post)`
    - Convert the Tags slice into a JSON string by `Marshaling` data. Example: `tagsJSON, err := json.Marshal(post.Tags)`
    - `Parse` or `Unmarshal` the tags JSON string into the Tags slice. Example: `err = json.Unmarshal([]byte(tagsJSON), &post.Tags)`
    - Generate unique, non-predictable `UUIDs` for each post using the github.com/google/uuid library in Go, which implements a good practice. With this instead of using Timestamp (`postID := time.Now().Format("200601021504105"`), we avoid collisions in high throughput scenarios and predictability
    - Extract `query parameter` from the `URL`, which appear after the ? symbol (e.g., http://localhost:8081/posts?term=Working). Example: `searchTerm := r.URL.Query().Get("term")`
    - Extracts `variables` from path segments in the `URL` that match placeholders defined in your gorilla/mux route patterns (e.g., /users/{id}). Example:   
      `vars := mux.Vars(r)`  
      `potID := vars["id"]`

- **Error Handling:** Enhanced skills in error propagation, using structured error messages to provide meaningful feedback.
    - With log/slog standard library package I have implemented structured logging by creating records containing the time of the call, level of Info, and the message. Example: `slog.Info("GetPost - Fetching post with ID: ", "id", id)` 
    - I have understood that i should not only log errors for the API Client but also for the server when handling the same error. Here's an example of that:  
      `slog.Error("post not found", "id", id)`  
      `http.Error(w, "Post not found", http.StatusNotFound)`  

### Database Integration with PostgresSQL
- **Local Database:** I had to install mysql service and start/stop it via brew services command. Also, always making sure it's running before start running the application.
- **MySQL Integration:** Implement reliable and efficient data storage with full support for relational database operations.
- **SQL Package:** Establish connections via database/sql package
- **MySQL driver:** Including the package (_ "github.com/go-sql-driver/mysql") is essential for enabling SQL database operations.
- **Database and Table creation:** Executing a query without returning any rows and performing database management.
- **Scanning tables:** Retrieve Data by scanning the database and parsing it's values into custom Golang struct types.  
- **Debugging:** Identify issues by running SQL commands directly in the database.
- **Access management:** Creation of users and password management.

### Containerization with Docker and Deploy
- **Colima:** A lightweight alternative to Docker Desktop. Want to avoid Docker Desktop's licensing requirements for commercial use.
- **Dockerized APP**: With the application containerized, we ensure it runs consistently across environments
    - **Dockerfile**: Defined it in `multi-stage` format for separation of concerns. I defined first the stage to Build the app and second one to Run the application, resulting in a smaller, more efficient, and secure final Docker image. 
    - **Credentials**: Read environment variables using `os` package. Example: `os.Getenv("MYSQL_DATABASE")`
- **Dockerized Database**: With the database containerized, we ensure it runs consistently across environments
    - **Docker compose**: Containerize Mysql defining a persistent storage for database data
    - **Credentials**: Reference .env file with credentials (local file).
- **Multi-container Setup:** Use Docker-Compose to define, manage and run both, app and mysql containers using a docker-compose.yml file.
- **Commands**:
    - **Start services (build if necessary)**:  
    `docker-compose up`

    - **Start services in detached mode**:  
    `docker-compose up -d`

    - **Stop services**:  
    `docker-compose down`
    
    - **List running services**:  
    `docker-compose ps`
    
    - **Rebuild services after changes**:  
    `docker-compose up --build`

## Challenges Faced
   Debugging errors in SQL queries.  
   Managing dependencies in Go (using tools like go mod).  
   Ensuring seamless communication between services in a containerized environment.  
   Designing APIs with scalability and maintainability in mind.  

## Achievements
   Built a fully functional blogging platform API.  
   Containerized the application for cross-environment compatibility.  
   Gained hands-on experience with Docker, PostgreSQL, and Go.  
   Developed an understanding of RESTful API best practices and deployment strategies.  

## Key Takeaways
   The importance of clean architecture for scalability.  
   Benefits of containerization in managing complex applications.  
   How to write efficient, testable, and maintainable Go code.  
   Real-world problem-solving in API development and integration.  