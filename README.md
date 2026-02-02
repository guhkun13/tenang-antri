# Queue Management System

A comprehensive web-based queue management system built with Go, Gin, HTMX, and PostgreSQL.

## Features

### Authentication & Authorization
- JWT token-based authentication
- Role-based access control (Admin/Staff)
- Secure password hashing
- Profile management

### Customer Features
- Self-service ticket generation kiosk
- Category selection
- Queue position display
- Estimated wait time

### Staff Features
- Counter operations dashboard
- Call next ticket
- Complete/No-show marking
- Pause/Resume counter
- Real-time queue visibility

### Admin Features
- Dashboard with real-time statistics
- Ticket management
- Category management (CRUD)
- Counter management (CRUD)
- Staff management (CRUD)
- Reports and analytics

### Display Board
- Real-time currently serving tickets
- Queue statistics
- Counter status
- Auto-refresh via WebSocket

## Tech Stack

- **Backend**: Go 1.21+, Gin framework
- **Database**: PostgreSQL 15+
- **Frontend**: HTML templates, HTMX, Alpine.js, Tailwind CSS
- **Real-time**: WebSocket
- **Deployment**: Docker, Docker Compose

## Quick Start

### Using Docker Compose

1. Clone the repository
2. Run: `docker-compose up -d`
3. Access the application at http://localhost

### Manual Setup

1. Install Go 1.21+ and PostgreSQL 15+
2. Create database: `createdb queue_system`
3. Run migrations: `psql -d queue_system -f migrations/001_initial_schema.up.sql`
4. Run seed data: `psql -d queue_system -f migrations/002_seed_data.up.sql`
5. Install dependencies: `go mod download`
6. Run server: `go run cmd/server/main.go`
7. Access at http://localhost:8080

## Default Credentials

- **Admin**: admin / admin123
- **Staff**: staff1 / staff123

## Project Structure

```
queue-system/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Authentication middleware
│   ├── models/          # Database models and repository
│   └── websocket/       # WebSocket hub
├── migrations/          # Database migrations
├── templates/           # HTML templates
├── static/              # Static assets
├── docker/              # Docker configuration
├── Dockerfile
└── docker-compose.yml
```

## API Endpoints

### Authentication
- `POST /login` - Login
- `GET /logout` - Logout
- `GET /api/profile` - Get profile
- `PUT /api/profile` - Update profile
- `POST /api/change-password` - Change password

### Admin
- `GET /admin/dashboard` - Dashboard
- `GET /admin/api/stats` - Get statistics
- `CRUD /admin/api/users` - User management
- `CRUD /admin/api/categories` - Category management
- `CRUD /admin/api/counters` - Counter management
- `CRUD /admin/api/tickets` - Ticket management

### Staff
- `GET /staff/dashboard` - Staff dashboard
- `POST /staff/call-next` - Call next ticket
- `POST /staff/complete` - Complete current ticket
- `POST /staff/no-show` - Mark as no-show
- `POST /staff/pause` - Pause counter
- `POST /staff/resume` - Resume counter

### Kiosk
- `GET /kiosk` - Kiosk interface
- `POST /kiosk/ticket` - Generate ticket

### Display
- `GET /display` - Display board
- `GET /display/serving` - Currently serving
- `GET /display/stats` - Queue statistics

### WebSocket
- `GET /ws` - WebSocket connection for real-time updates

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | Server port | 8080 |
| SERVER_MODE | Gin mode (debug/release) | debug |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | postgres |
| DB_NAME | Database name | queue_system |
| JWT_SECRET | JWT secret key | your-secret-key |
| JWT_ACCESS_TOKEN_EXPIRY | Token expiry | 24h |

## License

MIT License
