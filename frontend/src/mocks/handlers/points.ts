import { http, HttpResponse, delay } from 'msw'
import { mockTransactions, currentUser } from '../data'

export const pointsHandlers = [
  // GET /api/v1/points/balance
  http.get('/api/v1/points/balance', async () => {
    await delay(100)
    return HttpResponse.json({
      data: {
        points: currentUser.points,
        lifetime_earned: currentUser.points + 150,
        lifetime_spent: 50,
      },
    })
  }),

  // GET /api/v1/points/transactions
  http.get('/api/v1/points/transactions', async ({ request }) => {
    await delay(150)
    const url = new URL(request.url)
    const limit = parseInt(url.searchParams.get('limit') || '20')
    const offset = parseInt(url.searchParams.get('offset') || '0')

    const transactions = mockTransactions.slice(offset, offset + limit)

    return HttpResponse.json({
      data: transactions,
      meta: {
        total: mockTransactions.length,
        page: Math.floor(offset / limit) + 1,
        page_size: limit,
      },
    })
  }),

  // POST /api/v1/points/transfer
  http.post('/api/v1/points/transfer', async ({ request }) => {
    await delay(300)
    const body = await request.json() as { user_id: string; amount: number }

    if (body.amount <= 0) {
      return HttpResponse.json({ error: 'Amount must be positive' }, { status: 400 })
    }

    if (currentUser.points < body.amount) {
      return HttpResponse.json({ error: 'Insufficient points' }, { status: 400 })
    }

    const newTransaction = {
      id: 'tx-' + Date.now(),
      user_id: body.user_id,
      amount: -body.amount,
      type: 'transferred',
      description: `Transfer to user ${body.user_id}`,
      created_at: new Date().toISOString(),
    }

    return HttpResponse.json({
      data: {
        success: true,
        transaction: newTransaction,
        new_balance: currentUser.points - body.amount,
      },
    })
  }),

  // GET /api/v1/points/leaderboard
  http.get('/api/v1/points/leaderboard', async () => {
    await delay(150)
    return HttpResponse.json({
      data: [
        { rank: 1, user_id: 'user-2', username: 'janedoe', display_name: 'Jane Doe', points: 500 },
        { rank: 2, user_id: 'user-1', username: 'johndoe', display_name: 'John Doe', points: 250 },
        { rank: 3, user_id: 'user-3', username: 'alexr', display_name: 'Alex R', points: 175 },
        { rank: 4, user_id: 'user-4', username: 'sarahw', display_name: 'Sarah W', points: 120 },
        { rank: 5, user_id: 'user-5', username: 'mikeb', display_name: 'Mike B', points: 85 },
      ],
    })
  }),
]
