import { User, Bounty, PointTransaction } from '../types'

export const mockUsers: User[] = [
  {
    id: 'user-1',
    email: 'john@example.com',
    username: 'johndoe',
    display_name: 'John Doe',
    avatar_url: undefined,
    bio: 'Author of sci-fi novels',
    points: 250,
    reputation_score: 4.5,
    subscription_tier: 'free',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 'user-2',
    email: 'jane@example.com',
    username: 'janedoe',
    display_name: 'Jane Doe',
    avatar_url: undefined,
    bio: 'Writing coach and reviewer',
    points: 500,
    reputation_score: 4.8,
    subscription_tier: 'premium',
    created_at: '2024-01-15T00:00:00Z',
    updated_at: '2024-01-15T00:00:00Z',
  },
]

export const mockBounties: Bounty[] = [
  {
    id: 'bounty-1',
    user_id: 'user-1',
    product_id: 'product-1',
    bounty_points: 50,
    bounty_cash: 10,
    status: 'open',
    requirements: 'Honest review of at least 500 words',
    created_at: '2024-01-20T10:00:00Z',
    updated_at: '2024-01-20T10:00:00Z',
  },
  {
    id: 'bounty-2',
    user_id: 'user-2',
    product_id: 'product-2',
    bounty_points: 75,
    bounty_cash: 15,
    status: 'open',
    requirements: 'Detailed review covering plot, characters, and writing style',
    created_at: '2024-01-21T10:00:00Z',
    updated_at: '2024-01-21T10:00:00Z',
  },
  {
    id: 'bounty-3',
    user_id: 'user-1',
    product_id: 'product-3',
    bounty_points: 100,
    status: 'claimed',
    claimed_by: 'user-2',
    claimed_at: '2024-01-22T10:00:00Z',
    requirements: 'Review for my new fantasy novel',
    created_at: '2024-01-19T10:00:00Z',
    updated_at: '2024-01-22T10:00:00Z',
  },
]

export const mockTransactions: PointTransaction[] = [
  {
    id: 'tx-1',
    user_id: 'user-1',
    amount: 50,
    type: 'earned',
    reference_type: 'review',
    reference_id: 'review-1',
    description: 'Points for submitted review',
    created_at: '2024-01-20T10:00:00Z',
  },
  {
    id: 'tx-2',
    user_id: 'user-1',
    amount: -25,
    type: 'spent',
    reference_type: 'bounty',
    reference_id: 'bounty-1',
    description: 'Bounty creation',
    created_at: '2024-01-19T10:00:00Z',
  },
  {
    id: 'tx-3',
    user_id: 'user-1',
    amount: 100,
    type: 'bonus',
    description: 'Welcome bonus',
    created_at: '2024-01-01T10:00:00Z',
  },
]

export const currentUser = mockUsers[0]
