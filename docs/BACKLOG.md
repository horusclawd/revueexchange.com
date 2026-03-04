# RevUExchange Backlog

Features and bug fixes to implement. Will be organized into future sprints.

---

## Critical (Must Fix)

### Authentication
- [ ] User gets logged out on page refresh - investigate token persistence
- [ ] MSW should NOT be enabled by default (disable in main.tsx) - DONE

### Bounty Marketplace
- [ ] Create Bounty button does nothing - needs form/modal
- [ ] Claim Bounty button not connected to API
- [ ] No bounty detail page

### Reviews
- [ ] My Reviews page not connected to claimed bounties
- [ ] No review detail page

---

## High Priority

### Social Features
- [ ] Follow/Unfollow buttons not connected to API
- [ ] Comments UI on reviews - not implemented

### Products
- [ ] Products page UI - not implemented

### Gamification
- [ ] Badges display - not connected to API
- [ ] Streak indicator - not connected to API

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
| 2026-03-03 | User logged out on refresh | Investigating - likely MSW issue |
| 2026-03-03 | Create Bounty button does nothing | To Do |
| 2026-03-03 | Make local-api broken (wrong path) | Fixed |

---

## Notes

- Sprint plan updated 2026-03-03 to accurately reflect completion status
- Many features marked "done" are actually stubs or incomplete
- Focus should be on making core flows work first (register -> create bounty -> claim -> review)
