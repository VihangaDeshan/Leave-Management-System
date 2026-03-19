import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { leaveApi } from '../api';
import { toast } from 'react-toastify';
import { LogOut, Calendar, Plus, LayoutDashboard } from 'lucide-react';

const EmployeeDashboardPage = () => {
  const { user, logout, isManager } = useAuth();
  const navigate = useNavigate();
  const [leaves, setLeaves] = useState([]);
  const [balances, setBalances] = useState([]);
  const [leaveTypes, setLeaveTypes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    leave_type_id: '',
    start_date: '',
    end_date: '',
    reason: '',
  });

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [leavesRes, balancesRes, typesRes] = await Promise.all([
        leaveApi.getUserLeaves(),
        leaveApi.getUserBalance(),
        leaveApi.getLeaveTypes(),
      ]);
      setLeaves(leavesRes.data.leave_requests || []);
      setBalances(balancesRes.data || []);
      setLeaveTypes(typesRes.data || []);
    } catch (error) {
      toast.error('Failed to fetch data');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await leaveApi.createLeave(formData);
      toast.success('Leave request submitted successfully!');
      setShowModal(false);
      setFormData({ leave_type_id: '', start_date: '', end_date: '', reason: '' });
      fetchData();
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to submit leave request');
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
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
              <div>
                <h1 className="section-title text-lg font-bold text-[var(--ink-strong)] sm:text-xl">Leave Management</h1>
                <p className="text-xs text-[var(--ink-soft)]">Employee Workspace</p>
              </div>
            </div>

            <div className="flex flex-wrap items-center gap-2 sm:gap-3">
              {isManager() && (
                <button
                  onClick={() => navigate('/admin')}
                  className="ghost-btn rounded-xl px-3 py-2 text-xs font-semibold sm:text-sm"
                >
                  Admin Panel
                </button>
              )}
              <span className="rounded-full border border-[var(--line)] bg-white/80 px-3 py-1.5 text-xs font-semibold text-[var(--ink-soft)] sm:text-sm">
                {user?.first_name} {user?.last_name}
              </span>
              <button
                onClick={handleLogout}
                className="ghost-btn inline-flex items-center rounded-xl px-3 py-2 text-xs font-semibold sm:text-sm"
              >
                <LogOut className="w-4 h-4 mr-1" />
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
        {/* Dashboard Header */}
        <div className="mb-6 rounded-3xl border border-[var(--line-soft)] bg-white/70 p-5 shadow-[0_16px_36px_-26px_rgba(16,41,46,0.45)] sm:p-7">
          <h2 className="section-title text-2xl font-extrabold text-[var(--ink-strong)] sm:text-3xl">Welcome back, {user?.first_name}.</h2>
          <p className="mt-2 text-sm text-[var(--ink-soft)] sm:text-base">Track balances, request time off, and monitor approval status from one place.</p>
        </div>

        {/* Leave Balance Cards */}
        <div className="mb-8 grid grid-cols-1 gap-4 sm:gap-5 md:grid-cols-2 xl:grid-cols-3">
          {balances.map((balance) => (
            <div key={balance.id} className="surface-panel fade-in-up rounded-2xl p-5">
              <div className="flex items-center justify-between">
                <h3 className="text-base font-bold text-[var(--ink-strong)] sm:text-lg">{balance.leave_type?.name}</h3>
                <Calendar className="h-5 w-5 text-[var(--brand)] sm:h-6 sm:w-6" />
              </div>
              <div className="mt-4">
                <p className="section-title text-3xl font-extrabold text-[var(--brand)]">{balance.available_days}</p>
                <p className="mt-1 text-sm text-[var(--ink-soft)]">
                  Available out of {balance.total_days} days
                </p>
                <p className="mt-2 text-xs font-semibold uppercase tracking-[0.08em] text-[var(--ink-soft)]">Used: {balance.used_days} days</p>
              </div>
            </div>
          ))}
        </div>

        {/* Actions */}
        <div className="mb-6">
          <button
            onClick={() => setShowModal(true)}
            className="brand-btn inline-flex items-center gap-2 rounded-xl px-5 py-3 text-sm font-semibold"
          >
            <Plus className="h-4 w-4" />
            Request Leave
          </button>
        </div>

        {/* Leave Requests Table */}
        <div className="surface-panel overflow-hidden rounded-2xl">
          <div className="border-b border-[var(--line-soft)] px-5 py-4 sm:px-6">
            <h3 className="section-title text-lg font-bold text-[var(--ink-strong)]">My Leave Requests</h3>
          </div>
          <div className="hidden overflow-x-auto md:block soft-scrollbar">
            <table className="min-w-full divide-y divide-[var(--line-soft)]">
              <thead className="bg-[var(--bg-muted)]/55">
                <tr>
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
                    Reason
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--line-soft)] bg-white/90">
                {leaves.length === 0 ? (
                  <tr>
                    <td colSpan="6" className="px-6 py-5 text-center text-sm text-[var(--ink-soft)]">
                      No leave requests found
                    </td>
                  </tr>
                ) : (
                  leaves.map((leave) => (
                    <tr key={leave.id}>
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
                      <td className="max-w-xs truncate px-6 py-4 text-sm text-[var(--ink-soft)]">
                        {leave.reason}
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
                    <p className="font-semibold text-[var(--ink-strong)]">{leave.leave_type?.name}</p>
                    <span className={`fancy-badge ${getStatusColor(leave.status)}`}>{leave.status}</span>
                  </div>
                  <p className="text-xs text-[var(--ink-soft)]">{leave.start_date} to {leave.end_date}</p>
                  <p className="mt-1 text-sm text-[var(--ink-soft)]">{leave.total_days} day(s)</p>
                  <p className="mt-2 text-sm text-[var(--ink-soft)]">{leave.reason}</p>
                </div>
              ))
            )}
          </div>
        </div>
      </div>

      {/* Leave Request Modal */}
      {showModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-[#10292e]/45 p-4 backdrop-blur-sm">
          <div className="surface-panel max-h-[90vh] w-full max-w-md overflow-y-auto rounded-2xl p-6 sm:p-7">
            <h3 className="section-title mb-4 text-xl font-bold text-[var(--ink-strong)]">Request Leave</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-semibold text-[var(--ink-soft)]">Leave Type</label>
                <select
                  required
                  value={formData.leave_type_id}
                  onChange={(e) => setFormData({ ...formData, leave_type_id: parseInt(e.target.value, 10) || '' })}
                  className="input-base mt-1"
                >
                  <option value="">Select leave type</option>
                  {leaveTypes.map((type) => (
                    <option key={type.id} value={type.id}>
                      {type.name}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-semibold text-[var(--ink-soft)]">Start Date</label>
                <input
                  type="date"
                  required
                  value={formData.start_date}
                  onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                  className="input-base mt-1"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-[var(--ink-soft)]">End Date</label>
                <input
                  type="date"
                  required
                  value={formData.end_date}
                  onChange={(e) => setFormData({ ...formData, end_date: e.target.value })}
                  className="input-base mt-1"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-[var(--ink-soft)]">Reason</label>
                <textarea
                  required
                  minLength={10}
                  rows={3}
                  value={formData.reason}
                  onChange={(e) => setFormData({ ...formData, reason: e.target.value })}
                  className="input-base mt-1"
                  placeholder="Please provide a reason for your leave"
                />
              </div>

              <div className="flex space-x-3 pt-4">
                <button
                  type="submit"
                  className="brand-btn flex-1 rounded-xl px-4 py-2.5 text-sm font-semibold"
                >
                  Submit Request
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setShowModal(false);
                    setFormData({ leave_type_id: '', start_date: '', end_date: '', reason: '' });
                  }}
                  className="ghost-btn flex-1 rounded-xl px-4 py-2.5 text-sm font-semibold"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default EmployeeDashboardPage;

