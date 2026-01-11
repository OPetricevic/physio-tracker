import { useState } from 'react'
import type { Anamnesis } from '../types'

type Props = {
  patientName: string
  anamneses: Anamnesis[]
  loadSelectionOptions: () => Promise<Anamnesis[]>
  searchTerm: string
  onSearchChange: (term: string) => void
  page: number
  hasNext: boolean
  onPageChange: (page: number) => void
  disabled: boolean
  onAdd: (input: { note: string; status: string; diagnosis: string; therapy: string; otherInfo: string }) => void
  onDelete: (uuid: string) => void
  onUpdateIncludes: (uuid: string, includeVisitUuids: string[]) => Promise<void>
  onEdit: (uuid: string, payload: { note: string; status: string; diagnosis: string; therapy: string; otherInfo: string }) => Promise<void>
  onGeneratePdf: (anamnesisUuid: string, includes?: string[], onlyCurrent?: boolean) => void
  onBackup: () => void
}

export function AnamnesisPanel({
  patientName,
  anamneses,
  searchTerm,
  onSearchChange,
  page,
  hasNext,
  onPageChange,
  disabled,
  onAdd,
  onDelete,
  onUpdateIncludes,
  onEdit,
  onGeneratePdf,
  onBackup,
  loadSelectionOptions,
}: Props) {
  const [note, setNote] = useState('')
  const [status, setStatus] = useState('')
  const [diagnosis, setDiagnosis] = useState('')
  const [therapy, setTherapy] = useState('')
  const [otherInfo, setOtherInfo] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [selectionModal, setSelectionModal] = useState<{ open: boolean; currentId: string | null; selected: Set<string> }>({
    open: false,
    currentId: null,
    selected: new Set(),
  })
  const [selectionSearch, setSelectionSearch] = useState('')
  const [selectionList, setSelectionList] = useState<Anamnesis[]>([])
  const [onlyCurrent, setOnlyCurrent] = useState(false)
  const [editModal, setEditModal] = useState<{
    open: boolean
    currentId: string | null
    note: string
    status: string
    diagnosis: string
    therapy: string
    otherInfo: string
  }>({
    open: false,
    currentId: null,
    note: '',
    status: '',
    diagnosis: '',
    therapy: '',
    otherInfo: '',
  })

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (disabled) return
    onAdd({
      note: note.trim(),
      status: status.trim(),
      diagnosis: diagnosis.trim(),
      therapy: therapy.trim(),
      otherInfo: otherInfo.trim(),
    })
    setNote('')
    setStatus('')
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
          <input
            type="search"
            className="select"
            style={{ minWidth: 200 }}
            placeholder="Traži po dijagnozi"
            value={searchTerm}
            onChange={(e) => onSearchChange(e.target.value)}
            disabled={disabled}
          />
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
          return (
            <article
              key={entry.uuid}
              className="note"
            >
              <div className="note__header">
                <div>
                  <p className="note__eyebrow">Posjet</p>
                </div>
                <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                  <button
                    type="button"
                    className="btn text"
                    onClick={() => {
                    // open selection modal with current includes
                    setSelectionModal({
                      open: true,
                      currentId: entry.uuid,
                      selected: new Set(entry.includeVisitUuids ?? []),
                    })
                    setSelectionList(anamneses)
                    void (async () => {
                      const all = await loadSelectionOptions()
                      setSelectionList(all)
                    })()
                  }}
                  disabled={disabled}
                  >
                    Generiraj PDF
                  </button>
                  <button
                    type="button"
                    className="btn text"
                    onClick={(e) => {
                      e.stopPropagation()
                      setEditModal({
                        open: true,
                        currentId: entry.uuid,
                        note: entry.note,
                        status: entry.status,
                        diagnosis: entry.diagnosis,
                        therapy: entry.therapy,
                        otherInfo: entry.otherInfo,
                      })
                    }}
                    disabled={disabled}
                  >
                    Uredi
                  </button>
                  <button
                    type="button"
                    className="btn text danger"
                    onClick={(e) => {
                      e.stopPropagation()
                      onDelete(entry.uuid)
                    }}
                    disabled={disabled}
                  >
                    Obriši
                  </button>
                </div>
              </div>
              <p className="note__body"><strong>Datum posjete:</strong> {new Date(entry.createdAt).toLocaleDateString('hr-HR')}</p>
              <p className="note__body"><strong>Dijagnoza:</strong> {entry.diagnosis}</p>
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
            <span className="muted-small">Stranica {page}</span>
            <button
              type="button"
              className="btn ghost small"
              onClick={() => onPageChange(page + 1)}
              disabled={!hasNext}
            >
              Sljedeće ▶
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

              <label htmlFor="status">Status</label>
              <textarea
                id="status"
                name="status"
                placeholder="Status..."
                value={status}
                onChange={(e) => setStatus(e.target.value)}
                rows={2}
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

      {selectionModal.open && (
        <div className="modal-backdrop">
          <div className="modal">
            <p className="eyebrow">Odaberi posjete</p>
            <h3>Uključi prethodne posjete u PDF</h3>
            <input
              type="search"
              className="select"
              placeholder="Traži po dijagnozi ili datumu"
              value={selectionSearch}
              onChange={(e) => setSelectionSearch(e.target.value)}
              style={{ marginBottom: 12 }}
            />
            <div className="list" style={{ maxHeight: 300, overflowY: 'auto' }}>
              <label
                htmlFor="only-current"
                style={{
                  display: 'grid',
                  gridTemplateColumns: '1fr auto',
                  alignItems: 'center',
                  gap: 12,
                  padding: '10px 0',
                }}
              >
                <span style={{ fontSize: 15, color: '#4a4a4a' }}>Generiraj samo ovaj posjet</span>
                <input
                  type="checkbox"
                  id="only-current"
                  checked={onlyCurrent}
                  onChange={(e) => {
                    const checked = e.target.checked
                    setOnlyCurrent(checked)
                    if (checked) {
                      setSelectionModal((prev) => ({ ...prev, selected: new Set() }))
                    }
                  }}
                  style={{ justifySelf: 'end' }}
                />
              </label>
              {(selectionList.length ? selectionList : anamneses)
                .filter((a) => a.uuid !== selectionModal.currentId)
                .filter((a) => {
                  const term = selectionSearch.trim().toLowerCase()
                  if (!term) return true
                  return (
                    a.diagnosis.toLowerCase().includes(term) ||
                    new Date(a.createdAt).toLocaleDateString('hr-HR').includes(term)
                  )
                })
                .map((a) => {
                  const checked = selectionModal.selected.has(a.uuid)
                    return (
                      <label
                        key={a.uuid}
                        className="checkbox-row"
                        style={{
                          display: 'grid',
                          gridTemplateColumns: '1fr auto',
                          alignItems: 'center',
                          gap: 12,
                          padding: '10px 0',
                          opacity: onlyCurrent ? 0.4 : 1,
                        }}
                      >
                        <span style={{ fontSize: 15 }}>
                          {new Date(a.createdAt).toLocaleDateString('hr-HR')} — {a.diagnosis}
                        </span>
                        <input
                          type="checkbox"
                          checked={checked}
                          disabled={onlyCurrent}
                          onChange={(e) => {
                            setSelectionModal((prev) => {
                              const next = new Set(prev.selected)
                              if (e.target.checked) next.add(a.uuid)
                              else next.delete(a.uuid)
                            return { ...prev, selected: next }
                          })
                        }}
                        style={{ justifySelf: 'end' }}
                      />
                    </label>
                  )
                })}
              {anamneses.filter((a) => a.uuid !== selectionModal.currentId).length === 0 && (
                <p className="muted-small">Nema prethodnih posjeta.</p>
              )}
            </div>
            <div className="actions">
              <button
                className="btn ghost"
                onClick={() => setSelectionModal({ open: false, currentId: null, selected: new Set() })}
              >
                Odustani
              </button>
              <button
                className="btn primary"
                onClick={async () => {
                  if (!selectionModal.currentId) return
                  const selectedIds = Array.from(selectionModal.selected)
                  if (onlyCurrent) {
                    onGeneratePdf(selectionModal.currentId, [], true)
                  } else if (selectedIds.length > 0) {
                    await onUpdateIncludes(selectionModal.currentId, selectedIds)
                    onGeneratePdf(selectionModal.currentId, selectedIds)
                  } else {
                    // If nothing selected, keep previously saved includes on the record.
                    onGeneratePdf(selectionModal.currentId)
                  }
                  setSelectionModal({ open: false, currentId: null, selected: new Set() })
                  setOnlyCurrent(false)
                }}
              >
                Spremi i generiraj
              </button>
            </div>
          </div>
        </div>
      )}

      {editModal.open && (
        <div className="modal-backdrop">
          <div className="modal">
            <p className="eyebrow">Uređivanje</p>
            <h3>Uredi posjet</h3>
            <form
              className="composer"
              onSubmit={(event) => {
                event.preventDefault()
                if (!editModal.currentId) return
                onEdit(editModal.currentId, {
                  note: editModal.note.trim(),
                  status: editModal.status.trim(),
                  diagnosis: editModal.diagnosis.trim(),
                  therapy: editModal.therapy.trim(),
                  otherInfo: editModal.otherInfo.trim(),
                })
                setEditModal({
                  open: false,
                  currentId: null,
                  note: '',
                  status: '',
                  diagnosis: '',
                  therapy: '',
                  otherInfo: '',
                })
              }}
            >
              <label htmlFor="edit-note">Anamneza</label>
              <textarea
                id="edit-note"
                name="edit-note"
                value={editModal.note}
                onChange={(e) => setEditModal((prev) => ({ ...prev, note: e.target.value }))}
                rows={3}
                disabled={disabled}
              />

              <label htmlFor="edit-status">Status</label>
              <textarea
                id="edit-status"
                name="edit-status"
                value={editModal.status}
                onChange={(e) => setEditModal((prev) => ({ ...prev, status: e.target.value }))}
                rows={2}
                disabled={disabled}
              />

              <label htmlFor="edit-diagnosis">Dijagnoza</label>
              <textarea
                id="edit-diagnosis"
                name="edit-diagnosis"
                value={editModal.diagnosis}
                onChange={(e) => setEditModal((prev) => ({ ...prev, diagnosis: e.target.value }))}
                rows={2}
                disabled={disabled}
              />

              <label htmlFor="edit-therapy">Terapija</label>
              <textarea
                id="edit-therapy"
                name="edit-therapy"
                value={editModal.therapy}
                onChange={(e) => setEditModal((prev) => ({ ...prev, therapy: e.target.value }))}
                rows={2}
                disabled={disabled}
              />

              <label htmlFor="edit-other">Ostale informacije</label>
              <textarea
                id="edit-other"
                name="edit-other"
                value={editModal.otherInfo}
                onChange={(e) => setEditModal((prev) => ({ ...prev, otherInfo: e.target.value }))}
                rows={2}
                disabled={disabled}
              />
              <div className="actions">
                <button
                  type="button"
                  className="btn ghost"
                  onClick={() =>
                    setEditModal({
                      open: false,
                      currentId: null,
                      note: '',
                      status: '',
                      diagnosis: '',
                      therapy: '',
                      otherInfo: '',
                    })
                  }
                >
                  Odustani
                </button>
                <button type="submit" className="btn primary" disabled={disabled}>
                  Spremi promjene
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
