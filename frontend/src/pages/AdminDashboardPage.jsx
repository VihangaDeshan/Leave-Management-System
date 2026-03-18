import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { adminApi } from '../api';
import { toast } from 'react-toastify';
import { LogOut, Users, Clock, CheckCircle, XCircle, Calendar } from 'lucide-react';

const AdminDashboardPage = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [leaves, setLeaves] = useState([]);
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('pending');
  const [reviewModal, setReviewModal] = useState({ show: false, leave: null, action: '' });
  const [reviewNotes, setReviewNotes] = useState('');

  useEffect(() => {
    fetchData();
  }, [filter]);

  const fetchData = async () => {
    try {
      const [leavesRes, statsRes] = await Promise.all([
        adminApi.getAllLeaves(filter),
        adminApi.getDashboardStats(),
      ]);
      setLeaves(leavesRes.data.leave_requests || []);
      setStats(statsRes.data);
    } catch (error) {
      toast.error('Failed to fetch data');
    } finally {
      setLoading(false);
    }
  };

  const handleReview = async () => {
    try {
      if (reviewModal.action === 'approve') {
        await adminApi.approveLeave(reviewModal.leave.id, reviewNotes);
        toast.success('Leave request approved!');
      } else {
        await adminApi.rejectLeave(reviewModal.leave.id, reviewNotes);
        toast.success('Leave request rejected!');
      }
      setReviewModal({ show: false, leave: null, action: '' });
      setReviewNotes('');
      fetchData();
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to process request');
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'approved':
        return 'bg-green-100 text-green-800';
      case 'rejected':
        return 'bg-red-100 text-red-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  if (loading) {
    return <div className="min-h-screen flex items-center justify-center">Loading...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-bold text-gray-900">Admin Dashboard</h1>
            </div>
            <div className="flex items-center space-x-4">
              <button
                onClick={() => navigate('/dashboard')}
                className="text-sm text-blue-600 hover:text-blue-800"
              >
                Employee View
              </button>
              <span className="text-sm text-gray-700">
                {user?.first_name} {user?.last_name} ({user?.role})
              </span>
              <button
                onClick={() => {
                  logout();
                  navigate('/login');
                }}
                className="flex items-center text-sm text-gray-700 hover:text-gray-900"
              >
                <LogOut className="w-4 h-4 mr-1" />
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Stats Cards */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-medium text-gray-600">Pending Requests</h3>
                <Clock className="w-6 h-6 text-yellow-500" />
              </div>
              <p className="text-3xl font-bold text-gray-900 mt-2">{stats.pending_requests}</p>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-medium text-gray-600">Approved</h3>
                <CheckCircle className="w-6 h-6 text-green-500" />
              </div>
              <p className="text-3xl font-bold text-gray-900 mt-2">{stats.approved_requests}</p>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-medium text-gray-600">Rejected</h3>
                <XCircle className="w-6 h-6 text-red-500" />
              </div>
              <p className="text-3xl font-bold text-gray-900 mt-2">{stats.rejected_requests}</p>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-medium text-gray-600">On Leave Today</h3>
                <Calendar className="w-6 h-6 text-blue-500" />
              </div>
              <p className="text-3xl font-bold text-gray-900 mt-2">{stats.employees_on_leave}</p>
            </div>
          </div>
        )}

        {/* Filter Tabs */}
        <div className="bg-white rounded-lg shadow mb-6">
          <div className="border-b border-gray-200">
            <nav className="flex -mb-px">
              {['pending', 'approved', 'rejected', ''].map((status) => (
                <button
                  key={status}
                  onClick={() => setFilter(status)}
                  className={`${
                    filter === status
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  } whitespace-nowrap py-4 px-6 border-b-2 font-medium text-sm`}
                >
                  {status === '' ? 'All' : status.charAt(0).toUpperCase() + status.slice(1)}
                </button>
              ))}
            </nav>
          </div>
        </div>

        {/* Leave Requests Table */}
        <div className="bg-white shadow rounded-lg overflow-hidden">
          <div className="px-6 py-4 border-b border-gray-200">
            <h3 className="text-lg font-semibold text-gray-900">Leave Requests</h3>
          </div>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Employee
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Type
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Start Date
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    End Date
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Days
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {leaves.length === 0 ? (
                  <tr>
                    <td colSpan="7" className="px-6 py-4 text-center text-gray-500">
                      No leave requests found
                    </td>
                  </tr>
                ) : (
                  leaves.map((leave) => (
                    <tr key={leave.id}>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm font-medium text-gray-900">
                          {leave.user?.full_name}
                        </div>
                        <div className="text-sm text-gray-500">{leave.user?.email}</div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {leave.leave_type?.name}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {leave.start_date}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {leave.end_date}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {leave.total_days}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span
                          className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(
                            leave.status
                          )}`}
                        >
                          {leave.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        {leave.status === 'pending' && (
                          <div className="flex space-x-2">
                            <button
                              onClick={() =>
                                setReviewModal({ show: true, leave, action: 'approve' })
                              }
                              className="text-green-600 hover:text-green-900"
                            >
                              Approve
                            </button>
                            <button
                              onClick={() => setReviewModal({ show: true, leave, action: 'reject' })}
                              className="text-red-600 hover:text-red-900"
                            >
                              Reject
                            </button>
                          </div>
                        )}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* Review Modal */}
      {reviewModal.show && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-8 max-w-md w-full mx-4">
            <h3 className="text-xl font-bold text-gray-900 mb-4">
              {reviewModal.action === 'approve' ? 'Approve' : 'Reject'} Leave Request
            </h3>
            <div className="mb-4 text-sm text-gray-600">
              <p>
                <strong>Employee:</strong> {reviewModal.leave?.user?.full_name}
              </p>
              <p>
                <strong>Type:</strong> {reviewModal.leave?.leave_type?.name}
              </p>
              <p>
                <strong>Duration:</strong> {reviewModal.leave?.start_date} to{' '}
                {reviewModal.leave?.end_date} ({reviewModal.leave?.total_days} days)
              </p>
              <p>
                <strong>Reason:</strong> {reviewModal.leave?.reason}
              </p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Review Notes (Optional)
              </label>
              <textarea
                value={reviewNotes}
                onChange={(e) => setReviewNotes(e.target.value)}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Add any notes about your decision..."
              />
            </div>
            <div className="flex space-x-3 mt-6">
              <button
                onClick={handleReview}
                className={`flex-1 px-4 py-2 rounded-md text-white ${
                  reviewModal.action === 'approve'
                    ? 'bg-green-600 hover:bg-green-700'
                    : 'bg-red-600 hover:bg-red-700'
                }`}
              >
                Confirm {reviewModal.action === 'approve' ? 'Approval' : 'Rejection'}
              </button>
              <button
                onClick={() => {
                  setReviewModal({ show: false, leave: null, action: '' });
                  setReviewNotes('');
                }}
                className="flex-1 bg-gray-200 text-gray-800 px-4 py-2 rounded-md hover:bg-gray-300"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminDashboardPage;
