# Snippetbox - Go Web Application

A web application for creating and sharing code snippets, built with Go following modern web development practices.

## Features

- Create, view, and browse code snippets with expiration dates
- User authentication and account management
- Session-based security with CSRF protection
- Clean, responsive web interface
- HTTPS support

## Technology Stack

- **Go 1.24.4** with httprouter and Alice middleware
- **MySQL** with session storage
- **bcrypt** password hashing and **nosurf** CSRF protection
- HTML templates with embedded static files

## Project Structure

```
├── cmd/web/                 # Application entry point and web handlers
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

## Quick Start

### Prerequisites
- Go 1.24.4+
- MySQL 5.7+

### Setup
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

The application supports HTTPS using TLS certificates. The `tls/` directory is not included in the repository for security reasons, so you'll need to create certificates:

### Option 1: Generate Self-Signed Certificates (Development)
```bash
# Create tls directory
mkdir tls

# Generate self-signed certificates
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

# Move certificates to tls directory
mv cert.pem tls/
mv key.pem tls/

# Run with HTTPS
go run ./cmd/web
# Visit https://localhost:4000
```

### Option 2: Use OpenSSL (Alternative for Development)
```bash
# Create tls directory
mkdir tls

# Generate private key and certificate
openssl req -new -x509 -sha256 -key <(openssl genrsa 2048) -out tls/cert.pem -days 365 -nodes \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
openssl genrsa -out tls/key.pem 2048

# Or use a simpler one-liner
openssl req -x509 -newkey rsa:2048 -keyout tls/key.pem -out tls/cert.pem -days 365 -nodes \
  -subj "/CN=localhost"

# Run the application
go run ./cmd/web
```

### Option 3: Use Let's Encrypt (Production)
For production deployment with a domain:
```bash
# Create tls directory
mkdir tls

# Install certbot
# Ubuntu/Debian: sudo apt install certbot
# macOS: brew install certbot

# Obtain certificates (replace your-domain.com)
sudo certbot certonly --standalone -d your-domain.com

# Copy certificates to tls directory
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem tls/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem tls/key.pem
sudo chown $USER:$USER tls/*.pem

# Run with your domain
go run ./cmd/web -addr=":443"
```

### HTTP-Only Mode (Development)
To run without HTTPS (certificates not required):
```bash
# The application will automatically detect missing certificates and use HTTP
go run ./cmd/web
# Visit http://localhost:4000
```

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

## Database Schema

```sql
-- Snippets with expiration
CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);
CREATE INDEX idx_snippets_created ON snippets(created);

-- User accounts
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

-- Session storage
CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
```

## Testing

```bash
go test ./...           # Run all tests
go test -v ./cmd/web    # Verbose output for web package
```

## Development

### Project Structure
```
cmd/web/          # Web application (handlers, routes, middleware)
internal/models/  # Data models and database operations
ui/html/          # HTML templates
ui/static/        # CSS, JS, images
tls/              # TLS certificates
```

### Adding Features
1. Add handlers in `cmd/web/handlers.go`
2. Define routes in `cmd/web/routes.go`
3. Create templates in `ui/html/pages/`
4. Add models in `internal/models/`

### Security Features
- bcrypt password hashing
- CSRF protection via nosurf
- Secure session cookies
- Prepared SQL statements
- HTTPS support
