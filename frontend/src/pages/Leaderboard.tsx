import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { api } from '../services/api'
import { Trophy, Medal, Crown, Flame } from 'lucide-react'

function getRankIcon(rank: number) {
  switch (rank) {
    case 1:
      return <Crown className="w-6 h-6 text-amber-400" />
    case 2:
      return <Medal className="w-6 h-6 text-gray-400" />
    case 3:
      return <Medal className="w-6 h-6 text-amber-600" />
    default:
      return <span className="text-lg font-bold text-gray-500">#{rank}</span>
  }
}

function getRankBg(rank: number) {
  switch (rank) {
    case 1:
      return 'bg-gradient-to-r from-amber-50 to-yellow-50 border-amber-200'
    case 2:
      return 'bg-gradient-to-r from-gray-50 to-slate-50 border-gray-200'
    case 3:
      return 'bg-gradient-to-r from-amber-50 to-orange-50 border-amber-300'
    default:
      return 'bg-white border-gray-100'
  }
}

export default function Leaderboard() {
  const { data: leaderboard, isLoading } = useQuery({
    queryKey: ['leaderboard'],
    queryFn: () => api.getLeaderboard(50),
  })

  const { data: streak } = useQuery({
    queryKey: ['streak'],
    queryFn: api.getStreak,
  })

  return (
    <div>
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-amber-500 to-orange-500 rounded-2xl p-8 mb-8">
        <div className="flex items-center gap-4">
          <div className="w-16 h-16 bg-white/20 rounded-2xl flex items-center justify-center">
            <Trophy className="w-8 h-8 text-white" />
          </div>
          <div>
            <h1 className="text-3xl font-bold text-white mb-1">Leaderboard</h1>
            <p className="text-amber-100">Top reviewers by points</p>
          </div>
        </div>

        {/* Current user streak */}
        {streak && streak.current_streak > 0 && (
          <div className="mt-6 flex items-center gap-2 bg-white/20 rounded-xl px-4 py-3 w-fit">
            <Flame className="w-5 h-5 text-amber-300" />
            <span className="text-white font-semibold">{streak.current_streak} day streak!</span>
            <span className="text-amber-200 text-sm">(Best: {streak.longest_streak})</span>
          </div>
        )}
      </div>

      {/* Leaderboard */}
      {isLoading ? (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500"></div>
          <p className="mt-2 text-gray-500">Loading leaderboard...</p>
        </div>
      ) : !leaderboard || leaderboard.length === 0 ? (
        <div className="bg-white rounded-2xl p-12 text-center">
          <Trophy className="w-16 h-16 mx-auto text-gray-300 mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">No rankings yet</h3>
          <p className="text-gray-500 mb-6">Be the first to earn points!</p>
          <Link
            to="/bounties"
            className="inline-flex items-center px-6 py-3 bg-amber-500 text-white rounded-lg hover:bg-amber-600 transition-colors"
          >
            Browse Bounties
          </Link>
        </div>
      ) : (
        <div className="space-y-3">
          {leaderboard.map((entry) => (
            <Link
              key={entry.user_id}
              to={`/profile?id=${entry.user_id}`}
              className={`flex items-center gap-4 p-4 rounded-xl border hover:shadow-md transition-all ${getRankBg(entry.rank)}`}
            >
              <div className="w-12 h-12 flex items-center justify-center">
                {getRankIcon(entry.rank)}
              </div>
              <div className="flex-1">
                <p className="font-semibold text-gray-900">{entry.display_name || entry.username}</p>
                <p className="text-sm text-gray-500">@{entry.username}</p>
              </div>
              <div className="text-right">
                <p className="text-xl font-bold text-amber-600">{entry.points.toLocaleString()}</p>
                <p className="text-xs text-gray-500">{entry.review_count} reviews</p>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}
