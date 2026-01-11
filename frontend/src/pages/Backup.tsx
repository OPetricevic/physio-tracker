import { useState } from 'react'
import { useAuth } from '../auth'
import { apiRequest } from '../api/client'

export function BackupPage() {
  const { user } = useAuth()
  const token = user?.token ?? null
  const [saving, setSaving] = useState(false)
  const [restoring, setRestoring] = useState(false)
  const [message, setMessage] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  const handleBackup = async () => {
    if (!token) return
    setSaving(true)
    setError(null)
    setMessage(null)
    try {
      const res = await fetch('/api/backup', {
        method: 'GET',
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) {
        const text = await res.text()
        throw new Error(text || 'Neuspješno preuzimanje.')
      }
      const blob = await res.blob()
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `physio-backup-${new Date().toLocaleDateString('hr-HR')}.dump`
      a.click()
      URL.revokeObjectURL(url)
      setMessage('Sigurnosna kopija je preuzeta.')
    } catch (err) {
      setError('Sigurnosna kopija nije uspjela.')
      console.error(err)
    } finally {
      setSaving(false)
    }
  }

  const handleRestore = async (file: File | null) => {
    if (!token || !file) return
    setRestoring(true)
    setError(null)
    setMessage(null)
    const form = new FormData()
    form.append('file', file)
    try {
      const res = await fetch('/api/backup/restore', {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: form,
      })
      if (!res.ok) {
        const text = await res.text()
        throw new Error(text || 'Neuspješno vraćanje sigurnosne kopije.')
      }
      await fetch('/api/health', { method: 'GET' })
      setMessage('Sigurnosna kopija je vraćena.')
    } catch (err) {
      setError('Vraćanje sigurnosne kopije nije uspjelo.')
      console.error(err)
    } finally {
      setRestoring(false)
    }
  }

  return (
    <main className="page profile">
      <div className="page-header">
        <div>
          <p className="eyebrow">Postavke</p>
          <h1>Sigurnosna kopija</h1>
          <p className="muted">Preuzmite ili vratite podatke aplikacije.</p>
        </div>
      </div>

      {message && <div className="alert success">{message}</div>}
      {error && <div className="alert error">{error}</div>}

      <div className="card section">
        <div className="section-header">
          <h2>Preuzmi kopiju</h2>
          <p className="muted">Preporučujemo redovno spremanje na USB ili cloud.</p>
        </div>
        <button className="btn primary" onClick={handleBackup} disabled={saving}>
          {saving ? 'Spremanje...' : 'Preuzmi sigurnosnu kopiju'}
        </button>
      </div>

      <div className="card section">
        <div className="section-header">
          <h2>Vrati kopiju</h2>
          <p className="muted">Odaberite .dump datoteku koja je prethodno spremljena.</p>
        </div>
        <input
          type="file"
          accept=".dump"
          onChange={(e) => handleRestore(e.target.files?.[0] ?? null)}
          disabled={restoring}
        />
      </div>
    </main>
  )
}
