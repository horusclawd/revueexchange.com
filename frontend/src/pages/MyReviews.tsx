import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '../services/api'
import { Star, Clock, MessageSquare, ThumbsUp, Eye, Edit2, Send, BookOpen, ArrowRight } from 'lucide-react'

type ReviewStatus = 'draft' | 'submitted' | 'published'

const statusStyles: Record<ReviewStatus, { bg: string; text: string; label: string }> = {
  draft: { bg: 'bg-slate-100', text: 'text-slate-600', label: 'Draft' },
  submitted: { bg: 'bg-blue-100', text: 'text-blue-700', label: 'Submitted' },
  published: { bg: 'bg-green-100', text: 'text-green-700', label: 'Published' },
}

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

  // Mock data - in real app would fetch from API
  const mockReviews = [
    { id: '1', bounty_id: 'bounty-1', rating: 4, title: 'Great fantasy world-building', content: 'This was a fantastic read...', status: 'published' as ReviewStatus, word_count: 250, created_at: '2024-01-15T10:00:00Z' },
    { id: '2', bounty_id: 'bounty-2', rating: 3, title: '', content: '', status: 'draft' as ReviewStatus, word_count: 0, created_at: '2024-01-20T10:00:00Z' },
  ]

  const { data: bounties } = useQuery({
    queryKey: ['bounties'],
    queryFn: () => api.getBounties({ status: 'claimed' }),
  })

  const createMutation = useMutation({
    mutationFn: (data: { bounty_id: string; rating: number; title: string; content: string }) =>
      api.createReview(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reviews'] })
      setEditingId(null)
      setRating(0)
      setTitle('')
      setContent('')
    },
  })

  // const updateMutation = useMutation({
  //   mutationFn: (data: { id: string; rating: number; title: string; content: string }) =>
  //     api.updateReview(data.id, data),
  //   onSuccess: () => {
  //     queryClient.invalidateQueries({ queryKey: ['reviews'] })
  //     setEditingId(null)
  //   },
  // })

  const submitMutation = useMutation({
    mutationFn: (id: string) => api.submitReview(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reviews'] })
    },
  })

  const claimedBounties = bounties?.bounties || []

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

  return (
    <div className="min-h-screen -mx-4 -my-8 px-4 py-8 bg-gradient-to-b from-slate-50 to-blue-50/50">
      <div className="max-w-5xl mx-auto">
        {/* Header */}
        <div className="mb-10">
          <h1 className="text-4xl font-bold text-slate-800 mb-3">
            My <span className="text-blue-600">Reviews</span>
          </h1>
          <p className="text-slate-600 text-lg">
            Manage your reviews and track your progress
          </p>
        </div>

        {/* Create New Review Section */}
        <div className="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden mb-8">
          <div className="bg-gradient-to-r from-blue-600 to-indigo-600 px-6 py-4">
            <h2 className="text-white font-semibold text-lg flex items-center gap-2">
              <Edit2 className="w-5 h-5" />
              Write a Review
            </h2>
          </div>

          <div className="p-6">
            {editingId ? (
              <div className="space-y-6">
                {/* Rating */}
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-2">Rating</label>
                  <StarRating rating={rating} onChange={setRating} />
                </div>

                {/* Title */}
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-2">Title</label>
                  <input
                    type="text"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    placeholder="Give your review a title..."
                    className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500"
                  />
                </div>

                {/* Content */}
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-2">
                    Review Content
                    <span className={`ml-2 text-sm ${wordCount >= 10 ? 'text-green-600' : 'text-slate-400'}`}>
                      ({wordCount} words {wordCount < 10 && '- minimum 10 required'})
                    </span>
                  </label>
                  <textarea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    rows={8}
                    placeholder="Share your honest thoughts about the work..."
                    className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 resize-none"
                  />
                </div>

                <div className="flex justify-between items-center pt-2">
                  <button
                    onClick={() => setEditingId(null)}
                    className="text-slate-500 hover:text-slate-700 font-medium"
                  >
                    Cancel
                  </button>
                  <div className="flex gap-3">
                    <button
                      onClick={saveDraft}
                      disabled={createMutation.isPending}
                      className="px-5 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-700 font-medium rounded-lg transition-colors"
                    >
                      Save Draft
                    </button>
                    <button
                      onClick={saveDraft}
                      disabled={wordCount < 10 || rating === 0 || createMutation.isPending}
                      className="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                    >
                      <Send className="w-4 h-4" />
                      {createMutation.isPending ? 'Submitting...' : 'Submit Review'}
                    </button>
                  </div>
                </div>
              </div>
            ) : (
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-3">
                  Select a claimed bounty to review:
                </label>
                {claimedBounties.length === 0 ? (
                  <div className="text-center py-8 bg-slate-50 rounded-xl">
                    <BookOpen className="w-12 h-12 mx-auto text-slate-300 mb-3" />
                    <p className="text-slate-500">No claimed bounties yet.</p>
                    <a href="/bounties" className="text-blue-600 hover:underline font-medium">
                      Browse bounties to claim
                    </a>
                  </div>
                ) : (
                  <div className="grid gap-3">
                    {claimedBounties.map((bounty) => (
                      <button
                        key={bounty.id}
                        onClick={() => startNewReview(bounty.id)}
                        className="flex items-center justify-between p-4 bg-slate-50 hover:bg-blue-50 border border-slate-200 hover:border-blue-300 rounded-xl transition-all text-left group"
                      >
                        <div>
                          <p className="font-medium text-slate-800">Bounty #{bounty.id.slice(0, 8)}</p>
                          <p className="text-sm text-slate-500">{bounty.bounty_points} points</p>
                        </div>
                        <ArrowRight className="w-5 h-5 text-slate-400 group-hover:text-blue-600 group-hover:translate-x-1 transition-all" />
                      </button>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
        </div>

        {/* Existing Reviews */}
        <div className="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
          <div className="px-6 py-4 border-b border-slate-100">
            <h2 className="font-semibold text-slate-800">Your Reviews</h2>
          </div>

          <div className="divide-y divide-slate-100">
            {mockReviews.map((review) => {
              const status = statusStyles[review.status]
              return (
                <div key={review.id} className="p-6">
                  <div className="flex items-start justify-between mb-3">
                    <div className="flex items-center gap-3">
                      <span className={`px-3 py-1 rounded-full text-xs font-semibold ${status.bg} ${status.text}`}>
                        {status.label}
                      </span>
                      <StarRating rating={review.rating} readonly />
                    </div>
                    <span className="text-sm text-slate-400 flex items-center gap-1">
                      <Clock className="w-4 h-4" />
                      {new Date(review.created_at).toLocaleDateString()}
                    </span>
                  </div>

                  {review.title && (
                    <h3 className="font-semibold text-slate-800 text-lg mb-2">{review.title}</h3>
                  )}

                  {review.content ? (
                    <p className="text-slate-600 mb-3">{review.content}</p>
                  ) : (
                    <p className="text-slate-400 italic mb-3">No content yet...</p>
                  )}

                  <div className="flex items-center gap-4 text-sm text-slate-400">
                    <span className="flex items-center gap-1">
                      <MessageSquare className="w-4 h-4" />
                      {review.word_count} words
                    </span>
                    <span className="flex items-center gap-1">
                      <ThumbsUp className="w-4 h-4" />
                      0 helpful
                    </span>
                    <span className="flex items-center gap-1">
                      <Eye className="w-4 h-4" />
                      0 views
                    </span>
                  </div>

                  {review.status === 'draft' && (
                    <button
                      onClick={() => submitMutation.mutate(review.id)}
                      disabled={submitMutation.isPending}
                      className="mt-4 px-4 py-2 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-lg transition-colors disabled:opacity-50"
                    >
                      {submitMutation.isPending ? 'Submitting...' : 'Submit for Review'}
                    </button>
                  )}
                </div>
              )
            })}

            {mockReviews.length === 0 && (
              <div className="p-12 text-center">
                <MessageSquare className="w-12 h-12 mx-auto text-slate-300 mb-3" />
                <p className="text-slate-500">No reviews yet. Claim a bounty to get started!</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
