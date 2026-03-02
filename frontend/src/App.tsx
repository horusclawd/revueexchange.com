import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './context/AuthContext'
import Layout from './components/Layout'
import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import Bounties from './pages/Bounties'
import MyReviews from './pages/MyReviews'
import Profile from './pages/Profile'
import Points from './pages/Points'
import PurchasePoints from './pages/PurchasePoints'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth()

  if (loading) {
    return <div>Loading...</div>
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Home />} />
        <Route path="login" element={<Login />} />
        <Route path="register" element={<Register />} />
        <Route
          path="dashboard"
          element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          }
        />
        <Route
          path="bounties"
          element={
            <ProtectedRoute>
              <Bounties />
            </ProtectedRoute>
          }
        />
        <Route
          path="reviews"
          element={
            <ProtectedRoute>
              <MyReviews />
            </ProtectedRoute>
          }
        />
        <Route
          path="profile"
          element={
            <ProtectedRoute>
              <Profile />
            </ProtectedRoute>
          }
        />
        <Route
          path="points"
          element={
            <ProtectedRoute>
              <Points />
            </ProtectedRoute>
          }
        />
        <Route
          path="purchase-points"
          element={
            <ProtectedRoute>
              <PurchasePoints />
            </ProtectedRoute>
          }
        />
      </Route>
    </Routes>
  )
}

export default App
