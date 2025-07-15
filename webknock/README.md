# Library Management System Documentation

## Overview
This is a single-page HTTPS application built with a Golang backend and a Knockout.js frontend styled with Bootstrap. It implements OAuth2 authentication using JWT and a self-signed certificate for HTTPS. The application manages a library system with user authentication, book CRUD operations, member management, and book borrowing/returning functionality.

## Generating Self-Signed Certificates
To enable HTTPS, a self-signed certificate is used. Follow these steps to generate the certificate and key files:

1. **Install OpenSSL** (if not already installed):
   - On Ubuntu: `sudo apt-get install openssl`
   - On macOS: `brew install openssl`
   - On Windows: Download OpenSSL from a trusted source or use WSL.

2. **Generate the Certificate and Key**:
   Run the following command to create a self-signed certificate (`cert.pem`) and private key (`key.pem`):
   ```bash
   openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
   ```
   - `-x509`: Creates a self-signed certificate.
   - `-newkey rsa:4096`: Generates a 4096-bit RSA key.
   - `-keyout key.pem`: Outputs the private key to `key.pem`.
   - `-out cert.pem`: Outputs the certificate to `cert.pem`.
   - `-days 365`: Sets the certificate validity to 365 days.
   - `-nodes`: Skips passphrase protection for the private key.

3. **Place Files**:
   Place `cert.pem` and `key.pem` in the root directory of the project.

4. **Trust the Certificate** (optional, for development):
   - On browsers, you may need to manually trust the certificate by visiting `https://localhost:8080` and accepting the security warning.
   - For production, use a certificate from a trusted Certificate Authority (CA).

## Project Structure
- **main.go**: Golang backend with RESTful API endpoints, SQLite database, and JWT authentication.
- **static/index.html**: Frontend HTML with Bootstrap and Knockout.js bindings.
- **static/app.js**: Knockout.js view model for handling frontend logic.
- **library.db**: SQLite database (created automatically on first run).
- **cert.pem**, **key.pem**: Self-signed certificate and key for HTTPS.

## Web Services
All API endpoints are prefixed with `/api/` and require authentication (except `/api/login` and `/api/register`). Authentication is handled via JWT in the `Authorization` header as `Bearer <token>`.

### Authentication Endpoints
- **POST /api/register**
  - **Request Body**: `{ "username": string, "password": string }`
  - **Response**: `201 Created` on success, `409 Conflict` if username exists, `400 Bad Request` for invalid input.
  - **Description**: Registers a new user with a hashed password.

- **POST /api/login**
  - **Request Body**: `{ "username": string, "password": string }`
  - **Response**: `{ "token": string }` on success, `401 Unauthorized` for invalid credentials.
  - **Description**: Authenticates a user and returns a JWT token.

### Book Endpoints
- **GET /api/books**
  - **Query Parameters**: `status` (optional, values: "available", "borrowed")
  - **Response**: Array of books `[{ "id": int, "title": string, "author": string, "status": string, "borrower_id": int|null }]`
  - **Description**: Retrieves all books or filtered by status.

- **POST /api/books**
  - **Request Body**: `{ "title": string, "author": string }`
  - **Response**: `201 Created` on success, `400 Bad Request` for invalid input.
  - **Description**: Creates a new book with status "available".

- **PUT /api/books/{id}**
  - **Request Body**: `{ "title": string, "author": string, "status": string }`
  - **Response**: `200 OK` on success, `400 Bad Request` for invalid input.
  - **Description**: Updates a book by ID.

- **DELETE /api/books/{id}**
  - **Response**: `200 OK` on success, `500 Internal Server Error` on failure.
  - **Description**: Deletes a book by ID.

### Member Endpoints
- **GET /api/members**
  - **Response**: Array of members `[{ "id": int, "name": string }]`
  - **Description**: Retrieves all members.

- **POST /api/members**
  - **Request Body**: `{ "name": string }`
  - **Response**: `201 Created` on success, `400 Bad Request` for invalid input.
  - **Description**: Creates a new member.

- **PUT /api/members/{id}**
  - **Request Body**: `{ "name": string }`
  - **Response**: `200 OK` on success, `400 Bad Request` for invalid input.
  - **Description**: Updates a member by ID.

- **DELETE /api/members/{id}**
  - **Response**: `200 OK` on success, `500 Internal Server Error` on failure.
  - **Description**: Deletes a member by ID.

### Borrow/Return Endpoints
- **POST /api/borrow**
  - **Request Body**: `{ "book_id": int, "member_id": int }`
  - **Response**: `200 OK` on success, `400 Bad Request` for invalid input.
  - **Description**: Records a book borrow, updates book status to "borrowed" and sets borrower_id.

- **POST /api/return**
  - **Request Body**: `{ "book_id": int, "member_id": int }`
  - **Response**: `200 OK` on success, `400 Bad Request` for invalid input.
  - **Description**: Records a book return, updates book status to "available" and clears borrower_id.

## Pages
### Login Screen
- **URL**: `/` (initial view)
- **Description**: Allows users to log in with a username and password. Provides a link to the registration screen.
- **Fields**: Username, Password
- **Actions**: Login (POST to `/api/login`), Switch to Register

### Registration Screen
- **Description**: Allows new users to register with a username and password.
- **Fields**: Username, Password
- **Actions**: Register (POST to `/api/register`), Back to Login

### Books Screen
- **Description**: Displays a list of books with CRUD operations and a status filter (All, Available, Borrowed).
- **Fields**: Title, Author (for add/edit)
- **Actions**:
  - View books (GET `/api/books` with optional `status` query)
  - Add book (POST `/api/books`)
  - Edit book (PUT `/api/books/{id}`)
  - Delete book (DELETE `/api/books/{id}`)
  - Switch to Members
  - Logout

### Members Screen
- **Description**: Displays a list of members with CRUD operations and options to borrow/return books.
- **Fields**: Name (for add/edit), Book and Member dropdowns (for borrow/return)
- **Actions**:
  - View members (GET `/api/members`)
  - Add member (POST `/api/members`)
  - Edit member (PUT `/api/members/{id}`)
  - Delete member (DELETE `/api/members/{id}`)
  - Borrow book (POST `/api/borrow`)
  - Return book (POST `/api/return`)
  - Switch to Books
  - Logout

## Setup Instructions
1. **Install Go**:
   - Download and install Go from https://golang.org/dl/.
2. **Install Dependencies**:
   ```bash
   go get github.com/gorilla/mux
   go get github.com/mattn/go-sqlite3
   go get github.com/dgrijalva/jwt-go
   go get golang.org/x/crypto/bcrypt
   ```
3. **Generate Certificates**:
   Follow the certificate generation steps above.
4. **Run the Application**:
   ```bash
   go run main.go
   ```
5. **Access the Application**:
   Open `https://localhost:8080` in a browser and accept the self-signed certificate warning.

## Security Notes
- The JWT secret (`your_secret_key`) should be replaced with a secure, randomly generated key in production.
- The self-signed certificate is for development only. Use a trusted CA for production.
- Passwords are hashed using bcrypt for security.
- SQLite is used for simplicity; consider a more robust database for production.

## Limitations
- The application uses a simple SQLite database, which may not scale for large datasets.
- Error handling is basic; enhance with more detailed messages for production.
- The self-signed certificate may cause browser warnings; use a trusted certificate in production.

## Notes

- **Volume:** The library.db file is mounted as a volume to persist the SQLite database on the host machine.
- **Port Mapping:** The application runs on port 8080 inside the container, mapped to the host’s port 8080.
- **Environment:** GIN_MODE=release is set for better performance, though not strictly required for this app.
- **Security:** The self-signed certificate is included in the Docker image. For production, replace with a trusted CA certificate.
- **Stopping the Container:** Use Ctrl+C to stop, or run docker-compose down to remove the container.This setup ensures the application runs in a Docker container with all dependencies and configurations intact, including HTTPS support.