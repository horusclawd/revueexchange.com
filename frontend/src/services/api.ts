import axios from 'axios'
import type { User, Bounty, PointTransaction, Payment, Comment, Activity, Badge, LeaderboardEntry, Streak } from '../types'

const API_BASE = import.meta.env.VITE_API_URL || '/api'

const client = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add auth token to requests
client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export const api = {
  // Auth
  async login(email: string, password: string) {
    const { data } = await client.post<{ data: { user: User; token: string } }>('/v1/auth/login', { email, password })
    return { user: data.data.user, token: data.data.token }
  },

  async register(email: string, username: string, password: string) {
    const { data } = await client.post<{ data: { user: User; token: string } }>('/v1/auth/register', { email, username, password })
    return { user: data.data.user, token: data.data.token }
  },

  async getMe() {
    const { data } = await client.get<{ data: User }>('/v1/auth/me')
    return data.data
  },

  // User
  async getUser(id: string) {
    const { data } = await client.get<{ data: User }>(`/v1/users/${id}`)
    return data.data
  },

  async updateUser(id: string, updates: Partial<User>) {
    const { data } = await client.put<{ data: User }>(`/v1/users/${id}`, updates)
    return data.data
  },

  // Bounties
  async getBounties(params?: { status?: string; genre?: string; type?: string }) {
    const { data } = await client.get<{ data: Bounty[]; meta: { total: number } }>('/v1/bounties', { params })
    return { bounties: data.data, meta: data.meta }
  },

  async getBounty(id: string) {
    const { data } = await client.get<{ data: Bounty }>(`/v1/bounties/${id}`)
    return data.data
  },

  async createBounty(bounty: Partial<Bounty>) {
    const { data } = await client.post<{ data: Bounty }>('/v1/bounties', bounty)
    return data.data
  },

  async claimBounty(id: string) {
    const { data } = await client.post<{ data: Bounty }>(`/v1/bounties/${id}/claim`)
    return data.data
  },

  // Reviews
  async createReview(review: { bounty_id: string; rating: number; title: string; content: string }) {
    const { data } = await client.post<{ data: unknown }>('/v1/reviews', review)
    return data.data
  },

  async getReview(id: string) {
    const { data } = await client.get<{ data: unknown }>(`/v1/reviews/${id}`)
    return data.data
  },

  async updateReview(id: string, updates: { rating?: number; title?: string; content?: string }) {
    const { data } = await client.put<{ data: unknown }>(`/v1/reviews/${id}`, updates)
    return data.data
  },

  async submitReview(id: string) {
    const { data } = await client.post<{ data: unknown }>(`/v1/reviews/${id}/submit`)
    return data.data
  },

  // Points
  async getBalance() {
    const { data } = await client.get<{ data: { points: number } }>('/v1/points/balance')
    return data.data.points
  },

  async getTransactions(params?: { limit?: number; offset?: number }) {
    const { data } = await client.get<{ data: PointTransaction[]; meta: { total: number } }>('/v1/points/transactions', { params })
    return { transactions: data.data, meta: data.meta }
  },

  // Payments
  async createCheckout(amountCents: number) {
    const { data } = await client.post<{ data: { checkout_url: string; points_award: number; payment_id: string } }>('/v1/payments/checkout', { amount_cents: amountCents })
    return data.data
  },

  async getPaymentHistory() {
    const { data } = await client.get<{ data: Payment[] }>('/v1/payments/history')
    return data.data
  },

  // Social - Follow
  async followUser(userId: string) {
    const { data } = await client.post<{ message: string }>(`/v1/social/follow/${userId}`)
    return data
  },

  async unfollowUser(userId: string) {
    const { data } = await client.delete<{ message: string }>(`/v1/social/follow/${userId}`)
    return data
  },

  async getFollowers(userId: string) {
    const { data } = await client.get<{ data: User[] }>(`/v1/social/followers/${userId}`)
    return data.data
  },

  async getFollowing(userId: string) {
    const { data } = await client.get<{ data: User[] }>(`/v1/social/following/${userId}`)
    return data.data
  },

  async isFollowing(userId: string) {
    const { data } = await client.get<{ data: { following: boolean } }>(`/v1/social/following/${userId}`)
    return data.data.following
  },

  // Social - Activity Feed
  async getActivityFeed() {
    const { data } = await client.get<{ data: Activity[] }>('/v1/social/feed')
    return data.data
  },

  // Social - Comments
  async addComment(reviewId: string, content: string, parentId?: string) {
    const { data } = await client.post<{ data: Comment }>('/v1/comments', {
      review_id: reviewId,
      content,
      parent_id: parentId,
    })
    return data.data
  },

  async deleteComment(commentId: string) {
    const { data } = await client.delete<{ message: string }>(`/v1/comments/${commentId}`)
    return data
  },

  async getComments(reviewId: string) {
    const { data } = await client.get<{ data: Comment[] }>(`/v1/comments?review_id=${reviewId}`)
    return data.data
  },

  // Gamification - Badges
  async getBadges() {
    const { data } = await client.get<{ data: Badge[] }>('/v1/badges')
    return data.data
  },

  async checkBadges() {
    const { data } = await client.post<{ data: Badge[]; message: string }>('/v1/badges/check')
    return data
  },

  // Gamification - Leaderboard
  async getLeaderboard(limit = 50) {
    const { data } = await client.get<{ data: LeaderboardEntry[] }>('/v1/gamification/leaderboard', { params: { limit } })
    return data.data
  },

  // Gamification - Streaks
  async getStreak() {
    const { data } = await client.get<{ data: Streak }>('/v1/gamification/streak')
    return data.data
  },

  async updateStreak() {
    const { data } = await client.post<{ data: Streak; message: string }>('/v1/gamification/streak/update')
    return data.data
  },
}
