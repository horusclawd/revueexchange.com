import { http, HttpResponse, delay } from 'msw'
import { mockUsers, currentUser } from '../data'

export const authHandlers = [
  // POST /api/v1/auth/login
  http.post('/api/v1/auth/login', async ({ request }) => {
    await delay(200)
    const body = await request.json() as { email: string; password: string }

    const user = mockUsers.find(u => u.email === body.email)

    if (!user || body.password !== 'password123') {
      return HttpResponse.json(
        { error: 'Invalid credentials' },
        { status: 401 }
      )
    }

    return HttpResponse.json({
      data: {
        token: 'mock-jwt-token-' + user.id,
        user,
      },
    })
  }),

  // POST /api/v1/auth/register
  http.post('/api/v1/auth/register', async ({ request }) => {
    await delay(200)
    const body = await request.json() as { email: string; username: string; password: string }

    const existingUser = mockUsers.find(u => u.email === body.email || u.username === body.username)

    if (existingUser) {
      return HttpResponse.json(
        { error: 'User already exists' },
        { status: 400 }
      )
    }

    const newUser = {
      id: 'user-' + Date.now(),
      email: body.email,
      username: body.username,
      display_name: body.username,
      avatar_url: null,
      bio: null,
      points: 100,
      reputation_score: 0,
      subscription_tier: 'free',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }

    return HttpResponse.json({
      data: {
        token: 'mock-jwt-token-' + newUser.id,
        user: newUser,
      },
    }, { status: 201 })
  }),

  // GET /api/v1/auth/me
  http.get('/api/v1/auth/me', async ({ request }) => {
    await delay(100)
    const authHeader = request.headers.get('Authorization')

    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return HttpResponse.json(
        { error: 'Unauthorized' },
        { status: 401 }
      )
    }

    return HttpResponse.json({ data: { user: currentUser } })
  }),

  // POST /api/v1/auth/logout
  http.post('/api/v1/auth/logout', async () => {
    await delay(100)
    return HttpResponse.json({ success: true })
  }),
]
