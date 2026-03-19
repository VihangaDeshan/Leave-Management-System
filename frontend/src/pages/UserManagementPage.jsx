import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { adminApi, authApi } from '../api';
import { toast } from 'react-toastify';
import { LogOut, Users, Edit, ArrowLeft, Plus, Trash2 } from 'lucide-react';

const UserManagementPage = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [users, setUsers] = useState([]);
  const [departments, setDepartments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [createModal, setCreateModal] = useState(false);
  const [editModal, setEditModal] = useState({ show: false, user: null });
  const [createFormData, setCreateFormData] = useState({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    role: 'employee',
    department: '',
    manager_id: null,
  });
  const [formData, setFormData] = useState({
    first_name: '',
    last_name: '',
    department: '',
    manager_id: null,
  });

  useEffect(() => {
    fetchInitialData();
  }, []);

  const fetchInitialData = async () => {
    setLoading(true);
    try {
      const [usersRes, deptRes] = await Promise.all([
        adminApi.getAllUsers(1, 100),
        authApi.getDepartments(),
      ]);

      const fetchedUsers = usersRes.data?.users || [];
      setUsers(fetchedUsers);

      const apiDepartments = Array.isArray(deptRes.data) ? deptRes.data : [];
      const userDepartments = fetchedUsers
        .map((u) => u.department)
        .filter((d) => typeof d === 'string' && d.trim() !== '');
      const merged = Array.from(new Set([...apiDepartments, ...userDepartments]));
      setDepartments(merged);
    } catch (error) {
      toast.error('Failed to load user management data');
    } finally {
      setLoading(false);
    }
  };

  const fetchUsers = async () => {
    try {
      const response = await adminApi.getAllUsers(1, 100);
      const fetchedUsers = response.data.users || [];
      setUsers(fetchedUsers);

      const userDepartments = fetchedUsers
        .map((u) => u.department)
        .filter((d) => typeof d === 'string' && d.trim() !== '');
      setDepartments((prev) => Array.from(new Set([...prev, ...userDepartments])));
    } catch (error) {
      toast.error('Failed to fetch users');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateUser = async (e) => {
    e.preventDefault();

    if (!createFormData.email || !createFormData.password || !createFormData.first_name || !createFormData.last_name) {
      toast.error('Email, password, first name, and last name are required');
      return;
    }

    try {
      const payload = {
        email: createFormData.email,
        password: createFormData.password,
        first_name: createFormData.first_name,
        last_name: createFormData.last_name,
        role: createFormData.role,
        department: createFormData.department || null,
        manager_id: createFormData.manager_id || null,
      };

      await adminApi.createUser(payload);
      toast.success('User created successfully!');
      setCreateModal(false);
      setCreateFormData({
        email: '',
        password: '',
        first_name: '',
        last_name: '',
        role: 'employee',
        department: '',
        manager_id: null,
      });
      fetchUsers();
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to create user');
    }
  };

  const handleDeleteUser = async (userToDelete) => {
    if (!window.confirm(`Deactivate user ${userToDelete.first_name} ${userToDelete.last_name}?`)) {
      return;
    }

    try {
      await adminApi.deleteUser(userToDelete.id);
      toast.success('User deactivated successfully!');
      fetchUsers();
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to deactivate user');
    }
  };

  const handleEditClick = (userToEdit) => {
    setEditModal({ show: true, user: userToEdit });
    setFormData({
      first_name: userToEdit.first_name || '',
      last_name: userToEdit.last_name || '',
      department: userToEdit.department || '',
      manager_id: userToEdit.manager_id || null,
    });
  };

  const handleUpdateUser = async (e) => {
    e.preventDefault();
    try {
      // Only send fields that have values
      const updateData = {};
      if (formData.first_name) updateData.first_name = formData.first_name;
      if (formData.last_name) updateData.last_name = formData.last_name;
      if (formData.department) updateData.department = formData.department;
      if (formData.manager_id) updateData.manager_id = parseInt(formData.manager_id);

      await adminApi.updateUser(editModal.user.id, updateData);
      toast.success('User updated successfully!');
      setEditModal({ show: false, user: null });
      fetchUsers();
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to update user');
    }
  };

  const getAvailableManagers = () => {
    // Filter out the current user being edited and only show managers/admins
    return users.filter(
      (u) =>
        (u.role === 'manager' || u.role === 'admin') &&
        u.id !== editModal.user?.id
    );
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center text-sm font-semibold text-[var(--ink-soft)]">
        Loading users...
      </div>
    );
  }

  return (
    <div className="min-h-screen pb-10">
      {/* Header */}
      <nav className="glass-panel sticky top-0 z-30 border-b border-[var(--line-soft)]">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex min-h-16 flex-col justify-center gap-3 py-3 sm:min-h-20 sm:flex-row sm:items-center sm:justify-between sm:py-0">
            <div className="flex flex-wrap items-center gap-2 sm:gap-3">
              <button
                onClick={() => navigate('/admin')}
                className="ghost-btn inline-flex items-center rounded-xl px-3 py-2 text-xs font-semibold sm:text-sm"
              >
                <ArrowLeft className="w-4 h-4 mr-1" />
                Back to Dashboard
              </button>
              <h1 className="section-title text-lg font-bold text-[var(--ink-strong)] sm:text-xl">User Management</h1>
            </div>

            <div className="flex flex-wrap items-center gap-2 sm:gap-3">
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
        {/* Users Table */}
        <div className="surface-panel overflow-hidden rounded-2xl">
          <div className="flex flex-wrap items-center justify-between gap-3 border-b border-[var(--line-soft)] px-5 py-4 sm:px-6">
            <div className="flex items-center">
              <Users className="w-5 h-5 mr-2 text-[var(--brand)]" />
              <h3 className="section-title text-lg font-bold text-[var(--ink-strong)]">All Users</h3>
            </div>
            <div className="flex items-center gap-2 sm:gap-3">
              <span className="text-xs font-semibold uppercase tracking-[0.08em] text-[var(--ink-soft)] sm:text-sm">{users.length} total users</span>
              <button
                onClick={() => setCreateModal(true)}
                className="brand-btn inline-flex items-center rounded-xl px-3 py-2 text-xs font-semibold sm:text-sm"
              >
                <Plus className="w-4 h-4 mr-1" />
                Add User
              </button>
            </div>
          </div>
          <div className="hidden overflow-x-auto md:block soft-scrollbar">
            <table className="min-w-full divide-y divide-[var(--line-soft)]">
              <thead className="bg-[var(--bg-muted)]/55">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-bold uppercase tracking-[0.08em] text-[var(--ink-soft)]">
                    Name
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-bold uppercase tracking-[0.08em] text-[var(--ink-soft)]">
                    Email
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-bold uppercase tracking-[0.08em] text-[var(--ink-soft)]">
                    Role
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-bold uppercase tracking-[0.08em] text-[var(--ink-soft)]">
                    Department
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-bold uppercase tracking-[0.08em] text-[var(--ink-soft)]">
                    Manager
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
                {users.length === 0 ? (
                  <tr>
                    <td colSpan="7" className="px-6 py-5 text-center text-sm text-[var(--ink-soft)]">
                      No users found
                    </td>
                  </tr>
                ) : (
                  users.map((u) => {
                    const manager = users.find((m) => m.id === u.manager_id);
                    return (
                      <tr key={u.id}>
                        <td className="whitespace-nowrap px-6 py-4">
                          <div className="text-sm font-semibold text-[var(--ink-strong)]">
                            {u.first_name} {u.last_name}
                          </div>
                        </td>
                        <td className="whitespace-nowrap px-6 py-4 text-sm text-[var(--ink-soft)]">
                          {u.email}
                        </td>
                        <td className="whitespace-nowrap px-6 py-4">
                          <span
                            className={`fancy-badge inline-flex leading-5 ${
                              u.role === 'admin'
                                ? 'bg-violet-100 text-violet-800'
                                : u.role === 'manager'
                                ? 'bg-blue-100 text-blue-800'
                                : 'bg-green-100 text-green-800'
                            }`}
                          >
                            {u.role}
                          </span>
                        </td>
                        <td className="whitespace-nowrap px-6 py-4 text-sm text-[var(--ink-soft)]">
                          {u.department || '-'}
                        </td>
                        <td className="whitespace-nowrap px-6 py-4 text-sm text-[var(--ink-soft)]">
                          {manager ? `${manager.first_name} ${manager.last_name}` : '-'}
                        </td>
                        <td className="whitespace-nowrap px-6 py-4">
                          <span
                            className={`fancy-badge inline-flex leading-5 ${
                              u.is_active
                                ? 'bg-green-100 text-green-800'
                                : 'bg-red-100 text-red-800'
                            }`}
                          >
                            {u.is_active ? 'Active' : 'Inactive'}
                          </span>
                        </td>
                        <td className="whitespace-nowrap px-6 py-4 text-sm font-medium">
                          <div className="flex items-center space-x-4">
                            <button
                              onClick={() => handleEditClick(u)}
                              className="inline-flex items-center rounded-lg bg-blue-100 px-2.5 py-1 text-xs font-semibold text-blue-800 hover:bg-blue-200"
                            >
                              <Edit className="w-4 h-4 mr-1" />
                              Edit
                            </button>
                            <button
                              onClick={() => handleDeleteUser(u)}
                              className="inline-flex items-center rounded-lg bg-rose-100 px-2.5 py-1 text-xs font-semibold text-rose-800 hover:bg-rose-200"
                            >
                              <Trash2 className="w-4 h-4 mr-1" />
                              Delete
                            </button>
                          </div>
                        </td>
                      </tr>
                    );
                  })
                )}
              </tbody>
            </table>
          </div>

          <div className="space-y-3 p-4 md:hidden">
            {users.length === 0 ? (
              <div className="rounded-xl border border-[var(--line-soft)] bg-white/80 p-4 text-center text-sm text-[var(--ink-soft)]">
                No users found
              </div>
            ) : (
              users.map((u) => {
                const manager = users.find((m) => m.id === u.manager_id);
                return (
                  <div key={u.id} className="rounded-xl border border-[var(--line-soft)] bg-white/80 p-4">
                    <div className="flex items-start justify-between gap-2">
                      <div>
                        <p className="text-sm font-semibold text-[var(--ink-strong)]">{u.first_name} {u.last_name}</p>
                        <p className="text-xs text-[var(--ink-soft)]">{u.email}</p>
                      </div>
                      <span className={`fancy-badge ${u.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                        {u.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </div>

                    <div className="mt-2 flex flex-wrap gap-2 text-xs text-[var(--ink-soft)]">
                      <span className={`fancy-badge ${u.role === 'admin' ? 'bg-violet-100 text-violet-800' : u.role === 'manager' ? 'bg-blue-100 text-blue-800' : 'bg-green-100 text-green-800'}`}>
                        {u.role}
                      </span>
                      <span className="rounded-full border border-[var(--line)] px-2 py-0.5">{u.department || 'No department'}</span>
                    </div>

                    <p className="mt-2 text-xs text-[var(--ink-soft)]">Manager: {manager ? `${manager.first_name} ${manager.last_name}` : '-'}</p>

                    <div className="mt-3 flex items-center gap-2">
                      <button
                        onClick={() => handleEditClick(u)}
                        className="inline-flex items-center rounded-lg bg-blue-100 px-3 py-1.5 text-xs font-semibold text-blue-800"
                      >
                        <Edit className="mr-1 h-3.5 w-3.5" />
                        Edit
                      </button>
                      <button
                        onClick={() => handleDeleteUser(u)}
                        className="inline-flex items-center rounded-lg bg-rose-100 px-3 py-1.5 text-xs font-semibold text-rose-800"
                      >
                        <Trash2 className="mr-1 h-3.5 w-3.5" />
                        Delete
                      </button>
                    </div>
                  </div>
                );
              })
            )}
          </div>
        </div>
      </div>

      {/* Create User Modal */}
      {createModal && (
        <div className="fixed inset-0 z-50 overflow-y-auto bg-[#10292e]/45 p-4 backdrop-blur-sm">
          <div className="surface-panel mx-auto my-6 max-h-[90vh] w-full max-w-md overflow-y-auto rounded-2xl p-6 sm:p-8">
            <h3 className="section-title mb-4 text-xl font-bold text-[var(--ink-strong)]">Create User</h3>
            <form onSubmit={handleCreateUser}>
              <div className="space-y-4">
                <div>
                  <label className="mb-1 block text-sm font-semibold text-[var(--ink-soft)]">Email</label>
                  <input
                    type="email"
                    value={createFormData.email}
                    onChange={(e) => setCreateFormData({ ...createFormData, email: e.target.value })}
                    className="input-base"
                    required
                  />
                </div>

                <div>
                  <label className="mb-1 block text-sm font-semibold text-[var(--ink-soft)]">Password</label>
                  <input
                    type="password"
                    value={createFormData.password}
                    onChange={(e) => setCreateFormData({ ...createFormData, password: e.target.value })}
                    className="input-base"
                    minLength={8}
                    required
                  />
                </div>

                <div>
                  <label className="mb-1 block text-sm font-semibold text-[var(--ink-soft)]">First Name</label>
                  <input
                    type="text"
                    value={createFormData.first_name}
                    onChange={(e) => setCreateFormData({ ...createFormData, first_name: e.target.value })}
                    className="input-base"
                    required
                  />
                </div>

                <div>
                  <label className="mb-1 block text-sm font-semibold text-[var(--ink-soft)]">Last Name</label>
                  <input
                    type="text"
                    value={createFormData.last_name}
                    onChange={(e) => setCreateFormData({ ...createFormData, last_name: e.target.value })}
                    className="input-base"
                    required
                  />
                </div>

                <div>
                  <label className="mb-1 block text-sm font-semibold text-[var(--ink-soft)]">Role</label>
                  <select
                    value={createFormData.role}
                    onChange={(e) => setCreateFormData({ ...createFormData, role: e.target.value })}
                    className="input-base"
                  >
                    <option value="employee">Employee</option>
                    <option value="manager">Manager</option>
                    <option value="admin">Admin</option>
                  </select>
                </div>

                <div>
                  <label className="mb-1 block text-sm font-semibold text-[var(--ink-soft)]">Department</label>
                  <select
                    value={createFormData.department}
                    onChange={(e) => setCreateFormData({ ...createFormData, department: e.target.value })}
                    className="input-base"
                  >
                    <option value="">Select department</option>
                    {departments.map((dept) => (
                      <option key={dept} value={dept}>
                        {dept}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="mb-1 block text-sm font-semibold text-[var(--ink-soft)]">Assign Manager</label>
                  <select
                    value={createFormData.manager_id || ''}
                    onChange={(e) => setCreateFormData({ ...createFormData, manager_id: e.target.value ? parseInt(e.target.value, 10) : null })}
                    className="input-base"
                  >
                    <option value="">No Manager</option>
                    {users
                      .filter((u) => u.role === 'manager' || u.role === 'admin')
                      .map((manager) => (
                        <option key={manager.id} value={manager.id}>
                          {manager.first_name} {manager.last_name} ({manager.role})
                        </option>
                      ))}
                  </select>
                </div>
              </div>

              <div className="flex space-x-3 mt-6">
                <button type="submit" className="brand-btn flex-1 rounded-xl px-4 py-2.5 text-sm font-semibold">
                  Create User
                </button>
                <button
                  type="button"
                  onClick={() => setCreateModal(false)}
                  className="ghost-btn flex-1 rounded-xl px-4 py-2.5 text-sm font-semibold"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Edit User Modal */}
      {editModal.show && (
        <div className="fixed inset-0 z-50 overflow-y-auto bg-[#10292e]/45 p-4 backdrop-blur-sm">
          <div className="surface-panel mx-auto my-6 max-h-[90vh] w-full max-w-md overflow-y-auto rounded-2xl p-6 sm:p-8">
            <h3 className="section-title mb-4 text-xl font-bold text-[var(--ink-strong)]">Edit User</h3>
            <form onSubmit={handleUpdateUser}>
              <div className="space-y-4">
                {/* First Name */}
                <div>
                  <label className="mb-1 block text-sm font-semibold text-[var(--ink-soft)]">
                    First Name
                  </label>
                  <input
                    type="text"
                    value={formData.first_name}
                    onChange={(e) =>
                      setFormData({ ...formData, first_name: e.target.value })
                    }
                    className="input-base"
                  />
                </div>

                {/* Last Name */}
                <div>
                  <label className="mb-1 block text-sm font-semibold text-[var(--ink-soft)]">
                    Last Name
                  </label>
                  <input
                    type="text"
                    value={formData.last_name}
                    onChange={(e) =>
                      setFormData({ ...formData, last_name: e.target.value })
                    }
                    className="input-base"
                  />
                </div>

                {/* Department */}
                <div>
                  <label className="mb-1 block text-sm font-semibold text-[var(--ink-soft)]">
                    Department
                  </label>
                  <select
                    value={formData.department}
                    onChange={(e) =>
                      setFormData({ ...formData, department: e.target.value })
                    }
                    className="input-base"
                  >
                    <option value="">Select department</option>
                    {departments.map((dept) => (
                      <option key={dept} value={dept}>
                        {dept}
                      </option>
                    ))}
                  </select>
                </div>

                {/* Manager */}
                <div>
                  <label className="mb-1 block text-sm font-semibold text-[var(--ink-soft)]">
                    Assign Manager
                  </label>
                  <select
                    value={formData.manager_id || ''}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        manager_id: e.target.value ? parseInt(e.target.value) : null,
                      })
                    }
                    className="input-base"
                  >
                    <option value="">No Manager</option>
                    {getAvailableManagers().map((manager) => (
                      <option key={manager.id} value={manager.id}>
                        {manager.first_name} {manager.last_name} ({manager.role})
                      </option>
                    ))}
                  </select>
                </div>

                {/* Current User Info */}
                <div className="rounded-xl border border-[var(--line-soft)] bg-white/75 p-3 text-sm text-[var(--ink-soft)]">
                  <p>
                    <strong>Email:</strong> {editModal.user?.email}
                  </p>
                  <p>
                    <strong>Role:</strong> {editModal.user?.role}
                  </p>
                </div>
              </div>

              <div className="flex space-x-3 mt-6">
                <button
                  type="submit"
                  className="brand-btn flex-1 rounded-xl px-4 py-2.5 text-sm font-semibold"
                >
                  Update User
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setEditModal({ show: false, user: null });
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

export default UserManagementPage;
