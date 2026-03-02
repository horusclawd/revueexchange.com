import { setupWorker } from 'msw/browser'
import { authHandlers } from './handlers/auth'
import { userHandlers } from './handlers/users'
import { bountyHandlers } from './handlers/bounties'
import { pointsHandlers } from './handlers/points'

export const worker = setupWorker(
  ...authHandlers,
  ...userHandlers,
  ...bountyHandlers,
  ...pointsHandlers,
)
