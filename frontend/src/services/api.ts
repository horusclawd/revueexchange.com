import axios from 'axios'
import type { User, Bounty, PointTransaction } from '../types'

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

  // Points
  async getBalance() {
    const { data } = await client.get<{ data: { points: number } }>('/v1/points/balance')
    return data.data.points
  },

  async getTransactions(params?: { limit?: number; offset?: number }) {
    const { data } = await client.get<{ data: PointTransaction[]; meta: { total: number } }>('/v1/points/transactions', { params })
    return { transactions: data.data, meta: data.meta }
  },
}
