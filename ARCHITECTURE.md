# RevUExchange - Architecture Plan

## Quick Summary

| Aspect | Option A | Option B |
|--------|----------|----------|
| **Frontend** | Next.js 14 (React) | React + Vite |
| **Backend API** | Next.js API Routes / Lambda | Go (Golang) |
| **Styling** | Tailwind CSS | Tailwind CSS |
| **Database** | Aurora PostgreSQL | Aurora PostgreSQL |
| **Infrastructure** | AWS Serverless | AWS Containers |

---

## Hard Requirements

- Runs on AWS
- Scalable architecture
- Clear distinction between UI and API
- Uses API Gateway and/or messaging (SQS, EventBridge)

---

# Option A: Serverless (Next.js + Lambda)

## Stack

| Layer | Technology |
|-------|------------|
| **Frontend** | Next.js 14 (React) - TypeScript |
| **Backend** | AWS Lambda - Node.js/TypeScript |
| **API** | API Gateway REST |
| **Database** | Aurora PostgreSQL (Serverless v2) |
| **Cache** | ElastiCache Redis |
| **NoSQL** | DynamoDB |
| **Storage** | S3 |
| **Auth** | Amazon Cognito + JWT |
| **Messaging** | EventBridge + SQS |
| **Payments** | Stripe |
| **Email** | SendGrid |
| **IaC** | AWS CDK (TypeScript) |
| **CI/CD** | GitHub Actions |

### Monthly Cost: ~$134/month

---

# Option B: Go API + React Frontend (Recommended)

## Stack

| Layer | Technology |
|-------|------------|
| **Frontend** | React 18 + Vite |
| **Styling** | Tailwind CSS |
| **Backend API** | Go (Golang) |
| **API Gateway** | AWS API Gateway (HTTP API) |
| **Compute** | AWS Lambda OR ECS Fargate |
| **Database** | Aurora PostgreSQL |
| **Cache** | ElastiCache Redis |
| **NoSQL** | DynamoDB |
| **Storage** | S3 |
| **Auth** | Amazon Cognito + JWT |
| **Messaging** | EventBridge + SQS |
| **Payments** | Stripe |
| **Email** | SendGrid |
| **IaC** | AWS CDK (TypeScript) |
| **CI/CD** | GitHub Actions |

### Monthly Cost: ~$150-200/month

---

## Why Go?

| Factor | Benefit |
|--------|---------|
| **Performance** | 10-20x faster than Node.js for compute-heavy tasks |
| **Binary size** | Small Lambda packages (5-10MB vs 50MB+ Node.js) |
| **Cold starts** | Much faster startup (10-50ms vs 100-500ms) |
| **Concurrency** | Native goroutines for parallel processing |
| **Type safety** | Static typing catches errors at compile time |
| **Production** | Used by Google, Uber, Twitch, CloudFlare |

---

## Architecture Diagram (Option B)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            AWS Cloud                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                    CloudFront (Global CDN)                             │   │
│  │               200+ Edge Locations Worldwide                          │   │
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
│  │   Lambda (Go)   │    │  ECS Fargate    │    │   Lambda (Go)   │        │
│  │   - Auth        │    │   - Core API    │    │   - Async       │        │
│  │   - Lightweight │    │   - Business    │    │   - Workers     │        │
│  │                 │    │     Logic       │    │                 │        │
│  └─────────────────┘    └────────┬────────┘    └─────────────────┘        │
│                                   │                                          │
├───────────────────────────────────┼──────────────────────────────────────────┤
│                                   ▼                                          │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                     Amazon EventBridge                                │   │
│  │   user.registered  │  review.submitted  │  bounty.created           │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│           │                        │                        │                │
│           ▼                        ▼                        ▼                │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐        │
│  │   Lambda:      │    │   Lambda:      │    │   Lambda:      │        │
│  │  Notifications │    │  Analytics     │    │  Workflows     │        │
│  │  (SendGrid)    │    │  (Aggregations)│    │  (State Mach.)│        │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘        │
│                                                                              │
  ┌────────────────────────────────│──────────────────────────────────────┐   │
│  │                     Amazon SQS                                        │   │
│  │   email-queue  │  webhook-queue  │  export-queue                   │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                              Data Layer                                       │
│                                                                              │
│  ┌─────────────────────┐  ┌─────────────────────┐                          │
│  │  Aurora PostgreSQL  │  │  ElastiCache Redis  │                          │
│  │   (Serverless v2)   │  │   (Cluster Mode)    │                          │
│  │                     │  │                      │                          │
│  │  Primary Data:      │  │  • Sessions         │                          │
│  │  • Users            │  │  • API Cache         │                          │
│  │  • Bounties         │  │  • Rate Limiting    │                          │
│  │  • Reviews         │  │  • Real-time stats  │                          │
│  │  • Point TXs       │  │  • Leaderboard      │                          │
│  │  • Follows          │  │                      │                          │
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
  /repository      # Database queries
  /service          # Business logic
  /transport        # HTTP, GRPC, events
/pkg
  /logger           # Structured logging
  /errors           # Error handling
  /auth             # JWT utilities
  /stripe           # Stripe client
/migrations         # Database migrations
```

### HTTP Handlers

```go
// Handlers structure
type Handler struct {
    userService    *service.UserService
    bountyService  *service.BountyService
    reviewService  *service.ReviewService
    pointsService  *service.PointsService
    paymentService *service.PaymentService
}

// Routes setup
func (h *Handler) Routes(router *chi.Mux) {
    // Auth
    router.Post("/api/v1/auth/register", h.Register)
    router.Post("/api/v1/auth/login", h.Login)

    // Users
    router.Get("/api/v1/users/{id}", h.GetUser)
    router.Put("/api/v1/users/{id}", h.UpdateUser)

    // Bounties
    router.Get("/api/v1/bounties", h.ListBounties)
    router.Post("/api/v1/bounties", h.CreateBounty)
    router.Get("/api/v1/bounties/{id}", h.GetBounty)
    router.Post("/api/v1/bounties/{id}/claim", h.ClaimBounty)

    // Reviews
    router.Post("/api/v1/reviews", h.CreateReview)
    router.Get("/api/v1/reviews/{id}", h.GetReview)
    router.Put("/api/v1/reviews/{id}", h.UpdateReview)
    router.Post("/api/v1/reviews/{id}/submit", h.SubmitReview)

    // Points
    router.Get("/api/v1/points/balance", h.GetBalance)
    router.Get("/api/v1/points/transactions", h.GetTransactions)

    // Payments
    router.Post("/api/v1/payments/checkout", h.CreateCheckout)
    router.Post("/api/v1/payments/webhook", h.HandleWebhook)
}
```

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

### Tailwind CSS Setup

```javascript
// tailwind.config.js
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f0f9ff',
          100: '#e0f2fe',
          500: '#0ea5e9',
          600: '#0284c7',
          700: '#0369a1',
        }
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}
```

### API Client

```typescript
// services/api.ts
const API_BASE = import.meta.env.VITE_API_URL;

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(endpoint: string, options: RequestInit): Promise<T> {
    const token = localStorage.getItem('token');
    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...options.headers,
      },
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`);
    }

    return response.json();
  }

  async getBounties(filters: BountyFilters): Promise<Bounty[]> {
    const params = new URLSearchParams(filters as any);
    return this.request(`/bounties?${params}`, { method: 'GET' });
  }

  async createBounty(data: CreateBountyDto): Promise<Bounty> {
    return this.request('/bounties', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }
}

export const api = new ApiClient(API_BASE);
```

---

## Database Schema (Both Options)

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

## Comparison

| Aspect | Option A (Next.js) | Option B (Go + React) |
|--------|-------------------|----------------------|
| **Frontend** | Next.js | React + Vite |
| **Backend** | Node.js Lambda | Go Lambda/Fargate |
| **Cold Start** | 100-500ms | 10-50ms |
| **Dev Experience** | Single repo, full-stack | Clear separation |
| **Learning Curve** | Lower | Higher (Go) |
| **Startup Time** | Faster to build | More setup |
| **Best For** | MVPs, solo devs | Production, teams |

---

## Recommendation

**Option B (Go + React)** is recommended because:

1. **Clear UI/API separation** - Your hard requirement
2. **Go performance** - Faster cold starts, better throughput
3. **React + Vite** - You didn't specify a frontend framework, this is the standard choice
4. **Tailwind CSS** - Most popular, efficient stylesheet framework
5. **TypeScript** - Full-stack type safety

---

## Next Steps

1. Choose Option A or B
2. Set up AWS account
3. Initialize project (CDK, frontend, backend)
4. Build authentication
5. Implement core features
6. Deploy
