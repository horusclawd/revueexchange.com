import { Link } from 'react-router-dom'

export default function Home() {
  return (
    <div className="text-center py-12">
      <h1 className="text-4xl font-bold text-gray-900 mb- Points by Reviewing,4">
        Earn Get Reviews by Participating
      </h1>
      <p className="text-xl text-gray-600 mb-8">
        The reciprocal review platform for self-published authors and digital creators.
      </p>
      <div className="flex gap-4 justify-center">
        <Link
          to="/register"
          className="bg-primary-600 text-white px-6 py-3 rounded-lg text-lg hover:bg-primary-700"
        >
          Get Started
        </Link>
        <Link
          to="/bounties"
          className="border border-primary-600 text-primary-600 px-6 py-3 rounded-lg text-lg hover:bg-primary-50"
        >
          Browse Bounties
        </Link>
      </div>

      <div className="mt-16 grid md:grid-cols-3 gap-8">
        <div className="p-6 bg-white rounded-lg shadow">
          <h3 className="text-xl font-semibold mb-2">Earn Points</h3>
          <p className="text-gray-600">
            Write quality reviews and earn points that can be used to get your own work reviewed.
          </p>
        </div>
        <div className="p-6 bg-white rounded-lg shadow">
          <h3 className="text-xl font-semibold mb-2">Get Reviewed</h3>
          <p className="text-gray-600">
            Create bounties to attract reviewers for your books, courses, podcasts, and newsletters.
          </p>
        </div>
        <div className="p-6 bg-white rounded-lg shadow">
          <h3 className="text-xl font-semibold mb-2">Build Reputation</h3>
          <p className="text-gray-600">
            Earn badges and climb the leaderboard as you help the community.
          </p>
        </div>
      </div>
    </div>
  )
}
