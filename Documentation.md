# AffPilot Auth Service Documentation

## Introduction

### Overview
AffPilot Auth Service is a robust authentication and authorization system built in Go that provides secure user management, role-based access control, and email verification functionality.

## Getting Started

### Prerequisites

1. **Docker Installation**
   - Install Docker Engine
   - Install Docker Compose

### Clone and Run
1. **Clone Repository**
   ```bash
   git clone https://github.com/shahid-affpilot/Developer-Assignment.git
   cd Developer-Assignment
   go mod tidy
   ```

2. **Configure Environment**
   ```bash
   cp .env.example .env
   # Edit .env file with your configuration
   ```

3. **Build and Start Services**
   ```bash
   docker compose up --build // from linux os
   docker-compose up --build // from linux windows
   ```

4. **Verify Installation**
   Using curl:
   ```bash
   curl http://localhost:8080/register
   ```
   Or using Postman (recommended):
   - Import collection from Postman section
   - Execute "Health Check" request


    Import the complete API collection by import `Auth System.postman_collection.json` (in my root directory) into postman application. Or you can get the collection using this link:
    ```
    https://affpilot-2941.postman.co/workspace/Affpilot-Workspace~2220a0d0-ef3f-410d-84c5-1fd95022775b/collection/44639385-0a9112f5-2843-4de4-ad95-2e094fcb129e
    ```
    The collection includes all API endpoints with example requests and environment variables.

### Environment Setup

1. **Configuration**
   ```env
   # Create .env file with required variables
   DB_HOST=postgres
   DB_PORT=5432
   DB_USER=youruser
   DB_PASSWORD=yourpassword
   DB_NAME=affpilot
   JWT_SECRET=your-secret-key
   SMTP_HOST=smtp.example.com
   SMTP_PORT=587
   SMTP_USER=your-email
   SMTP_PASSWORD=your-password
   ```

#### Initialization Process
The system automatically handles:
- Loading environment variables
- Establishing database connection
- Creating initial system admin account
- Setting up required permissions


## Database Architecture
The system uses PostgreSQL with the following table structure:

1. **Users Table**
   - Stores user information
   - Columns: id, username, email, first_name, last_name, etc.

2. **Roles Table**
   - Manages system roles
   - Columns: role_id, name, description, created_at, updated_at

3. **Permissions Table**
   - Defines system permissions
   - Columns: id, name, description, resource, action, created_at, updated_at

4. **UserRoles Table**
   - Maps users to roles
   - Columns: user_id, role_id, assigned_by, timestamp
   - Composite primary key: (user_id, role_id)

5. **RolePermissions Table**
   - Maps roles to permissions
   - Columns: id, role_id, permission_id, timestamp

## Role-Based Access Control

The system implements a hierarchical permission structure:

1. **System Admin**
   - Has all permissions of lower roles
   - Exclusive ability to promote users to admin

2. **Admin**
   - User management (list, update, delete)
   - Role management (CRUD operations)
   - Permission management
   - Can promote users to moderator

3. **Moderator**
   - View user details
   - Delete users

4. **Authenticated Users**
   - Manage own profile
   - View own permissions
   - Request account deletion
   - Logout

5. **Public Access**
   - Register
   - Login
   - Email verification
   - Password reset


## Project Structure

### Key Components

1. **Entry Point**
   - `/cmd/main.go`: Application initialization and route setup

2. **Core Modules**
   - `/internal/http/routes`: API route definitions
   - `/internal/http/middleware`: Request middleware
   - `/internal/http/handlers`: Request handlers organized by domain
     - `/auth`: Authentication handlers
     - `/permission`: Permission management
     - `/role`: Role management
     - `/user`: User management

## Security Features

The service implements multiple security layers:

1. **JWT Authentication**
   - Secure token-based authentication
   - Token validation and refresh

2. **Email Verification**
   - Required for account activation
   - Secure verification links

3. **Password Security**
   - Bcrypt hashing
   - Secure reset mechanism


## Project Initialization

The application's entry point is `cmd/main.go`, which handles the following initialization sequence:

1. **Configuration Loading**
   ```go
   // filepath: /cmd/main.go
   func init() {
       config.LoadConfig()  // Loads environment variables from .env
   }
   ```

2. **Database Setup**
   ```go
   // filepath: /internal/database/db.go
   func ConnDB() {
       // Establishes PostgreSQL connection
       // Sets up connection pool
   }
   ```

3. **System Admin Initialization**
   ```go
   // filepath: /internal/database/admin.go
   func InitAdminUser() {
       // Creates initial system admin if not exists
       // Assigns all permissions
   }
   ```


## API Documentation

The project contains 24 API endpoints grouped into 4 categories: /auth, /users, /roles, and /permissions. All endpoints are prefixed with `/api/v1/`. Here's the detailed documentation:

### Authentication APIs (/auth)

* **POST - auth/register** - Register new user - Public
  - Handled by `Register()` in `/internal/http/handlers/auth`
  - Validates request body and converts to struct
  - Checks for existing username/email
  - Creates user with default 'user' role
  - Sends verification email
  - Response includes user ID, username, email and success message

* **POST - auth/login** - Authenticate user - Public
  - Validates username/password credentials
  - Issues JWT token on successful authentication
  - Returns token and user profile information

* **POST - auth/logout** - End user session - Authenticated
  - Invalidates current JWT token
  - Requires valid authentication token

* **GET - auth/verify** - Verify email address - Public
  - Validates verification token from email
  - Updates user's email_verified status
  - Expires tokens after 5 minutes

* **POST - auth/resend-verification** - Resend verification email - Public
  - Generates new verification token
  - Sends fresh verification email
  - Updates token expiry timestamp

* **POST - auth/password-reset** - Reset user password - Public
  - Sends password reset link to registered email
  - Token expires after 5 minutes

### User Management APIs (/users)

* **GET - users** - List all users - Admin+
  - Returns paginated list of all users
  - Includes basic user information

* **GET - users/{user_id}** - Get user details - User(self)/Moderator+
  - Returns detailed user profile
  - Users can only access their own profile
  - Moderators and above can access any profile

* **PUT - users/{user_id}** - Update user - User(self)/Admin+
  - Update user profile information
  - Users can only update their own profile
  - Admins can update any user except System Admins

* **POST - users/{user_id}/request-deletion** - Request account deletion - User(self)
  - Marks account for deletion
  - Requires moderator approval

* **DELETE - users/{user_id}** - Delete user account - Moderator+
  - Permanently removes user account
  - Cannot delete System Admins

* **POST - users/{user_id}/role** - Change user role - Admin+
  - Assigns new role to user
  - Records who made the change

* **POST - users/{user_id}/promote/admin** - Promote to admin - System Admin
  - Promotes user to Admin role
  - Only available to System Admins

* **POST - users/{user_id}/promote/moderator** - Promote to moderator - Admin+
  - Promotes user to Moderator role
  - Available to Admins and System Admins

* **POST - users/{user_id}/demote** - Demote user role - Admin+
  - Demotes user to lower role
  - Cannot demote System Admins

### Role Management APIs (/roles)

* **GET - roles** - List all roles - Admin+
  - Returns all system roles
  - Includes role permissions

* **GET - roles/{role_id}** - Get role details - Admin+
  - Returns detailed role information
  - Includes assigned permissions

* **POST - roles** - Create new role - Admin+
  - Creates custom role
  - Assigns permissions to role

* **PUT - roles/{role_id}** - Update role - Admin+
  - Modifies role details
  - Updates role permissions

* **DELETE - roles/{role_id}** - Delete role - Admin+
  - Removes role from system
  - Cannot delete built-in roles

### Permission Management APIs (/permissions)

* **GET - permissions** - List all permissions - Admin+
  - Returns all system permissions
  - Grouped by resource type

* **GET - permissions/{permission_id}** - Get permission details - Admin+
  - Returns detailed permission information

### User Profile APIs (/me)

* **GET - me** - Get own profile - Authenticated
  - Returns current user's profile
  - Includes role information

* **GET - me/permissions** - Get own permissions - Authenticated
  - Lists all permissions for current user
  - Derived from assigned roles



## Troubleshooting

### PostgreSQL Port Conflicts

1. **Check for Active PostgreSQL Processes**
   ```bash
   sudo lsof -i :5432
   sudo lsof -i :8080
   ```
   This command shows:
   - Process ID (PID)
   - User running the process
   - Connection details
   
2. **Kill Conflicting Process**
   ```bash
   sudo kill <PID>
   ```
   Replace `<PID>` with the process ID from step 1

### Docker Container Issues

1. **Stop Running Containers**
   ```bash
   docker stop $(docker ps -aq)
   ```

2. **Remove Containers**
   ```bash
   docker rm $(docker ps -aq)
   ```

3. **Clean Docker System**
   ```bash
   docker system prune -a
   ```

### Common Issues

1. **Database Connection Failures**
   - Check PostgreSQL service status:
     ```bash
     sudo systemctl status postgresql
     ```
   - Verify port availability:
     ```bash
     netstat -tulpn | grep 5432
     ```

2. **Docker Permission Issues**
   - Add user to docker group:
     ```bash
     sudo usermod -aG docker $USER
     ```
   - Apply changes:
     ```bash
     newgrp docker
     ```

3. **Network Issues**
   - Check Docker network:
     ```bash
     docker network ls
     docker network inspect affpilot_network
     ```

## Technology Stack

| Component    | Technology                           |
|--------------|--------------------------------------|
| Language     | Go (Golang)                          |
| Database     | PostgreSQL                           |
| Security     | JWT, bcrypt                          |
| Deployment   | Docker                               |