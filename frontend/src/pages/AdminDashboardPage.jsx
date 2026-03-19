import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { adminApi } from '../api';
import { toast } from 'react-toastify';
import { LogOut, Users, Clock, CheckCircle, XCircle, Calendar, LayoutDashboard } from 'lucide-react';

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
    return (
      <div className="min-h-screen flex items-center justify-center text-sm font-semibold text-[var(--ink-soft)]">
        Loading dashboard...
      </div>
    );
  }

  return (
    <div className="min-h-screen pb-10">
      {/* Header */}
      <nav className="glass-panel sticky top-0 z-30 border-b border-[var(--line-soft)]">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex min-h-16 flex-col justify-center gap-3 py-3 sm:min-h-20 sm:flex-row sm:items-center sm:justify-between sm:py-0">
            <div className="flex items-center gap-2">
              <div className="rounded-xl bg-[var(--brand)]/15 p-2 text-[var(--brand)]">
                <LayoutDashboard className="h-5 w-5" />
              </div>
              <h1 className="section-title text-lg font-bold text-[var(--ink-strong)] sm:text-xl">
                {user?.role === 'manager' ? 'Manager Dashboard' : 'Admin Dashboard'}
              </h1>
            </div>

            <div className="flex flex-wrap items-center gap-2 sm:gap-3">
              <button
                onClick={() => navigate('/dashboard')}
                className="ghost-btn rounded-xl px-3 py-2 text-xs font-semibold sm:text-sm"
              >
                Employee View
              </button>
              {user?.role === 'admin' && (
                <button
                  onClick={() => navigate('/admin/users')}
                  className="ghost-btn inline-flex items-center rounded-xl px-3 py-2 text-xs font-semibold sm:text-sm"
                >
                  <Users className="w-4 h-4 mr-1" />
                  Manage Users
                </button>
              )}
              <span className="rounded-full border border-[var(--line)] bg-white/80 px-3 py-1.5 text-xs font-semibold text-[var(--ink-soft)] sm:text-sm">
                {user?.first_name} {user?.last_name} ({user?.role})
              </span>
              <button
                onClick={() => {
                  logout();
                  navigate('/login');
                }}
                className="ghost-btn inline-flex items-center rounded-xl px-3 py-2 text-xs font-semibold sm:text-sm"
              >
                <LogOut className="w-4 h-4 mr-1" />
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 py-6 sm:px-6 sm:py-8 lg:px-8">
        {/* Stats Cards */}
        {stats && (
          <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
            <div className="surface-panel rounded-2xl p-5">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold text-[var(--ink-soft)]">Pending Requests</h3>
                <Clock className="h-6 w-6 text-amber-500" />
              </div>
              <p className="section-title mt-2 text-3xl font-extrabold text-[var(--ink-strong)]">{stats.pending_requests}</p>
            </div>

            <div className="surface-panel rounded-2xl p-5">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold text-[var(--ink-soft)]">Approved</h3>
                <CheckCircle className="h-6 w-6 text-emerald-500" />
              </div>
              <p className="section-title mt-2 text-3xl font-extrabold text-[var(--ink-strong)]">{stats.approved_requests}</p>
            </div>

            <div className="surface-panel rounded-2xl p-5">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold text-[var(--ink-soft)]">Rejected</h3>
                <XCircle className="h-6 w-6 text-rose-500" />
              </div>
              <p className="section-title mt-2 text-3xl font-extrabold text-[var(--ink-strong)]">{stats.rejected_requests}</p>
            </div>

            <div className="surface-panel rounded-2xl p-5">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold text-[var(--ink-soft)]">On Leave Today</h3>
                <Calendar className="h-6 w-6 text-[var(--brand)]" />
              </div>
              <p className="section-title mt-2 text-3xl font-extrabold text-[var(--ink-strong)]">{stats.employees_on_leave}</p>
            </div>
          </div>
        )}

        {/* Filter Tabs */}
        <div className="surface-panel mb-6 rounded-2xl p-2">
            <nav className="flex flex-wrap gap-2">
              {['pending', 'approved', 'rejected', ''].map((status) => (
                <button
                  key={status}
                  onClick={() => setFilter(status)}
                  className={`${
                    filter === status
                      ? 'brand-btn'
                      : 'ghost-btn'
                  } whitespace-nowrap rounded-xl px-4 py-2 text-xs font-semibold sm:px-5 sm:text-sm`}
                >
                  {status === '' ? 'All' : status.charAt(0).toUpperCase() + status.slice(1)}
                </button>
              ))}
            </nav>
        </div>

        {/* Leave Requests Table */}
        <div className="surface-panel overflow-hidden rounded-2xl">
          <div className="border-b border-[var(--line-soft)] px-5 py-4 sm:px-6">
            <h3 className="section-title text-lg font-bold text-[var(--ink-strong)]">Leave Requests</h3>
          </div>
          <div className="hidden overflow-x-auto md:block soft-scrollbar">
            <table className="min-w-full divide-y divide-[var(--line-soft)]">
              <thead className="bg-[var(--bg-muted)]/55">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-bold uppercase tracking-[0.08em] text-[var(--ink-soft)]">
                    Employee
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-bold uppercase tracking-[0.08em] text-[var(--ink-soft)]">
                    Type
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-bold uppercase tracking-[0.08em] text-[var(--ink-soft)]">
                    Start Date
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-bold uppercase tracking-[0.08em] text-[var(--ink-soft)]">
                    End Date
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-bold uppercase tracking-[0.08em] text-[var(--ink-soft)]">
                    Days
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-bold uppercase tracking-[0.08em] text-[var(--ink-soft)]">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-bold uppercase tracking-[0.08em] text-[var(--ink-soft)]">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--line-soft)] bg-white/90">
                {leaves.length === 0 ? (
                  <tr>
                    <td colSpan="7" className="px-6 py-5 text-center text-sm text-[var(--ink-soft)]">
                      No leave requests found
                    </td>
                  </tr>
                ) : (
                  leaves.map((leave) => (
                    <tr key={leave.id}>
                      <td className="whitespace-nowrap px-6 py-4">
                        <div className="text-sm font-semibold text-[var(--ink-strong)]">
                          {leave.user?.full_name}
                        </div>
                        <div className="text-sm text-[var(--ink-soft)]">{leave.user?.email}</div>
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm font-semibold text-[var(--ink-strong)]">
                        {leave.leave_type?.name}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm text-[var(--ink-soft)]">
                        {leave.start_date}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm text-[var(--ink-soft)]">
                        {leave.end_date}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm font-semibold text-[var(--ink-strong)]">
                        {leave.total_days}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4">
                        <span
                          className={`fancy-badge inline-flex leading-5 ${getStatusColor(
                            leave.status
                          )}`}
                        >
                          {leave.status}
                        </span>
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm font-medium">
                        {leave.status === 'pending' && (
                          <div className="flex space-x-2">
                            <button
                              onClick={() =>
                                setReviewModal({ show: true, leave, action: 'approve' })
                              }
                              className="rounded-lg bg-emerald-100 px-2.5 py-1 text-xs font-semibold text-emerald-800 transition hover:bg-emerald-200"
                            >
                              Approve
                            </button>
                            <button
                              onClick={() => setReviewModal({ show: true, leave, action: 'reject' })}
                              className="rounded-lg bg-rose-100 px-2.5 py-1 text-xs font-semibold text-rose-800 transition hover:bg-rose-200"
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

          <div className="space-y-3 p-4 md:hidden">
            {leaves.length === 0 ? (
              <div className="rounded-xl border border-[var(--line-soft)] bg-white/80 p-4 text-center text-sm text-[var(--ink-soft)]">
                No leave requests found
              </div>
            ) : (
              leaves.map((leave) => (
                <div key={leave.id} className="rounded-xl border border-[var(--line-soft)] bg-white/80 p-4">
                  <div className="mb-2 flex items-start justify-between gap-2">
                    <div>
                      <p className="text-sm font-semibold text-[var(--ink-strong)]">{leave.user?.full_name}</p>
                      <p className="text-xs text-[var(--ink-soft)]">{leave.user?.email}</p>
                    </div>
                    <span className={`fancy-badge ${getStatusColor(leave.status)}`}>{leave.status}</span>
                  </div>
                  <p className="text-sm text-[var(--ink-soft)]">{leave.leave_type?.name}</p>
                  <p className="mt-1 text-xs text-[var(--ink-soft)]">{leave.start_date} to {leave.end_date} ({leave.total_days} day(s))</p>

                  {leave.status === 'pending' && (
                    <div className="mt-3 flex items-center gap-2">
                      <button
                        onClick={() => setReviewModal({ show: true, leave, action: 'approve' })}
                        className="rounded-lg bg-emerald-100 px-3 py-1.5 text-xs font-semibold text-emerald-800"
                      >
                        Approve
                      </button>
                      <button
                        onClick={() => setReviewModal({ show: true, leave, action: 'reject' })}
                        className="rounded-lg bg-rose-100 px-3 py-1.5 text-xs font-semibold text-rose-800"
                      >
                        Reject
                      </button>
                    </div>
                  )}
                </div>
              ))
            )}
          </div>
        </div>
      </div>

      {/* Review Modal */}
      {reviewModal.show && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-[#10292e]/45 p-4 backdrop-blur-sm">
          <div className="surface-panel max-h-[90vh] w-full max-w-md overflow-y-auto rounded-2xl p-6 sm:p-7">
            <h3 className="section-title mb-4 text-xl font-bold text-[var(--ink-strong)]">
              {reviewModal.action === 'approve' ? 'Approve' : 'Reject'} Leave Request
            </h3>
            <div className="mb-4 rounded-xl border border-[var(--line-soft)] bg-white/75 p-4 text-sm text-[var(--ink-soft)]">
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
              <label className="mb-2 block text-sm font-semibold text-[var(--ink-soft)]">
                Review Notes (Optional)
              </label>
              <textarea
                value={reviewNotes}
                onChange={(e) => setReviewNotes(e.target.value)}
                rows={3}
                className="input-base"
                placeholder="Add any notes about your decision..."
              />
            </div>
            <div className="flex space-x-3 mt-6">
              <button
                onClick={handleReview}
                className={`flex-1 px-4 py-2 rounded-md text-white ${
                  reviewModal.action === 'approve'
                    ? 'bg-emerald-600 hover:bg-emerald-700'
                    : 'bg-rose-600 hover:bg-rose-700'
                }`}
              >
                Confirm {reviewModal.action === 'approve' ? 'Approval' : 'Rejection'}
              </button>
              <button
                onClick={() => {
                  setReviewModal({ show: false, leave: null, action: '' });
                  setReviewNotes('');
                }}
                className="ghost-btn flex-1 rounded-xl px-4 py-2"
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
