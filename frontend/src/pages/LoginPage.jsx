import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { toast } from 'react-toastify';
import { Building2, ShieldCheck, Sparkles } from 'lucide-react';

const LoginPage = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const user = await login(email, password);
      toast.success('Login successful!');

      // Redirect based on role
      if (user.role === 'admin' || user.role === 'manager') {
        navigate('/admin');
      } else {
        navigate('/dashboard');
      }
    } catch (error) {
      toast.error(error.response?.data?.message || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen px-4 py-6 sm:px-6 lg:px-8">
      <div className="mx-auto grid min-h-[calc(100vh-3rem)] w-full max-w-6xl grid-cols-1 overflow-hidden rounded-3xl border border-white/40 shadow-[0_28px_70px_-42px_rgba(16,41,46,0.55)] lg:grid-cols-[1.08fr_1fr]">
        <section className="glass-panel relative hidden overflow-hidden p-10 lg:block">
          <div className="absolute -right-10 top-6 h-36 w-36 rounded-full bg-white/35 blur-2xl" />
          <div className="absolute bottom-8 left-4 h-24 w-24 rounded-full bg-amber-300/40 blur-2xl" />
          <div className="relative z-10 flex h-full flex-col justify-between">
            <div className="fade-in-up">
              <div className="inline-flex items-center gap-2 rounded-full border border-[var(--line)] bg-white/70 px-3 py-1 text-xs font-semibold uppercase tracking-[0.12em] text-[var(--ink-soft)]">
                <Sparkles className="h-3.5 w-3.5 text-[var(--accent)]" />
                Smart Leave Ops
              </div>
              <h1 className="section-title mt-5 max-w-md text-4xl font-extrabold leading-tight text-[var(--ink-strong)] xl:text-5xl">
                Better leave workflows for modern teams.
              </h1>
              <p className="mt-4 max-w-md text-sm leading-6 text-[var(--ink-soft)] xl:text-base">
                Approvals, balances, and visibility in one calm dashboard. Built for employees, managers, and admins.
              </p>
            </div>

            <div className="space-y-3 text-sm">
              <div className="surface-panel flex items-center gap-3 rounded-2xl p-4">
                <Building2 className="h-5 w-5 text-[var(--brand)]" />
                <span className="font-medium text-[var(--ink-soft)]">Department-aware leave policies</span>
              </div>
              <div className="surface-panel flex items-center gap-3 rounded-2xl p-4">
                <ShieldCheck className="h-5 w-5 text-[var(--brand)]" />
                <span className="font-medium text-[var(--ink-soft)]">Role-based access and secure review flow</span>
              </div>
            </div>
          </div>
        </section>

        <section className="surface-panel soft-scrollbar flex items-center p-6 sm:p-10 lg:p-12">
          <div className="mx-auto w-full max-w-md fade-in-up">
            <div className="mb-6 space-y-2 text-center lg:text-left">
              <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--ink-soft)]">Welcome Back</p>
              <h2 className="section-title text-3xl font-extrabold text-[var(--ink-strong)]">Sign in to continue</h2>
              <p className="text-sm text-[var(--ink-soft)]">Use your company account to access the leave portal.</p>
            </div>

            <form className="space-y-5" onSubmit={handleSubmit}>
              <div className="space-y-2">
                <label htmlFor="email" className="text-sm font-semibold text-[var(--ink-soft)]">
                  Work Email
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
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
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="input-base"
                  placeholder="Enter your password"
                />
              </div>

              <button
                type="submit"
                disabled={loading}
                className="brand-btn w-full rounded-xl px-4 py-3 text-sm font-semibold"
              >
                {loading ? 'Signing in...' : 'Sign in'}
              </button>

              <div className="text-center text-sm text-[var(--ink-soft)]">
                Don&apos;t have an account?{' '}
                <Link to="/register" className="font-bold text-[var(--brand)] transition hover:text-[var(--brand-strong)]">
                  Register here
                </Link>
              </div>

              <div className="rounded-2xl border border-[var(--line-soft)] bg-white/75 p-4">
                <p className="mb-2 text-xs font-bold uppercase tracking-[0.12em] text-[var(--ink-soft)]">Test Accounts</p>
                <div className="space-y-1.5 text-xs text-[var(--ink-soft)] sm:text-sm">
                  <p><span className="font-semibold">Admin:</span> admin@abccompany.com / admin123</p>
                  <p><span className="font-semibold">Manager:</span> manager@abccompany.com / manager123</p>
                  <p><span className="font-semibold">Employee:</span> employee@abccompany.com / employee123</p>
                </div>
              </div>
            </form>
          </div>
        </section>
      </div>
    </div>
  );
};

export default LoginPage;
