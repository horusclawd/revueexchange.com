import { useQuery } from '@tanstack/react-query'
import { api } from '../services/api'
import { Users, BookOpen, FileText, DollarSign, TrendingUp, Activity } from 'lucide-react'

export default function Analytics() {
  const { data: overview } = useQuery({
    queryKey: ['analytics-overview'],
    queryFn: api.getAnalyticsOverview,
  })

  const { data: reviewMetrics } = useQuery({
    queryKey: ['analytics-reviews'],
    queryFn: api.getReviewMetrics,
  })

  const { data: activity } = useQuery({
    queryKey: ['analytics-activity'],
    queryFn: () => api.getUserActivity(7),
  })

  return (
    <div>
      {/* Hero */}
      <div className="bg-gradient-to-r from-indigo-600 to-blue-600 rounded-2xl p-8 mb-8">
        <div className="flex items-center gap-4">
          <div className="w-14 h-14 bg-white/20 rounded-xl flex items-center justify-center">
            <TrendingUp className="w-7 h-7 text-white" />
          </div>
          <div>
            <h1 className="text-3xl font-bold text-white">Analytics</h1>
            <p className="text-blue-100">Platform overview and metrics</p>
          </div>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-white rounded-xl p-6 shadow-sm">
          <div className="flex items-center gap-3 mb-2">
            <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
              <Users className="w-5 h-5 text-blue-600" />
            </div>
            <span className="text-gray-500 text-sm">Total Users</span>
          </div>
          <p className="text-2xl font-bold text-gray-900">{overview?.total_users ?? '-'}</p>
        </div>

        <div className="bg-white rounded-xl p-6 shadow-sm">
          <div className="flex items-center gap-3 mb-2">
            <div className="w-10 h-10 bg-amber-100 rounded-lg flex items-center justify-center">
              <BookOpen className="w-5 h-5 text-amber-600" />
            </div>
            <span className="text-gray-500 text-sm">Total Bounties</span>
          </div>
          <p className="text-2xl font-bold text-gray-900">{overview?.total_bounties ?? '-'}</p>
        </div>

        <div className="bg-white rounded-xl p-6 shadow-sm">
          <div className="flex items-center gap-3 mb-2">
            <div className="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
              <FileText className="w-5 h-5 text-green-600" />
            </div>
            <span className="text-gray-500 text-sm">Total Reviews</span>
          </div>
          <p className="text-2xl font-bold text-gray-900">{overview?.total_reviews ?? '-'}</p>
        </div>

        <div className="bg-white rounded-xl p-6 shadow-sm">
          <div className="flex items-center gap-3 mb-2">
            <div className="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center">
              <DollarSign className="w-5 h-5 text-purple-600" />
            </div>
            <span className="text-gray-500 text-sm">Revenue</span>
          </div>
          <p className="text-2xl font-bold text-gray-900">
            ${((overview?.total_points_spent ?? 0) / 100).toFixed(2)}
          </p>
        </div>
      </div>

      {/* Activity & Metrics */}
      <div className="grid md:grid-cols-2 gap-6">
        {/* Recent Activity */}
        <div className="bg-white rounded-xl p-6 shadow-sm">
          <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
            <Activity className="w-5 h-5 text-gray-400" />
            Recent Activity
          </h2>
          {activity && activity.length > 0 ? (
            <div className="space-y-3">
              {activity.slice(-7).map((day: any) => (
                <div key={day.date} className="flex items-center justify-between text-sm">
                  <span className="text-gray-500">{new Date(day.date).toLocaleDateString()}</span>
                  <div className="flex gap-4">
                    <span className="text-green-600">+{day.new_users} users</span>
                    <span className="text-blue-600">+{day.new_reviews} reviews</span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-400 text-center py-4">No recent activity</p>
          )}
        </div>

        {/* Review Metrics */}
        <div className="bg-white rounded-xl p-6 shadow-sm">
          <h2 className="text-lg font-semibold mb-4">Review Metrics</h2>
          {reviewMetrics ? (
            <div className="space-y-3">
              {reviewMetrics.map((m: any) => (
                <div key={m.status} className="flex items-center justify-between">
                  <span className="capitalize text-gray-600">{m.status}</span>
                  <div className="text-right">
                    <span className="font-semibold">{m.count}</span>
                    <span className="text-gray-400 text-sm ml-2">
                      (avg {m.avg_rating?.toFixed(1)} stars)
                    </span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-400 text-center py-4">Loading...</p>
          )}
        </div>
      </div>
    </div>
  )
}
