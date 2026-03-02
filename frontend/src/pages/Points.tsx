import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '../services/api'
import { PointTransaction } from '../types'

export default function Points() {
  const [filter, setFilter] = useState<'all' | 'earned' | 'spent'>('all')

  const { data: balance, isLoading: balanceLoading } = useQuery({
    queryKey: ['balance'],
    queryFn: () => api.getBalance(),
  })

  const { data: transactionsData, isLoading: transactionsLoading } = useQuery({
    queryKey: ['transactions'],
    queryFn: () => api.getTransactions({ limit: 50 }),
  })

  const transactions = transactionsData?.transactions || []

  const filteredTransactions = filter === 'all'
    ? transactions
    : transactions.filter(t => t.type === filter || (filter === 'earned' && ['earned', 'bonus', 'refund', 'received'].includes(t.type)) || (filter === 'spent' && ['spent', 'penalty', 'transferred'].includes(t.type)))

  const earnedTotal = transactions
    .filter(t => ['earned', 'bonus', 'refund', 'received'].includes(t.type))
    .reduce((sum, t) => sum + t.amount, 0)

  const spentTotal = transactions
    .filter(t => ['spent', 'penalty', 'transferred'].includes(t.type))
    .reduce((sum, t) => sum + Math.abs(t.amount), 0)

  const getTypeLabel = (type: string) => {
    const labels: Record<string, string> = {
      earned: 'Review Bonus',
      bonus: 'Bonus',
      refund: 'Refund',
      received: 'Received',
      spent: 'Bounty Claim',
      penalty: 'Penalty',
      transferred: 'Transfer',
    }
    return labels[type] || type
  }

  if (balanceLoading || transactionsLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-teal-500"></div>
      </div>
    )
  }

  return (
    <div className="min-h-screen pb-12">
      {/* Hero Section - Teal gradient with modern feel */}
      <div className="relative overflow-hidden bg-gradient-to-br from-teal-500 via-teal-600 to-cyan-600 rounded-3xl mb-8">
        <div className="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGNpcmNsZSBjeD0iMzAiIGN5PSIzMCIgcj0iMiIgZmlsbD0iI2ZmZiIgZmlsbC1vcGFjaXR5PSIwLjEiLz48L3N2Zz4=')] opacity-30"></div>
        <div className="absolute -right-20 -top-20 w-80 h-80 bg-white/10 rounded-full blur-3xl"></div>
        <div className="absolute -left-20 -bottom-20 w-60 h-60 bg-cyan-300/20 rounded-full blur-3xl"></div>

        <div className="relative px-8 py-12">
          <div className="flex items-center gap-3 mb-2">
            <div className="w-10 h-10 bg-white/20 rounded-xl flex items-center justify-center">
              <svg className="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.736 6.979C9.208 6.193 9.696 6 10 6c.304 0 .792.193 1.264.979a1 1 0 001.715-1.029C12.279 4.784 11.232 4 10 4s-2.279.784-2.979 1.95c-.285.475-.507 1-.67 1.55H6a1 1 0 000 2h.013a9.358 9.358 0 000 1H6a1 1 0 100 2h.351c.163.55.385 1.075.67 1.55C7.721 15.216 8.768 16 10 16s2.279-.784 2.979-1.95a1 1 0 10-1.715-1.029c-.472.786-.96.979-1.264.979-.304 0-.792-.193-1.264-.979a1 1 0 00-1.715 1.029c.7 1.167 1.747 1.951 2.979 1.951s2.279-.784 2.979-1.95c.285-.475.507-1 .67-1.55H14a1 1 0 100-2h-.013a9.358 9.358 0 00-1.351-4.5H14a1 1 0 100-2h-.013a9.358 9.358 0 00-1.351-4.5H10z" />
              </svg>
            </div>
            <span className="text-teal-100 font-medium">Your Points Balance</span>
          </div>

          <div className="text-6xl font-bold text-white mb-4 tracking-tight">
            {balance?.toLocaleString() || 0}
            <span className="text-2xl text-teal-200 ml-2">pts</span>
          </div>

          <div className="flex gap-6">
            <div className="bg-white/10 backdrop-blur-sm rounded-2xl px-6 py-3">
              <div className="text-teal-200 text-sm mb-1">Total Earned</div>
              <div className="text-white font-bold text-xl">+{earnedTotal.toLocaleString()}</div>
            </div>
            <div className="bg-white/10 backdrop-blur-sm rounded-2xl px-6 py-3">
              <div className="text-teal-200 text-sm mb-1">Total Spent</div>
              <div className="text-white font-bold text-xl">-{spentTotal.toLocaleString()}</div>
            </div>
          </div>
        </div>
      </div>

      {/* Stats Cards - Modern card design */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-gradient-to-br from-emerald-50 to-teal-50 rounded-2xl p-6 border border-emerald-100">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 bg-emerald-500 rounded-xl flex items-center justify-center">
              <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </div>
            <span className="text-emerald-800 font-semibold">Earned This Month</span>
          </div>
          <div className="text-3xl font-bold text-emerald-600">{earnedTotal.toLocaleString()}</div>
        </div>

        <div className="bg-gradient-to-br from-rose-50 to-orange-50 rounded-2xl p-6 border border-rose-100">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 bg-rose-500 rounded-xl flex items-center justify-center">
              <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
              </svg>
            </div>
            <span className="text-rose-800 font-semibold">Spent This Month</span>
          </div>
          <div className="text-3xl font-bold text-rose-600">{spentTotal.toLocaleString()}</div>
        </div>

        <div className="bg-gradient-to-br from-cyan-50 to-blue-50 rounded-2xl p-6 border border-cyan-100">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 bg-cyan-500 rounded-xl flex items-center justify-center">
              <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            </div>
            <span className="text-cyan-800 font-semibold">Net Points</span>
          </div>
          <div className={`text-3xl font-bold ${earnedTotal - spentTotal >= 0 ? 'text-emerald-600' : 'text-rose-600'}`}>
            {earnedTotal - spentTotal >= 0 ? '+' : ''}{(earnedTotal - spentTotal).toLocaleString()}
          </div>
        </div>
      </div>

      {/* Transaction History */}
      <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
        <div className="px-6 py-5 border-b border-gray-100 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-teal-100 rounded-lg flex items-center justify-center">
              <svg className="w-4 h-4 text-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <h2 className="text-xl font-bold text-gray-900">Transaction History</h2>
          </div>

          <div className="flex gap-2">
            {(['all', 'earned', 'spent'] as const).map((f) => (
              <button
                key={f}
                onClick={() => setFilter(f)}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                  filter === f
                    ? 'bg-teal-500 text-white shadow-md shadow-teal-200'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {f.charAt(0).toUpperCase() + f.slice(1)}
              </button>
            ))}
          </div>
        </div>

        {filteredTransactions.length === 0 ? (
          <div className="py-16 text-center">
            <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <p className="text-gray-500">No transactions yet</p>
            <p className="text-sm text-gray-400 mt-1">Start reviewing to earn points!</p>
          </div>
        ) : (
          <div className="divide-y divide-gray-50">
            {filteredTransactions.map((tx: PointTransaction) => (
              <div key={tx.id} className="px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors">
                <div className="flex items-center gap-4">
                  <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${
                    tx.amount > 0 ? 'bg-emerald-100' : 'bg-rose-100'
                  }`}>
                    {tx.amount > 0 ? (
                      <svg className="w-5 h-5 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 11l5-5m0 0l5 5m-5-5v12" />
                      </svg>
                    ) : (
                      <svg className="w-5 h-5 text-rose-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 13l-5 5m0 0l-5-5m5 5V6" />
                      </svg>
                    )}
                  </div>
                  <div>
                    <div className="font-medium text-gray-900">{getTypeLabel(tx.type)}</div>
                    <div className="text-sm text-gray-500">{tx.description || 'No description'}</div>
                  </div>
                </div>
                <div className="text-right">
                  <div className={`font-bold text-lg ${
                    tx.amount > 0 ? 'text-emerald-600' : 'text-rose-600'
                  }`}>
                    {tx.amount > 0 ? '+' : ''}{tx.amount.toLocaleString()}
                  </div>
                  <div className="text-xs text-gray-400">
                    {new Date(tx.created_at).toLocaleDateString()}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Tips Section */}
      <div className="mt-8 bg-gradient-to-r from-teal-50 to-cyan-50 rounded-2xl p-6 border border-teal-100">
        <div className="flex items-start gap-4">
          <div className="w-10 h-10 bg-teal-500 rounded-xl flex items-center justify-center flex-shrink-0">
            <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <h3 className="font-semibold text-teal-900 mb-1">How to Earn More Points</h3>
            <p className="text-teal-700 text-sm">
              Complete quality reviews to earn points. Reviews must be at least 100 words to qualify.
              The better your review, the more likely authors are to leave you positive feedback!
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
