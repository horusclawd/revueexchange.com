# RevUExchange - Architecture Plan

## Stack

| Layer | Technology |
|-------|------------|
| **Frontend** | React 18 + Vite + TypeScript |
| **Styling** | Tailwind CSS |
| **Backend API** | Go (Golang) |
| **API Gateway** | AWS API Gateway (HTTP API) |
| **Compute** | AWS ECS Fargate |
| **Database** | Aurora PostgreSQL |
| **Cache** | ElastiCache Redis |
| **NoSQL** | DynamoDB |
| **Storage** | S3 |
| **Auth** | Amazon Cognito + JWT |
| **Messaging** | EventBridge + SQS |
| **Payments** | Stripe |
| **Email** | SendGrid |
| **IaC** | Terraform |
| **CI/CD** | GitHub Actions |

### Monthly Cost: ~$150-200/month

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            AWS Cloud                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                    CloudFront (Global CDN)                             │   │
│  │               200+ Edge Locations Worldwide                            │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                    │                                          │
│           ┌────────────────────────┴────────────────────────┐               │
│           ▼                                                ▼               │
│  ┌─────────────────────┐                      ┌─────────────────────┐        │
│  │   S3 Static Website │                      │   API Gateway       │        │
│  │   (React + Vite)   │                      │   (HTTP API)        │        │
│  │                     │                      │                     │        │
│  │   /                │                      │   /api/v1/*         │        │
│  │   /bounties        │                      │   - auth            │        │
│  │   /bounties/:id    │                      │   - users           │        │
│  │   /dashboard       │                      │   - bounties        │        │
│  │   /profile         │                      │   - reviews         │        │
│  │                     │                      │   - points          │        │
│  └─────────────────────┘                      │   - payments        │        │
│                                                └──────────┬──────────┘        │
│                                                           │                    │
│           ┌───────────────────────────────────────────────┼────────┐         │
│           ▼                                               ▼        ▼         │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐        │
│  │   ECS Fargate   │    │  ECS Fargate    │    │  ECS Fargate    │        │
│  │   - Core API     │    │  - Workers      │    │  - Scheduled    │        │
│  │   - Business     │    │  - Email        │    │  - Analytics    │        │
│  │     Logic        │    │  - Webhooks     │    │  - Cleanup      │        │
│  └────────┬────────┘    └─────────────────┘    └─────────────────┘        │
│           │                                                                  │
├───────────┼──────────────────────────────────────────────────────────────────┤
│           ▼                                                                  │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                     Amazon EventBridge                                │   │
│  │   user.registered  │  review.submitted  │  bounty.created           │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│           │                        │                        │                │
│           ▼                        ▼                        ▼                │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐        │
│  │   SQS:         │    │   SQS:         │    │   SQS:         │        │
│  │  notifications │    │  webhooks       │    │  exports        │        │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘        │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                              Data Layer                                       │
│                                                                              │
│  ┌─────────────────────┐  ┌─────────────────────┐                          │
│  │  Aurora PostgreSQL  │  │  ElastiCache Redis  │                          │
│  │   (Serverless v2)   │  │   (Cluster Mode)    │                          │
│  │                     │  │                      │                          │
│  │  Primary Data:      │  │  • Sessions         │                          │
│  │  • Users            │  │  • API Cache        │                          │
│  │  • Bounties         │  │  • Rate Limiting    │                          │
│  │  • Reviews          │  │  • Real-time stats  │                          │
│  │  • Point TXs        │  │  • Leaderboard      │                          │
│  │  • Follows           │  │                      │                          │
│  └─────────────────────┘  └─────────────────────┘                          │
│                                                                              │
│  ┌─────────────────────┐  ┌─────────────────────┐                          │
│  │        S3           │  │     DynamoDB       │                          │
│  │   (Documents)      │  │    (Global)        │                          │
│  │                     │  │                     │                          │
│  │  • Book metadata   │  │  • Badges          │                          │
│  │  • User uploads    │  │  • Counters        │                          │
│  │  • Exports         │  │  • Streaks         │                          │
│  └─────────────────────┘  └─────────────────────┘                          │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                        External Integrations                                 │
│                                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │    Stripe    │  │  Amazon KDP │  │  Goodreads  │  │   SendGrid  │        │
│  │  (Payments)  │  │    (Books)  │  │  (Reviews)  │  │  (Emails)   │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Why Go?

| Factor | Benefit |
|--------|---------|
| **Performance** | 10-20x faster than Node.js for compute-heavy tasks |
| **Binary size** | Small Docker images (~20MB) |
| **Startup** | Fast container startup |
| **Concurrency** | Native goroutines for parallel processing |
| **Type safety** | Static typing catches errors at compile time |
| **Production** | Used by Google, Uber, Twitch, CloudFlare |

---

## Go API Structure

### Directory Layout

```
/cmd
  /api              # Main API entry point
    main.go
  /worker           # Async worker process
    main.go
/internal
  /config           # Configuration
  /handler          # HTTP handlers
  /middleware       # Auth, logging, cors
  /model            # Database models
  /repository       # Database queries
  /service          # Business logic
/pkg
  /logger           # Structured logging
  /errors           # Error handling
  /auth             # JWT utilities
  /stripe           # Stripe client
/migrations         # Database migrations
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/auth/register | Register new user |
| POST | /api/v1/auth/login | Login |
| GET | /api/v1/auth/me | Get current user |
| GET | /api/v1/users/{id} | Get user profile |
| PUT | /api/v1/users/{id} | Update user |
| GET | /api/v1/bounties | List bounties |
| POST | /api/v1/bounties | Create bounty |
| GET | /api/v1/bounties/{id} | Get bounty |
| POST | /api/v1/bounties/{id}/claim | Claim bounty |
| POST | /api/v1/reviews | Create review |
| GET | /api/v1/reviews/{id} | Get review |
| PUT | /api/v1/reviews/{id} | Update review |
| POST | /api/v1/reviews/{id}/submit | Submit review |
| GET | /api/v1/points/balance | Get points balance |
| GET | /api/v1/points/transactions | Get point history |
| POST | /api/v1/payments/checkout | Create Stripe session |
| POST | /api/v1/payments/webhook | Stripe webhook |
| GET | /api/v1/analytics/overview | Dashboard stats |
| GET | /api/v1/social/feed | Activity feed |
| POST | /api/v1/social/follow/{id} | Follow user |
| GET | /api/v1/gamification/badges | Get user badges |
| GET | /api/v1/gamification/leaderboard | Get leaderboard |

---

## Frontend (React + Vite)

### Project Structure

```
/src
  /components       # Reusable UI components
    /ui             # Base components (Button, Input, Card)
    /bounty         # Bounty-specific components
    /review         # Review-specific components
  /pages            # Page components
    /Home.tsx
    /Bounties.tsx
    /BountyDetail.tsx
    /Dashboard.tsx
    /Profile.tsx
  /hooks            # Custom React hooks
  /services         # API client functions
  /types            # TypeScript interfaces
  /utils            # Helper functions
  /context          # React context (Auth, Theme)
```

---

## Database Schema

### PostgreSQL Tables

```sql
-- Users
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255),
  username VARCHAR(50) UNIQUE NOT NULL,
  display_name VARCHAR(100),
  avatar_url TEXT,
  bio TEXT,
  points INTEGER DEFAULT 0,
  reputation_score DECIMAL(3,2) DEFAULT 0,
  stripe_customer_id VARCHAR(255),
  subscription_tier VARCHAR(20) DEFAULT 'free',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Products (books, courses, podcasts, newsletters)
CREATE TABLE products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  type VARCHAR(50) NOT NULL,
  title VARCHAR(500) NOT NULL,
  description TEXT,
  url TEXT,
  cover_image_url TEXT,
  genre VARCHAR(100),
  word_count INTEGER,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Bounties
CREATE TABLE bounties (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  product_id UUID REFERENCES products(id),
  bounty_points INTEGER NOT NULL,
  bounty_cash DECIMAL(10,2),
  status VARCHAR(50) DEFAULT 'open',
  requirements TEXT,
  claimed_by UUID REFERENCES users(id),
  claimed_at TIMESTAMP,
  completed_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Reviews
CREATE TABLE reviews (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  bounty_id UUID REFERENCES bounties(id),
  reviewer_id UUID REFERENCES users(id),
  rating INTEGER CHECK (rating >= 1 AND rating <= 5),
  title VARCHAR(500),
  content TEXT,
  word_count INTEGER,
  verified_purchase BOOLEAN DEFAULT FALSE,
  status VARCHAR(50) DEFAULT 'draft',
  amazon_review_url TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Point Transactions
CREATE TABLE point_transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  amount INTEGER NOT NULL,
  type VARCHAR(50) NOT NULL,
  reference_type VARCHAR(50),
  reference_id UUID,
  description VARCHAR(255),
  created_at TIMESTAMP DEFAULT NOW()
);

-- Follows
CREATE TABLE follows (
  follower_id UUID REFERENCES users(id),
  following_id UUID REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (follower_id, following_id)
);

-- Comments
CREATE TABLE comments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  review_id UUID REFERENCES reviews(id),
  parent_id UUID REFERENCES comments(id),
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Payments
CREATE TABLE payments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  stripe_session_id VARCHAR(255),
  stripe_payment_intent VARCHAR(255),
  amount_cents INTEGER NOT NULL,
  currency VARCHAR(3) DEFAULT 'usd',
  type VARCHAR(50),
  status VARCHAR(50),
  created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_bounties_status ON bounties(status);
CREATE INDEX idx_bounties_user ON bounties(user_id);
CREATE INDEX idx_reviews_bounty ON reviews(bounty_id);
CREATE INDEX idx_reviews_reviewer ON reviews(reviewer_id);
CREATE INDEX idx_point_transactions_user ON point_transactions(user_id);
```

---

## Local Development

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Node.js 18+
- Terraform

### Setup

```bash
# Start local infrastructure (PostgreSQL, Redis)
docker-compose up -d

# Start Go API
cd api && go run cmd/api/main.go

# Start React frontend
cd frontend && npm run dev
```

---

## Next Steps

1. Initialize Go API project
2. Initialize React frontend project
3. Set up Terraform infrastructure
4. Build authentication
5. Implement core features
6. Deploy
