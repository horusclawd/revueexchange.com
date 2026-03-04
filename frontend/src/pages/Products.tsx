import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '../services/api'
import { Plus, BookOpen, Headphones, FileText, Sparkles, X } from 'lucide-react'

type ProductType = 'book' | 'course' | 'podcast' | 'newsletter'

const typeColors: Record<ProductType, { bg: string; text: string; icon: React.ReactNode }> = {
  book: { bg: 'bg-amber-50', text: 'text-amber-700', icon: <BookOpen className="w-4 h-4" /> },
  course: { bg: 'bg-indigo-50', text: 'text-indigo-700', icon: <FileText className="w-4 h-4" /> },
  podcast: { bg: 'bg-rose-50', text: 'text-rose-700', icon: <Headphones className="w-4 h-4" /> },
  newsletter: { bg: 'bg-emerald-50', text: 'text-emerald-700', icon: <Sparkles className="w-4 h-4" /> },
}

export default function Products() {
  const [showCreateModal, setShowCreateModal] = useState(false)
  const queryClient = useQueryClient()

  // Fetch user's products
  const { data: user } = useQuery({
    queryKey: ['user'],
    queryFn: () => api.getMe(),
  })

  const { data: products, isLoading } = useQuery({
    queryKey: ['products', user?.id],
    queryFn: () => user?.id ? api.getProducts(user.id) : Promise.resolve([]),
    enabled: !!user?.id,
  })

  const createMutation = useMutation({
    mutationFn: (data: { title: string; description: string; type: ProductType; genre: string }) =>
      api.createProduct(data as any),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      setShowCreateModal(false)
    },
  })

  if (isLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="w-12 h-12 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin" />
      </div>
    )
  }

  return (
    <div className="max-w-6xl mx-auto py-8 px-4">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-slate-800 mb-2">My Products</h1>
          <p className="text-slate-500">Manage your books, courses, podcasts, and newsletters</p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium px-4 py-2 rounded-lg transition-colors"
        >
          <Plus className="w-5 h-5" />
          Add Product
        </button>
      </div>

      {products && products.length > 0 ? (
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {products.map((product: any) => {
            const productType = (product.type as ProductType) || 'book'
            const typeStyle = typeColors[productType]
            return (
              <div key={product.id} className="bg-white rounded-xl border border-slate-200 overflow-hidden hover:shadow-lg transition-shadow">
                <div className="p-6">
                  <div className="flex items-start justify-between mb-3">
                    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-medium ${typeStyle.bg} ${typeStyle.text}`}>
                      {typeStyle.icon}
                      {productType}
                    </span>
                  </div>
                  <h3 className="font-semibold text-slate-800 mb-2">{product.title}</h3>
                  <p className="text-sm text-slate-500 line-clamp-2">{product.description}</p>
                  {product.genre && (
                    <span className="inline-block mt-3 text-xs bg-slate-100 text-slate-600 px-2 py-1 rounded">
                      {product.genre}
                    </span>
                  )}
                </div>
              </div>
            )
          })}
        </div>
      ) : (
        <div className="text-center py-20">
          <BookOpen className="w-16 h-16 mx-auto text-slate-300 mb-4" />
          <h3 className="text-xl font-semibold text-slate-700 mb-2">No products yet</h3>
          <p className="text-slate-500 mb-6">Add your first product to start getting reviews</p>
          <button
            onClick={() => setShowCreateModal(true)}
            className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium px-6 py-3 rounded-lg transition-colors"
          >
            <Plus className="w-5 h-5" />
            Add Your First Product
          </button>
        </div>
      )}

      {showCreateModal && (
        <CreateProductModal
          onClose={() => setShowCreateModal(false)}
          onSubmit={(data) => createMutation.mutate(data as any)}
          isPending={createMutation.isPending}
          error={createMutation.error?.message}
        />
      )}
    </div>
  )
}

function CreateProductModal({ onClose, onSubmit, isPending, error }: {
  onClose: () => void
  onSubmit: (data: { title: string; description: string; type: ProductType; genre: string }) => void
  isPending: boolean
  error?: string
}) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [type, setType] = useState<ProductType>('book')
  const [genre, setGenre] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit({ title, description, type, genre })
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl max-w-md w-full p-6 relative">
        <button onClick={onClose} className="absolute top-4 right-4 text-slate-400 hover:text-slate-600">
          <X className="w-5 h-5" />
        </button>

        <h2 className="text-2xl font-bold text-slate-800 mb-6">Add Product</h2>

        {error && (
          <div className="bg-red-50 text-red-600 p-3 rounded-lg mb-4 text-sm">{error}</div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">Title</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              className="w-full px-4 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
              className="w-full px-4 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">Type</label>
            <select
              value={type}
              onChange={(e) => setType(e.target.value as ProductType)}
              className="w-full px-4 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              <option value="book">Book</option>
              <option value="course">Course</option>
              <option value="podcast">Podcast</option>
              <option value="newsletter">Newsletter</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">Genre</label>
            <input
              type="text"
              value={genre}
              onChange={(e) => setGenre(e.target.value)}
              placeholder="e.g., Fiction, Technology, Business"
              className="w-full px-4 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div className="flex gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border border-slate-200 text-slate-600 rounded-lg hover:bg-slate-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isPending}
              className="flex-1 px-4 py-2 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 disabled:opacity-50"
            >
              {isPending ? 'Creating...' : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
