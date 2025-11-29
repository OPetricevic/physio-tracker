import type { Patient } from '../types'

type Props = {
  patients: Patient[]
  selectedUuid: string | null
  onSelect: (uuid: string) => void
  searchTerm: string
  onSearchChange: (value: string) => void
}

export function PatientList({
  patients,
  selectedUuid,
  onSelect,
  searchTerm,
  onSearchChange,
}: Props) {
  return (
    <div className="panel">
      <div className="panel-header">
        <div>
          <p className="eyebrow">Patients</p>
          <h2>Roster</h2>
        </div>
        <span className="badge">{patients.length}</span>
      </div>
      <div className="field">
        <label htmlFor="search">Search</label>
        <input
          id="search"
          type="search"
          placeholder="Search by name or phone"
          value={searchTerm}
          onChange={(e) => onSearchChange(e.target.value)}
        />
      </div>
      <div className="list">
        {patients.map((patient) => {
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
                <p className="list-subtitle">{patient.phone || 'No phone added'}</p>
              </div>
              <span className="pill">
                {new Date(patient.createdAt).toLocaleDateString()}
              </span>
            </button>
          )
        })}
        {patients.length === 0 && (
          <div className="empty">No patients match your search.</div>
        )}
      </div>
    </div>
  )
}
