import { useQuery } from '@tanstack/react-query'
import { useAuth } from '../context/AuthContext'
import { api } from '../services/api'

export default function Dashboard() {
  const { user } = useAuth()

  const { data: balance } = useQuery({
    queryKey: ['balance'],
    queryFn: api.getBalance,
  })

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Dashboard</h1>

      <div className="grid md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-sm font-medium text-gray-500">Points Balance</h3>
          <p className="text-3xl font-bold text-primary-600">{balance ?? 0}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-sm font-medium text-gray-500">Reputation Score</h3>
          <p className="text-3xl font-bold">{user?.reputation_score ?? 0}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-sm font-medium text-gray-500">Subscription</h3>
          <p className="text-3xl font-bold capitalize">{user?.subscription_tier ?? 'free'}</p>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-xl font-semibold mb-4">Profile</h2>
        <div className="space-y-2">
          <p><span className="font-medium">Email:</span> {user?.email}</p>
          <p><span className="font-medium">Username:</span> {user?.username}</p>
          <p><span className="font-medium">Display Name:</span> {user?.display_name}</p>
          <p><span className="font-medium">Member Since:</span> {new Date(user?.created_at || '').toLocaleDateString()}</p>
        </div>
      </div>
    </div>
  )
}
