# WillowSuite Vault

A comprehensive inventory and organization management system that helps you track and organize your belongings across multiple hierarchical levels - from buildings down to individual items.

## 🏗️ Overview

WillowSuite Vault is a full-stack web application that provides a hierarchical organization system for managing physical spaces and items. The system supports a 6-level hierarchy:

- **Buildings** → **Rooms** → **Shelving Units** → **Shelves** → **Containers** → **Items**

Each level can contain multiple children of the level below it, creating a flexible and scalable organization structure.

## ✨ Features

- **Hierarchical Organization**: 6-level hierarchy for comprehensive space management
- **User Authentication**: Secure authentication using AWS Cognito
- **QR Code Generation**: Generate QR codes for quick item identification
- **Search & Filtering**: Advanced search and filtering capabilities across all entity types
- **Pagination**: Efficient pagination for large datasets
- **Real-time Updates**: Live updates with Redis caching
- **Responsive Design**: Modern, responsive UI built with SvelteKit and Tailwind CSS
- **RESTful API**: Clean, well-documented API endpoints
- **Docker Support**: Containerized deployment with Docker Compose
- **Infrastructure as Code**: Terraform configuration for AWS resources

## 🏛️ Architecture

### Backend (Go)
- **Framework**: Chi router with GORM ORM
- **Database**: PostgreSQL with Redis caching
- **Authentication**: AWS Cognito integration
- **File Storage**: AWS S3 for file uploads
- **Testing**: Comprehensive unit and integration tests with mocks

### Frontend (SvelteKit)
- **Framework**: SvelteKit with TypeScript
- **Styling**: Tailwind CSS with Skeleton UI components
- **Testing**: Playwright for E2E tests, Vitest for unit tests
- **State Management**: Svelte stores and context

### Infrastructure
- **Containerization**: Docker with multi-stage builds
- **Reverse Proxy**: Nginx for load balancing and SSL termination
- **Infrastructure**: Terraform-managed AWS resources (Cognito, S3, IAM)

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development)
- Node.js 18+ (for local development)
- AWS CLI configured (for production deployment)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd WillowSuite-Vault
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start the application**
   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:3000
   - Database: localhost:5432
   - Redis: localhost:6379

### Environment Variables

Create a `.env` file with the following variables:

```env
# Database
MASTER_DB_NAME=willowsuite_vault
MASTER_DB_USER=postgres
MASTER_DB_PASSWORD=your_password
MASTER_DB_HOST=db
MASTER_DB_PORT=5432
MASTER_DB_LOG_MODE=true
MASTER_SSL_MODE=disable

# Server
SECRET=your_secret_key
DEBUG=true
ALLOWED_HOSTS=localhost,127.0.0.1
SERVER_HOST=0.0.0.0
SERVER_PORT=3000
FRONT_END_URL=http://localhost:5173

# Redis
REDIS_HOST=redis
REDIS_USER=default
REDIS_PASSWORD=password

# AWS (for production)
AWS_ACCESS_KEY=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1
AWS_CLIENT_ID=your_cognito_client_id
AWS_CLIENT_SECRET=your_cognito_client_secret
AWS_USER_POOL_ID=your_user_pool_id
AWS_S3_BUCKET_NAME=your_s3_bucket

# Frontend
API_URL=http://localhost:3000
```

## 🧪 Testing

### Backend Tests
```bash
cd Backend
go test ./...
```

### Frontend Tests
```bash
cd Frontend
npm run test:unit      # Unit tests
npm run test:integration  # E2E tests
npm run test           # All tests
```

## 🏗️ Project Structure

```
WillowSuite-Vault/
├── Backend/                 # Go backend application
│   ├── config/             # Configuration management
│   ├── controllers/        # HTTP request handlers
│   ├── models/            # Database models
│   ├── repository/        # Data access layer
│   ├── routers/           # Route definitions
│   ├── infra/             # Infrastructure components
│   └── tests/             # Backend tests
├── Frontend/              # SvelteKit frontend application
│   ├── src/
│   │   ├── lib/           # Shared components and utilities
│   │   ├── routes/        # SvelteKit routes
│   │   └── stores/        # State management
│   ├── tests/             # Frontend tests
│   └── playwright-tests/  # E2E tests
├── Nginx/                 # Nginx configuration
├── Terraform/             # Infrastructure as Code
└── docker-compose.yml     # Local development setup
```

## 🔧 Development

### Backend Development

The backend is built with Go and follows clean architecture principles:

- **Controllers**: Handle HTTP requests and responses
- **Models**: Define database schema and business logic
- **Repository**: Abstract data access layer
- **Infrastructure**: External service integrations (Redis, S3, Cognito)

### Frontend Development

The frontend uses SvelteKit with modern tooling:

- **Components**: Reusable UI components in `src/lib/`
- **Routes**: Page components in `src/routes/`
- **Stores**: State management with Svelte stores
- **API Integration**: Server-side API calls in route handlers

## 🚀 Deployment

### Production Deployment

1. **Set up AWS infrastructure**
   ```bash
   cd Terraform/local
   terraform init
   terraform plan
   terraform apply
   ```

2. **Build and deploy**
   ```bash
   docker-compose -f docker-compose.prd.yml up -d
   ```

### Environment Configuration

- **Development**: Uses local Docker containers
- **Production**: Uses AWS services (Cognito, S3, RDS)

## 📚 API Documentation

### Authentication Endpoints
- `POST /api/v1/user/signup` - User registration
- `POST /api/v1/user/signin` - User login
- `POST /api/v1/user/logout` - User logout
- `POST /api/v1/user/refresh` - Token refresh

### Entity Management
- `GET /api/v1/entities` - Get paginated entities
- `POST /api/v1/entity` - Create new entity
- `GET /api/v1/entity/{category}/{id}` - Get specific entity
- `PUT /api/v1/entity` - Update entity
- `DELETE /api/v1/entity/{category}/{id}` - Delete entity

### Hierarchy Management
- `GET /api/v1/parents/{category}` - Get available parents
- `GET /api/v1/children/{category}/{id}` - Get entity children

### QR Code Generation
- `GET /api/v1/qr/{category}/{id}` - Generate QR code for entity

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions, please open an issue in the GitHub repository.
