# RevUExchange

A reciprocal review platform for self-published authors and digital creators. Users earn points by reviewing others' work and spend points to get their own work reviewed.

## Tech Stack

- **Backend**: Go with chi/v5 router
- **Frontend**: React 18 with TypeScript, Vite, Tailwind CSS
- **Database**: PostgreSQL (AWS RDS)
- **Cache**: Redis (AWS ElastiCache)
- **Storage**: AWS S3
- **CDN**: AWS CloudFront
- **Infrastructure**: Terraform
- **CI/CD**: GitHub Actions
- **Authentication**: JWT

## Features

- User authentication (register/login)
- Product management (books, courses, etc.)
- Bounty marketplace for review requests
- Review system with ratings
- Points economy (earn/spend/transfer)
- Stripe payments for point purchases
- Social features (follow, comments, activity feed)
- Gamification (badges, leaderboards, streaks)
- Analytics dashboard
- Anti-fraud detection
- Email notifications (SendGrid)

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 20+
- Docker (for local development)
- LocalStack (for local AWS services)

### Local Development

```bash
# Install dependencies
cd api && go mod download
cd frontend && npm install

# Start infrastructure
make local-start

# Start API
make local-api

# Start frontend (in another terminal)
make local-ui
```

### Running Tests

```bash
# All tests
make test

# Backend only
make test-backend

# Frontend only
make test-frontend
```

## Project Structure

```
revueexchange.com/
├── api/                    # Go backend
│   ├── cmd/server/        # Entry point
│   └── internal/
│       ├── config/        # Configuration
│       ├── handler/       # HTTP handlers
│       ├── middleware/    # Auth middleware
│       ├── model/        # Data models
│       ├── repository/   # Database operations
│       └── service/      # Business logic
├── frontend/              # React frontend
│   ├── src/
│   │   ├── components/   # Reusable components
│   │   ├── pages/        # Page components
│   │   ├── services/     # API client
│   │   └── hooks/        # Custom hooks
│   └── public/
├── terraform/            # Infrastructure
│   ├── env/prod/       # Production config
│   └── modules/         # Reusable modules
└── docs/               # Documentation
```

## API Documentation

See [API Endpoints](docs/API.md) for detailed endpoint documentation.

## Deployment

See [Production Deployment Guide](docs/PRODUCTION_DEPLOYMENT.md) for deployment instructions.

## License

MIT
