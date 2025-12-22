import { useState } from 'react'
import type { Anamnesis } from '../types'

type Props = {
  patientName: string
  anamneses: Anamnesis[]
  page: number
  totalPages: number
  onPageChange: (page: number) => void
  disabled: boolean
  onAdd: (input: {
    note: string
    diagnosis?: string
    therapy?: string
    otherInfo?: string
  }) => void
  onGeneratePdf: (anamnesisUuid: string) => void
  onBackup: () => void
  selectedVisits: Set<string>
  onToggleVisit: (uuid: string) => void
  onBulkPdf: () => void
}

export function AnamnesisPanel({
  patientName,
  anamneses,
  page,
  totalPages,
  onPageChange,
  disabled,
  onAdd,
  onGeneratePdf,
  onBackup,
  selectedVisits,
  onToggleVisit,
  onBulkPdf,
}: Props) {
  const [note, setNote] = useState('')
  const [diagnosis, setDiagnosis] = useState('')
  const [therapy, setTherapy] = useState('')
  const [otherInfo, setOtherInfo] = useState('')
  const [showForm, setShowForm] = useState(false)

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!note.trim() || disabled) return
    onAdd({
      note: note.trim(),
      diagnosis: diagnosis.trim() || undefined,
      therapy: therapy.trim() || undefined,
      otherInfo: otherInfo.trim() || undefined,
    })
    setNote('')
    setDiagnosis('')
    setTherapy('')
    setOtherInfo('')
  }

  return (
    <div className="panel detail">
      <div className="panel-header">
        <div>
          <p className="eyebrow">Anamneze</p>
          <h2>{patientName || 'Odaberite pacijenta'}</h2>
        </div>
        <div className="actions">
          <button type="button" className="btn primary small" onClick={() => setShowForm(true)} disabled={disabled}>
            Novi zapis
          </button>
          <button type="button" className="btn ghost" onClick={onBackup} disabled={disabled}>
            Sigurnosna kopija
          </button>
        </div>
      </div>

      <div className="stack">
        {anamneses.length === 0 && <div className="empty">Nema unosa.</div>}
        {anamneses.map((entry) => {
          const isSelected = selectedVisits.has(entry.uuid)
          return (
            <article
              key={entry.uuid}
              className={`note ${isSelected ? 'is-selected' : ''}`}
              role="button"
              tabIndex={0}
              onClick={() => !disabled && onToggleVisit(entry.uuid)}
              onKeyDown={(e) => {
                if (!disabled && (e.key === 'Enter' || e.key === ' ')) {
                  e.preventDefault()
                  onToggleVisit(entry.uuid)
                }
              }}
            >
              <div className="note__header">
                <div>
                  <p className="note__eyebrow">Posjet</p>
                  <strong>{new Date(entry.createdAt).toLocaleDateString()}</strong>
                </div>
              {isSelected && <span className="pill">Za PDF</span>}
              <button
                type="button"
                className="btn text"
                onClick={() => onGeneratePdf(entry.uuid)}
                disabled={disabled}
              >
                Generiraj PDF
              </button>
            </div>
            <p className="note__body">
              {entry.note}
            </p>
          </article>
          )
        })}
        {anamneses.length > 0 && (
          <div className="pager">
            <button
              type="button"
              className="btn ghost small"
              onClick={() => onPageChange(page - 1)}
              disabled={page <= 1}
            >
              ◀ Prethodne
            </button>
            <span className="muted-small">
              Stranica {page} / {totalPages}
            </span>
            <button
              type="button"
              className="btn ghost small"
              onClick={() => onPageChange(page + 1)}
              disabled={page >= totalPages}
            >
              Sljedeće ▶
            </button>
            <button type="button" className="btn primary small" onClick={onBulkPdf} disabled={disabled || selectedVisits.size === 0}>
              Odaberi za PDF
            </button>
          </div>
        )}
      </div>

      {showForm && (
        <div className="modal-backdrop">
          <div className="modal">
            <p className="eyebrow">Novi zapis</p>
            <h3>Dodaj posjet</h3>
            <form className="composer" onSubmit={handleSubmit}>
              <label htmlFor="note">Anamneza</label>
              <textarea
                id="note"
                name="note"
                placeholder="Bilješke sa tretmana, vježbe, napredak..."
                value={note}
                onChange={(e) => setNote(e.target.value)}
                rows={3}
                disabled={disabled}
              />

              <label htmlFor="diagnosis">Dijagnoza</label>
              <textarea
                id="diagnosis"
                name="diagnosis"
                placeholder="Dijagnoza..."
                value={diagnosis}
                onChange={(e) => setDiagnosis(e.target.value)}
                rows={2}
                disabled={disabled}
              />

              <label htmlFor="therapy">Terapija</label>
              <textarea
                id="therapy"
                name="therapy"
                placeholder="Terapija..."
                value={therapy}
                onChange={(e) => setTherapy(e.target.value)}
                rows={2}
                disabled={disabled}
              />

              <label htmlFor="otherInfo">Ostale informacije</label>
              <textarea
                id="otherInfo"
                name="otherInfo"
                placeholder="Drugi dolazak, stanje..."
                value={otherInfo}
                onChange={(e) => setOtherInfo(e.target.value)}
                rows={2}
                disabled={disabled}
              />
              <div className="actions">
                <button type="button" className="btn ghost" onClick={() => setShowForm(false)}>
                  Odustani
                </button>
                <button type="submit" className="btn primary" disabled={disabled}>
                  Spremi anamnezu
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
