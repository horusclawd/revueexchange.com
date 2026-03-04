import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { api } from '../services/api'
import type { Activity } from '../types'

function getActivityIcon(type: Activity['type']) {
  switch (type) {
    case 'review_submitted':
      return (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
        </svg>
      )
    case 'bounty_created':
      return (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      )
    case 'follow':
      return (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      )
    case 'comment':
      return (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
      )
    case 'points_earned':
      return (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      )
    case 'bounty_completed':
      return (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      )
    default:
      return (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      )
  }
}

function getActivityColor(type: Activity['type']) {
  switch (type) {
    case 'review_submitted':
      return 'bg-blue-100 text-blue-600'
    case 'bounty_created':
      return 'bg-amber-100 text-amber-600'
    case 'follow':
      return 'bg-purple-100 text-purple-600'
    case 'comment':
      return 'bg-teal-100 text-teal-600'
    case 'points_earned':
      return 'bg-emerald-100 text-emerald-600'
    case 'bounty_completed':
      return 'bg-green-100 text-green-600'
    default:
      return 'bg-gray-100 text-gray-600'
  }
}

function getActivityText(activity: Activity): string {
  switch (activity.type) {
    case 'review_submitted':
      return 'submitted a review'
    case 'bounty_created':
      return 'created a new bounty'
    case 'follow':
      return 'started following someone'
    case 'comment':
      return 'left a comment'
    case 'points_earned':
      return 'earned points'
    case 'bounty_completed':
      return 'completed a bounty'
    default:
      return 'did something'
  }
}

function formatTimeAgo(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000)

  if (seconds < 60) return 'just now'
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`
  if (seconds < 604800) return `${Math.floor(seconds / 86400)}d ago`
  return date.toLocaleDateString()
}

export default function Feed() {
  const { data: activities, isLoading } = useQuery({
    queryKey: ['feed'],
    queryFn: api.getActivityFeed,
  })

  return (
    <div>
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-violet-600 to-purple-600 rounded-2xl p-8 mb-8">
        <h1 className="text-3xl font-bold text-white mb-2">Activity Feed</h1>
        <p className="text-violet-100">See what people you follow have been up to</p>
      </div>

      {/* Feed */}
      {isLoading ? (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-violet-600"></div>
          <p className="mt-2 text-gray-500">Loading feed...</p>
        </div>
      ) : !activities || activities.length === 0 ? (
        <div className="bg-white rounded-2xl p-12 text-center">
          <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-gray-900 mb-2">No activity yet</h3>
          <p className="text-gray-500 mb-6">Follow some authors to see their activity here!</p>
          <Link
            to="/bounties"
            className="inline-flex items-center px-6 py-3 bg-violet-600 text-white rounded-lg hover:bg-violet-700 transition-colors"
          >
            Explore Bounties
          </Link>
        </div>
      ) : (
        <div className="space-y-4">
          {activities.map((activity) => (
            <div
              key={activity.id}
              className="bg-white rounded-xl p-6 border border-gray-100 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start gap-4">
                <div className={`w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0 ${getActivityColor(activity.type)}`}>
                  {getActivityIcon(activity.type)}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <Link
                      to={`/profile?id=${activity.user_id}`}
                      className="font-semibold text-gray-900 hover:text-violet-600 transition-colors"
                    >
                      {activity.user?.display_name || activity.user?.username || 'Someone'}
                    </Link>
                    <span className="text-gray-500">{getActivityText(activity)}</span>
                  </div>
                  <p className="text-sm text-gray-500">{formatTimeAgo(activity.created_at)}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
