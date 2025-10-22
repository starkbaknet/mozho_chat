# Mozho Chat Backend

A scalable, secure chat application backend built with Go, featuring real-time messaging, file attachments, end-to-end encryption, and modern architecture patterns.

## 🚀 Features

- **User Management**: Registration, authentication, profile management
- **Real-time Chat**: Direct messages and group conversations
- **File Attachments**: Secure file upload with S3-compatible storage
- **End-to-End Encryption**: AES encryption for message security
- **Message Status**: Read receipts, delivery confirmations
- **Session Management**: JWT-based authentication with refresh tokens
- **Scalable Architecture**: Clean separation of concerns with repository pattern

## 🏗️ Architecture

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── api/            # HTTP routing and middleware
│   ├── user/           # User management (auth, profiles)
│   ├── chatroom/       # Chat room operations
│   ├── message/        # Message handling and encryption
│   ├── models/         # Database models
│   ├── repository/     # Data access layer
│   └── config/         # Configuration management
├── pkg/
│   ├── auth/           # JWT and password utilities
│   ├── encryption/     # AES encryption service
│   ├── middleware/     # HTTP middleware (CORS, auth)
│   └── s3/            # File storage service
├── migrations/         # Database schema migrations
└── scripts/           # Database migration scripts
```

## 🛠️ Tech Stack

- **Language**: Go 1.23
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Storage**: S3-compatible (MinIO)
- **Authentication**: JWT tokens
- **Encryption**: AES-256
- **Migration**: golang-migrate

## 📋 Prerequisites

- Go 1.23+
- PostgreSQL 16+
- Redis
- MinIO (or AWS S3)
- golang-migrate CLI

## 🚀 Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd mozho_chat
go mod download
```

### 2. Environment Configuration

Create a `.env` file in the project root:

```env
# Database
POSTGRES_URL=postgres://username:password@localhost:5432/mozho_chat?sslmode=disable
POSTGRES_USER=your_username
POSTGRES_PASSWORD=your_password
POSTGRES_DB=mozho_chat
POSTGRES_PORT=5432

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASS=
REDIS_DB=0

# AWS S3 (or MinIO)
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
S3_BUCKET_NAME=mozho-chat-files
S3_ENDPOINT=http://localhost:9000  # For MinIO
```

### 3. Database Setup

```bash
# Start PostgreSQL and Redis using Docker Compose
docker-compose up -d postgres redis minio

# Run database migrations
chmod +x scripts/migrate-up.sh
./scripts/migrate-up.sh
```

### 4. Run the Application

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

Most endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

### User Endpoints

#### Register User

```http
POST /users/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

#### Login

```http
POST /users/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword"
}
```

#### Get Profile

```http
GET /users/me
Authorization: Bearer <token>
```

#### Update Profile

```http
PATCH /users/me
Authorization: Bearer <token>
Content-Type: application/json

{
  "username": "newusername",
  "profile": {
    "full_name": "John Doe",
    "bio": "Software Developer",
    "avatar_url": "https://example.com/avatar.jpg"
  }
}
```

### Chat Room Endpoints

#### Create Chat Room

```http
POST /chatrooms
Authorization: Bearer <token>
Content-Type: application/json

{
  "other_user_id": "uuid-of-other-user"
}
```

#### Get Chat Room

```http
GET /chatrooms/{room_id}
Authorization: Bearer <token>
```

#### List User's Chat Rooms

```http
GET /chatrooms
Authorization: Bearer <token>
```

#### Join Chat Room

```http
POST /chatrooms/{room_id}/join
Authorization: Bearer <token>
```

#### Leave Chat Room

```http
POST /chatrooms/{room_id}/leave
Authorization: Bearer <token>
```

### Message Endpoints

#### Send Message

```http
POST /messages/send
Authorization: Bearer <token>
Content-Type: multipart/form-data

receiver_id: uuid-of-receiver
content: Hello, this is a test message
algorithm: aes-256-gcm
encryption_key: base64-encoded-key
attachments: [file1, file2, ...]
```

#### Get Messages

```http
GET /messages/{chat_room_id}?limit=20&offset=0
Authorization: Bearer <token>
```

#### Mark Message as Read

```http
POST /messages/{message_id}/read
Authorization: Bearer <token>
```

#### Mark Message as Delivered

```http
POST /messages/{message_id}/delivered
Authorization: Bearer <token>
```

#### Generate Encryption Key

```http
POST /messages/generate-key
Authorization: Bearer <token>
```

## 🔐 Security Features

### Authentication

- JWT-based authentication
- Password hashing with bcrypt
- Session management with refresh tokens

### Encryption

- End-to-end message encryption using AES-256-GCM
- Client-side key generation
- Secure key exchange

### File Security

- Secure file upload to S3-compatible storage
- File metadata stored in database
- Access control for attachments

## 🗄️ Database Schema

The application uses PostgreSQL with the following main entities:

- **Users**: User accounts with profile information
- **Chat Rooms**: Conversation containers (direct messages or groups)
- **Chat Room Members**: User-room relationships
- **Messages**: Chat messages with encryption support
- **Message Status**: Read/delivery status tracking
- **Sessions**: User authentication sessions
- **User Public Keys**: Encryption key management
- **Message Attachments**: File attachment metadata

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run specific test packages
go test ./tests/
```

## 🐳 Docker Support

The project includes Docker Compose configuration for development:

```bash
# Start all services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f
```

Services included:

- PostgreSQL 16
- Redis
- MinIO (S3-compatible storage)

## 📝 Migration Commands

```bash
# Apply migrations
./scripts/migrate-up.sh

# Rollback migrations
./scripts/migrate-down.sh

# Create new migration
migrate create -ext sql -dir migrations -seq <migration_name>
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:

- Create an issue in the repository
- Check the API documentation above
- Review the code comments for implementation details
