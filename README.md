# Brute-Force Protected Login Application

A full-stack login application with advanced brute-force protection mechanisms, featuring both user-level suspension and IP-level blocking.

## 🎯 Features

- **User-Level Suspension**: After 5 failed login attempts within 5 minutes, the account is suspended for 15 minutes
- **IP-Level Block**: After 100 failed login attempts from any IP within 5 minutes, all login attempts from that IP are blocked
- **Real-time Feedback**: Immediate user feedback for blocked or suspended attempts
- **Persistent Storage**: All data persists across server restarts using PostgreSQL
- **Modern UI**: Beautiful, responsive React frontend
- **RESTful API**: Clean Golang backend with Gin framework

### Tech Stack

- **Frontend**: React 18 with Vite
- **Backend**: Golang with Gin framework
- **Database**: PostgreSQL 15
- **Deployment**: Docker & Docker Compose

### Project Structure

```
.
├── backend/
│   ├── config/          # Configuration management
│   ├── database/        # Database schema and migrations
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data models
│   ├── services/        # Business logic (auth, lockout)
│   ├── main.go          # Application entry point
│   └── Dockerfile       # Backend container config
├── frontend/
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── App.jsx      # Main app component
│   │   └── main.jsx     # React entry point
│   └── Dockerfile       # Frontend container config
├── docker-compose.yml   # Multi-container orchestration
└── README.md            # This file
```

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose installed
- OR Go 1.21+ and Node.js 18+ for local development

### Option 1: Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd brute-force-login
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - Frontend: http://<Server IP>:3000
   - Backend API: http://<Server IP>:8080
   - Database: <IP>:5432

4. **Stop services**
   ```bash
   docker-compose down
   ```

### Option 2: Local Development

#### Backend Setup

1. **Install PostgreSQL** and create a database:
   ```sql
   CREATE DATABASE brute_force_login;
   ```

2. **Run database migrations**:
   ```bash
   cd backend
   psql -U postgres -d brute_force_login -f database/schema.sql
   psql -U postgres -d brute_force_login -f database/init.sql
   ```

3. **Set environment variables** (optional, defaults provided):
   ```bash
   export DB_HOST=<IP>
   export HOST=<Server IP>
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=brute_force_login
   export JWT_SECRET=your-secret-key
   export PORT=8080
   ```

4. **Run the backend**:
   ```bash
   cd backend
   go mod download
   go run main.go
   ```

#### Frontend Setup

1. **Install dependencies**:
   ```bash
   cd frontend
   npm install
   ```

2. **Set API URL** (optional, defaults to <Server IP>:8080):
   ```bash
   export VITE_API_URL=http://<Server IP>:8080/api
   ```

3. **Run the frontend**:
   ```bash
   npm run dev
   ```

4. **Access the application**: http://<Server IP>:3000

## 🧪 Testing

### Run Unit Tests

```bash
cd backend
go test ./services/... -v
```

The tests verify:
- User suspension logic (5 failed attempts → suspension)
- IP blocking logic (100 failed attempts → IP block)

### Manual Testing

**Test User Suspension:**
1. Try logging in with wrong password 5 times
2. 6th attempt should show: "Account temporarily suspended due to too many failed attempts."

**Test IP Block:**
1. Make 100 failed login attempts from the same IP (different users)
2. 101st attempt should show: "IP temporarily blocked due to excessive failed login attempts."

## 📊 Database Schema

### Tables

- **users**: Stores user credentials (email, password hash)
- **user_failed_attempts**: Tracks failed login attempts per user
- **ip_failed_attempts**: Tracks failed login attempts per IP
- **user_suspensions**: Stores active user suspensions
- **ip_blocks**: Stores active IP blocks

### Indexes

Optimized indexes on email, IP address, and timestamps for fast lookups.

## 🔧 Configuration

### Environment Variables

**Backend:**
- `DB_HOST`: PostgreSQL host (default: <IP>)
- `HOST`: Server host (default: <Server IP>)
- `DB_PORT`: PostgreSQL port (default: 5432)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password (default: postgres)
- `DB_NAME`: Database name (default: brute_force_login)
- `DB_SSLMODE`: SSL mode (default: disable)
- `JWT_SECRET`: Secret key for JWT tokens
- `PORT`: Server port (default: 8080)

**Frontend:**
- `VITE_API_URL`: Backend API URL (default: http://<Server IP>:8080/api)

### Lockout Parameters

Currently hardcoded in `backend/config/config.go`:
- User attempt limit: 5
- User window: 5 minutes
- User suspension: 15 minutes
- IP attempt limit: 100
- IP window: 5 minutes

## 🌐 Deployment
   Deployed In AWS EC2 Instance.

### Deploy to Production

1. **Update environment variables** for production
2. **Change JWT_SECRET** to a secure random string
3. **Enable SSL** for database connections
4. **Configure CORS** origins in `backend/main.go`
5. **Build and deploy**:
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### Platform-Specific Deployment

**Vercel (Frontend):**
- Connect GitHub repository
- Set build command: `npm run build`
- Set output directory: `dist`

**Render/Railway (Backend):**
- Connect GitHub repository
- Set build command: `go build -o main`
- Set start command: `./main`
- Configure PostgreSQL addon
- Set environment variables

## 📝 API Endpoints

### POST /api/login

Login endpoint with brute-force protection.

**Request:**
```json
{
  "email": "alice@example.com",
  "password": "password123"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Login successful.",
  "token": "jwt-token-here"
}
```

**Error Responses (401):**
```json
{
  "success": false,
  "message": "Invalid email or password."
}
```

```json
{
  "success": false,
  "message": "Account temporarily suspended due to too many failed attempts."
}
```

```json
{
  "success": false,
  "message": "IP temporarily blocked due to excessive failed login attempts."
}
```

### GET /api/health

Health check endpoint.

**Response (200):**
```json
{
  "status": "ok"
}
```

## 🔐 Security Features

1. **Password Hashing**: Bcrypt with cost factor 10
2. **JWT Tokens**: Secure token-based authentication
3. **IP Detection**: Supports X-Forwarded-For and X-Real-IP headers
4. **Time-based Windows**: Sliding window for attempt tracking
5. **Database Indexing**: Optimized queries for performance

## 📋 Sample Credentials

For testing purposes, the following users are pre-configured:

- **Email**: alice@example.com
- **Password**: password123

- **Email**: bob@example.com
- **Password**: password123

## 🤖 AI Usage Report

This project was built with AI assistance:

- **Code Generation**: ~40% generated via AI
- **Time Spent**: ~6 hours fine-tuning and debugging
- **AI-Assisted Parts**:
    - Architecture Design
    - Docker configuration
     - Deployment configuration adjustments
     - README documentation
- **Manual Work**:
    - Initial project structure and boilerplate
  - Testing and debugging lockout logic
  - UI/UX refinements
  - Error handling improvements
  - Database schema design
  - Authentication service logic
  - React component structure






