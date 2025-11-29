import { useState } from 'react'
import type { Anamnesis } from '../types'

type Props = {
  patientName: string
  anamneses: Anamnesis[]
  disabled: boolean
  onAdd: (note: string) => void
  onGeneratePdf: (anamnesisUuid: string) => void
  onBackup: () => void
}

export function AnamnesisPanel({
  patientName,
  anamneses,
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
          <p className="eyebrow">Anamneses</p>
          <h2>{patientName || 'Select a patient'}</h2>
        </div>
        <div className="actions">
          <button type="button" className="btn ghost" onClick={onBackup} disabled={disabled}>
            Backup patient data
          </button>
        </div>
      </div>

      <div className="stack">
        {anamneses.length === 0 && <div className="empty">No anamneses yet.</div>}
        {anamneses.map((entry) => (
          <article key={entry.uuid} className="note">
            <div className="note__header">
              <div>
                <p className="note__eyebrow">Visit</p>
                <strong>{new Date(entry.createdAt).toLocaleDateString()}</strong>
              </div>
              <button
                type="button"
                className="btn text"
                onClick={() => onGeneratePdf(entry.uuid)}
                disabled={disabled}
              >
                Generate PDF
              </button>
            </div>
            <p className="note__body">{entry.note}</p>
          </article>
        ))}
      </div>

      <form className="composer" onSubmit={handleSubmit}>
        <label htmlFor="note">Add anamnesis</label>
        <textarea
          id="note"
          name="note"
          placeholder="Session notes, exercises, progress..."
          value={note}
          onChange={(e) => setNote(e.target.value)}
          rows={4}
          disabled={disabled}
        />
        <div className="actions">
          <button type="submit" className="btn primary" disabled={disabled}>
            Save anamnesis
          </button>
        </div>
      </form>
    </div>
  )
}
