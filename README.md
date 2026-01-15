# GoJobs API 🚀

A high-performance, production-grade REST API for a Job Board platform built with **Go**.

This project demonstrates **backend engineering best practices**, moving beyond heavy frameworks to focus on core principles of system design, security, and scalability. It features a clean, layered architecture, robust security mechanisms, and production-ready operational tooling.

---

## 🏗️ Architecture & Philosophy

Built using **Standard Library Go** (`net/http`) to maximize performance and minimize hidden abstractions.

- **Separation of Concerns:** Clear layers for routing, handlers, business logic, and data access.
- **No ORMs:** Raw SQL via `pgx` for fine-grained control over queries, performance, and joins.
- **Type Safety:** Strong typing from the database layer through to JSON responses.

---

## 📂 Project Structure

This project follows the **Standard Go Project Layout**, designed for scalability and long-term maintainability.

```bash
.
├── cmd/
│   └── api/
│       ├── main.go          # Entry point: config, DB connection, server startup
│       ├── handlers.go      # HTTP handlers: parse requests, call models, write responses
│       ├── middleware.go    # Middleware: JWT auth, CORS, logging, panic recovery
│       └── helpers.go       # Utilities: JSON helpers, error handling, query parsing
├── internal/
│   ├── data/                # Data access layer & business logic
│   │   ├── models.go        # Model registry and shared logic
│   │   ├── users.go         # User logic (registration, password hashing)
│   │   ├── jobs.go          # Job logic (CRUD, search, pagination)
│   │   ├── applications.go  # Application logic (many-to-many relationships)
│   │   └── filters.go       # Filtering, sorting, pagination metadata
│   └── validator/           # Custom request validation logic
├── migrations/              # SQL migrations (version-controlled schema)
├── go.mod                   # Dependency definitions
└── README.md                # Project documentation
```

---

## ✨ Key Features

### 🔐 Security & Identity

- **Stateless Authentication:** JWT-based authentication with expiration.
- **RBAC (Role-Based Access Control):**

  - **Recruiters:** Post and manage jobs.
  - **Candidates:** Browse and apply for jobs.

- **Context-Aware Requests:** Authenticated user information injected into request context via middleware.
- **Password Security:** Secure password hashing using `bcrypt` with salt.

### ⚙️ Production Operations (DevOps)

- **Graceful Shutdown:** Handles `SIGTERM` / `SIGINT` to complete in-flight requests (zero-downtime friendly).
- **Structured Logging:** JSON logs via `log/slog`, compatible with tools like Splunk and Datadog.
- **Database Migrations:** Versioned schema management using `golang-migrate`.
- **CORS Policy:** Custom middleware for secure cross-origin resource sharing.
- **Resiliency:** Configured `ReadTimeout` and `WriteTimeout` to mitigate Slowloris-style attacks.

### 💾 Data & Business Logic

- **Relational Data Model:** Users ↔ Jobs ↔ Applications.
- **Data Integrity:** Database-level constraints (foreign keys, unique indexes) prevent invalid or duplicate data.
- **Advanced Querying:** Optimized SQL for full-text search, filtering, sorting, and offset-based pagination.

---

## 🛠️ Tech Stack

- **Language:** Go 1.22+
- **Database:** PostgreSQL (Neon.tech)
- **Driver:** `pgx`
- **Migrations:** `golang-migrate`
- **Routing:** Standard `net/http` (`ServeMux`)
- **Logging:** `log/slog` (structured logging)

---

## 🚀 Getting Started

### 1. Prerequisites

- Go 1.22 or higher
- PostgreSQL database (local or cloud)
- `golang-migrate` CLI installed

### 2. Installation

```bash
git clone https://github.com/karnop/gojobs.git
cd gojobs
```

### 3. Environment Configuration

Create a `.env` file in the project root:

```env
PORT=:8080
DB_DSN=postgres://user:pass@host:port/dbname?sslmode=require
JWT_SECRET=your_super_secret_key_change_me_in_prod
```

### 4. Database Setup (Migrations)

Initialize the database schema:

```bash
# On Windows (PowerShell)
$env:DB_DSN="your_connection_string_here"
migrate -path ./migrations -database $env:DB_DSN up

# On macOS / Linux
export DB_DSN="your_connection_string_here"
migrate -path ./migrations -database $DB_DSN up
```

### 5. Run the Server

```bash
go run ./cmd/api
```

You should see a structured JSON log indicating the server has started.

---

## 📡 API Endpoints

### Public Routes

| Method | Endpoint     | Description                                          |
| -----: | ------------ | ---------------------------------------------------- |
|    GET | /health      | Health check                                         |
|    GET | /jobs        | List jobs (supports `?page=1&title=go&sort=-salary`) |
|    GET | /jobs/{id}   | Get job details                                      |
|   POST | /users       | Register a new user                                  |
|   POST | /users/login | Login and receive a Bearer token                     |

### Protected Routes (Requires JWT)

| Method | Endpoint         | Role      | Description     |
| -----: | ---------------- | --------- | --------------- |
|   POST | /jobs            | Recruiter | Post a new job  |
|   POST | /jobs/{id}/apply | Candidate | Apply for a job |

---

## 🧪 Testing

Example: apply for a job using `curl`:

```bash
curl -X POST http://localhost:8080/jobs/1/apply \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

---

## 📜 License

MIT License
