import { Outlet, Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function Layout() {
  const { user, logout } = useAuth()

  return (
    <div className="min-h-screen">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
          <Link to="/" className="text-2xl font-bold text-primary-600">
            RevUExchange
          </Link>
          <nav className="flex items-center gap-4">
            <Link to="/bounties" className="text-gray-600 hover:text-gray-900">
              Bounties
            </Link>
            {user ? (
              <>
                <Link to="/dashboard" className="text-gray-600 hover:text-gray-900">
                  Dashboard
                </Link>
                <Link to="/reviews" className="text-gray-600 hover:text-gray-900">
                  Reviews
                </Link>
                <Link to="/profile" className="text-gray-600 hover:text-gray-900">
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
                <button
                  onClick={logout}
                  className="text-gray-600 hover:text-gray-900"
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
      <main className="max-w-7xl mx-auto px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}
