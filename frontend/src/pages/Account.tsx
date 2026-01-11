import { useEffect, useMemo, useState } from 'react'
import { apiRequest } from '../api/client'
import { useAuth } from '../auth'

type DoctorAccount = {
  uuid?: string
  email: string
  username: string
  first_name: string
  last_name: string
}

type DoctorResponse = DoctorAccount

export function AccountPage() {
  const { user } = useAuth()
  const token = user?.token ?? null
  const [account, setAccount] = useState<DoctorAccount>({
    email: '',
    username: '',
    first_name: '',
    last_name: '',
  })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [saved, setSaved] = useState(false)
  const [touched, setTouched] = useState<Record<string, boolean>>({})

  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [passwordError, setPasswordError] = useState<string | null>(null)
  const [passwordSaved, setPasswordSaved] = useState(false)

  const dirty = useMemo(() => Object.values(touched).some(Boolean), [touched])
  const isValid = useMemo(() => {
    return (
      account.email.trim() !== '' &&
      account.username.trim() !== '' &&
      account.first_name.trim() !== '' &&
      account.last_name.trim() !== ''
    )
  }, [account])

  useEffect(() => {
    let mounted = true
    if (!token) return
    setLoading(true)
    apiRequest<DoctorResponse>('/api/doctors/me', { method: 'GET', token })
      .then((res) => {
        if (!mounted) return
        setAccount({
          uuid: res.uuid,
          email: res.email ?? '',
          username: res.username ?? '',
          first_name: res.first_name ?? '',
          last_name: res.last_name ?? '',
        })
      })
      .catch((err) => {
        const message = err instanceof Error ? err.message : 'Greška pri dohvaćanju profila.'
        setError(message)
      })
      .finally(() => {
        if (mounted) setLoading(false)
      })
    return () => {
      mounted = false
    }
  }, [token])

  const handleChange = (field: keyof DoctorAccount) => (e: React.ChangeEvent<HTMLInputElement>) => {
    setAccount((prev) => ({ ...prev, [field]: e.target.value }))
    setTouched((prev) => ({ ...prev, [field]: true }))
    setSaved(false)
    setError(null)
  }

  const handleSave = async () => {
    if (!token) return
    if (!isValid) {
      setError('Molimo ispunite obavezna polja.')
      return
    }
    setSaving(true)
    setError(null)
    setSaved(false)
    try {
      const res = await apiRequest<DoctorResponse>('/api/doctors/me', {
        method: 'PATCH',
        token,
        body: {
          email: account.email,
          username: account.username,
          first_name: account.first_name,
          last_name: account.last_name,
        },
      })
      setAccount({
        uuid: res.uuid,
        email: res.email ?? '',
        username: res.username ?? '',
        first_name: res.first_name ?? '',
        last_name: res.last_name ?? '',
      })
      setSaved(true)
      setTouched({})
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Greška pri spremanju profila.'
      setError(message)
    } finally {
      setSaving(false)
    }
  }

  const handlePasswordSave = async () => {
    if (!token) return
    setPasswordError(null)
    setPasswordSaved(false)
    if (!currentPassword || !newPassword || !confirmPassword) {
      setPasswordError('Molimo ispunite sva polja.')
      return
    }
    if (newPassword !== confirmPassword) {
      setPasswordError('Lozinke se ne podudaraju.')
      return
    }
    try {
      await apiRequest<null>('/api/auth/change-password', {
        method: 'POST',
        token,
        body: {
          current_password: currentPassword,
          new_password: newPassword,
        },
      })
      setCurrentPassword('')
      setNewPassword('')
      setConfirmPassword('')
      setPasswordSaved(true)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Greška pri promjeni lozinke.'
      setPasswordError(message)
    }
  }

  return (
    <main className="page profile">
      <div className="page-header">
        <div>
          <p className="eyebrow">Postavke</p>
          <h1>Račun</h1>
          <p className="muted">Uredite osobne podatke ili promijenite lozinku.</p>
        </div>
        <button className="btn primary" onClick={handleSave} disabled={saving || loading || !dirty}>
          {saving ? 'Spremanje...' : 'Spremi'}
        </button>
      </div>

      {saved && <div className="alert success">Podaci spremljeni.</div>}
      {error && <div className="alert error">{error}</div>}

      {loading ? <div className="card">Učitavanje...</div> : (
        <div className="profile-layout">
          <div className="card section">
            <div className="section-header">
              <h2>Osobni podaci</h2>
              <p className="muted">Ime, prezime i kontakt za prijavu.</p>
            </div>
            <div className="grid two-cols">
              <Field
                label="Ime"
                required
                value={account.first_name}
                onChange={handleChange('first_name')}
              />
              <Field
                label="Prezime"
                required
                value={account.last_name}
                onChange={handleChange('last_name')}
              />
              <Field
                label="Email"
                required
                value={account.email}
                onChange={handleChange('email')}
              />
              <Field
                label="Korisničko ime"
                required
                value={account.username}
                onChange={handleChange('username')}
              />
            </div>
          </div>

          <div className="card section">
            <div className="section-header">
              <h2>Promjena lozinke</h2>
              <p className="muted">Upišite trenutnu i novu lozinku.</p>
            </div>
            <div className="grid two-cols">
              <Field
                label="Trenutna lozinka"
                type="password"
                value={currentPassword}
                onChange={(e) => setCurrentPassword(e.target.value)}
              />
              <Field
                label="Nova lozinka"
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
              />
              <Field
                label="Potvrdite lozinku"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
              />
            </div>
            <div className="section-actions">
              <button className="btn primary" onClick={handlePasswordSave}>
                Spremi lozinku
              </button>
            </div>
            {passwordSaved && <div className="alert success">Lozinka je promijenjena.</div>}
            {passwordError && <div className="alert error">{passwordError}</div>}
          </div>
        </div>
      )}
    </main>
  )
}

function Field({
  label,
  value,
  onChange,
  required,
  type = 'text',
}: {
  label: string
  value: string
  required?: boolean
  type?: string
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void
}) {
  return (
    <label className="field">
      <span className="field-label">
        {label}
        {required ? ' *' : ''}
      </span>
      <input type={type} value={value} onChange={onChange} />
    </label>
  )
}
