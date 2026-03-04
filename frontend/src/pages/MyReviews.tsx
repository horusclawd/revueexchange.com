import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '../services/api'
import { Star, BookOpen, ArrowRight } from 'lucide-react'

function StarRating({ rating, onChange, readonly = false }: { rating: number; onChange?: (r: number) => void; readonly?: boolean }) {
  return (
    <div className="flex gap-1">
      {[1, 2, 3, 4, 5].map((star) => (
        <button
          key={star}
          type="button"
          disabled={readonly}
          onClick={() => onChange?.(star)}
          className={`transition-transform ${!readonly && 'hover:scale-110'} ${star <= rating ? 'text-amber-400' : 'text-slate-300'}`}
        >
          <Star className={`w-6 h-6 ${star <= rating ? 'fill-current' : ''}`} />
        </button>
      ))}
    </div>
  )
}

export default function MyReviews() {
  const [editingId, setEditingId] = useState<string | null>(null)
  const [rating, setRating] = useState(0)
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const queryClient = useQueryClient()

  // Fetch claimed bounties (bounties user has claimed to review)
  const { data: claimedBountiesResult, isLoading: bountiesLoading } = useQuery({
    queryKey: ['bounties', 'claimed'],
    queryFn: () => api.getBounties({ status: 'claimed' }),
  })

  const claimedBounties = claimedBountiesResult?.bounties || []

  // Create review mutation
  const createMutation = useMutation({
    mutationFn: (data: { bounty_id: string; rating: number; title: string; content: string }) =>
      api.createReview(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['bounties'] })
      setEditingId(null)
      setRating(0)
      setTitle('')
      setContent('')
    },
  })

  // Submit review mutation
  const submitMutation = useMutation({
    mutationFn: (id: string) => api.submitReview(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['bounties'] })
    },
  })

  const startNewReview = (bountyId: string) => {
    setEditingId(bountyId)
    setRating(0)
    setTitle('')
    setContent('')
  }

  const saveDraft = () => {
    if (editingId) {
      createMutation.mutate({ bounty_id: editingId, rating, title, content })
    }
  }

  const wordCount = content.trim().split(/\s+/).filter(Boolean).length

  if (bountiesLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full animate-spin" />
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto py-8 px-4">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-slate-800 mb-2">My Reviews</h1>
        <p className="text-slate-500">Manage your reviews and track your progress</p>
      </div>

      {/* Claimed Bounties - Start Reviews */}
      <div className="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden mb-8">
        <div className="px-6 py-4 border-b border-slate-100">
          <h2 className="font-semibold text-slate-800">Bounties To Review</h2>
        </div>

        <div className="divide-y divide-slate-100">
          {claimedBounties.length === 0 ? (
            <div className="p-12 text-center">
              <BookOpen className="w-12 h-12 mx-auto text-slate-300 mb-3" />
              <p className="text-slate-500">No bounties claimed yet. Claim a bounty to start reviewing!</p>
            </div>
          ) : (
            claimedBounties.map((bounty) => (
              <div key={bounty.id} className="p-6">
                {editingId === bounty.id ? (
                  // Review form
                  <div className="space-y-4">
                    <h3 className="font-semibold text-slate-800">Writing Review for Bounty {bounty.id.slice(0, 8)}...</h3>

                    <div>
                      <label className="block text-sm font-medium text-slate-700 mb-2">Rating</label>
                      <StarRating rating={rating} onChange={setRating} />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-slate-700 mb-2">Title (optional)</label>
                      <input
                        type="text"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        placeholder="Give your review a title"
                        className="w-full px-4 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-slate-700 mb-2">
                        Review Content
                        <span className="text-slate-400 font-normal ml-2">({wordCount} words, min 10)</span>
                      </label>
                      <textarea
                        value={content}
                        onChange={(e) => setContent(e.target.value)}
                        placeholder="Write your review here..."
                        rows={6}
                        className="w-full px-4 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                    </div>

                    <div className="flex gap-3">
                      <button
                        onClick={saveDraft}
                        disabled={createMutation.isPending}
                        className="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-medium rounded-lg transition-colors disabled:opacity-50"
                      >
                        {createMutation.isPending ? 'Saving...' : 'Save Draft'}
                      </button>
                      <button
                        onClick={() => submitMutation.mutate(bounty.id)}
                        disabled={submitMutation.isPending || wordCount < 10}
                        className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white font-medium rounded-lg transition-colors disabled:opacity-50"
                      >
                        {submitMutation.isPending ? 'Submitting...' : 'Submit Review'}
                      </button>
                      <button
                        onClick={() => setEditingId(null)}
                        className="px-4 py-2 text-slate-500 hover:text-slate-700"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                ) : (
                  // Bounty info
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium text-slate-800">Bounty #{bounty.id.slice(0, 8)}</p>
                      <p className="text-sm text-slate-500">{bounty.bounty_points} points</p>
                      <p className="text-sm text-slate-400 mt-1">
                        Claimed {bounty.claimed_at ? new Date(bounty.claimed_at).toLocaleDateString() : 'recently'}
                      </p>
                    </div>
                    <button
                      onClick={() => startNewReview(bounty.id)}
                      className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors"
                    >
                      Write Review
                      <ArrowRight className="w-4 h-4" />
                    </button>
                  </div>
                )}
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  )
}
