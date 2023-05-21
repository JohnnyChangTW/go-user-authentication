# HTTP APIs for Account and Password Management

- This project provides two RESTful HTTP APIs for creating and verifying an account and password.

- Still have some issues related to database initialization when using docker-compose to spin up this application

## APIs document

### 1. Create Account

**Endpoint:** POST /accounts

**Request Payload:**
```json
{
  "username": "john_doe",
  "password": "Password123"
}
```
**Response Payload (Success):**
```json
{
  "success": true
}
```
**Response Payload (Failure):**
```json
{
  "success": false,
  "reason": "Username already exists"
}
```

### 2. Verify Account and Password

**Endpoint:** POST /accounts

**Request Payload:**
```json
{
  "username": "john_doe",
  "password": "Password123"
}
```
**Response Payload (Success):**
```json
{
  "success": true
}
```
**Response Payload (Failure):**
```json
{
  "success": false,
  "reason": "Invalid username or password"
}
```

## User Guide: Running the Application with Docker

### Prerequisites
- Docker

### Steps 
#### 1. Clone the GitHub repository:
```bash
git clone https://github.com/JohnnyChangTW/senao-coding-assessment.git
```
#### 2. Build the Docker image and start the containers:
```bash
docker-compose up -d
```
Please note that the MySQL database will be persisted in a Docker volume named "mysql_data".
#### 3. The application will be accessible at http://localhost:8000.
#### 4. Use the provided API endpoints to create and verify accounts.
#### 5. To stop the application and containers, run:
```bash
docker-compose down
```




