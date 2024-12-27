# Blogging Platform API App
A Go-based API application for managing blog posts, built to provide a robust and scalable backend for a blogging platform. The application integrates with MySQL for data persistence and offers essential blogging features, making it a solid foundation for modern web applications.

## Features
- **Blog Post Management**: Supports creating, retrieving, updating, and deleting blog posts through RESTful endpoints.
- **MySQL Integration**: Ensures reliable and efficient data storage with full support for relational database operations.
- **Scalable Deployment**: Easily scalable and containerized for seamless deployment using Docker.
- **Lightweight Server Setup**: Built with a straightforward and maintainable server architecture using Go and the Gorilla Mux router.
- **Port Configuration**: Runs on port 8081 by default, configurable for flexibility.
- **Future Extensibility**: Designed with a modular structure to support additional features like user authentication, categories, and tags.

## Prerequisites
Before running the application, ensure you have the following installed on your machine:
- **Go**: Version 1.22.3 (darwin/arm64)
- **Docker**: Make sure Docker is installed and running.
- **Docker Compose**: This is included with Docker Desktop.
- **cURL**: Ensure cURL is installed for testing API requests.

## Configuration
- Make sure docker or colima is running
- Store the MySQL credentials in a .env file and reference them in your docker-compose.yml
- Example .env File:  
  ```
    MYSQL_ROOT_PASSWORD=securepassword123
    MYSQL_USER=blog_user
    MYSQL_PASSWORD=strongpassword456
    MYSQL_DATABASE=blog_db
  ```

## Installation
To install and run the app locally, clone the repository and build the Go binary, following the steps below:
```
git clone https://github.com/lugomas/blogging-platform-api.git
cd blogging-platform-api
go build -o blogging-platform-api
```

## Running the Application
1. Start the app:  
   ```docker-compose up --build```

## Usage
1. Access the app:
   Open your web browser and navigate to:  
   ```
    http://localhost:8081/posts
   ```  
   ```
    http://localhost:8081/posts/{postId}
   ```    

2. Open a new terminal and test it via cURL:  

   GET post by ID:  
    ```
    curl -v http://localhost:8081/posts/{postId}
   ```    

   Create a post:  
    ```
    curl -X POST http://localhost:8081/posts \                                    
    -H "Content-Type: application/json" \
    -d '{"title":"My First Blog Post","content":"Here is my first post","category":"Programming","tags":["Golang","API","Database","Docker"]}'
    ```
   
   Update a post:  
    ```
    curl -X PUT http://localhost:8081/posts/{postID} \
    -H "Content-Type: application/json" \
    -d '{"title":"Updated Blog Post","content":"This is the updated content for the first post.","category":"Astronomy","tags":["Cosmology","Astrophysics"]}'
    ```
   
    Delete a post:  
    ```
    curl -X DELETE http://localhost:8081/posts/{postID}
   ```  

   Search for a post containing a string that matches the specified term:  
   ```
   curl -X GET "http://localhost:8081/posts?term=Working"  
   ```
3. Stop and remove all containers:
   ```
    docker-compose down
   ```

## License
This project is licensed under the MIT License.

## Project Inspiration
This project was developed based on the guidelines provided by [roadmap.sh's Blogging Platform API project](https://roadmap.sh/projects/blogging-platform-api)

## Backlog
- user can also filter posts by a search term
- Include MKDOCS to this project
  - Explain what i've learned from this project