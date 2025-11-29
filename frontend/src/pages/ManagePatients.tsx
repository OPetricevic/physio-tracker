import { useMemo, useState } from 'react'
import { usePatients } from '../data/PatientsContext'
import type { Patient } from '../types'
import '../App.css'

type EditingState = {
  [uuid: string]: {
    firstName: string
    lastName: string
    phone?: string
  }
}

export function ManagePatientsPage() {
  const { patients, anamneses, updatePatient, deletePatient } = usePatients()
  const [editing, setEditing] = useState<EditingState>({})
  const [searchTerm, setSearchTerm] = useState('')

  const filteredPatients = useMemo(() => {
    const term = searchTerm.trim().toLowerCase()
    if (!term) return patients
    return patients.filter((p) => {
      const name = `${p.firstName} ${p.lastName}`.toLowerCase()
      return name.includes(term) || (p.phone ?? '').toLowerCase().includes(term)
    })
  }, [patients, searchTerm])

  const startEdit = (p: Patient) => {
    setEditing((prev) => ({
      ...prev,
      [p.uuid]: { firstName: p.firstName, lastName: p.lastName, phone: p.phone },
    }))
  }

  const cancelEdit = (uuid: string) => {
    setEditing((prev) => {
      const next = { ...prev }
      delete next[uuid]
      return next
    })
  }

  const saveEdit = (uuid: string) => {
    const draft = editing[uuid]
    if (!draft) return
    if (!draft.firstName.trim() || !draft.lastName.trim()) return
    updatePatient(uuid, {
      firstName: draft.firstName.trim(),
      lastName: draft.lastName.trim(),
      phone: draft.phone?.trim() || undefined,
    })
    cancelEdit(uuid)
  }

  const handleDelete = (uuid: string) => {
    const hasAnamneses = (anamneses[uuid] ?? []).length > 0
    const confirmed = window.confirm(
      hasAnamneses
        ? 'Jeste li sigurni? Pacijent i sve anamneze bit će trajno obrisani.'
        : 'Jeste li sigurni da želite obrisati pacijenta?',
    )
    if (!confirmed) return
    deletePatient(uuid)
    cancelEdit(uuid)
  }

  return (
    <div className="page">
      <header className="app-header">
        <div>
          <p className="eyebrow">Pacijenti</p>
          <h1>Upravljanje</h1>
          <p className="lede">Uredite ili obrišite pacijente. Brisanje je trajno.</p>
        </div>
      </header>

      <main className="panel">
        <div className="field" style={{ maxWidth: 360 }}>
          <label htmlFor="search-manage">Pretraga</label>
          <input
            id="search-manage"
            type="search"
            placeholder="Traži po imenu ili telefonu"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>

        <div className="table">
          <div className="table-head">
            <span>Pacijent</span>
            <span>Telefon</span>
            <span>Anamneze</span>
            <span>Akcije</span>
          </div>
          {filteredPatients.map((p) => {
            const draft = editing[p.uuid]
            const count = (anamneses[p.uuid] ?? []).length
            return (
              <div key={p.uuid} className="table-row">
                <div className="table-cell">
                  {draft ? (
                    <div className="field">
                      <input
                        value={draft.firstName}
                        onChange={(e) =>
                          setEditing((prev) => ({
                            ...prev,
                            [p.uuid]: { ...prev[p.uuid], firstName: e.target.value },
                          }))
                        }
                      />
                      <input
                        value={draft.lastName}
                        onChange={(e) =>
                          setEditing((prev) => ({
                            ...prev,
                            [p.uuid]: { ...prev[p.uuid], lastName: e.target.value },
                          }))
                        }
                      />
                    </div>
                  ) : (
                    <div>
                      <strong>
                        {p.firstName} {p.lastName}
                      </strong>
                      <div className="muted-small">{new Date(p.createdAt).toLocaleDateString()}</div>
                    </div>
                  )}
                </div>
                <div className="table-cell">
                  {draft ? (
                    <input
                      value={draft.phone ?? ''}
                      onChange={(e) =>
                        setEditing((prev) => ({
                          ...prev,
                          [p.uuid]: { ...prev[p.uuid], phone: e.target.value },
                        }))
                      }
                    />
                  ) : (
                    <span>{p.phone || 'Telefon nije unesen'}</span>
                  )}
                </div>
                <div className="table-cell">
                  <span className="pill">{count}</span>
                </div>
                <div className="table-cell actions">
                  {draft ? (
                    <>
                      <button className="btn primary small" onClick={() => saveEdit(p.uuid)}>
                        Spremi
                      </button>
                      <button className="btn ghost small" onClick={() => cancelEdit(p.uuid)}>
                        Odustani
                      </button>
                    </>
                  ) : (
                    <>
                      <button className="btn ghost small" onClick={() => startEdit(p)}>
                        Uredi
                      </button>
                      <button className="btn ghost small" onClick={() => handleDelete(p.uuid)}>
                        Obriši
                      </button>
                    </>
                  )}
                </div>
              </div>
            )
          })}
          {filteredPatients.length === 0 && (
            <div className="empty" style={{ marginTop: 8 }}>
              Nema pacijenata za prikaz.
            </div>
          )}
        </div>
      </main>
    </div>
  )
}
