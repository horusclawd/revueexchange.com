import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '../services/api'
import { Search, Filter, Plus, Star, Clock, BookOpen, Headphones, FileText, Sparkles, X } from 'lucide-react'

type ProductType = 'book' | 'course' | 'podcast' | 'newsletter'
type BountyStatus = 'open' | 'claimed' | 'completed' | 'cancelled'

// Creative color palette for product types
const typeColors: Record<ProductType, { bg: string; text: string; icon: React.ReactNode }> = {
  book: { bg: 'bg-amber-50', text: 'text-amber-700', icon: <BookOpen className="w-4 h-4" /> },
  course: { bg: 'bg-indigo-50', text: 'text-indigo-700', icon: <FileText className="w-4 h-4" /> },
  podcast: { bg: 'bg-rose-50', text: 'text-rose-700', icon: <Headphones className="w-4 h-4" /> },
  newsletter: { bg: 'bg-emerald-50', text: 'text-emerald-700', icon: <Sparkles className="w-4 h-4" /> },
}

const statusStyles: Record<BountyStatus, { bg: string; text: string; label: string }> = {
  open: { bg: 'bg-emerald-500/10', text: 'text-emerald-600', label: 'Open' },
  claimed: { bg: 'bg-amber-500/10', text: 'text-amber-600', label: 'In Progress' },
  completed: { bg: 'bg-blue-500/10', text: 'text-blue-600', label: 'Completed' },
  cancelled: { bg: 'bg-slate-500/10', text: 'text-slate-500', label: 'Closed' },
}

export default function Bounties() {
  const [filter, setFilter] = useState<string>('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [showCreateModal, setShowCreateModal] = useState(false)
  const queryClient = useQueryClient()

  const { data: result, isLoading } = useQuery({
    queryKey: ['bounties', filter],
    queryFn: () => api.getBounties(filter !== 'all' ? { status: filter } : undefined),
  })

  const claimMutation = useMutation({
    mutationFn: (bountyId: string) => api.claimBounty(bountyId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['bounties'] })
    },
  })

  const createMutation = useMutation({
    mutationFn: (bounty: { product_id: string; bounty_points: number; requirements?: string }) =>
      api.createBounty(bounty),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['bounties'] })
      setShowCreateModal(false)
    },
  })

  const bounties = result?.bounties || []

  const filteredBounties = bounties.filter(b =>
    searchQuery === '' || b.id.includes(searchQuery) || b.requirements?.toLowerCase().includes(searchQuery.toLowerCase())
  )

  if (isLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="flex flex-col items-center gap-4">
          <div className="w-12 h-12 border-4 border-amber-500 border-t-transparent rounded-full animate-spin" />
          <p className="text-slate-500 font-medium">Loading bounties...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen -mx-4 -my-8 px-4 py-8 bg-gradient-to-br from-slate-50 via-white to-amber-50/30">
      {/* Hero Section */}
      <div className="max-w-6xl mx-auto mb-12">
        <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 p-8 md:p-12 text-white">
          <div className="absolute inset-0 opacity-10">
            <div className="absolute top-0 right-0 w-96 h-96 bg-amber-500 rounded-full blur-3xl transform translate-x-1/2 -translate-y-1/2" />
            <div className="absolute bottom-0 left-0 w-64 h-64 bg-emerald-500 rounded-full blur-3xl transform -translate-x-1/2 translate-y-1/2" />
          </div>

          <div className="relative z-10">
            <h1 className="text-4xl md:text-5xl font-black mb-4 tracking-tight">
              Review <span className="text-transparent bg-clip-text bg-gradient-to-r from-amber-400 to-emerald-400">Bounties</span>
            </h1>
            <p className="text-slate-300 text-lg md:text-xl max-w-2xl mb-8">
              Earn points and cash by writing honest reviews for self-published creators.
              Your feedback helps authors improve and readers discover great content.
            </p>

            <div className="flex flex-wrap gap-4">
              <button
                onClick={() => setShowCreateModal(true)}
                className="group flex items-center gap-2 bg-amber-500 hover:bg-amber-400 text-slate-900 font-bold px-6 py-3 rounded-xl transition-all hover:shadow-lg hover:shadow-amber-500/25"
              >
                <Plus className="w-5 h-5 group-hover:scale-110 transition-transform" />
                Create Bounty
              </button>
              <button className="flex items-center gap-2 bg-white/10 hover:bg-white/20 text-white font-medium px-6 py-3 rounded-xl backdrop-blur-sm transition-all">
                <Star className="w-5 h-5" />
                Leaderboard
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Filters & Search */}
      <div className="max-w-6xl mx-auto mb-8">
        <div className="flex flex-col md:flex-row gap-4 items-start md:items-center justify-between">
          {/* Filter Chips */}
          <div className="flex flex-wrap gap-2">
            {[
              { value: 'all', label: 'All Bounties' },
              { value: 'open', label: 'Open' },
              { value: 'claimed', label: 'In Progress' },
              { value: 'completed', label: 'Completed' },
            ].map((f) => (
              <button
                key={f.value}
                onClick={() => setFilter(f.value)}
                className={`px-4 py-2 rounded-full font-medium text-sm transition-all ${
                  filter === f.value
                    ? 'bg-slate-800 text-white shadow-lg shadow-slate-500/20'
                    : 'bg-white text-slate-600 hover:bg-slate-100 border border-slate-200'
                }`}
              >
                {f.label}
              </button>
            ))}
          </div>

          {/* Search */}
          <div className="relative w-full md:w-auto">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
            <input
              type="text"
              placeholder="Search bounties..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full md:w-72 pl-10 pr-4 py-2.5 bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-amber-500/50 focus:border-amber-500 transition-all"
            />
          </div>
        </div>
      </div>

      {/* Bounties Grid */}
      <div className="max-w-6xl mx-auto">
        {filteredBounties.length === 0 ? (
          <div className="text-center py-20">
            <div className="w-20 h-20 mx-auto mb-6 bg-slate-100 rounded-full flex items-center justify-center">
              <Filter className="w-8 h-8 text-slate-400" />
            </div>
            <h3 className="text-xl font-bold text-slate-700 mb-2">No bounties found</h3>
            <p className="text-slate-500">Try adjusting your filters or create the first bounty!</p>
          </div>
        ) : (
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredBounties.map((bounty) => {
              const status = statusStyles[bounty.status as BountyStatus] || statusStyles.open
              const productType: ProductType = 'book' // Mock - would come from product data
              const typeStyle = typeColors[productType]

              return (
                <div
                  key={bounty.id}
                  className="group bg-white rounded-2xl border border-slate-200 overflow-hidden hover:shadow-xl hover:shadow-slate-200/50 transition-all duration-300 hover:-translate-y-1"
                >
                  {/* Card Header */}
                  <div className="p-6 pb-4">
                    <div className="flex items-start justify-between mb-4">
                      {/* Status Badge */}
                      <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold ${status.bg} ${status.text}`}>
                        <span className={`w-1.5 h-1.5 rounded-full ${bounty.status === 'open' ? 'bg-emerald-500' : 'bg-amber-500'}`} />
                        {status.label}
                      </span>

                      {/* Points Badge */}
                      <div className="flex items-center gap-1.5 bg-gradient-to-r from-amber-500 to-orange-500 text-white px-3 py-1.5 rounded-lg font-bold text-sm shadow-md">
                        <Star className="w-4 h-4 fill-current" />
                        {bounty.bounty_points} pts
                      </div>
                    </div>

                    {/* Product Type */}
                    <div className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-medium mb-3 ${typeStyle.bg} ${typeStyle.text}`}>
                      {typeStyle.icon}
                      {productType}
                    </div>

                    {/* Title */}
                    <h3 className="text-lg font-bold text-slate-800 mb-2 group-hover:text-amber-600 transition-colors">
                      Bounty #{bounty.id.slice(0, 8)}
                    </h3>

                    {/* Requirements */}
                    <p className="text-slate-500 text-sm line-clamp-2 mb-4">
                      {bounty.requirements || 'Review requirements will be provided upon claim'}
                    </p>
                  </div>

                  {/* Card Footer */}
                  <div className="px-6 py-4 bg-slate-50 border-t border-slate-100 flex items-center justify-between">
                    <div className="flex items-center gap-1.5 text-slate-400 text-sm">
                      <Clock className="w-4 h-4" />
                      {new Date(bounty.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
                    </div>

                    {bounty.status === 'open' && (
                      <button
                        onClick={() => claimMutation.mutate(bounty.id)}
                        disabled={claimMutation.isPending}
                        className="flex items-center gap-2 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white text-sm font-medium rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        {claimMutation.isPending ? 'Claiming...' : 'Claim'}
                      </button>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>

      {/* Create Bounty Modal */}
      {showCreateModal && (
        <CreateBountyModal
          onClose={() => setShowCreateModal(false)}
          onSubmit={(data) => createMutation.mutate(data)}
          isPending={createMutation.isPending}
          error={createMutation.error?.message}
        />
      )}
    </div>
  )
}

// Create Bounty Modal Component
function CreateBountyModal({ onClose, onSubmit, isPending, error }: {
  onClose: () => void
  onSubmit: (data: { product_id: string; bounty_points: number; requirements?: string }) => void
  isPending: boolean
  error?: string
}) {
  const [productId, setProductId] = useState('')
  const [points, setPoints] = useState(50)
  const [requirements, setRequirements] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit({
      product_id: productId,
      bounty_points: points,
      requirements: requirements || undefined,
    })
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl max-w-md w-full p-6 relative">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-slate-400 hover:text-slate-600"
        >
          <X className="w-5 h-5" />
        </button>

        <h2 className="text-2xl font-bold text-slate-800 mb-6">Create Bounty</h2>

        {error && (
          <div className="bg-red-50 text-red-600 p-3 rounded-lg mb-4 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              Product ID
            </label>
            <input
              type="text"
              value={productId}
              onChange={(e) => setProductId(e.target.value)}
              placeholder="Enter product ID"
              required
              className="w-full px-4 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-amber-500"
            />
            <p className="text-xs text-slate-500 mt-1">
              You need to create a product first to get its ID
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              Points to Award
            </label>
            <input
              type="number"
              value={points}
              onChange={(e) => setPoints(parseInt(e.target.value) || 0)}
              min={10}
              max={500}
              required
              className="w-full px-4 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-amber-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              Requirements (optional)
            </label>
            <textarea
              value={requirements}
              onChange={(e) => setRequirements(e.target.value)}
              placeholder="What should the review include?"
              rows={3}
              className="w-full px-4 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-amber-500"
            />
          </div>

          <div className="flex gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border border-slate-200 text-slate-600 rounded-lg hover:bg-slate-50 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isPending}
              className="flex-1 px-4 py-2 bg-amber-500 text-white font-medium rounded-lg hover:bg-amber-400 transition-colors disabled:opacity-50"
            >
              {isPending ? 'Creating...' : 'Create Bounty'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
