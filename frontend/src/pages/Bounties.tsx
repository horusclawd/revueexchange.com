import { useQuery } from '@tanstack/react-query'
import { api } from '../services/api'

export default function Bounties() {
  const { data: result, isLoading } = useQuery({
    queryKey: ['bounties'],
    queryFn: () => api.getBounties(),
  })

  const bounties = result?.bounties

  if (isLoading) {
    return <div>Loading...</div>
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Review Bounties</h1>
        <button className="bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700">
          Create Bounty
        </button>
      </div>

      {!bounties || bounties.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          No bounties available. Be the first to create one!
        </div>
      ) : (
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {bounties.map((bounty) => (
            <div key={bounty.id} className="bg-white rounded-lg shadow p-6">
              <div className="flex justify-between items-start mb-4">
                <span className={`px-2 py-1 rounded text-sm ${
                  bounty.status === 'open' ? 'bg-green-100 text-green-800' :
                  bounty.status === 'claimed' ? 'bg-yellow-100 text-yellow-800' :
                  'bg-gray-100 text-gray-800'
                }`}>
                  {bounty.status}
                </span>
                <span className="text-2xl font-bold text-primary-600">
                  {bounty.bounty_points} pts
                </span>
              </div>
              <h3 className="text-lg font-semibold mb-2">Bounty #{bounty.id.slice(0, 8)}</h3>
              <p className="text-gray-600 text-sm mb-4">
                {bounty.requirements || 'No requirements specified'}
              </p>
              <div className="flex justify-between items-center text-sm text-gray-500">
                <span>Created {new Date(bounty.created_at).toLocaleDateString()}</span>
                {bounty.status === 'open' && (
                  <button className="text-primary-600 hover:underline">
                    Claim
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
