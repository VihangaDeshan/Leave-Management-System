# Leave Management System API Documentation

## 1. Overview

This document describes the REST API for the Leave Management System backend.

- Base URL: `http://localhost:8080`
- API Prefix: `/api/v1`
- Content-Type: `application/json`
- Auth Type: JWT Bearer token

Health endpoint:

- `GET /health`

## 2. Authentication and Authorization

### 2.1 Bearer Token

Protected endpoints require this header:

`Authorization: Bearer <access_token>`

### 2.2 Roles

- `employee`: can manage own leave requests and view own balance
- `manager`: employee permissions + review leave requests for managed employees
- `admin`: full access, including user management

## 3. Standard Response Format

All API responses use this shape:

```json
{
  "success": true,
  "message": "string",
  "data": {},
  "error": null
}
```

Error example:

```json
{
  "success": false,
  "message": "Validation error",
  "error": "Key: 'CreateLeaveRequest.Reason' Error:Field validation for 'Reason' failed on the 'min' tag"
}
```

## 4. Status Codes (Common)

- `200 OK`: Request successful
- `201 Created`: Resource created
- `400 Bad Request`: Invalid input or business rule violation
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Duplicate/conflicting request
- `500 Internal Server Error`: Unexpected server-side error

## 5. Public Endpoints

### 5.1 Register

- Method: `POST`
- URL: `/api/v1/auth/register`
- Auth: No

Request body:

```json
{
  "email": "new.user@abccompany.com",
  "password": "Password123",
  "first_name": "New",
  "last_name": "User",
  "department": "Engineering"
}
```

Success response (`201`):

```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "access_token": "<jwt>",
    "refresh_token": "<jwt>",
    "user": {
      "id": 7,
      "email": "new.user@abccompany.com",
      "first_name": "New",
      "last_name": "User",
      "role": "employee",
      "department": "Engineering",
      "is_active": true
    }
  }
}
```

### 5.2 Login

- Method: `POST`
- URL: `/api/v1/auth/login`
- Auth: No

Request body:

```json
{
  "email": "admin@abccompany.com",
  "password": "admin123"
}
```

Success response (`200`): same shape as register response.

### 5.3 Refresh Token

- Method: `POST`
- URL: `/api/v1/auth/refresh`
- Auth: No

Request body:

```json
{
  "refresh_token": "<refresh_jwt>"
}
```

Success response (`200`):

```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "<new_access_jwt>"
  }
}
```

### 5.4 Get Departments

- Method: `GET`
- URL: `/api/v1/departments`
- Auth: No

Success response (`200`):

```json
{
  "success": true,
  "message": "Departments retrieved successfully",
  "data": ["Engineering", "HR", "Finance"]
}
```

## 6. Protected Endpoints (Authenticated Users)

### 6.1 Get Current User

- Method: `GET`
- URL: `/api/v1/auth/me`
- Auth: Yes

### 6.2 Create Leave Request

- Method: `POST`
- URL: `/api/v1/leaves`
- Auth: Yes

Request body:

```json
{
  "leave_type_id": 1,
  "start_date": "2026-04-01",
  "end_date": "2026-04-03",
  "reason": "Family event and travel arrangements"
}
```

Business rules:

- Dates must use `YYYY-MM-DD`
- End date must be greater than or equal to start date
- Start date cannot be in the past
- Must include at least one working day
- Overlapping leave requests are rejected
- Available balance must be sufficient

### 6.3 Get My Leave Requests

- Method: `GET`
- URL: `/api/v1/leaves?page=1&page_size=10`
- Auth: Yes

Response `data` shape:

```json
{
  "leave_requests": [],
  "total_count": 0,
  "page": 1,
  "page_size": 10,
  "total_pages": 0
}
```

### 6.4 Get Leave Request by ID

- Method: `GET`
- URL: `/api/v1/leaves/{id}`
- Auth: Yes

Notes:

- Employees can view only their own leave request.
- Admin and manager can view any request they are permitted to access.

### 6.5 Cancel Leave Request

- Method: `DELETE`
- URL: `/api/v1/leaves/{id}`
- Auth: Yes

Notes:

- Only pending leave requests can be cancelled.
- Employee can cancel own request.
- Admin and manager can cancel where authorized by business rules.

### 6.6 Get Leave Types

- Method: `GET`
- URL: `/api/v1/leave-types`
- Auth: Yes

### 6.7 Get My Leave Balance

- Method: `GET`
- URL: `/api/v1/profile/balance`
- Auth: Yes

Success `data` sample:

```json
[
  {
    "id": 1,
    "user_id": 3,
    "leave_type_id": 1,
    "leave_type": {
      "id": 1,
      "name": "Annual Leave",
      "description": "Yearly paid leave",
      "is_active": true
    },
    "total_days": 20,
    "used_days": 2,
    "available_days": 18,
    "year": 2026
  }
]
```

## 7. Admin and Manager Endpoints

Auth required: yes
Allowed roles: `admin`, `manager`

### 7.1 Get All Leave Requests

- Method: `GET`
- URL: `/api/v1/admin/leaves?page=1&page_size=10&status=pending`

Notes:

- Admin can view all leave requests.
- Manager gets requests scoped to managed employees.
- `status` filter is optional.

### 7.2 Approve Leave Request

- Method: `PUT`
- URL: `/api/v1/admin/leaves/{id}/approve`

Optional request body:

```json
{
  "review_notes": "Approved for project schedule flexibility"
}
```

Rules:

- Only pending requests can be approved.
- Managers can approve only requests for employees assigned to them.
- Approval updates used leave balance.

### 7.3 Reject Leave Request

- Method: `PUT`
- URL: `/api/v1/admin/leaves/{id}/reject`

Optional request body:

```json
{
  "review_notes": "Please re-submit with revised dates"
}
```

Rules:

- Only pending requests can be rejected.
- Managers can reject only requests for employees assigned to them.

### 7.4 Dashboard Statistics

- Method: `GET`
- URL: `/api/v1/admin/dashboard`

Success `data` shape:

```json
{
  "pending_requests": 3,
  "approved_requests": 12,
  "rejected_requests": 2,
  "employees_on_leave": 4,
  "total_employees": 20,
  "leaves_by_type": [],
  "recent_leave_requests": [],
  "employees_on_leave_today": []
}
```

### 7.5 Allocate Leave Balance

- Method: `POST`
- URL: `/api/v1/admin/balances`

Request body:

```json
{
  "user_id": 3,
  "leave_type_id": 1,
  "total_days": 20,
  "year": 2026
}
```

## 8. Admin-Only Endpoints

Auth required: yes
Allowed roles: `admin`

### 8.1 Get All Users

- Method: `GET`
- URL: `/api/v1/admin/users?page=1&page_size=10`

Success `data` shape:

```json
{
  "users": [],
  "total_count": 0,
  "page": 1,
  "page_size": 10,
  "total_pages": 0
}
```

### 8.2 Create User

- Method: `POST`
- URL: `/api/v1/admin/users`

Request body:

```json
{
  "email": "manager2@abccompany.com",
  "password": "Password123",
  "first_name": "Mina",
  "last_name": "Manager",
  "role": "manager",
  "department": "Engineering",
  "manager_id": null
}
```

`role` must be one of:

- `employee`
- `manager`
- `admin`

### 8.3 Update User

- Method: `PUT`
- URL: `/api/v1/admin/users/{id}`

Request body (all fields optional):

```json
{
  "first_name": "Updated",
  "last_name": "Name",
  "department": "Finance",
  "manager_id": 2
}
```

### 8.4 Deactivate User

- Method: `DELETE`
- URL: `/api/v1/admin/users/{id}`

Notes:

- Performs soft delete (deactivate).
- Returns `400` if user is already deactivated.

## 9. cURL Quick Start

### 9.1 Login and save token (PowerShell)

```powershell
$loginBody = @{
  email = "admin@abccompany.com"
  password = "admin123"
} | ConvertTo-Json

$loginRes = Invoke-RestMethod -Method POST -Uri "http://localhost:8080/api/v1/auth/login" -ContentType "application/json" -Body $loginBody
$token = $loginRes.data.access_token
```

### 9.2 Call protected endpoint

```powershell
Invoke-RestMethod -Method GET -Uri "http://localhost:8080/api/v1/auth/me" -Headers @{ Authorization = "Bearer $token" }
```

### 9.3 Create leave request

```powershell
$leaveBody = @{
  leave_type_id = 1
  start_date = "2026-04-01"
  end_date = "2026-04-03"
  reason = "Family event and travel arrangements"
} | ConvertTo-Json

Invoke-RestMethod -Method POST -Uri "http://localhost:8080/api/v1/leaves" -Headers @{ Authorization = "Bearer $token" } -ContentType "application/json" -Body $leaveBody
```

## 10. Endpoint Summary Table

| Method | Endpoint | Access |
|---|---|---|
| GET | /health | Public |
| POST | /api/v1/auth/register | Public |
| POST | /api/v1/auth/login | Public |
| POST | /api/v1/auth/refresh | Public |
| GET | /api/v1/departments | Public |
| GET | /api/v1/auth/me | Authenticated |
| POST | /api/v1/leaves | Authenticated |
| GET | /api/v1/leaves | Authenticated |
| GET | /api/v1/leaves/{id} | Authenticated |
| DELETE | /api/v1/leaves/{id} | Authenticated |
| GET | /api/v1/leave-types | Authenticated |
| GET | /api/v1/profile/balance | Authenticated |
| GET | /api/v1/admin/leaves | Admin/Manager |
| PUT | /api/v1/admin/leaves/{id}/approve | Admin/Manager |
| PUT | /api/v1/admin/leaves/{id}/reject | Admin/Manager |
| GET | /api/v1/admin/dashboard | Admin/Manager |
| POST | /api/v1/admin/balances | Admin/Manager |
| GET | /api/v1/admin/users | Admin |
| POST | /api/v1/admin/users | Admin |
| PUT | /api/v1/admin/users/{id} | Admin |
| DELETE | /api/v1/admin/users/{id} | Admin |
