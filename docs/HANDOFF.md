# Project Handoff Documentation

## Overview

RevUExchange is a reciprocal review platform for self-published authors and digital creators. Users earn points by reviewing others' work and spend points to get their own work reviewed.

## Architecture

### Technology Stack
- **Backend**: Go with chi/v5 router
- **Frontend**: React with TypeScript, Vite, Tailwind CSS
- **Database**: PostgreSQL (RDS)
- **Cache**: Redis (ElastiCache)
- **Storage**: S3 for documents/images
- **CDN**: CloudFront
- **Infrastructure**: Terraform
- **CI/CD**: GitHub Actions

### Project Structure
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
│   ├── env/prod/        # Production config
│   └── modules/         # Reusable modules
└── docs/               # Documentation
```

## Key Features

1. **User Authentication**
   - Email/password registration and login
   - JWT-based authentication
   - Bcrypt password hashing

2. **Product Management**
   - Create/edit/delete products (books, courses, etc.)
   - Product categories and metadata

3. **Bounty Marketplace**
   - Create bounties for review requests
   - Claim and complete reviews
   - Points-based compensation

4. **Reviews System**
   - Word count validation (100+ words)
   - Star ratings (1-5)
   - Draft/submitted states

5. **Points System**
   - Earn points by reviewing
   - Spend points to request reviews
   - Transfer points between users

6. **Payments (Stripe)**
   - Purchase points with real money
   - Secure checkout

7. **Social Features**
   - Follow/unfollow users
   - Comment on reviews
   - Activity feed

8. **Gamification**
   - Badges (DynamoDB)
   - Leaderboards
   - Streaks

9. **Analytics Dashboard**
   - Overview stats
   - Bounty/revenue metrics
   - User activity

10. **Security**
    - Anti-fraud detection
    - Email notifications (SendGrid)
    - Genre matching

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/me` - Get current user

### Users
- `GET /api/v1/users/{id}` - Get profile
- `PUT /api/v1/users/{id}` - Update profile

### Products
- `GET /api/v1/products/{id}` - Get product
- `POST /api/v1/products` - Create product
- `PUT /api/v1/products/{id}` - Update product
- `DELETE /api/v1/products/{id}` - Delete product

### Bounties
- `GET /api/v1/bounties` - List bounties
- `POST /api/v1/bounties` - Create bounty
- `GET /api/v1/bounties/{id}` - Get bounty
- `POST /api/v1/bounties/{id}/claim` - Claim bounty

### Reviews
- `POST /api/v1/reviews` - Create review
- `GET /api/v1/reviews/{id}` - Get review
- `PUT /api/v1/reviews/{id}` - Update review
- `POST /api/v1/reviews/{id}/submit` - Submit review

### Points
- `GET /api/v1/points/balance` - Get balance
- `GET /api/v1/points/transactions` - Transaction history
- `POST /api/v1/points/transfer` - Transfer points

### Social
- `POST /api/v1/users/{id}/follow` - Follow user
- `DELETE /api/v1/users/{id}/follow` - Unfollow user
- `GET /api/v1/users/{id}/followers` - Get followers
- `GET /api/v1/users/{id}/following` - Get following

### Comments
- `POST /api/v1/reviews/{id}/comments` - Add comment
- `GET /api/v1/reviews/{id}/comments` - Get comments

### Analytics
- `GET /api/v1/analytics/overview` - Overview stats
- `GET /api/v1/analytics/bounties` - Bounty metrics
- `GET /api/v1/analytics/reviews` - Review metrics
- `GET /api/v1/analytics/revenue` - Revenue stats
- `GET /api/v1/analytics/activity` - User activity

## Environment Variables

### Backend
```
DB_HOST=...
DB_PORT=5432
DB_NAME=revueexchange
DB_USER=...
DB_PASSWORD=...

REDIS_HOST=...
REDIS_PORT=6379

JWT_SECRET=...

STRIPE_SECRET_KEY=...
SENDGRID_API_KEY=...
```

### Frontend
```
VITE_API_URL=https://api.revueexchange.com/api
```

## Common Tasks

### Running Locally
```bash
# Start infrastructure
make local-start

# Start API
make local-api

# Start frontend
make local-ui

# Run tests
make test
```

### Deploying to Production
```bash
# Via GitHub Actions (push to main)
git checkout main
git merge develop
git push origin main
```

### Database Migrations
```bash
cd api
go run cmd/server/main.go migrate
```

### Scaling Services
```bash
# Update ECS service desired count
aws ecs update-service \
  --cluster revueexchange-prod \
  --service revueexchange-backend \
  --desired-count 4
```

## Support & Troubleshooting

### Check Logs
```bash
aws logs tail /ecs/revueexchange-prod-api --follow
```

### Restart Service
```bash
aws ecs update-service \
  --cluster revueexchange-prod \
  --service revueexchange-backend \
  --force-new-deployment
```

### Rollback
```bash
# Revert to previous task definition
aws ecs update-service \
  --cluster revueexchange-prod \
  --service revueexchange-backend \
  --task-definition revueexchange-backend:PREVIOUS-REVISION
```

## Future Improvements

1. Mobile app (React Native)
2. Real-time notifications (WebSocket)
3. Advanced analytics
4. Social sharing
5. Content moderation
6. Multi-language support
7. API rate limiting improvements
8. Advanced fraud detection (ML)

## Contact

For questions or issues, refer to:
- Documentation: `/docs/`
- GitHub Issues: https://github.com/horusclawd/revueexchange.com/issues
- API Health: `/health` endpoint
