import type { Patient } from '../types'

type Props = {
  patients: Patient[]
  selectedUuid: string | null
  onSelect: (uuid: string) => void
  searchTerm: string
  onSearchChange: (value: string) => void
  recentPatients?: Patient[]
}

export function PatientList({
  patients,
  selectedUuid,
  onSelect,
  searchTerm,
  onSearchChange,
  recentPatients = [],
}: Props) {
  const showResults = searchTerm.trim().length > 0
  const limitedPatients = showResults ? patients.slice(0, 10) : []

  return (
    <div className="panel">
      <div className="panel-header">
        <div>
          <p className="eyebrow">Pacijenti</p>
          <h2>{showResults ? 'Rezultati' : 'Pretraga'}</h2>
        </div>
        {showResults && <span className="badge">{patients.length}</span>}
      </div>
      <div className="field">
        <label htmlFor="search">Pretraga</label>
        <input
          id="search"
          type="search"
          placeholder="Traži po imenu ili telefonu"
          value={searchTerm}
          onChange={(e) => onSearchChange(e.target.value)}
        />
      </div>
      {!showResults && recentPatients.length > 0 && (
        <div className="list">
          <p className="muted-small" style={{ marginTop: 0 }}>
            Nedavni pacijenti
          </p>
          {recentPatients.map((patient) => {
            const isSelected = patient.uuid === selectedUuid
            return (
              <button
                key={patient.uuid}
                className={`list-item ${isSelected ? 'is-selected' : ''}`}
                onClick={() => onSelect(patient.uuid)}
              >
                <div>
                  <p className="list-title">
                    {patient.firstName} {patient.lastName}
                  </p>
                  <p className="list-subtitle">{patient.phone || 'Telefon nije unesen'}</p>
                </div>
              </button>
            )
          })}
        </div>
      )}

      {showResults && (
        <div className="list">
          {limitedPatients.map((patient) => {
            const isSelected = patient.uuid === selectedUuid
            return (
              <button
                key={patient.uuid}
                className={`list-item ${isSelected ? 'is-selected' : ''}`}
                onClick={() => onSelect(patient.uuid)}
              >
                <div>
                  <p className="list-title">
                    {patient.firstName} {patient.lastName}
                  </p>
                  <p className="list-subtitle">{patient.phone || 'Telefon nije unesen'}</p>
                </div>
                <span className="pill">
                  {new Date(patient.createdAt).toLocaleDateString()}
                </span>
              </button>
            )
          })}
          {limitedPatients.length === 0 && (
            <div className="empty">
              {searchTerm.trim().length === 0 ? 'Upišite pojam za pretragu.' : 'Nema pacijenata za ovaj upit.'}
            </div>
          )}
          {patients.length > limitedPatients.length && (
            <div className="muted-small">Prikazano prvih 10; suzite pretragu za više.</div>
          )}
        </div>
      )}
    </div>
  )
}
