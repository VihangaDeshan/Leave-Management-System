# Leave Management System - ABC Company

A modern, full-stack Leave Management System built with **Golang (Backend)**, **PostgreSQL (Database)**, and **React (Frontend)** with **Tailwind CSS**. This system enables employees to manage their leave requests digitally while allowing managers/admins to review and approve them efficiently with a professional, responsive user interface.

---

## 🎯 Project Overview

ABC Company currently manages employee leave requests manually and needed a digital solution. This system provides:

- **Employee Portal**: Submit leave requests, view history, and track leave balances
- **Admin/Manager Portal**: Review, approve/reject leave requests, view company-wide dashboard statistics
- **Professional UI**: Modern, clean design that works seamlessly across desktop, tablet, and mobile devices
- **Secure Authentication**: JWT-based authentication with role-based access control (RBAC)
- **Real-time Dashboard**: View pending requests, employees on leave, and leave statistics

---

## 🛠️ Technology Stack

### Backend
- **Language**: Golang (Go 1.21+)
- **Web Framework**: Gorilla Mux (HTTP routing)
- **Database**: PostgreSQL 14+
- **Authentication**: JWT (JSON Web Tokens)
- **Security**: Bcrypt password hashing
- **Architecture**: Layered architecture (Controller → Service → Repository)

### Frontend
- **Framework**: React 18+
- **Build Tool**: Vite
- **Styling**: Tailwind CSS 3+
- **Routing**: React Router DOM v6
- **HTTP Client**: Axios
- **Icons**: Lucide React
- **Notifications**: React Toastify

### Database
- **DBMS**: PostgreSQL 14+
- **Features**: ACID compliance, Foreign keys, Indexes, Triggers

---

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

1. **Go (Golang)** - Version 1.21 or higher
   - Download: [https://golang.org/dl/](https://golang.org/dl/)
   - Verify: `go version`

2. **PostgreSQL** - Version 14 or higher
   - Download: [https://www.postgresql.org/download/](https://www.postgresql.org/download/)
   - Verify: `psql --version`

3. **Node.js & npm** - Version 18+ (for frontend)
   - Download: [https://nodejs.org/](https://nodejs.org/)
   - Verify: `node --version` and `npm --version`

4. **Git** (optional, for cloning)
   - Download: [https://git-scm.com/downloads](https://git-scm.com/downloads)

5. **Docker & Docker Compose** (optional, for running PostgreSQL quickly)
   - Download: [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)

---

## 🚀 Installation & Setup Guide

### Step 1: Clone or Extract the Project

```bash
# If using Git
git clone <repository-url>
cd Leave-Management-System

# Or simply navigate to the extracted folder
cd Leave-Management-System
```

### Step 2: Database Setup

You can set up the database using either **local PostgreSQL** or **Docker Compose**.

#### Option A: Local PostgreSQL

1. **Create PostgreSQL Database:**

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE leave_management;

# Exit psql
\q
```

2. **Run Database Migrations:**

```bash
# Navigate to migrations directory
cd backend/internal/database/migrations

# Run migrations in order
psql -U postgres -d leave_management -f 001_create_users.sql
psql -U postgres -d leave_management -f 002_create_leave_types.sql
psql -U postgres -d leave_management -f 003_create_leave_balances.sql
psql -U postgres -d leave_management -f 004_create_leave_requests.sql
psql -U postgres -d leave_management -f 005_create_audit_logs.sql
psql -U postgres -d leave_management -f 006_seed_data.sql
```

#### Option B: Docker Compose (PostgreSQL)

```bash
# Navigate to backend directory
cd backend

# Copy environment file (choose command for your OS)
# Windows PowerShell
Copy-Item .env.example .env
# macOS/Linux
cp .env.example .env

# Start PostgreSQL container
docker compose up -d
```

After the container is running, execute the same migration files listed in Option A.

**Note:** The seed data file creates test accounts with the following credentials:
- **Admin**: admin@abccompany.com / admin123
- **Manager**: manager@abccompany.com / manager123
- **Employee**: employee@abccompany.com / employee123

### Step 3: Backend Setup

1. **Navigate to backend directory:**

```bash
cd ../../..  # Back to project root if in migrations folder
cd backend
```

2. **Create environment file:**

```bash
# Copy the example file (choose command for your OS)
# Windows PowerShell
Copy-Item .env.example .env
# macOS/Linux
cp .env.example .env

# Edit .env with your settings (use your preferred text editor)
# Update database credentials if necessary
```

3. **Install Go dependencies:**

```bash
go mod download
```

4. **Run the backend server:**

```bash
go run cmd/server/main.go
```

The backend server will start on **http://localhost:8080**

You should see:
```
Configuration loaded successfully
Database connection established successfully
Server starting on port 8080
Environment: development
```

### Step 4: Frontend Setup

Open a **new terminal window/tab** and:

1. **Navigate to frontend directory:**

```bash
cd Leave-Management-System/frontend
```

2. **Create environment file:**

```bash
# Copy the example file (choose command for your OS)
# Windows PowerShell
Copy-Item .env.example .env
# macOS/Linux
cp .env.example .env

# The default settings should work if backend is on localhost:8080
```

3. **Install Node.js dependencies:**

```bash
npm install
```

4. **Run the frontend development server:**

```bash
npm run dev
```

The frontend will start on **http://localhost:5173**

You should see:
```
VITE v5.x.x  ready in xxx ms
➜  Local:   http://localhost:5173/
```

---

## 🎮 Using the Application

### Accessing the Application

1. Open your browser and navigate to **http://localhost:5173**
2. You'll be redirected to the login page

### Test Accounts

Use these pre-seeded accounts to test different roles:

| Role | Email | Password | Capabilities |
|------|-------|----------|--------------|
| Admin | admin@abccompany.com | admin123 | Full system access, user management, approve/reject leaves, view all reports |
| Manager | manager@abccompany.com | manager123 | Approve/reject leaves, view team reports, submit own leaves |
| Employee | employee@abccompany.com | employee123 | Submit leave requests, view own history and balance |

### Employee Workflow

1. **Login** with employee credentials
2. **View Dashboard**: See your leave balance for different leave types
3. **Request Leave**:
   - Click "Request Leave" button
   - Select leave type (Annual, Sick, Casual, etc.)
   - Choose start and end dates
   - Provide a reason (minimum 10 characters)
   - Submit request
4. **View History**: See all your leave requests with their status
5. **Track Status**: Monitor if your requests are pending, approved, or rejected

### Manager/Admin Workflow

1. **Login** with manager or admin credentials
2. **View Admin Dashboard**:
   - See pending requests count
   - View employees currently on leave
   - See approval/rejection statistics
3. **Review Leave Requests**:
   - Filter by status (Pending, Approved, Rejected, All)
   - View employee details, leave duration, and reason
   - **Approve** or **Reject** pending requests
   - Add optional review notes
4. **Switch to Employee View**: Use "Employee View" button to access your own leave portal
5. **User Management** (Admin only): Create new users, manage roles

---

## 📁 Project Structure

```
Leave-Management-System/
├── backend/
│   ├── cmd/server/
│   │   └── main.go                 # Application entry point
│   ├── internal/
│   │   ├── config/                 # Configuration management
│   │   ├── database/               # Database connection & migrations
│   │   ├── models/                 # Data models
│   │   ├── repository/             # Data access layer
│   │   ├── service/                # Business logic layer
│   │   ├── handler/                # HTTP handlers (controllers)
│   │   ├── middleware/             # Auth, RBAC, CORS, logging
│   │   ├── dto/                    # Request/response DTOs
│   │   └── utils/                  # Utilities (JWT, password, date)
│   ├── pkg/errors/                 # Custom error types
│   ├── go.mod                      # Go module definition
│   └── .env.example                # Environment variables template
│
└── frontend/
    ├── src/
    │   ├── api/                    # API service layer (Axios)
    │   ├── components/             # React components
    │   │   ├── auth/               # Authentication components
    │   │   ├── employee/           # Employee-specific features
    │   │   ├── admin/              # Admin-specific features
    │   │   └── common/             # Reusable UI components
    │   ├── context/                # React Context (AuthContext)
    │   ├── pages/                  # Page components
    │   ├── styles/                 # Global styles
    │   └── utils/                  # Helper functions
    ├── package.json                # Node.js dependencies
    ├── vite.config.js              # Vite configuration
    ├── tailwind.config.js          # Tailwind CSS configuration
    └── .env.example                # Frontend environment variables
```

---

## 🔧 Configuration

### Backend Configuration (.env)

```env
# Server
PORT=8080
ENV=development

# Optional for Docker Compose database service
POSTGRES_USER=my_secure_user
POSTGRES_PASSWORD=my_secure_password
POSTGRES_DB=leave_management

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=leave_management
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key-min-32-chars
JWT_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# CORS
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# App
DEFAULT_LEAVE_DAYS=20
```

### Frontend Configuration (.env)

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

---

## 📊 Database Schema

### Key Tables

1. **users**: Stores user accounts with roles (employee, manager, admin)
2. **leave_types**: Defines leave categories (Annual, Sick, Casual, etc.)
3. **leave_balances**: Tracks leave allocation and usage per user per year
4. **leave_requests**: Stores all leave submissions and their approval status
5. **audit_logs**: Tracks critical actions for accountability

### Relationships

- Users can have many leave requests (1:N)
- Users can have many leave balances (1:N)
- Leave requests reference leave types (N:1)
- Leave requests track who reviewed them (N:1 to users)

---

## 🔐 Security Features

1. **Password Security**: Bcrypt hashing with cost factor 12
2. **JWT Authentication**: Stateless token-based authentication
   - Access tokens expire in 15 minutes
   - Refresh tokens expire in 7 days
3. **Role-Based Access Control (RBAC)**: Middleware enforces permissions
4. **Input Validation**: All inputs validated on backend
5. **SQL Injection Prevention**: Parameterized queries
6. **CORS Protection**: Whitelist-based origin checking

---

## 🧪 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/me` - Get current user
- `GET /api/v1/departments` - Get available departments

### Employee Endpoints
- `POST /api/v1/leaves` - Submit leave request
- `GET /api/v1/leaves` - Get own leave history
- `GET /api/v1/leaves/:id` - Get leave request details
- `DELETE /api/v1/leaves/:id` - Cancel leave request
- `GET /api/v1/leave-types` - Get available leave types
- `GET /api/v1/profile/balance` - Get leave balance

### Admin/Manager Endpoints
- `GET /api/v1/admin/leaves` - Get all leave requests
- `PUT /api/v1/admin/leaves/:id/approve` - Approve leave
- `PUT /api/v1/admin/leaves/:id/reject` - Reject leave
- `GET /api/v1/admin/dashboard` - Get dashboard statistics
- `POST /api/v1/admin/balances` - Allocate or update leave balance

### Admin Only Endpoints
- `GET /api/v1/admin/users` - Get all users
- `POST /api/v1/admin/users` - Create new user
- `PUT /api/v1/admin/users/:id` - Update user
- `DELETE /api/v1/admin/users/:id` - Delete user

### System Endpoints
- `GET /health` - Health check endpoint

---

## 🐛 Troubleshooting

### Backend Issues

**Database connection failed:**
- Verify PostgreSQL is running: `pg_isready`
- Check credentials in `.env` file
- Ensure database exists: `psql -U postgres -l`

**Port already in use:**
- Change `PORT` in backend `.env` file
- Kill existing process (Windows PowerShell): `Get-NetTCPConnection -LocalPort 8080 | ForEach-Object { Stop-Process -Id $_.OwningProcess -Force }`
- Kill existing process (Mac/Linux): `lsof -ti:8080 | xargs kill`

### Frontend Issues

**Cannot connect to backend:**
- Verify backend is running on port 8080
- Check `VITE_API_BASE_URL` in frontend `.env`
- Check browser console for CORS errors

**Blank page after login:**
- Check browser console for JavaScript errors
- Verify API responses in Network tab
- Clear browser cache and localStorage

---

## 🚀 Deployment

### Backend Deployment

1. Build for production:
```bash
go build -o leave-management-server cmd/server/main.go
```

2. Set production environment variables
3. Run: `./leave-management-server`

### Frontend Deployment

1. Build for production:
```bash
npm run build
```

2. Deploy `dist` folder to static hosting (Netlify, Vercel, etc.)
3. Update `VITE_API_BASE_URL` to production backend URL

---

## 💡 Features Implemented

✅ User registration and authentication (JWT)
✅ Role-based access control (Employee, Manager, Admin)
✅ Leave request submission with validation
✅ Leave balance tracking and management
✅ Leave approval/rejection workflow
✅ Dashboard with statistics
✅ Responsive design (mobile, tablet, desktop)
✅ Real-time leave balance calculation
✅ Working days calculation (excludes weekends)
✅ Leave history with pagination
✅ Employee on-leave tracking
✅ Audit logging
✅ Error handling and input validation

---

## 📝 License

This project was developed as part of an assessment for ABC Company.

---

## 👨‍💻 Support

For issues or questions:
1. Check the troubleshooting section above
2. Review the PROJECT_EXPLANATION.md for system design details
3. Verify all prerequisites are correctly installed

---

## 🎉 Conclusion

You now have a fully functional Leave Management System! The system demonstrates clean architecture, proper separation of concerns, and professional UI/UX design. All core features for managing employee leaves are implemented and ready for use.

**Happy Testing! 🚀**
