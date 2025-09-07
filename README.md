# Chirpy 🐦

A simple backend API for posting and managing short messages (chirps).

## Features

* User registration & authentication (JWT-based).
* Post, fetch, and delete chirps.
* Filter chirps by user.
* Secure password hashing & refresh tokens.
* PostgreSQL database integration.
* JSON-based REST API.

## Tech Stack

* **Go** (Golang) – backend
* **PostgreSQL** – database
## Getting Started

### Prerequisites

* Go 1.22+
* PostgreSQL 14+
* Docker (optional, for local database)

### Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/austinthieu/chirpy.git
   cd chirpy
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Set up your environment variables (create a `.env` file):

   ```env
   DATABASE_URL="postgres://user:password@localhost:5432/chirpy"
   PLATFORM="dev"
   SECRET="supersecretkey"
   POLKA_KEY="secretapikey"
   ```

4. Run database migrations:

   ```bash
   goose postgres postgres://user:password@localhost:5432/chirpy up
   ```

5. Start the server:

   ```bash
   go run .
   ```

The API will be available at:

```
http://localhost:8080
```

## API Endpoints

### Users

* `POST /api/users` – Register a new user
* `PUT /api/users` – Update a user's red status

### Chirps

* `GET /api/chirps` – Get all chirps
* `GET /api/chirps?author_id={id}` – Get chirps from a specific user
* `POST /api/chirps` – Create a new chirp
* `DELETE /api/chirps/{id}` – Delete a chirp (if owned)

### Auth
* `GET /api/login` – Login and get a refresh token and access token
* `GET /api/refresh` – Refresh to get a new access token
* `POST /api/revoke` – Revoke a user's refresh token

## Development

* Run tests:

  ```bash
  go test ./...
  ```
