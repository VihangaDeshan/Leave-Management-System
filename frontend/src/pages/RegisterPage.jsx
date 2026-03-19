import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { authApi } from '../api';
import { toast } from 'react-toastify';
import { BadgeCheck, Briefcase, UserPlus } from 'lucide-react';

const RegisterPage = () => {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    department: '',
  });
  const [departments, setDepartments] = useState([]);
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchDepartments = async () => {
      try {
        const response = await authApi.getDepartments();
        if (response.data) {
          setDepartments(Array.isArray(response.data) ? response.data : []);
        }
      } catch (err) {
        console.error('Failed to fetch departments:', err);
      }
    };
    fetchDepartments();
  }, []);

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      await register(formData);
      toast.success('Registration successful!');
      navigate('/dashboard');
    } catch (error) {
      toast.error(error.response?.data?.message || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen px-4 py-6 sm:px-6 lg:px-8">
      <div className="mx-auto grid min-h-[calc(100vh-3rem)] w-full max-w-6xl grid-cols-1 overflow-hidden rounded-3xl border border-white/40 shadow-[0_28px_70px_-42px_rgba(16,41,46,0.55)] lg:grid-cols-[1.08fr_1fr]">
        <section className="glass-panel relative hidden overflow-hidden p-10 lg:block">
          <div className="absolute right-0 top-0 h-52 w-52 rounded-full bg-white/35 blur-3xl" />
          <div className="absolute -left-4 bottom-10 h-40 w-40 rounded-full bg-teal-200/40 blur-3xl" />

          <div className="relative z-10 flex h-full flex-col justify-between">
            <div className="fade-in-up">
              <div className="inline-flex items-center gap-2 rounded-full border border-[var(--line)] bg-white/70 px-3 py-1 text-xs font-semibold uppercase tracking-[0.12em] text-[var(--ink-soft)]">
                <UserPlus className="h-3.5 w-3.5 text-[var(--brand)]" />
                Quick Onboarding
              </div>
              <h1 className="section-title mt-5 max-w-md text-4xl font-extrabold leading-tight text-[var(--ink-strong)] xl:text-5xl">
                Create your workspace account in minutes.
              </h1>
              <p className="mt-4 max-w-md text-sm leading-6 text-[var(--ink-soft)] xl:text-base">
                Set your role, link your department, and start managing leave requests with full transparency.
              </p>
            </div>

            <div className="space-y-3 text-sm">
              <div className="surface-panel flex items-center gap-3 rounded-2xl p-4">
                <BadgeCheck className="h-5 w-5 text-[var(--brand)]" />
                <span className="font-medium text-[var(--ink-soft)]">Secure account setup and role-based access</span>
              </div>
              <div className="surface-panel flex items-center gap-3 rounded-2xl p-4">
                <Briefcase className="h-5 w-5 text-[var(--brand)]" />
                <span className="font-medium text-[var(--ink-soft)]">Department-aligned user profiles</span>
              </div>
            </div>
          </div>
        </section>

        <section className="surface-panel soft-scrollbar max-h-[calc(100vh-3rem)] overflow-y-auto p-6 sm:p-10 lg:p-12">
          <div className="mx-auto w-full max-w-md fade-in-up">
            <div className="mb-6 space-y-2 text-center lg:text-left">
              <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--ink-soft)]">Get Started</p>
              <h2 className="section-title text-3xl font-extrabold text-[var(--ink-strong)]">Create your account</h2>
              <p className="text-sm text-[var(--ink-soft)]">Join ABC Company&apos;s leave management platform.</p>
            </div>

            <form className="space-y-5" onSubmit={handleSubmit}>
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label htmlFor="first_name" className="text-sm font-semibold text-[var(--ink-soft)]">
                    First Name
                  </label>
                  <input
                    id="first_name"
                    name="first_name"
                    type="text"
                    required
                    value={formData.first_name}
                    onChange={handleChange}
                    className="input-base"
                    placeholder="First name"
                  />
                </div>

                <div className="space-y-2">
                  <label htmlFor="last_name" className="text-sm font-semibold text-[var(--ink-soft)]">
                    Last Name
                  </label>
                  <input
                    id="last_name"
                    name="last_name"
                    type="text"
                    required
                    value={formData.last_name}
                    onChange={handleChange}
                    className="input-base"
                    placeholder="Last name"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <label htmlFor="email" className="text-sm font-semibold text-[var(--ink-soft)]">
                  Work Email
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  required
                  value={formData.email}
                  onChange={handleChange}
                  className="input-base"
                  placeholder="name@company.com"
                />
              </div>

              <div className="space-y-2">
                <label htmlFor="password" className="text-sm font-semibold text-[var(--ink-soft)]">
                  Password
                </label>
                <input
                  id="password"
                  name="password"
                  type="password"
                  required
                  minLength={8}
                  value={formData.password}
                  onChange={handleChange}
                  className="input-base"
                  placeholder="Minimum 8 characters"
                />
              </div>

              <div className="space-y-2">
                <label htmlFor="department" className="text-sm font-semibold text-[var(--ink-soft)]">
                  Department (Optional)
                </label>
                <select
                  id="department"
                  name="department"
                  value={formData.department}
                  onChange={handleChange}
                  className="input-base"
                >
                  <option value="">Select a department...</option>
                  {departments.map((dept, idx) => (
                    <option key={idx} value={dept}>
                      {dept}
                    </option>
                  ))}
                </select>
              </div>

              <button
                type="submit"
                disabled={loading}
                className="brand-btn w-full rounded-xl px-4 py-3 text-sm font-semibold"
              >
                {loading ? 'Creating account...' : 'Create account'}
              </button>

              <div className="text-center text-sm text-[var(--ink-soft)]">
                Already have an account?{' '}
                <Link to="/login" className="font-bold text-[var(--brand)] transition hover:text-[var(--brand-strong)]">
                  Sign in here
                </Link>
              </div>
            </form>
          </div>
        </section>
      </div>
    </div>
  );
};

export default RegisterPage;
