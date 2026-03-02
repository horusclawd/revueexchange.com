import { http, HttpResponse, delay } from 'msw'
import { mockUsers } from '../data'

export const userHandlers = [
  // GET /api/v1/users
  http.get('/api/v1/users', async ({ request }) => {
    await delay(150)
    const url = new URL(request.url)
    const username = url.searchParams.get('username')

    if (username) {
      const user = mockUsers.find(u => u.username === username)
      if (!user) {
        return HttpResponse.json({ error: 'User not found' }, { status: 404 })
      }
      return HttpResponse.json({ data: user })
    }

    return HttpResponse.json({ data: mockUsers })
  }),

  // GET /api/v1/users/:id
  http.get('/api/v1/users/:id', async ({ params }) => {
    await delay(100)
    const user = mockUsers.find(u => u.id === params.id)

    if (!user) {
      return HttpResponse.json({ error: 'User not found' }, { status: 404 })
    }

    return HttpResponse.json({ data: user })
  }),

  // PUT /api/v1/users/:id
  http.put('/api/v1/users/:id', async ({ params, request }) => {
    await delay(200)
    const user = mockUsers.find(u => u.id === params.id)

    if (!user) {
      return HttpResponse.json({ error: 'User not found' }, { status: 404 })
    }

    const body = await request.json() as Partial<typeof user>
    const updatedUser = { ...user, ...body, updated_at: new Date().toISOString() }

    return HttpResponse.json({ data: updatedUser })
  }),

  // GET /api/v1/users/:id/profile
  http.get('/api/v1/users/:id/profile', async ({ params }) => {
    await delay(100)
    const user = mockUsers.find(u => u.id === params.id)

    if (!user) {
      return HttpResponse.json({ error: 'User not found' }, { status: 404 })
    }

    return HttpResponse.json({
      data: {
        id: user.id,
        username: user.username,
        display_name: user.display_name,
        avatar_url: user.avatar_url,
        bio: user.bio,
        reputation_score: user.reputation_score,
        created_at: user.created_at,
        stats: {
          reviews_given: 15,
          reviews_received: 8,
          bounties_created: 3,
          followers: 42,
          following: 25,
        },
      },
    })
  }),
]
