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

export interface Payment {
  id: string
  user_id: string
  stripe_session_id?: string
  stripe_payment_intent?: string
  amount_cents: number
  currency: string
  type?: string
  status: 'pending' | 'completed' | 'failed' | 'refunded' | 'expired'
  created_at: string
}

export interface Comment {
  id: string
  user_id: string
  review_id: string
  parent_id?: string
  content: string
  created_at: string
  user?: User
}

export interface Activity {
  id: string
  user_id: string
  type: 'review_submitted' | 'bounty_created' | 'follow' | 'comment' | 'points_earned' | 'bounty_completed'
  reference_type?: string
  reference_id?: string
  created_at: string
  user?: User
}

export interface Badge {
  id: string
  user_id: string
  badge_type: string
  badge_name: string
  description: string
  tier: 'bronze' | 'silver' | 'gold' | 'platinum'
  icon_url?: string
  awarded_at: string
}

export interface LeaderboardEntry {
  user_id: string
  username: string
  display_name: string
  points: number
  review_count: number
  rank: number
  last_updated: string
}

export interface Streak {
  user_id: string
  current_streak: number
  longest_streak: number
  last_activity_at: string
  streak_started_at: string
}

export interface ApiResponse<T> {
  data?: T
  error?: string
  message?: string
}
