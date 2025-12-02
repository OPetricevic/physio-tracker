import { useEffect, useMemo, useState } from 'react'
import { PatientList } from '../components/PatientList'
import { PatientForm } from '../components/PatientForm'
import { AnamnesisPanel } from '../components/AnamnesisPanel'
import { usePatients } from '../data/PatientsContext'
import type { Patient } from '../types'
import '../App.css'

export function WorkspacePage() {
  const { patients, anamneses, createPatient, addAnamnesis } = usePatients()
  const [selectedUuid, setSelectedUuid] = useState<string | null>(null)
  const [searchTerm, setSearchTerm] = useState('')
  const [page, setPage] = useState(1)
  const [showForm, setShowForm] = useState(false)
  const [recent, setRecent] = useState<string[]>([])
  const [reasonFilter, setReasonFilter] = useState('')
  const [selectedVisits, setSelectedVisits] = useState<Set<string>>(new Set())
  const pageSize = 5

  const filteredPatients = useMemo(() => {
    const term = searchTerm.trim().toLowerCase()
    if (!term) return patients
    return patients.filter((p) => {
      const name = `${p.firstName} ${p.lastName}`.toLowerCase()
      return name.includes(term) || (p.phone ?? '').toLowerCase().includes(term)
    })
  }, [patients, searchTerm])

  const recentPatients = useMemo(
    () => recent.map((id) => patients.find((p) => p.uuid === id)).filter(Boolean) as Patient[],
    [recent, patients],
  )

  const selectedPatient = patients.find((p) => p.uuid === selectedUuid) ?? null
  const selectedAnamneses = (selectedUuid && anamneses[selectedUuid]) || []
  useEffect(() => {
    setPage(1)
  }, [selectedUuid, selectedAnamneses.length])

  const sortedAnamneses = useMemo(() => {
    let list = [...selectedAnamneses].sort(
      (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
    )
    if (reasonFilter) {
      list = list.filter((a) => (a.visitReason || '').toLowerCase() === reasonFilter.toLowerCase())
    }
    return list
  }, [selectedAnamneses, reasonFilter])

  const reasonOptions = useMemo(() => {
    const set = new Set<string>()
    selectedAnamneses.forEach((a) => {
      if (a.visitReason) set.add(a.visitReason)
    })
    return Array.from(set)
  }, [selectedAnamneses])
  const totalPages = Math.max(1, Math.ceil(sortedAnamneses.length / pageSize))
  const currentPage = Math.min(page, totalPages)
  const pagedAnamneses = sortedAnamneses.slice((currentPage - 1) * pageSize, currentPage * pageSize)

  const handleCreatePatient = (input: { firstName: string; lastName: string; phone?: string }) => {
    const next = createPatient(input)
    setSelectedUuid(next.uuid)
    setRecent((prev) => [next.uuid, ...prev.filter((id) => id !== next.uuid)].slice(0, 5))
  }

  const handleAddAnamnesis = (input: { note: string; diagnosis?: string; therapy?: string; otherInfo?: string; visitReason?: string }) => {
    if (!selectedPatient) return
    addAnamnesis(selectedPatient.uuid, input)
    setPage(1)
  }

  const handleSelectPatient = (uuid: string) => {
    setSelectedUuid(uuid)
    setRecent((prev) => [uuid, ...prev.filter((id) => id !== uuid)].slice(0, 5))
    setSelectedVisits(new Set())
  }

  const handleGeneratePdf = (anamnesisUuid: string) => {
    alert(`Generirao bi se PDF za anamnezu ${anamnesisUuid}`)
  }

  const handleBackup = () => {
    alert(`Pokrenula bi se sigurnosna kopija za pacijenta ${selectedPatient?.uuid ?? ''}`)
  }

  const handleToggleVisit = (visitUuid: string) => {
    setSelectedVisits((prev) => {
      const next = new Set(prev)
      if (next.has(visitUuid)) next.delete(visitUuid)
      else next.add(visitUuid)
      return next
    })
  }

  const handleBulkPdf = () => {
    if (!selectedPatient || selectedVisits.size === 0) return
    alert(`Generirao bi se PDF za posjete: ${Array.from(selectedVisits).join(', ')}`)
  }

  return (
    <div className="page">
      <header className="app-header">
        <div>
          <p className="eyebrow">Physio Tracker</p>
          <h1>Radni prostor</h1>
          <p className="lede">
            Dodajte pacijente, bilježite anamneze i generirajte PDF zapise tretmana.
          </p>
        </div>
      </header>

      <main className="layout">
        <div className="column">
          <div className="toggle-row">
            <button className="btn ghost small" onClick={() => setShowForm((v) => !v)}>
              {showForm ? 'Sakrij obrazac' : 'Novi pacijent'}
            </button>
          </div>
          {showForm && <PatientForm onCreate={handleCreatePatient} />}
          <PatientList
            patients={filteredPatients}
            recentPatients={recentPatients}
            selectedUuid={selectedUuid}
            onSelect={handleSelectPatient}
            searchTerm={searchTerm}
            onSearchChange={setSearchTerm}
          />
        </div>
        <div className="column wide">
          {selectedPatient ? (
            <AnamnesisPanel
              patientName={`${selectedPatient.firstName} ${selectedPatient.lastName}`}
              anamneses={pagedAnamneses}
              page={currentPage}
              totalPages={totalPages}
              onPageChange={setPage}
              disabled={!selectedPatient}
              onAdd={handleAddAnamnesis}
              onGeneratePdf={handleGeneratePdf}
              onBackup={handleBackup}
              reasonFilter={reasonFilter}
              reasonOptions={reasonOptions}
              onReasonFilterChange={setReasonFilter}
              selectedVisits={selectedVisits}
              onToggleVisit={handleToggleVisit}
              onBulkPdf={handleBulkPdf}
            />
          ) : (
            <div className="panel empty">
              <h2>Odaberite pacijenta</h2>
              <p className="muted-small">Pronađite pacijenta pretragom ili dodajte novog.</p>
            </div>
          )}
        </div>
      </main>
    </div>
  )
}
