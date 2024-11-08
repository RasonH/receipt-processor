
# Receipt Processor

## Introduction
This is a coding challenge @Fetch (done by Juisheng Hung), more requirements can be found in this repository: https://github.com/fetch-rewards/receipt-processor-challenge/ 

---
---
## Solution Summary

This solution is a backend service written in Go that provides endpoints for calculating points for receipts based on specific business rules. The service uses Gorilla Mux for routing and is structured to ensure modularity, maintainability and scalibility, with key functionality divided into separate packages (e.g., api, models, services, and storage).

---
---
## Key Features

### 1. API Endpoints:

The service exposes RESTful endpoints for submitting receipts (***/receipts/process***) and retrieving points by receipt ID (***/receipts/{id}/points***).

### 2. Points Calculation:
The service implements a set of rules to calculate points based on details in the receipt, such as the retailer’s name, purchase date, and item prices. These rules are encapsulated within helper functions, making them easily testable and extendable for future requirements.
> Points Calculation Rules:
> - One point for every alphanumeric character in the retailer name.
> - 50 points if the total is a round dollar amount with no cents.
> - 25 points if the total is a multiple of 0.25.
> - 5 points for every two items on the receipt.
> - If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
> - 6 points if the day in the purchase date is odd.
> - 10 points if the time of purchase is after 2:00pm and before 4:00pm.


### 3. Storage and Data Management:
The solution uses an in-memory storage mechanism to store receipt data and their corresponding points. This storage is managed through a singleton instance, ensuring efficient retrieval and management of data during the runtime.

### 4. Duplicate Receipt Prevention:
Each receipt is assigned a unique ID using a SHA-256 hash based on receipt content, ensuring consistent results and avoiding duplicate entries.

### 5. Error Handling and Validation
The solution includes basic validation checks for input data, returning appropriate HTTP status codes, and can show detailed error messages on the server side, making it easier to debug and troubleshoot issues.

### 6.  Unit Testing
The solution includes unit tests for key components, such as models, services, and handlers, to ensure the correctness of the implementation and facilitate future changes and refactoring.

---
---
## Areas for Further Discussion and Improvement

### 1. Edge Case Behavior Discussion:
Some edge cases need further clarification to ensure the system behaves as expected. This structure has left room for adaptation and scaling, but ongoing discussions on edge cases and requirements are necessary to ensure the service meets all business needs comprehensively.

### 2.	Persistence Layer:
Currently, the solution uses in-memory storage, which works for a lightweight setup but may require migration to a persistent database for a production environment.

### 3. Security Considerations
The service can introduce more middlewares such as input validation, rate limiting, and authentication to prevent abuse and unauthorized access.


---
---
## File Structure
```bash
.
├── Dockerfile
├── README.md
├── api
│   ├── handlers.go
│   ├── handlers_test.go
│   └── routes.go
├── go.mod
├── go.sum
├── main.go
├── models
│   ├── models.go
│   └── models_test.go
├── services
│   ├── points.go
│   ├── points_helpers.go
│   └── points_test.go
└── storage
    └── storage.go
```

---
---
## Building and Running
### 0. Clone the repository and navigate to the directory
```bash
$ git clone https://github.com/RasonH/receipt-processor.git
$ cd receipt-processor
```

## Approach 1: Using Go directly
### 1. Download dependencies
```bash
$ go mod download
```

### 2. Build the application
```bash
$ go build -o main
```

### 3. Run the application
```bash
$ ./main
```

## Approach 2: Using Docker
### 1. Build the Docker image
```bash
$ docker build -t receipt-processor .
```
### 2. Run the Docker container
```bash
$ docker run -p 8080:8080 receipt-processor
```

---
---
## API Documentation

### 1. Process Receipts
#### POST /receipts/process

- Function: Submits a receipt for processing.
- Request Body: JSON object representing the receipt.
- Response:
    - Status: 200 OK - Receipt processed successfully.
    - Status: 400 Bad Request - Invalid request body (receipt data).
    - Status: 409 Conflict - ID collision detected (with different receipt data).
    - Status: 500 Internal Server Error - Server error during processing.

### 2. Get Points by Receipt ID
#### GET /receipts/{id}/points

- Function: Retrieves the points calculated for a specific receipt.
- Response:
    - Status: 200 OK - Points retrieved successfully.
    - Status: 404 Not Found - Receipt ID not found.
    
---
---
## Sample Requests and Responses

### 1. Process Receipt
#### Request
```bash
$ curl -X POST -H "Content-Type: application/json" -d '{                                                     
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    { "shortDescription": "Gatorade", "price": "2.25" },
    { "shortDescription": "Gatorade", "price": "2.25" },
    { "shortDescription": "Gatorade", "price": "2.25" },
    { "shortDescription": "Gatorade", "price": "2.25" }
  ],
  "total": "9.00"
}' http://localhost:8080/receipts/process
```
#### Response
```json
{"id":"78cfa3241cc6e557b1531a3bdc19b69bb9e309cdad605d4af6e93fc1a9482e66"}
```

### 2. Get Points by Receipt ID
#### Request
```bash
$ curl http://localhost:8080/receipts/78cfa3241cc6e557b1531a3bdc19b69bb9e309cdad605d4af6e93fc1a9482e66/points
```
#### Response
```json
{"points":109}
```

---
---
## Running Unit Tests and Coverage
### 1. Running on Go
#### Run all tests
```bash
$ go test ./...
```

#### Run tests with coverage
```bash
$ go test -cover ./...
```

#### Run tests for a specific package
```bash
$ go test ./<package_name>
```

### 2. Running on Docker
#### (Prerequisite: make sure the Docker image has been built)
#### Run all tests 
```bash
$ docker run --rm receipt-processor go test ./...
```

#### Running tests with coverage
```bash
$ docker run --rm receipt-processor go test -cover ./...
```
#### Running tests for a specific package
```bash
$ docker run --rm receipt-processor go test ./<package_name>
```