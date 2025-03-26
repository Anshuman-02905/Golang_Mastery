# File-Based Database in Go

## Overview
This project implements a lightweight, file-based NoSQL-style database in Go using JSON files for storage. It ensures thread safety with **mutex locks** and provides CRUD (Create, Read, Update, Delete) operations on JSON data.

The approach is **similar to MongoDB**, as it stores structured data in JSON format and allows for easy retrieval and modification.

## Features
- **Data Handling**: Supports structured data storage in JSON format.
- **Thread Safety**: Uses `sync.Mutex` for safe concurrent access.
- **Atomic Writes**: Ensures data consistency with temporary file writes before renaming.
- **CRUD Operations**:
  - **Write**: Save JSON data into a specified collection.
  - **Read**: Retrieve a single or all records from a collection.
  - **Delete**: Remove single or all records.
- **Simple and Fast**: No external dependencies, purely based on file storage.

## Installation
To use this project, clone the repository and run the Go program:

```sh
git clone <repository_url>
cd <project_directory>
go run main.go
```

## Usage

### 1. Initializing the Database
```go
  dir := "./data"
  db, err := New(dir, nil)
  if err != nil {
      fmt.Println("Error initializing database:", err)
  }
```

### 2. Writing Data
```go
type User struct {
    Name    string
    Age     json.Number
    Contact string
    Company string
    Address Address
}

type Address struct {
    City    string
    State   string
    Country string
    Pincode json.Number
}

db.Write("users", "John", User{
    Name: "John", Age: "25", Contact: "1234567890",
    Company: "TechCorp", Address: Address{"New York", "NY", "USA", "10001"},
})
```

### 3. Reading Data
```go
var user User
db.Read("users", "John", &user)
fmt.Println(user)
```

### 4. Reading All Records
```go
records, err := db.ReadAll("users")
if err == nil {
    fmt.Println(records)
}
```

### 5. Deleting Data
```go
db.Delete("users", "John")
```

## Concurrency Handling
This implementation prevents race conditions by assigning a **mutex per collection**, ensuring that multiple goroutines do not write to the same file at the same time.

```go
func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
    d.mutex.Lock()
    defer d.mutex.Unlock()
    
    if m, ok := d.mutexes[collection]; ok {
        return m
    }
    d.mutexes[collection] = &sync.Mutex{}
    return d.mutexes[collection]
}
```

## Similarity to MongoDB
- Data is stored in **JSON documents**.
- Collections group similar documents together.
- CRUD operations are similar to **MongoDB’s insert, find, and delete** methods.

## License
This project is open-source and available under the MIT License.

## Author
[Your Name]

