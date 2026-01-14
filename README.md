# GoJobs API

A production-grade REST API for a Job Board platform, built with Go.
It features secure JWT authentication, Role-Based Access Control (RBAC), and robust data validation.

## 🚀 Tech Stack

- **Language:** Go 1.22+
- **Database:** PostgreSQL (via Neon.tech)
- **Router:** `net/http` (Standard Library)
- **Authentication:** JWT (JSON Web Tokens)
- **Security:** bcrypt password hashing
- **Architecture:** Layered (Handlers → Models → Database)

## ✨ Features

### Public API

- `GET /jobs` — List all active job postings
- `GET /jobs/{id}` — View detailed job description
- `POST /users` — Register a new account (Candidate / Recruiter)
- `POST /users/login` — Authenticate and receive a Bearer token

### Protected API (Recruiters Only)

- `POST /jobs` — Create a new job posting (Requires JWT and `recruiter` role)

## 🛠️ Installation & Setup

### Prerequisites

- Go 1.22 or higher
- PostgreSQL database (local or cloud)

### 1. Clone the Repository

```bash
git clone https://github.com/karnop/gojobs.git
cd gojobs
```

### 2. Configure Environment Variables

Create a `.env` file in the root directory:

```env
# Database Connection String (PostgreSQL)
DB_DSN=postgres://user:pass@host:port/dbname?sslmode=require

# Server Port
PORT=:8080

# JWT Secret Key (keep this secure)
JWT_SECRET=your_super_secret_random_string
```

### 3. Run the Server

```bash
go run ./cmd/api
```

## 🧪 Testing the API

### 1. Register a User

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Manav","email":"manav@example.com","password":"password123"}'
```

### 2. Login (Get Token)

```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"manav@example.com","password":"password123"}'
```

### 3. Create a Job (As Recruiter)

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Go Dev","company":"Tech","description":"Remote","salary":100000}'
```

## 📂 Project Structure

```text
cmd/api/            # Application entry point and HTTP handlers
internal/
  ├── data/         # Database models and DAO layer
  └── validator/   # Request validation logic
```

---
