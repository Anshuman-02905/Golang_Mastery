BOOK MANAGEMENT PROJECT
Sprint 1: Project Setup & Database Configuration
Objective: Set up the foundational structure and ensure database connectivity.
Tasks:
1. Initialize the Project
    * Set up a new Go module (go mod init bookstore)
    * Create the project folder structure.
2. Set Up Configurations
    * Create a config package to manage environment variables.
    * Use godotenv to load .env file (store DB credentials, etc.).
    * Implement a function to read config values.
3. Set Up PostgreSQL with Docker
    * Write a docker-compose.yml file for PostgreSQL.
    * Ensure PostgreSQL is running inside a container.
4. Define Database Connection
    * Use GORM as the ORM (gorm.io/gorm, gorm.io/driver/postgres).
    * Implement a ConnectDB() function to establish a connection.
    * Handle errors and logging.
5. Define the Book Model
    * Create a models package.
    * Define a Book struct with fields (ID, Title, Author, Description, etc.).
    * Use GORM struct tags to map fields correctly.
6. Auto-Migrate the Database
    * Implement auto-migration logic in main.go or a separate migration file.

✅ Deliverables:
* Proper folder structure
* Dockerized PostgreSQL instance
* Database connection logic
* Book model defined
* Auto-migration implemented

-go mod init 
.
├── ReadMe.md
├── config
├── controllers
├── go.mod
├── main.go
├── models
├── routes
└── utils