# RevUExchange - Full Sprint Plan (Zero to MVP)

## Product Overview

**RevUExchange** - A reciprocal review platform for self-published authors and digital creators.

---

## Individual Sprints

---

### Sprint 1: Terraform Infrastructure ✅
**Duration**: ~14 hours

**Goal**: Set up complete Terraform infrastructure for AWS

**Deliverables**:
- ✅ Core module (VPC, IAM, Cognito, KMS, Secrets) - *replaced with separate modules pattern*
- ✅ Database module (Aurora PostgreSQL, ElastiCache Redis)
- ✅ Storage module (S3 buckets)
- ✅ Events module (EventBridge, SQS) - *simplified, combined with other modules*
- ✅ API module (ECS Fargate, API Gateway)
- ✅ CDN module (CloudFront, S3 frontend)
- ✅ Dev environment configuration
- ✅ Prod environment configuration
- ✅ Route53 module (uses existing hosted zone)
- ✅ Resource naming: `revueexchange-{env}-{resource}`

---

### Sprint 2: Local Development Environment
**Duration**: ~8 hours

**Goal**: Set up local development with Docker, LocalStack

**Deliverables**:
- ❗️ docker-compose.yml (LocalStack, PostgreSQL, Redis)
- ❗️ start-localstack.sh script
- ❗️ deploy-localstack.sh script
- ❗️ stop-localstack.sh script
- ❗️ Makefile with local commands
- ❗️ Test LocalStack deployment

---

### Sprint 3: Go API Project Setup
**Duration**: ~10 hours

**Goal**: Initialize Go API with project structure and base config

**Deliverables**:
- ❗️ Go module initialization
- ❗️ Project directory structure (cmd, internal, pkg, migrations)
- ❗️ Configuration from environment variables
- ❗️ Structured logging (zerolog)
- ❗️ Database connection (pgx)
- ❗️ Router setup (chi)
- ❗️ Health check endpoint
- ❗️ Basic error handling

---

### Sprint 4: React Frontend Setup
**Duration**: ~10 hours

**Goal**: Initialize React frontend with Vite, TypeScript, Tailwind

**Deliverables**:
- ❗️ Vite + React + TypeScript project
- ❗️ Tailwind CSS configuration
- ❗️ Project structure (components, pages, hooks, services, types)
- ❗️ React Query setup
- ❗️ Base UI components (Button, Input, Card)
- ❗️ API client with axios

---

### Sprint 5: Mock API for UI Development
**Duration**: ~12 hours

**Goal**: Create MSW mock handlers for UI development without backend

**Deliverables**:
- ✅ MSW setup and configuration (msw init, worker setup in main.tsx)
- ✅ Mock handlers for:
  - ✅ Auth (register, login, me, logout)
  - ✅ Users (list, get, update, profile)
  - ✅ Bounties (list, get, create, update, delete, claim)
  - ❌ Reviews (not implemented in mock - backend pending)
  - ✅ Points (balance, transactions, transfer, leaderboard)
  - ❌ Payments (not implemented in mock - backend pending)
  - ❌ Social (not implemented in mock - backend pending)
  - ❌ Gamification (not implemented in mock - backend pending)
- ✅ Mock data (users, bounties, transactions)
- ✅ Environment toggle (DEV mode uses MSW automatically)

---

### Sprint 6: User Authentication (Backend)
**Duration**: ~12 hours

**Goal**: Implement user registration, login, JWT tokens

**Deliverables**:
- ❗️ User model and repository
- ❗️ Password hashing (bcrypt)
- ❗️ JWT token generation and validation
- ❗️ Register endpoint (POST /api/v1/auth/register)
- ❗️ Login endpoint (POST /api/v1/auth/login)
- ❗️ Get current user endpoint (GET /api/v1/auth/me)
- ❗️ Refresh token endpoint
- ❗️ Logout functionality
- ❗️ Auth middleware

---

### Sprint 7: User Authentication (Frontend)
**Duration**: ~8 hours

**Goal**: Build auth UI and integrate with backend

**Deliverables**:
- ✅ Auth context (React) - AuthContext with login, register, logout
- ✅ Register page - Full form with email, username, password
- ✅ Login page - Full form with email, password
- ✅ Protected routes - ProtectedRoute component in App.tsx
- ✅ Token storage (localStorage) - Token stored on login/register
- ✅ API client with auth headers - Axios interceptor for Bearer token
- ✅ Logout functionality - Clears token and user state
- ✅ Auth guards - Redirect to login if not authenticated

---

### Sprint 8: User Profiles
**Duration**: ~8 hours

**Goal**: User profile management

**Deliverables**:
- ✅ Get user profile endpoint (GET /api/v1/users/{id}) - Already exists
- ✅ Update user profile endpoint (PUT /api/v1/users/{id}) - Already exists
- ❌ Avatar upload functionality - Deferred (requires S3 integration)
- ✅ Profile page UI - Profile page with view/edit modes
- ✅ Edit profile UI - Form with display name and bio

---

### Sprint 9: Products Management
**Duration**: ~10 hours

**Goal**: Support for books, courses, podcasts, newsletters

**Deliverables**:
- ✅ Product model (type, title, description, URL, cover, genre) - Already exists
- ✅ Create product endpoint (POST /api/v1/products) - Implemented
- ✅ Get product endpoint (GET /api/v1/products/{id}) - Implemented
- ✅ List user's products (GET /api/v1/users/{id}/products) - Implemented
- ✅ Product types: book, course, podcast, newsletter - Supported
- ✅ Genre/tags support - Genre field in model

---

### Sprint 10: Bounty Marketplace - Backend
**Duration**: ~14 hours

**Goal**: Create, list, claim bounties

**Deliverables**:
- ✅ Bounty model - Already exists
- ✅ Create bounty endpoint (POST /api/v1/bounties) - Already exists
- ✅ List bounties endpoint (GET /api/v1/bounties) - Updated with filters
  - ✅ Filters: genre, type, status, min/max points - Implemented
  - ✅ Pagination - Implemented
- ✅ Get bounty detail (GET /api/v1/bounties/{id}) - Already exists
- ✅ Claim bounty endpoint (POST /api/v1/bounties/{id}/claim) - Implemented
- ✅ Cancel bounty endpoint - Implemented
- ✅ Anti-swap protection logic - Implemented

---

### Sprint 11: Bounty Marketplace - Frontend
**Duration**: ~10 hours

**Goal**: Build bounty UI

**Deliverables**:
- ✅ Bounties list page - Creative redesign with gradient hero, cards
- ❌ Bounty filters (genre, type, status) - Partial (status filter chips)
- ❌ Bounty detail page - Deferred
- ❌ Create bounty form - Deferred
- ✅ Claim bounty button - Implemented
- ❌ My bounties page (as author) - Deferred
- ❌ My claimed bounties page (as reviewer) - Deferred

---

### Sprint 12: Reviews System - Backend
**Duration**: ~12 hours

**Goal**: Submit and manage reviews

**Deliverables**:
- ✅ Review model - Already exists
- ✅ Create review endpoint (POST /api/v1/reviews) - Implemented
- ✅ Get review endpoint (GET /api/v1/reviews/{id}) - Implemented
- ✅ Update review endpoint (PUT /api/v1/reviews/{id}) - Implemented
- ✅ Submit review endpoint (POST /api/v1/reviews/{id}/submit) - Implemented
- ✅ Word count validation - Implemented (min 100 words)
- ✅ Rating system (1-5 stars) - Implemented
- ✅ Review status (draft, submitted, published) - Implemented

---

### Sprint 13: Reviews System - Frontend
**Duration**: ~8 hours

**Goal**: Build review UI

**Deliverables**:
- ✅ Review form (draft mode) - MyReviews page with form
- ❌ Review detail page - Deferred
- ❌ Edit review functionality - Deferred
- ✅ Submit review flow - Implemented
- ✅ Rating component - StarRating component
- ✅ Word count display - With 100 word minimum

---

### Sprint 14: Points System - Backend
**Duration**: ~10 hours

**Goal**: Point earning, spending, transactions

**Deliverables**:
- ✅ Point transaction model - Already exists
- ✅ Award points on review submission - In PointsService
- ✅ Deduct points on bounty claim - In PointsService
- ✅ Get balance endpoint (GET /api/v1/points/balance) - Implemented
- ✅ Get transactions endpoint (GET /api/v1/points/transactions) - Implemented
- ✅ Point transfer between users - Implemented (POST /api/v1/points/transfer)
- ❌ Bonus/penalty system - Not implemented

---

### Sprint 15: Points System - Frontend
**Duration**: ~6 hours

**Goal**: Build points UI

**Deliverables**:
- ✅ Points balance display (header/nav) - Linked in header with teal badge
- ✅ Transaction history page - Full page with filters
- ✅ Points earned/spent visualization - Stats cards and totals

---

### Sprint 16: Payments (Stripe) - Backend
**Duration**: ~12 hours

**Goal**: Integrate Stripe for payments

**Deliverables**:
- ✅ Stripe SDK integration - Added config for Stripe keys
- ✅ Create checkout session (POST /api/v1/payments/checkout) - Returns session URL and points award
- ✅ Stripe webhook handler (POST /api/v1/payments/webhook) - Processes checkout events
- ✅ Payment model - Already exists
- ✅ Convert cash to points logic - 100 points per dollar
- ❌ Refund handling - Not implemented (simplified)
- ❌ Stripe webhook signature verification - Skipped for MVP (would add in production)

---

### Sprint 17: Payments (Stripe) - Frontend
**Duration**: ~8 hours

**Goal**: Build payment UI

**Deliverables**:
- ✅ Points purchase page - Purple/violet gradient design, purchase packages
- ❌ Stripe Elements integration - Skipped (simplified for MVP)
- ✅ Checkout flow - Simulated checkout with success modal
- ✅ Success page - Modal with celebration animation
- ❌ Cancel page - Not needed (simplified flow)
- ✅ Purchase history - Display in purchase page

---

### Sprint 18: Social Features - Backend
**Duration**: ~12 hours

**Goal**: Follows, activity feed, comments

**Deliverables**:
- ✅ Follow model - Already exists
- ✅ Follow user endpoint (POST /api/v1/social/follow/{id}) - Implemented
- ✅ Unfollow user endpoint (DELETE /api/v1/social/follow/{id}) - Implemented
- ✅ Get followers endpoint (GET /api/v1/social/followers/{id}) - Implemented
- ✅ Get following endpoint (GET /api/v1/social/following/{id}) - Implemented
- ✅ Activity feed endpoint (GET /api/v1/social/feed) - Implemented
- ✅ Comment model - Added
- ✅ Add comment endpoint - POST /api/v1/comments
- ✅ Delete comment endpoint - DELETE /api/v1/comments/{id}

---

### Sprint 19: Social Features - Frontend
**Duration**: ~8 hours

**Goal**: Build social UI

**Deliverables**:
- ❗️ User profile page
- ❗️ Follow/unfollow button
- ❗️ Followers/following lists
- ❗️ Activity feed page
- ❗️ Comments on reviews
- ❗️ Add/delete comments

---

### Sprint 20: Gamification (Badges) - Backend
**Duration**: ~10 hours

**Goal**: Badge system using DynamoDB

**Deliverables**:
- ❗️ Badge definitions (types, tiers)
- ❗️ Award badge logic
- ❗️ Check badge conditions:
  - ❗️ First review
  - ❗️ 10 reviews
  - ❗️ 50 reviews
  - ❗️ Top reviewer
  - ❗️ Streak milestones
- ❗️ Get user badges endpoint

---

### Sprint 21: Gamification (Leaderboard & Streaks) - Backend
**Duration**: ~8 hours

**Goal**: Leaderboards and streaks

**Deliverables**:
- ❗️ Leaderboard model in DynamoDB
- ❗️ Update leaderboard on points change
- ❗️ Get leaderboard endpoint (GET /api/v1/gamification/leaderboard)
- ❗️ Streak model in DynamoDB
- ❗️ Update streak on activity
- ❗️ Get streak endpoint

---

### Sprint 22: Gamification - Frontend
**Duration**: ~8 hours

**Goal**: Build gamification UI

**Deliverables**:
- ❗️ Badges display on profile
- ❗️ Badge modal/details
- ❗️ Leaderboard page
- ❗️ Rankings display
- ❗️ Streak indicator
- ❗️ Achievement notifications

---

### Sprint 23: Analytics Dashboard - Backend
**Duration**: ~12 hours

**Goal**: Analytics and insights

**Deliverables**:
- ❗️ Analytics aggregation service
- ❗️ Overview stats endpoint
- ❗️ Bounty performance metrics
- ❗️ Review metrics (views, helpful)
- ❗️ Revenue stats
- ❗️ User activity tracking
- ❗️ Daily/weekly/monthly aggregations

---

### Sprint 24: Analytics Dashboard - Frontend
**Duration**: ~10 hours

**Goal**: Build analytics UI

**Deliverables**:
- ❗️ Dashboard overview page
- ❗️ Charts (views over time)
- ❗️ Bounty performance table
- ❗️ Revenue analytics
- ❗️ Review quality metrics
- ❗️ Export data functionality

---

### Sprint 25: Anti-Fraud System
**Duration**: ~12 hours

**Goal**: Prevent gaming the system

**Deliverables**:
- ❗️ Review quality scoring
- ❗️ Suspicious activity detection
- ❗️ Rate limiting (per user, per endpoint)
- ❗️ IP fingerprinting
- ❗️ Device fingerprinting
- ❗️ Manual review queue
- ❗️ Report review endpoint
- ❗️ Flagged reviews handling

---

### Sprint 26: Email Notifications (SendGrid)
**Duration**: ~8 hours

**Goal**: Transactional emails

**Deliverables**:
- ❗️ SendGrid integration
- ❗️ Welcome email
- ❗️ Bounty claimed notification
- ❗️ Review submitted notification
- ❗️ Points awarded notification
- ❗️ Follower notification
- ❗️ Email templates
- ❗️ Email queue (async sending)

---

### Sprint 27: Genre & Expertise Matching
**Duration**: ~8 hours

**Goal**: Match reviewers to bounties

**Deliverables**:
- ❗️ User genre preferences
- ❗️ Product genre tags
- ❗️ Matching algorithm
- ❗️ Suggested bounties endpoint
- ❗️ Genre-based recommendations

---

### Sprint 28: Polish & Error Handling
**Duration**: ~12 hours

**Goal**: Final polish and robustness

**Deliverables**:
- ❗️ Loading states (all pages)
- ❗️ Error boundaries
- ❗️ Toast notifications
- ❗️ Form validation (frontend + backend)
- ❗️ Input sanitization
- ❗️ Global error handler
- ❗️ 404 pages
- ❗️ Empty states
- ❗️ Responsive design fixes

---

### Sprint 29: Testing
**Duration**: ~12 hours

**Goal**: Test coverage

**Deliverables**:
- ❗️ Unit tests for services
- ❗️ Integration tests for handlers
- ❗️ Frontend component tests
- ❗️ E2E tests (critical flows)
- ❗️ Auth flow tests
- ❗️ Bounty → Claim → Review flow tests

---

### Sprint 30: CI/CD Pipeline
**Duration**: ~10 hours

**Goal**: Automated deployments

**Deliverables**:
- ❗️ GitHub Actions workflow
- ❗️ Build Go API
- ❗️ Build React frontend
- ❗️ Run tests
- ❗️ Deploy to dev (auto)
- ❗️ Deploy to staging (on merge)
- ❗️ Deploy to prod (manual approval)
- ❗️ Database migration runner
- ❗️ Rollback procedure

---

### Sprint 31: Production Deployment
**Duration**: ~12 hours

**Goal**: Go live

**Deliverables**:
- ❗️ Production Terraform apply
- ❗️ Database migrations
- ❗️ Domain registration/setup
- ❗️ SSL certificates (ACM)
- ❗️ DNS configuration (Route53)
- ❗️ CloudFront distribution
- ❗️ Production environment variables
- ❗️ Health checks
- ❗️ Monitoring setup (CloudWatch)
- ❗️ Alerts (error rates, latency)
- ❗️ Log aggregation

---

### Sprint 32: Launch & Handoff
**Duration**: ~6 hours

**Goal**: Launch preparation

**Deliverables**:
- ❗️ Launch checklist
- ❗️ Runbook documentation
- ❗️ On-call rotation setup
- ❗️ Incident response plan
- ❗️ Feature flag list
- ❗️ Analytics tracking (internal)
- ❗️ Social media assets
- ❗️ Press release (optional)

---

## Summary

| Sprint | Name | Status |
|--------|------|--------|
| 1 | Terraform Infrastructure | ✅ |
| 2 | Local Development Environment | ❗️ |
| 3 | Go API Project Setup | ❗️ |
| 4 | React Frontend Setup | ✅ |
| 5 | Mock API for UI Development | ✅ |
| 6 | User Authentication (Backend) | ✅ |
| 7 | User Authentication (Frontend) | ✅ |
| 8 | User Profiles | ✅ |
| 9 | Products Management | ✅ |
| 10 | Bounty Marketplace (Backend) | ✅ |
| 11 | Bounty Marketplace (Frontend) | ✅ |
| 12 | Reviews System (Backend) | ✅ |
| 13 | Reviews System (Frontend) | ✅ |
| 14 | Points System (Backend) | ❗️ |
| 15 | Points System (Frontend) | ❗️ |
| 16 | Payments (Stripe Backend) | ❗️ |
| 17 | Payments (Stripe Frontend) | ❗️ |
| 18 | Social Features (Backend) | ❗️ |
| 19 | Social Features (Frontend) | ❗️ |
| 20 | Gamification (Badges Backend) | ❗️ |
| 21 | Gamification (Leaderboard & Streaks) | ❗️ |
| 22 | Gamification (Frontend) | ❗️ |
| 23 | Analytics Dashboard (Backend) | ❗️ |
| 24 | Analytics Dashboard (Frontend) | ❗️ |
| 25 | Anti-Fraud System | ❗️ |
| 26 | Email Notifications | ❗️ |
| 27 | Genre & Expertise Matching | ❗️ |
| 28 | Polish & Error Handling | ❗️ |
| 29 | Testing | ❗️ |
| 30 | CI/CD Pipeline | ❗️ |
| 31 | Production Deployment | ❗️ |
| 32 | Launch & Handoff | ❗️ |

**Completed**: 13/32 sprints

---

## Out of Scope (Post-MVP)

- Mobile apps (iOS/Android)
- White-label/enterprise
- AI-powered recommendations
- Marketplace for other services
- Subscription tiers beyond free
- API for third-parties
