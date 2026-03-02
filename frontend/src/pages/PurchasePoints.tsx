import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '../services/api'

interface PurchaseOption {
  id: string
  amountCents: number
  points: number
  bonusPoints: number
  popular?: boolean
}

const purchaseOptions: PurchaseOption[] = [
  { id: '1', amountCents: 500, points: 500, bonusPoints: 0 },
  { id: '2', amountCents: 1000, points: 1000, bonusPoints: 50, popular: true },
  { id: '3', amountCents: 2500, points: 2500, bonusPoints: 250 },
  { id: '4', amountCents: 5000, points: 5000, bonusPoints: 750 },
  { id: '5', amountCents: 10000, points: 10000, bonusPoints: 2000 },
]

export default function PurchasePoints() {
  const queryClient = useQueryClient()
  const [selectedOption, setSelectedOption] = useState<PurchaseOption | null>(null)
  const [showSuccess, setShowSuccess] = useState(false)
  const [purchasedPoints, setPurchasedPoints] = useState(0)

  const { data: balance } = useQuery({
    queryKey: ['balance'],
    queryFn: () => api.getBalance(),
  })

  const { data: payments } = useQuery({
    queryKey: ['payment-history'],
    queryFn: () => api.getPaymentHistory(),
  })

  const checkoutMutation = useMutation({
    mutationFn: (amountCents: number) => api.createCheckout(amountCents),
    onSuccess: (data) => {
      // In production, redirect to Stripe checkout
      // For MVP, simulate success
      setPurchasedPoints(data.points_award)
      setShowSuccess(true)
      queryClient.invalidateQueries({ queryKey: ['balance'] })
      queryClient.invalidateQueries({ queryKey: ['payment-history'] })
    },
  })

  const handlePurchase = () => {
    if (selectedOption) {
      checkoutMutation.mutate(selectedOption.amountCents)
    }
  }

  const formatCurrency = (cents: number) => {
    return `$${(cents / 100).toFixed(2)}`
  }

  if (showSuccess) {
    return (
      <div className="min-h-screen flex items-center justify-center p-4">
        <div className="bg-white rounded-3xl shadow-2xl p-8 max-w-md w-full text-center">
          <div className="w-20 h-20 bg-gradient-to-br from-violet-500 to-purple-600 rounded-full flex items-center justify-center mx-auto mb-6">
            <svg className="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Purchase Complete!</h2>
          <p className="text-gray-600 mb-6">
            You've successfully added <span className="font-bold text-violet-600">{purchasedPoints} points</span> to your account.
          </p>
          <div className="bg-violet-50 rounded-2xl p-4 mb-6">
            <div className="text-4xl font-bold text-violet-600">{balance?.toLocaleString() || 0}</div>
            <div className="text-sm text-violet-500">Total Points</div>
          </div>
          <button
            onClick={() => {
              setShowSuccess(false)
              setSelectedOption(null)
            }}
            className="w-full bg-gradient-to-r from-violet-500 to-purple-600 text-white py-3 rounded-xl font-semibold hover:from-violet-600 hover:to-purple-700 transition-all"
          >
            Purchase More Points
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen pb-12">
      {/* Hero Section - Purple gradient */}
      <div className="relative overflow-hidden bg-gradient-to-br from-violet-600 via-purple-600 to-fuchsia-600 rounded-3xl mb-8">
        <div className="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGNpcmNsZSBjeD0iMzAiIGN5PSIzMCIgcj0iMiIgZmlsbD0iI2ZmZiIgZmlsbC1vcGFjaXR5PSIwLjEiLz48L3N2Zz4=')] opacity-30"></div>
        <div className="absolute -right-20 -top-20 w-80 h-80 bg-white/10 rounded-full blur-3xl"></div>
        <div className="absolute -left-20 -bottom-20 w-60 h-60 bg-fuchsia-300/20 rounded-full blur-3xl"></div>

        <div className="relative px-8 py-12 text-center">
          <div className="inline-flex items-center gap-2 bg-white/20 backdrop-blur-sm rounded-full px-4 py-2 mb-4">
            <svg className="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
              <path d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.736 6.979C9.208 6.193 9.696 6 10 6c.304 0 .792.193 1.264.979a1 1 0 001.715-1.029C12.279 4.784 11.232 4 10 4s-2.279.784-2.979 1.95c-.285.475-.507 1-.67 1.55H6a1 1 0 000 2h.013a9.358 9.358 0 000 1H6a1 1 0 100 2h.351c.163.55.385 1.075.67 1.55C7.721 15.216 8.768 16 10 16s2.279-.784 2.979-1.95c.285-.475.507-1 .67-1.55H14a1 1 0 100-2h-.013a9.358 9.358 0 00-1.351-4.5H14a1 1 0 100-2h-.013a9.358 9.358 0 00-1.351-4.5H10z" />
            </svg>
            <span className="text-white font-medium">Points Store</span>
          </div>
          <h1 className="text-4xl font-bold text-white mb-2">Boost Your Review Power</h1>
          <p className="text-violet-100 text-lg max-w-md mx-auto">
            Get more points to claim bounties and grow your audience
          </p>
        </div>
      </div>

      {/* Current Balance */}
      <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 mb-8">
        <div className="flex items-center justify-between">
          <div>
            <div className="text-sm text-gray-500 mb-1">Your Current Balance</div>
            <div className="text-3xl font-bold text-gray-900">{balance?.toLocaleString() || 0} pts</div>
          </div>
          <div className="w-16 h-16 bg-violet-100 rounded-2xl flex items-center justify-center">
            <svg className="w-8 h-8 text-violet-600" fill="currentColor" viewBox="0 0 20 20">
              <path d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.736 6.979C9.208 6.193 9.696 6 10 6c.304 0 .792.193 1.264.979a1 1 0 001.715-1.029C12.279 4.784 11.232 4 10 4s-2.279.784-2.979 1.95c-.285.475-.507 1-.67 1.55H6a1 1 0 000 2h.013a9.358 9.358 0 000 1H6a1 1 0 100 2h.351c.163.55.385 1.075.67 1.55C7.721 15.216 8.768 16 10 16s2.279-.784 2.979-1.95c.285-.475.507-1 .67-1.55H14a1 1 0 100-2h-.013a9.358 9.358 0 00-1.351-4.5H14a1 1 0 100-2h-.013a9.358 9.358 0 00-1.351-4.5H10z" />
            </svg>
          </div>
        </div>
      </div>

      {/* Purchase Options */}
      <h2 className="text-xl font-bold text-gray-900 mb-4">Choose Points Package</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
        {purchaseOptions.map((option) => (
          <button
            key={option.id}
            onClick={() => setSelectedOption(option)}
            className={`relative bg-white rounded-2xl p-6 text-left transition-all ${
              selectedOption?.id === option.id
                ? 'ring-2 ring-violet-500 shadow-lg shadow-violet-100'
                : 'border border-gray-100 hover:border-violet-200 hover:shadow-md'
            }`}
          >
            {option.popular && (
              <div className="absolute -top-3 left-1/2 -translate-x-1/2">
                <span className="bg-gradient-to-r from-violet-500 to-purple-600 text-white text-xs font-bold px-3 py-1 rounded-full">
                  MOST POPULAR
                </span>
              </div>
            )}
            <div className="flex items-center justify-between mb-3">
              <div className="text-3xl font-bold text-gray-900">{option.points.toLocaleString()}</div>
              <div className="text-violet-600 font-semibold">pts</div>
            </div>
            {option.bonusPoints > 0 && (
              <div className="text-sm text-emerald-600 font-medium mb-3">
                +{option.bonusPoints} bonus points!
              </div>
            )}
            <div className="text-gray-500 text-lg">{formatCurrency(option.amountCents)}</div>
          </button>
        ))}
      </div>

      {/* Purchase Button */}
      <button
        onClick={handlePurchase}
        disabled={!selectedOption || checkoutMutation.isPending}
        className={`w-full py-4 rounded-2xl font-bold text-lg transition-all ${
          selectedOption
            ? 'bg-gradient-to-r from-violet-500 to-purple-600 text-white hover:from-violet-600 hover:to-purple-700 shadow-lg shadow-violet-200'
            : 'bg-gray-100 text-gray-400 cursor-not-allowed'
        }`}
      >
        {checkoutMutation.isPending ? (
          <span className="flex items-center justify-center gap-2">
            <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            Processing...
          </span>
        ) : selectedOption ? (
          `Purchase ${selectedOption.points.toLocaleString()} Points for ${formatCurrency(selectedOption.amountCents)}`
        ) : (
          'Select a package to purchase'
        )}
      </button>

      {/* Purchase History */}
      {payments && payments.length > 0 && (
        <div className="mt-12">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Purchase History</h2>
          <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
            <div className="divide-y divide-gray-50">
              {payments.slice(0, 5).map((payment) => (
                <div key={payment.id} className="px-6 py-4 flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${
                      payment.status === 'completed' ? 'bg-emerald-100' : 'bg-gray-100'
                    }`}>
                      {payment.status === 'completed' ? (
                        <svg className="w-5 h-5 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                        </svg>
                      ) : (
                        <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      )}
                    </div>
                    <div>
                      <div className="font-medium text-gray-900">
                        {formatCurrency(payment.amount_cents)} Points Purchase
                      </div>
                      <div className="text-sm text-gray-500">
                        {new Date(payment.created_at).toLocaleDateString()}
                      </div>
                    </div>
                  </div>
                  <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                    payment.status === 'completed' ? 'bg-emerald-100 text-emerald-700' :
                    payment.status === 'pending' ? 'bg-yellow-100 text-yellow-700' :
                    payment.status === 'failed' ? 'bg-red-100 text-red-700' :
                    'bg-gray-100 text-gray-700'
                  }`}>
                    {payment.status}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Info Section */}
      <div className="mt-12 bg-gradient-to-r from-violet-50 to-purple-50 rounded-2xl p-6 border border-violet-100">
        <div className="flex items-start gap-4">
          <div className="w-10 h-10 bg-violet-500 rounded-xl flex items-center justify-center flex-shrink-0">
            <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <h3 className="font-semibold text-violet-900 mb-1">How Points Work</h3>
            <p className="text-violet-700 text-sm">
              Points are used to claim review bounties. Each bounty has a point value that will be deducted from your balance when you claim it.
              Complete quality reviews to earn more points!
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
