import { useState } from 'react'
import type { Anamnesis } from '../types'

type Props = {
  patientName: string
  anamneses: Anamnesis[]
  page: number
  totalPages: number
  onPageChange: (page: number) => void
  disabled: boolean
  onAdd: (note: string) => void
  onGeneratePdf: (anamnesisUuid: string) => void
  onBackup: () => void
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
}: Props) {
  const [note, setNote] = useState('')

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!note.trim() || disabled) return
    onAdd(note.trim())
    setNote('')
  }

  return (
    <div className="panel detail">
      <div className="panel-header">
        <div>
          <p className="eyebrow">Anamneze</p>
          <h2>{patientName || 'Odaberite pacijenta'}</h2>
        </div>
        <div className="actions">
          <button type="button" className="btn ghost" onClick={onBackup} disabled={disabled}>
            Sigurnosna kopija
          </button>
        </div>
      </div>

      <div className="stack">
        {anamneses.length === 0 && <div className="empty">Nema unosa.</div>}
        {anamneses.map((entry) => (
          <article key={entry.uuid} className="note">
            <div className="note__header">
              <div>
                <p className="note__eyebrow">Posjet</p>
                <strong>{new Date(entry.createdAt).toLocaleDateString()}</strong>
              </div>
              <button
                type="button"
                className="btn text"
                onClick={() => onGeneratePdf(entry.uuid)}
                disabled={disabled}
              >
                Generiraj PDF
              </button>
            </div>
            <p className="note__body">{entry.note}</p>
          </article>
        ))}
        {totalPages > 1 && (
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
          </div>
        )}
      </div>

      <form className="composer" onSubmit={handleSubmit}>
        <label htmlFor="note">Dodaj anamnezu</label>
        <textarea
          id="note"
          name="note"
          placeholder="Bilješke sa tretmana, vježbe, napredak..."
          value={note}
          onChange={(e) => setNote(e.target.value)}
          rows={4}
          disabled={disabled}
        />
        <div className="actions">
          <button type="submit" className="btn primary" disabled={disabled}>
            Spremi anamnezu
          </button>
        </div>
      </form>
    </div>
  )
}
