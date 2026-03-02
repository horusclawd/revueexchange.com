export interface User {
  id: string
  email: string
  username: string
  display_name: string
  avatar_url?: string
  bio?: string
  points: number
  reputation_score: number
  subscription_tier: string
  created_at: string
  updated_at: string
}

export interface Product {
  id: string
  user_id: string
  type: 'book' | 'course' | 'podcast' | 'newsletter'
  title: string
  description?: string
  url?: string
  cover_image_url?: string
  genre?: string
  word_count?: number
  created_at: string
}

export interface Bounty {
  id: string
  user_id: string
  product_id: string
  bounty_points: number
  bounty_cash?: number
  status: 'open' | 'claimed' | 'under_review' | 'completed' | 'cancelled'
  requirements?: string
  claimed_by?: string
  claimed_at?: string
  completed_at?: string
  created_at: string
  updated_at: string
}

export interface Review {
  id: string
  bounty_id: string
  reviewer_id: string
  rating: number
  title?: string
  content?: string
  word_count?: number
  verified_purchase: boolean
  status: 'draft' | 'submitted' | 'published'
  amazon_review_url?: string
  created_at: string
  updated_at: string
}

export interface PointTransaction {
  id: string
  user_id: string
  amount: number
  type: 'earned' | 'spent' | 'bonus' | 'penalty' | 'refund'
  reference_type?: string
  reference_id?: string
  description?: string
  created_at: string
}

export interface ApiResponse<T> {
  data?: T
  error?: string
  message?: string
}
