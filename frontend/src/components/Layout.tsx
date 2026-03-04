import { useState } from 'react'
import { Outlet, Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { Menu, X, Star, LogOut, BarChart2, MessageCircle, Trophy, BookOpen, User, Home, Package } from 'lucide-react'

const navItems = [
  { to: '/dashboard', label: 'Dashboard', icon: Home },
  { to: '/products', label: 'Products', icon: Package },
  { to: '/bounties', label: 'Bounties', icon: BookOpen },
  { to: '/reviews', label: 'Reviews', icon: Star },
  { to: '/feed', label: 'Feed', icon: MessageCircle },
  { to: '/leaderboard', label: 'Leaderboard', icon: Trophy },
  { to: '/analytics', label: 'Analytics', icon: BarChart2 },
  { to: '/profile', label: 'Profile', icon: User },
]

export default function Layout() {
  const { user, logout } = useAuth()
  const [isMenuOpen, setIsMenuOpen] = useState(false)

  return (
    <div className="min-h-screen">
      {/* Mobile hamburger button - left side */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
          <div className="flex items-center gap-4">
            {user && (
              <button
                onClick={() => setIsMenuOpen(!isMenuOpen)}
                className="lg:hidden p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg"
              >
                {isMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
              </button>
            )}
            <Link to="/" className="text-2xl font-bold text-primary-600">
              RevUExchange
            </Link>
          </div>
          <nav className="flex items-center gap-4">
            <Link to="/bounties" className="text-gray-600 hover:text-gray-900 hidden md:block">
              Bounties
            </Link>
            {user ? (
              <>
                <Link to="/dashboard" className="text-gray-600 hover:text-gray-900 hidden md:block">
                  Dashboard
                </Link>
                <Link to="/reviews" className="text-gray-600 hover:text-gray-900 hidden md:block">
                  Reviews
                </Link>
                <Link to="/feed" className="text-gray-600 hover:text-gray-900 hidden md:block">
                  Feed
                </Link>
                <Link to="/leaderboard" className="text-gray-600 hover:text-gray-900 hidden md:block">
                  Leaderboard
                </Link>
                <Link to="/analytics" className="text-gray-600 hover:text-gray-900 hidden md:block">
                  Analytics
                </Link>
                <Link to="/profile" className="text-gray-600 hover:text-gray-900 hidden md:block">
                  Profile
                </Link>
                <Link
                  to="/points"
                  className="flex items-center gap-1 bg-teal-50 text-teal-700 px-3 py-1.5 rounded-full hover:bg-teal-100 transition-colors"
                >
                  <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.736 6.979C9.208 6.193 9.696 6 10 6c.304 0 .792.193 1.264.979a1 1 0 001.715-1.029C12.279 4.784 11.232 4 10 4s-2.279.784-2.979 1.95c-.285.475-.507 1-.67 1.55H6a1 1 0 000 2h.013a9.358 9.358 0 000 1H6a1 1 0 100 2h.351c.163.55.385 1.075.67 1.55C7.721 15.216 8.768 16 10 16s2.279-.784 2.979-1.95c.285-.475.507-1 .67-1.55H14a1 1 0 100-2h-.013a9.358 9.358 0 00-1.351-4.5H14a1 1 0 100-2h-.013a9.358 9.358 0 00-1.351-4.5H10z" />
                  </svg>
                  <span className="font-semibold">{user.points}</span>
                </Link>
                <Link
                  to="/purchase-points"
                  className="flex items-center gap-1 bg-gradient-to-r from-violet-500 to-purple-600 text-white px-3 py-1.5 rounded-full hover:from-violet-600 hover:to-purple-700 transition-all text-sm font-semibold"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                  </svg>
                  Buy
                </Link>
                <button
                  onClick={logout}
                  className="text-gray-600 hover:text-gray-900 hidden md:block"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link to="/login" className="text-gray-600 hover:text-gray-900">
                  Login
                </Link>
                <Link
                  to="/register"
                  className="bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700"
                >
                  Sign Up
                </Link>
              </>
            )}
          </nav>
        </div>
      </header>

      {/* Mobile hamburger menu - slides in from left */}
      {user && isMenuOpen && (
        <div className="lg:hidden fixed inset-0 z-50">
          {/* Backdrop */}
          <div
            className="absolute inset-0 bg-black/50"
            onClick={() => setIsMenuOpen(false)}
          />
          {/* Menu panel */}
          <div className="absolute left-0 top-0 h-full w-64 bg-white shadow-xl overflow-y-auto">
            <div className="p-4 border-b border-gray-100">
              <div className="flex items-center justify-between">
                <span className="font-semibold text-gray-800">Menu</span>
                <button
                  onClick={() => setIsMenuOpen(false)}
                  className="p-2 text-gray-500 hover:text-gray-700"
                >
                  <X className="w-5 h-5" />
                </button>
              </div>
            </div>
            <nav className="p-4 space-y-1">
              {navItems.map((item) => (
                <Link
                  key={item.to}
                  to={item.to}
                  onClick={() => setIsMenuOpen(false)}
                  className="flex items-center gap-3 px-4 py-3 text-gray-600 hover:text-gray-900 hover:bg-gray-50 rounded-lg transition-colors"
                >
                  <item.icon className="w-5 h-5" />
                  {item.label}
                </Link>
              ))}
              <hr className="my-4 border-gray-100" />
              <button
                onClick={() => {
                  logout()
                  setIsMenuOpen(false)
                }}
                className="flex items-center gap-3 px-4 py-3 w-full text-left text-red-600 hover:bg-red-50 rounded-lg transition-colors"
              >
                <LogOut className="w-5 h-5" />
                Logout
              </button>
            </nav>
          </div>
        </div>
      )}

      <main className="max-w-7xl mx-auto px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}
