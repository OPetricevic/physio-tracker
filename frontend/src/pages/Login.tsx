import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth'
import '../App.css'

export function LoginPage() {
  const navigate = useNavigate()
  const { login } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!email.trim() || !password.trim()) return
    try {
      setLoading(true)
      setError(null)
      await login(email.trim(), password.trim())
      navigate('/app', { replace: true })
    } catch (err) {
      setError('Prijava nije uspjela. Provjerite podatke i pokušajte ponovno.')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-card">
        <p className="eyebrow">Physio Tracker</p>
        <h1>Prijava</h1>
        <p className="lede">Prijavite se kako biste pristupili pacijentima i anamnezama.</p>
        <form className="auth-form" onSubmit={handleSubmit}>
          <div className="field">
            <label htmlFor="email">Email</label>
            <input
              id="email"
              name="email"
              type="email"
              autoComplete="email"
              placeholder="doktor@ordinacija.hr"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <div className="field">
            <label htmlFor="password">Lozinka</label>
            <input
              id="password"
              name="password"
              type="password"
              autoComplete="current-password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
          />
        </div>
        {error && <p className="error-text">{error}</p>}
        <button type="submit" className="btn primary full" disabled={loading}>
          Prijavi se
        </button>
      </form>
        <p className="muted-small">
          Nemate račun? <Link to="/register">Registrirajte se</Link>
        </p>
      </div>
    </div>
  )
}
