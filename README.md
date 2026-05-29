# 🚀 Go Todo API

> Production-grade Productivity + Personal Finance Backend API built with Go, Gin, PostgreSQL, JWT Authentication, Session Management, Background Schedulers, Swagger Documentation, and Integration Testing.

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)
![Gin](https://img.shields.io/badge/Gin-Web_Framework-00ADD8)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-336791?logo=postgresql)
![JWT](https://img.shields.io/badge/JWT-Authentication-orange)
![Swagger](https://img.shields.io/badge/Swagger-OpenAPI-green)
![License](https://img.shields.io/badge/License-MIT-blue)

---

# 📖 Overview

Go Todo API is a backend platform designed to manage:

- Productivity workflows
- Personal finance management
- Budget tracking
- Recurring transactions
- Financial analytics
- User authentication & session management

The system follows a **clean layered architecture** and is built as a **modular monolith** to maximize maintainability, development speed, and future scalability.

---

# ✨ Features

## 🔐 Authentication

- User Registration
- Login
- JWT Access Token
- Refresh Token Rotation
- Session-Based Authentication
- Multi-Device Sessions
- Session Revocation
- Logout
- Current User Endpoint

---

## 📝 Productivity Module

### Todos

- Create Todo
- Update Todo
- Delete Todo
- Complete Todo
- Todo Dashboard

### Notes

- Create Note
- Update Note
- Delete Note
- Pin Note
- Search Notes

---

## 💰 Finance Module

### Categories

- Income Categories
- Expense Categories

### Transactions

- Create Transaction
- Update Transaction
- Delete Transaction
- Filter Transactions
- Pagination
- Monthly Summary

### Budgets

- Monthly Budget
- Category Budget
- Budget Tracking
- Remaining Budget Calculation

### Recurring Transactions

- Monthly Transactions
- Automated Transaction Creation
- Background Scheduler

### Alerts

- Budget Alerts
- Reminder Notifications

### Analytics

- Income Summary
- Expense Summary
- Financial Dashboard

---

# 🛡 Production Features

- JWT Authentication
- Refresh Token Rotation
- Session Management
- Exact Session Logout
- Swagger Documentation
- Structured Logging
- Request Tracing
- Security Headers
- Panic Recovery
- Graceful Shutdown
- Health Checks
- Readiness Checks
- Integration Testing Foundation
- Scheduler Architecture

---

# 🏗 Architecture

```text
Mobile App (React Native / Flutter)
           │
           ▼
      REST API
       (Gin)
           │
           ▼
    Middleware Layer
           │
           ▼
      Handler Layer
           │
           ▼
      Service Layer
           │
           ▼
    Repository Layer
           │
           ▼
      PostgreSQL
```

Supporting Components:

```text
JWT Authentication
Session Management
Schedulers
Swagger Documentation
Structured Logging
Integration Tests
```

---

# 📂 Project Structure

```text
go-todo-api/
│
├── app/
│   └── router.go
│
├── config/
│   ├── config.go
│   └── database.go
│
├── docs/
│   └── swagger docs
│
├── dto/
│   └── request/response DTOs
│
├── handler/
│   └── HTTP handlers
│
├── middleware/
│   └── auth, logger, recovery, cors
│
├── model/
│   └── database models
│
├── repository/
│   └── database access layer
│
├── response/
│   └── response helpers
│
├── scheduler/
│   └── background jobs
│
├── service/
│   └── business logic
│
├── tests/
│   └── integration tests
│
├── utils/
│   └── helpers
│
├── main.go
│
└── README.md
```

---

# 🧠 Architecture Rules

## Handler Layer

Responsible for:

- Request Binding
- Validation
- Service Calls
- API Responses

Never:

- Query Database Directly
- Implement Business Logic

---

## Service Layer

Responsible for:

- Business Rules
- Domain Logic
- Workflow Orchestration

Never:

- Return HTTP Responses
- Access Gin Context

---

## Repository Layer

Responsible for:

- Database Queries
- Persistence
- Filtering
- Pagination

Never:

- Implement Business Logic

---

# ⚙️ Installation

## Clone Repository

```bash
git clone https://github.com/your-username/go-todo-api.git

cd go-todo-api
```

---

## Install Dependencies

```bash
go mod tidy
```

---

# 🔧 Environment Variables

Create:

```bash
.env
```

Example:

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=go_todo_api

TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USER=postgres
TEST_DB_PASSWORD=password
TEST_DB_NAME=go_todo_test

JWT_SECRET=super-secret-key

JWT_ACCESS_EXP=15m
JWT_REFRESH_EXP=168h
```

---

# 🗄 Database

Create databases:

```sql
CREATE DATABASE go_todo_api;
CREATE DATABASE go_todo_test;
```

---

# ▶ Running Application

```bash
go run main.go
```

Server:

```text
http://localhost:8080
```

---

# 📘 Swagger Documentation

Generate docs:

```bash
swag init
```

Swagger UI:

```text
http://localhost:8080/api/docs/index.html
```

---

# 🔐 Authentication Flow

## Login

```text
POST /auth/login
```

Returns:

```json
{
  "access_token": "...",
  "refresh_token": "..."
}
```

---

## Protected Endpoint

```http
Authorization: Bearer ACCESS_TOKEN
```

---

## Refresh Token

```text
POST /auth/refresh
```

Flow:

```text
Access Token Expired
        ↓
Refresh Endpoint
        ↓
New Access Token
        ↓
Retry Request
```

---

## Logout

```text
POST /auth/logout
```

Behavior:

```text
Current Session Revoked
Other Sessions Remain Active
```

---

# 📱 Mobile Integration Guide

## Recommended Storage

Access Token:

```text
Memory
```

Refresh Token:

```text
Secure Storage
```

Examples:

- react-native-keychain
- Expo SecureStore
- Flutter Secure Storage

---

## Auto Refresh Strategy

Use interceptor:

```text
401 Response
      ↓
Refresh Token
      ↓
Store New Token
      ↓
Retry Original Request
```

---

## Logout

Send:

```json
{
  "refresh_token": "..."
}
```

Backend revokes exact session.

---

# 🧪 Testing

Run integration tests:

```bash
go test ./tests/... -v
```

Current Coverage:

- Register
- Login

Planned Coverage:

- Refresh
- Logout
- Me Endpoint
- Transactions
- Budgets
- Analytics

---

# ❤️ Health Monitoring

## Health Check

```text
GET /health
```

Purpose:

```text
Application is alive
```

---

## Readiness Check

```text
GET /ready
```

Purpose:

```text
Application is ready to serve requests
```

---

# ⏰ Scheduler Architecture

Current schedulers:

- Recurring Scheduler
- Alert Scheduler
- Session Cleanup Scheduler

Pattern:

```text
scheduler/
    feature_scheduler.go
```

Required:

```go
type Scheduler struct{}

func NewScheduler() *Scheduler

func (s *Scheduler) Start()
```

---

# 🚀 Deployment

Recommended:

- Docker
- Nginx
- PostgreSQL
- SSL
- Reverse Proxy

Production Settings:

```env
APP_ENV=production
```

---

# 📚 Engineering Documentation

See:

```text
TECHNICAL_HANDBOOK.md
```

Contains:

- Architecture
- Coding Standards
- Mobile Integration
- Deployment Guides
- Future Roadmap

---

# 🛣 Roadmap

## Phase 2

- Wallets
- Transfers
- Savings Goals
- Debt Tracking
- Loan Tracking

## Phase 3

- Exports (PDF / Excel)
- Push Notifications
- Admin Dashboard
- RBAC

## Phase 4

- Redis Caching
- Queue Workers
- Event-Driven Architecture

---

# 🤝 Contributing

Before contributing:

- Follow architecture guidelines
- Follow service/repository separation
- Write Swagger annotations
- Add tests for new features
- Keep handlers thin

---

# 📄 License

MIT License

---

# 🔥 Current Status

```text
Sprint 1  ✅ Auth Foundation
Sprint 2  ✅ Productivity Core
Sprint 3  ✅ Dashboard
Sprint 4  ✅ Profile
Sprint 5  ✅ Hardening
Sprint 6  ✅ Finance Foundation
Sprint 7  ✅ Finance Expansion
Sprint 8  ✅ Architecture Stabilization
Sprint 9  ✅ Production Readiness
```

### Project Maturity

```text
Production Launch Candidate 🚀
```