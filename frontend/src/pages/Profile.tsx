import { useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '../context/AuthContext'
import { api } from '../services/api'
import { Award, X, Trophy } from 'lucide-react'

function getBadgeColor(tier: string) {
  switch (tier) {
    case 'bronze':
      return 'bg-amber-100 text-amber-700 border-amber-200'
    case 'silver':
      return 'bg-gray-100 text-gray-700 border-gray-200'
    case 'gold':
      return 'bg-yellow-100 text-yellow-700 border-yellow-200'
    case 'platinum':
      return 'bg-purple-100 text-purple-700 border-purple-200'
    default:
      return 'bg-gray-100 text-gray-700 border-gray-200'
  }
}

export default function Profile() {
  const [searchParams] = useSearchParams()
  const profileId = searchParams.get('id')
  const { user: currentUser } = useAuth()
  const queryClient = useQueryClient()
  const [isEditing, setIsEditing] = useState(false)
  const [displayName, setDisplayName] = useState(currentUser?.display_name || '')
  const [bio, setBio] = useState(currentUser?.bio || '')
  const [selectedBadge, setSelectedBadge] = useState<any>(null)

  // Determine which user to show
  const isOwnProfile = !profileId || profileId === currentUser?.id
  const userId = isOwnProfile ? currentUser!.id : profileId!

  // Fetch user profile
  const { data: profileUser, isLoading } = useQuery({
    queryKey: ['user', userId],
    queryFn: () => api.getUser(userId),
    enabled: !!userId,
  })

  // Check if following
  const { data: isFollowing } = useQuery({
    queryKey: ['isFollowing', userId],
    queryFn: () => api.isFollowing(userId),
    enabled: !!userId && !isOwnProfile,
  })

  // Fetch badges
  const { data: badges } = useQuery({
    queryKey: ['badges'],
    queryFn: api.getBadges,
  })

  // Follow/unfollow mutation
  const followMutation = useMutation({
    mutationFn: () => api.followUser(userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['isFollowing', userId] })
    },
  })

  const unfollowMutation = useMutation({
    mutationFn: () => api.unfollowUser(userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['isFollowing', userId] })
    },
  })

  // Update profile mutation
  const updateMutation = useMutation({
    mutationFn: (updates: { display_name?: string; bio?: string }) =>
      api.updateUser(currentUser!.id, updates),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] })
      setIsEditing(false)
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    updateMutation.mutate({ display_name: displayName, bio })
  }

  const handleFollow = () => {
    if (isFollowing) {
      unfollowMutation.mutate()
    } else {
      followMutation.mutate()
    }
  }

  if (isLoading) {
    return <div className="text-center py-12">Loading...</div>
  }

  if (!profileUser) {
    return <div>User not found</div>
  }

  // Viewing another user's profile
  if (!isOwnProfile) {
    return (
      <div className="max-w-2xl mx-auto">
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
          {/* Header */}
          <div className="bg-gradient-to-r from-violet-500 to-purple-600 h-32"></div>

          <div className="px-6 pb-6">
            <div className="flex items-end gap-4 -mt-12 mb-4">
              <div className="w-24 h-24 rounded-full bg-white p-1">
                <div className="w-full h-full rounded-full bg-violet-100 flex items-center justify-center text-3xl font-bold text-violet-600">
                  {profileUser.display_name?.[0]?.toUpperCase() || profileUser.username[0].toUpperCase()}
                </div>
              </div>
              <div className="flex-1 pb-2">
                <h2 className="text-2xl font-bold text-gray-900">{profileUser.display_name || profileUser.username}</h2>
                <p className="text-gray-500">@{profileUser.username}</p>
              </div>
              <button
                onClick={handleFollow}
                disabled={followMutation.isPending || unfollowMutation.isPending}
                className={`px-6 py-2.5 rounded-lg font-semibold transition-colors ${
                  isFollowing
                    ? 'bg-gray-100 text-gray-700 hover:bg-gray-200 border border-gray-200'
                    : 'bg-violet-600 text-white hover:bg-violet-700'
                } disabled:opacity-50`}
              >
                {followMutation.isPending || unfollowMutation.isPending
                  ? '...'
                  : isFollowing
                  ? 'Following'
                  : 'Follow'}
              </button>
            </div>

            {/* Stats */}
            <div className="flex gap-6 py-4 border-y border-gray-100 mb-4">
              <div className="text-center">
                <p className="text-2xl font-bold text-gray-900">{profileUser.points}</p>
                <p className="text-sm text-gray-500">Points</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-bold text-gray-900">{profileUser.reputation_score}</p>
                <p className="text-sm text-gray-500">Reputation</p>
              </div>
            </div>

            {/* Bio */}
            {profileUser.bio && (
              <div className="mb-4">
                <p className="text-gray-600">{profileUser.bio}</p>
              </div>
            )}

            <p className="text-sm text-gray-500">
              Member since {new Date(profileUser.created_at).toLocaleDateString()}
            </p>
          </div>
        </div>

        {/* Badges Section */}
        {badges && badges.length > 0 && (
          <div className="mt-6">
            <h3 className="font-semibold text-gray-900 mb-3">Badges</h3>
            <div className="flex flex-wrap gap-2">
              {badges.map((badge) => (
                <button
                  key={badge.id}
                  onClick={() => setSelectedBadge(badge)}
                  className={`px-3 py-1.5 rounded-full text-sm font-medium border ${getBadgeColor(badge.tier)}`}
                >
                  {badge.badge_name}
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Badge Modal */}
        {selectedBadge && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" onClick={() => setSelectedBadge(null)}>
            <div className="bg-white rounded-2xl p-6 max-w-md w-full" onClick={e => e.stopPropagation()}>
              <div className="flex items-start justify-between mb-4">
                <div className={`px-4 py-2 rounded-full text-lg font-bold ${getBadgeColor(selectedBadge.tier)}`}>
                  {selectedBadge.badge_name}
                </div>
                <button onClick={() => setSelectedBadge(null)} className="text-gray-400 hover:text-gray-600">
                  <X className="w-5 h-5" />
                </button>
              </div>
              <p className="text-gray-600 mb-4">{selectedBadge.description}</p>
              <div className="flex items-center gap-2 text-sm text-gray-500">
                <Award className="w-4 h-4" />
                <span className="capitalize">{selectedBadge.tier} tier</span>
                <span>•</span>
                <span>Earned {new Date(selectedBadge.awarded_at).toLocaleDateString()}</span>
              </div>
            </div>
          </div>
        )}
      </div>
    )
  }

  // Viewing own profile
  if (!currentUser) {
    return <div>Please login to view your profile</div>
  }

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">Profile</h1>

      {/* Badges for own profile */}
      {badges && badges.length > 0 && (
        <div className="mb-6">
          <div className="flex items-center gap-2 mb-3">
            <Trophy className="w-5 h-5 text-amber-500" />
            <h2 className="font-semibold text-gray-900">Your Badges</h2>
          </div>
          <div className="flex flex-wrap gap-2">
            {badges.map((badge) => (
              <button
                key={badge.id}
                onClick={() => setSelectedBadge(badge)}
                className={`px-3 py-1.5 rounded-full text-sm font-medium border ${getBadgeColor(badge.tier)}`}
              >
                {badge.badge_name}
              </button>
            ))}
          </div>
        </div>
      )}

      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <div className="flex items-center gap-4 mb-6">
          <div className="w-20 h-20 rounded-full bg-primary-100 flex items-center justify-center text-2xl font-bold text-primary-600">
            {currentUser.display_name?.[0]?.toUpperCase() || currentUser.username[0].toUpperCase()}
          </div>
          <div>
            <h2 className="text-xl font-semibold">{currentUser.display_name || currentUser.username}</h2>
            <p className="text-gray-500">@{currentUser.username}</p>
          </div>
        </div>

        {!isEditing ? (
          <>
            <div className="space-y-4 mb-6">
              <div>
                <label className="block text-sm font-medium text-gray-500">Email</label>
                <p>{currentUser.email}</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-500">Display Name</label>
                <p>{currentUser.display_name || 'Not set'}</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-500">Bio</label>
                <p>{currentUser.bio || 'No bio yet'}</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-500">Member Since</label>
                <p>{new Date(currentUser.created_at).toLocaleDateString()}</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-500">Reputation Score</label>
                <p>{currentUser.reputation_score || 0}</p>
              </div>
            </div>
            <button
              onClick={() => setIsEditing(true)}
              className="bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700"
            >
              Edit Profile
            </button>
          </>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">Display Name</label>
              <input
                type="text"
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Bio</label>
              <textarea
                value={bio || ''}
                onChange={(e) => setBio(e.target.value)}
                rows={4}
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
                placeholder="Tell us about yourself..."
              />
            </div>
            <div className="flex gap-4">
              <button
                type="submit"
                disabled={updateMutation.isPending}
                className="bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700 disabled:opacity-50"
              >
                {updateMutation.isPending ? 'Saving...' : 'Save'}
              </button>
              <button
                type="button"
                onClick={() => setIsEditing(false)}
                className="bg-gray-200 text-gray-700 px-4 py-2 rounded-md hover:bg-gray-300"
              >
                Cancel
              </button>
            </div>
            {updateMutation.isError && (
              <p className="text-red-600">Failed to update profile</p>
            )}
          </form>
        )}
      </div>

      {/* Badge Modal for own profile */}
      {selectedBadge && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" onClick={() => setSelectedBadge(null)}>
          <div className="bg-white rounded-2xl p-6 max-w-md w-full" onClick={e => e.stopPropagation()}>
            <div className="flex items-start justify-between mb-4">
              <div className={`px-4 py-2 rounded-full text-lg font-bold ${getBadgeColor(selectedBadge.tier)}`}>
                {selectedBadge.badge_name}
              </div>
              <button onClick={() => setSelectedBadge(null)} className="text-gray-400 hover:text-gray-600">
                <X className="w-5 h-5" />
              </button>
            </div>
            <p className="text-gray-600 mb-4">{selectedBadge.description}</p>
            <div className="flex items-center gap-2 text-sm text-gray-500">
              <Award className="w-4 h-4" />
              <span className="capitalize">{selectedBadge.tier} tier</span>
              <span>•</span>
              <span>Earned {new Date(selectedBadge.awarded_at).toLocaleDateString()}</span>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
