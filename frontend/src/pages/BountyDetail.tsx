import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { api } from '../services/api'
import { ArrowLeft, Star, Clock, BookOpen, Headphones, FileText, Sparkles, User } from 'lucide-react'

type ProductType = 'book' | 'course' | 'podcast' | 'newsletter'
type BountyStatus = 'open' | 'claimed' | 'completed' | 'cancelled'

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

export default function BountyDetail() {
  const { id } = useParams<{ id: string }>()
  const [claimError, setClaimError] = useState('')

  const { data: bounty, isLoading, error } = useQuery({
    queryKey: ['bounty', id],
    queryFn: () => api.getBounty(id!),
    enabled: !!id,
  })

  const claimMutation = {
    mutate: async (bountyId: string) => {
      try {
        setClaimError('')
        await api.claimBounty(bountyId)
        window.location.reload()
      } catch (err: any) {
        setClaimError(err.message || 'Failed to claim bounty')
      }
    },
    isPending: false,
  }

  if (isLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="w-12 h-12 border-4 border-amber-500 border-t-transparent rounded-full animate-spin" />
      </div>
    )
  }

  if (error || !bounty) {
    return (
      <div className="max-w-4xl mx-auto py-8">
        <Link to="/bounties" className="flex items-center gap-2 text-slate-600 hover:text-slate-900 mb-8">
          <ArrowLeft className="w-4 h-4" />
          Back to Bounties
        </Link>
        <div className="bg-red-50 text-red-600 p-4 rounded-lg">
          Failed to load bounty details
        </div>
      </div>
    )
  }

  const status = statusStyles[bounty.status as BountyStatus] || statusStyles.open
  const productType: ProductType = 'book' // Would come from product data
  const typeStyle = typeColors[productType]

  return (
    <div className="max-w-4xl mx-auto py-8">
      <Link to="/bounties" className="flex items-center gap-2 text-slate-600 hover:text-slate-900 mb-8">
        <ArrowLeft className="w-4 h-4" />
        Back to Bounties
      </Link>

      <div className="bg-white rounded-2xl border border-slate-200 overflow-hidden">
        {/* Header */}
        <div className="p-8 border-b border-slate-100">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center gap-3">
              <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold ${status.bg} ${status.text}`}>
                <span className={`w-1.5 h-1.5 rounded-full ${bounty.status === 'open' ? 'bg-emerald-500' : 'bg-amber-500'}`} />
                {status.label}
              </span>
              <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium ${typeStyle.bg} ${typeStyle.text}`}>
                {typeStyle.icon}
                {productType}
              </span>
            </div>
            <div className="flex items-center gap-1.5 bg-gradient-to-r from-amber-500 to-orange-500 text-white px-4 py-2 rounded-lg font-bold shadow-md">
              <Star className="w-5 h-5 fill-current" />
              {bounty.bounty_points} pts
            </div>
          </div>

          <h1 className="text-3xl font-bold text-slate-800 mb-2">
            {bounty.product_id}
          </h1>
          <p className="text-slate-500">Bounty ID: {bounty.id}</p>
        </div>

        {/* Details */}
        <div className="p-8">
          {bounty.requirements && (
            <div className="mb-8">
              <h3 className="font-semibold text-slate-800 mb-2">Requirements</h3>
              <p className="text-slate-600">{bounty.requirements}</p>
            </div>
          )}

          {/* Meta info */}
          <div className="grid grid-cols-2 gap-6 mb-8">
            <div className="flex items-center gap-3 text-slate-600">
              <Clock className="w-5 h-5 text-slate-400" />
              <div>
                <p className="text-sm text-slate-500">Created</p>
                <p className="font-medium">{new Date(bounty.created_at).toLocaleDateString()}</p>
              </div>
            </div>
            <div className="flex items-center gap-3 text-slate-600">
              <User className="w-5 h-5 text-slate-400" />
              <div>
                <p className="text-sm text-slate-500">Author</p>
                <p className="font-medium">{bounty.user_id.slice(0, 8)}...</p>
              </div>
            </div>
          </div>

          {/* Claim button */}
          {bounty.status === 'open' && (
            <div>
              {claimError && (
                <div className="bg-red-50 text-red-600 p-3 rounded-lg mb-4 text-sm">
                  {claimError}
                </div>
              )}
              <button
                onClick={() => claimMutation.mutate(bounty.id)}
                className="w-full py-3 bg-amber-500 hover:bg-amber-400 text-white font-bold rounded-xl transition-colors"
              >
                Claim This Bounty
              </button>
              <p className="text-center text-sm text-slate-500 mt-2">
                You will earn {bounty.bounty_points} points when you complete the review
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
