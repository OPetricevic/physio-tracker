import { useEffect, useMemo, useState } from 'react'
import { apiRequest } from '../api/client'
import { useAuth } from '../auth'

type DoctorProfile = {
  uuid?: string
  doctor_uuid?: string
  practice_name: string
  department?: string
  role_title?: string
  address: string
  phone: string
  email?: string
  website?: string
  logo_path?: string
  protocol_prefix?: string
  header_note?: string
  footer_note?: string
}

type ProfileResponse = {
  profile: DoctorProfile
}

export function ProfilePage() {
  const { user } = useAuth()
  const token = user?.token ?? null
  const [profile, setProfile] = useState<DoctorProfile>({
    practice_name: '',
    address: '',
    phone: '',
  })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [saved, setSaved] = useState(false)
  const [touched, setTouched] = useState<Record<string, boolean>>({})
  const dirty = useMemo(() => Object.values(touched).some(Boolean), [touched])

  const isValid = useMemo(() => {
    return profile.practice_name.trim() !== '' && profile.address.trim() !== '' && profile.phone.trim() !== ''
  }, [profile])

  useEffect(() => {
    let mounted = true
    setLoading(true)
    apiRequest<ProfileResponse>('/api/doctor/profile', { method: 'GET', token })
      .then((res) => {
        if (!mounted) return
        setProfile({
          practice_name: res.profile.practice_name ?? '',
          department: res.profile.department ?? '',
          role_title: res.profile.role_title ?? '',
          address: res.profile.address ?? '',
          phone: res.profile.phone ?? '',
          email: res.profile.email ?? '',
          website: res.profile.website ?? '',
          logo_path: res.profile.logo_path ?? '',
          protocol_prefix: res.profile.protocol_prefix ?? '',
          header_note: res.profile.header_note ?? '',
          footer_note: res.profile.footer_note ?? '',
          uuid: res.profile.uuid,
          doctor_uuid: res.profile.doctor_uuid,
        })
      })
      .catch(() => {
        // If 404, we just keep empty profile; apiRequest throws for non-2xx, so ignore quietly.
      })
      .finally(() => {
        if (mounted) setLoading(false)
      })
    return () => {
      mounted = false
    }
  }, [token])

  const handleChange = (field: keyof DoctorProfile) => (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setProfile((prev) => ({ ...prev, [field]: e.target.value }))
    setTouched((prev) => ({ ...prev, [field]: true }))
    setSaved(false)
    setError(null)
  }

  const handleSave = async () => {
    if (!isValid || !token) return
    setSaving(true)
    setError(null)
    setSaved(false)
    try {
      const res = await apiRequest<ProfileResponse>('/api/doctor/profile', {
        method: 'PUT',
        token,
        body: { profile },
      })
      setProfile(res.profile)
      setSaved(true)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Greška pri spremanju profila.'
      setError(message)
    } finally {
      setSaving(false)
    }
  }

  const uploadFile = async (field: 'logo_path', file: File) => {
    if (!token) return
    setSaving(true)
    setError(null)
    setSaved(false)
    setTouched((prev) => ({ ...prev, [field]: true }))
    const form = new FormData()
    form.append('file', file)
    try {
      const res = await fetch('/api/files/upload', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: form,
      })
      if (!res.ok) {
        const text = await res.text()
        throw new Error(text || 'Neuspješan upload.')
      }
      const data = (await res.json()) as { url: string }
      setProfile((prev) => ({ ...prev, [field]: data.url }))
      setSaved(false)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Neuspješan upload.'
      setError(message)
    } finally {
      setSaving(false)
    }
  }

  const handleFileInput = (field: 'logo_path') => (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) uploadFile(field, file)
  }

  return (
    <main className="page profile">
      <div className="page-header">
        <div>
          <p className="eyebrow">Postavke</p>
          <h1>Profil ordinacije</h1>
          <p className="muted">Podaci za zaglavlje PDF nalaza: naziv, kontakt, adresa i logo.</p>
        </div>
        <button className="btn primary" onClick={handleSave} disabled={!isValid || saving || loading || !dirty}>
          {saving ? 'Spremanje...' : 'Spremi'}
        </button>
      </div>

      {saved && <div className="alert success">Profil spremljen.</div>}
      {error && <div className="alert error">{error}</div>}

      {loading ? <div className="card">Učitavanje...</div> : (
        <div className="profile-layout">
          <div className="card section">
            <div className="section-header">
              <h2>Osnovni podaci</h2>
              <p className="muted">Naziv ordinacije i stručni podaci.</p>
            </div>
            <div className="grid two-cols">
              <Field
                label="Naziv ordinacije"
                required
                value={profile.practice_name}
                onChange={handleChange('practice_name')}
                ghost={!touched.practice_name && !!profile.practice_name}
                onFocus={() => setTouched((p) => ({ ...p, practice_name: true }))}
              />
              <Field
                label="Odjel / pododjel"
                value={profile.department ?? ''}
                onChange={handleChange('department')}
                ghost={!touched.department && !!profile.department}
                onFocus={() => setTouched((p) => ({ ...p, department: true }))}
              />
              <Field
                label="Titula / uloga"
                value={profile.role_title ?? ''}
                onChange={handleChange('role_title')}
                ghost={!touched.role_title && !!profile.role_title}
                onFocus={() => setTouched((p) => ({ ...p, role_title: true }))}
              />
              <Field
                label="Protokol prefiks"
                value={profile.protocol_prefix ?? ''}
                onChange={handleChange('protocol_prefix')}
                ghost={!touched.protocol_prefix && !!profile.protocol_prefix}
                onFocus={() => setTouched((p) => ({ ...p, protocol_prefix: true }))}
              />
            </div>
          </div>

          <div className="card section">
            <div className="section-header">
              <h2>Kontakt i adresa</h2>
              <p className="muted">Telefon je obavezan, ostalo po želji.</p>
            </div>
            <div className="grid two-cols">
              <Field
                label="Telefon"
                required
                value={profile.phone}
                onChange={handleChange('phone')}
                ghost={!touched.phone && !!profile.phone}
                onFocus={() => setTouched((p) => ({ ...p, phone: true }))}
              />
              <Field
                label="Email"
                value={profile.email ?? ''}
                onChange={handleChange('email')}
                ghost={!touched.email && !!profile.email}
                onFocus={() => setTouched((p) => ({ ...p, email: true }))}
              />
              <Field
                label="Web"
                value={profile.website ?? ''}
                onChange={handleChange('website')}
                ghost={!touched.website && !!profile.website}
                onFocus={() => setTouched((p) => ({ ...p, website: true }))}
              />
              <Field
                label="Adresa"
                required
                full
                value={profile.address}
                onChange={handleChange('address')}
                ghost={!touched.address && !!profile.address}
                onFocus={() => setTouched((p) => ({ ...p, address: true }))}
              />
            </div>
          </div>

          <div className="card section">
            <div className="section-header">
              <h2>Logo</h2>
              <p className="muted">Najbolje .png s transparentnom pozadinom.</p>
            </div>
            <div className="logo-row">
              {profile.logo_path ? (
                <div className="logo-preview">
                  <img src={profile.logo_path} alt="Logo" />
                  <span className="muted">{profile.logo_path}</span>
                </div>
              ) : (
                <p className="muted">Logo nije postavljen.</p>
              )}
              <div className="logo-actions">
                <input type="file" accept=".png,.jpg,.jpeg" onChange={handleFileInput('logo_path')} />
              </div>
            </div>
          </div>

          <div className="card section">
            <div className="section-header">
              <h2>Napomene</h2>
              <p className="muted">Dodatni tekst u zaglavlju ili podnožju nalaza.</p>
            </div>
            <div className="grid one-col">
              <Field
                label="Napomena u zaglavlju"
                textarea
                value={profile.header_note ?? ''}
                onChange={handleChange('header_note')}
                ghost={!touched.header_note && !!profile.header_note}
                onFocus={() => setTouched((p) => ({ ...p, header_note: true }))}
              />
              <Field
                label="Napomena u podnožju"
                textarea
                value={profile.footer_note ?? ''}
                onChange={handleChange('footer_note')}
                ghost={!touched.footer_note && !!profile.footer_note}
                onFocus={() => setTouched((p) => ({ ...p, footer_note: true }))}
              />
            </div>
          </div>

          {error && <p className="error">{error}</p>}
          {saved && <p className="success">Profil spremljen.</p>}
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
  full,
  textarea,
  children,
  ghost,
  onFocus,
}: {
  label: string
  value: string
  required?: boolean
  full?: boolean
  textarea?: boolean
  children?: React.ReactNode
  ghost?: boolean
  onFocus?: () => void
  onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void
}) {
  const className = full ? 'field full' : 'field'
  const inputClass = ghost ? 'ghost' : undefined
  return (
    <label className={className}>
      <span className="field-label">
        {label}
        {required ? ' *' : ''}
      </span>
      {textarea ? (
        <textarea value={value} onChange={onChange} onFocus={onFocus} rows={3} className={inputClass} />
      ) : (
        <input type="text" value={value} onChange={onChange} onFocus={onFocus} className={inputClass} />
      )}
      {children ? <div className="field-hint">{children}</div> : null}
    </label>
  )
}
