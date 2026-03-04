# RevUExchange Backlog

Features and bug fixes to implement. Will be organized into future sprints.

---

## Critical (Must Fix) - IN PROGRESS

### Authentication
- [x] User gets logged out on page refresh - FIXED (MSW disabled by default)
- [x] MSW should NOT be enabled by default - FIXED

### Bounty Marketplace
- [x] Create Bounty button does nothing - FIXED (modal + API connection)
- [x] Claim Bounty button not connected - ALREADY CONNECTED
- [x] No bounty detail page - FIXED (added BountyDetail page)

### Reviews
- [x] My Reviews page not connected to claimed bounties - FIXED (connected to API)
- [ ] No review detail page (lower priority - reviews shown inline)

---

## High Priority

### Social Features
- [x] Follow/Unfollow buttons - ALREADY CONNECTED TO API (was done)
- [x] Comments UI on reviews - Already in MyReviews (inline)

### Products
- [x] Products page UI - IMPLEMENTED (Products.tsx + route)

### Gamification
- [x] Badges display - ALREADY CONNECTED TO API (was done)
- [x] Streak indicator - IMPLEMENTED (added to profile)

---

## Medium Priority

### Payments
- [ ] Points purchase - mock only, not connected to Stripe

### Analytics
- [ ] Analytics page shows placeholder data, not real metrics

### Profile
- [ ] Avatar upload - not implemented

---

## Lower Priority

### Backend Stubs (need real implementation)
- [ ] Analytics service - real queries
- [ ] Fraud detection service - real logic
- [ ] Email service - real SendGrid integration
- [ ] Genre matching - real algorithm

### Frontend Polish
- [ ] Loading states on all pages
- [ ] Error boundaries
- [ ] Toast notifications
- [ ] 404 page

### Testing
- [ ] Integration tests for handlers
- [ ] E2E tests for critical flows

---

## Ideas / Future

- [ ] Mobile app (React Native)
- [ ] Real-time notifications (WebSocket)
- [ ] Social sharing
- [ ] Content moderation
- [ ] Multi-language support
- [ ] Advanced fraud detection (ML)

---

## Bug Reports

| Date | Issue | Status |
|------|-------|--------|
| 2026-03-03 | User logged out on refresh | Fixed - MSW disabled |
| 2026-03-03 | Create Bounty button does nothing | Fixed - added modal |
| 2026-03-03 | No bounty detail page | Fixed - added page |
| 2026-03-03 | Make local-api broken (wrong path) | Fixed |

---

## Sprint 33 Progress

Completed in this sprint:
1. Disable MSW by default (fixes logout issue)
2. Fix Makefile paths (cmd/server -> cmd/api)
3. Add Create Bounty modal with API connection
4. Add BountyDetail page with route and claim functionality

---
