# Snippetbox - Go Web Application

A web application for creating and sharing code snippets, built with Go following modern web development practices.

## Features

- Create, view, and browse code snippets with expiration dates
- User authentication and account management
- Session-based security with CSRF protection
- Secure session cookies and bcrypt password hashing
- Prepared SQL statements to prevent SQL injection

## Technology Stack

- **Go 1.24.4** with httprouter and Alice middleware
- **MySQL** with session storage
- **bcrypt** password hashing and **nosurf** CSRF protection
- HTML templates with embedded static files

## Setup
1. **Clone and install dependencies**
   ```bash
   git clone https://github.com/thienel/lets-go.git
   cd lets-go
   go mod download
   ```

2. **Database setup**
   ```bash
   mysql -u root -e "CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
   mysql -u root snippetbox < internal/models/testdata/setup.sql
   ```

3. **Run the application**
   ```bash
   go run ./cmd/web
   ```

   Visit `http://localhost:4000`

### Configuration
- `-addr`: Server address (default: ":4000")
- `-dsn`: MySQL DSN (default: "web:pass@/snippetbox?parseTime=true")
- `-debug`: Debug mode

## TLS/HTTPS Setup

The application supports HTTPS using TLS certificates. The `tls/` directory is not included in the repository for security reasons, so you'll need to create certificates.

**Note:** Add `tls/` to your `.gitignore` to avoid accidentally committing private keys.

## API Routes

**Public:**
- `GET /` - Home page with latest snippets
- `GET /snippet/view/:id` - View snippet
- `GET|POST /user/signup` - User registration
- `GET|POST /user/login` - User login
- `GET /about` - About page

**Protected (auth required):**
- `GET|POST /snippet/create` - Create snippet
- `GET /account/view` - User account
- `GET|POST /account/password/update` - Change password
- `POST /user/logout` - Logout

## Testing

```bash
go test ./...           # Run all tests
```

## Project Structure

```
├── cmd/web/                # Application entry point and web handlers
│   ├── main.go             # Main application setup and configuration
│   ├── handlers.go         # HTTP request handlers
│   ├── routes.go           # Route definitions and middleware setup
│   ├── middleware.go       # Custom middleware functions
│   ├── helpers.go          # Helper functions for handlers
│   ├── templates.go        # Template caching and rendering
│   └── context.go          # Request context utilities
├── internal/
│   ├── models/             # Data models and database layer
│   │   ├── snippets.go     # Snippet model and database operations
│   │   ├── users.go        # User model and authentication
│   │   ├── errors.go       # Custom error definitions
│   │   ├── mocks/          # Mock implementations for testing
│   │   └── testdata/       # Test database schemas and data
│   ├── assert/             # Testing utilities
│   └── validator/          # Input validation utilities
├── ui/
│   ├── html/               # HTML templates
│   │   ├── base.html       # Base template layout
│   │   ├── pages/          # Page templates
│   │   └── partials/       # Reusable template components
│   ├── static/             # Static assets (CSS, JS, images)
│   └── efs.go              # Embedded file system
├── tls/                    # TLS certificates
└── go.mod                  # Go module dependencies
```


## Security Features
- bcrypt password hashing
- CSRF protection via nosurf
- Secure session cookies
- Prepared SQL statements
- HTTPS support
