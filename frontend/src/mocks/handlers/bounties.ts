import { http, HttpResponse, delay } from 'msw'
import { mockBounties, mockUsers } from '../data'

export const bountyHandlers = [
  // GET /api/v1/bounties
  http.get('/api/v1/bounties', async ({ request }) => {
    await delay(150)
    const url = new URL(request.url)
    const status = url.searchParams.get('status')
    // Note: genre and type filters available but not implemented in mock

    let bounties = [...mockBounties]

    if (status && status !== 'all') {
      bounties = bounties.filter(b => b.status === status)
    }

    // Add product info to each bounty
    const bountiesWithProduct = bounties.map(bounty => ({
      ...bounty,
      product: {
        id: bounty.product_id,
        title: 'Product ' + bounty.product_id,
        type: 'book',
        genre: 'fantasy',
        cover_url: undefined,
      },
      user: mockUsers.find(u => u.id === bounty.user_id),
    }))

    return HttpResponse.json({
      data: bountiesWithProduct,
      meta: {
        total: bountiesWithProduct.length,
        page: 1,
        page_size: 20,
      },
    })
  }),

  // POST /api/v1/bounties
  http.post('/api/v1/bounties', async ({ request }) => {
    await delay(300)
    const body = await request.json() as {
      product_id: string
      bounty_points: number
      bounty_cash?: number
      requirements: string
    }

    const newBounty = {
      id: 'bounty-' + Date.now(),
      user_id: 'user-1',
      product_id: body.product_id,
      bounty_points: body.bounty_points,
      bounty_cash: body.bounty_cash || undefined,
      status: 'open',
      requirements: body.requirements,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }

    return HttpResponse.json({ data: newBounty }, { status: 201 })
  }),

  // GET /api/v1/bounties/:id
  http.get('/api/v1/bounties/:id', async ({ params }) => {
    await delay(100)
    const bounty = mockBounties.find(b => b.id === params.id)

    if (!bounty) {
      return HttpResponse.json({ error: 'Bounty not found' }, { status: 404 })
    }

    return HttpResponse.json({
      data: {
        ...bounty,
        product: {
          id: bounty.product_id,
          title: 'Product ' + bounty.product_id,
          type: 'book',
          genre: 'fantasy',
          cover_url: undefined,
        },
        user: mockUsers.find(u => u.id === bounty.user_id),
      },
    })
  }),

  // PUT /api/v1/bounties/:id
  http.put('/api/v1/bounties/:id', async ({ params, request }) => {
    await delay(200)
    const bounty = mockBounties.find(b => b.id === params.id)

    if (!bounty) {
      return HttpResponse.json({ error: 'Bounty not found' }, { status: 404 })
    }

    const body = await request.json() as Partial<typeof bounty>
    const updatedBounty = { ...bounty, ...body, updated_at: new Date().toISOString() }

    return HttpResponse.json({ data: updatedBounty })
  }),

  // DELETE /api/v1/bounties/:id
  http.delete('/api/v1/bounties/:id', async ({ params }) => {
    await delay(100)
    const bounty = mockBounties.find(b => b.id === params.id)

    if (!bounty) {
      return HttpResponse.json({ error: 'Bounty not found' }, { status: 404 })
    }

    return HttpResponse.json({ success: true })
  }),

  // POST /api/v1/bounties/:id/claim
  http.post('/api/v1/bounties/:id/claim', async ({ params }) => {
    await delay(200)
    const bounty = mockBounties.find(b => b.id === params.id)

    if (!bounty) {
      return HttpResponse.json({ error: 'Bounty not found' }, { status: 404 })
    }

    if (bounty.status !== 'open') {
      return HttpResponse.json({ error: 'Bounty is not available' }, { status: 400 })
    }

    const updatedBounty = {
      ...bounty,
      status: 'claimed',
      claimed_by: 'user-2',
      claimed_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }

    return HttpResponse.json({ data: updatedBounty })
  }),
]
