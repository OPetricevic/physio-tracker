import { useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { PatientList } from '../components/PatientList'
import { PatientForm } from '../components/PatientForm'
import { AnamnesisPanel } from '../components/AnamnesisPanel'
import { usePatients } from '../data/PatientsContext'
import type { Anamnesis, Patient } from '../types'
import '../App.css'

export function WorkspacePage() {
  const {
    patients,
    createPatient,
    fetchAnamneses,
    createAnamnesis,
    deleteAnamnesis,
    updateAnamnesis,
    generateAnamnesisPdf,
    searchTerm,
    setSearchTerm,
    loading,
    error,
  } = usePatients()
  const navigate = useNavigate()
  const [selectedUuid, setSelectedUuid] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [recent, setRecent] = useState<string[]>([])
  const [anamneses, setAnamneses] = useState<Anamnesis[]>([])
  const [anaPage, setAnaPage] = useState(1)
  const [anaHasNext, setAnaHasNext] = useState(false)
  const [anaQuery, setAnaQuery] = useState('')
  const pageSize = 5
  const selectionPageSize = 500

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

  const loadAnamneses = async (patientUuid: string, query: string, pageNum: number) => {
    const { items, hasNext } = await fetchAnamneses(patientUuid, { query, page: pageNum, pageSize })
    setAnamneses(items)
    setAnaHasNext(hasNext)
    setAnaPage(pageNum)
  }

  const loadSelectionOptions = async () => {
    if (!selectedUuid) return []
    const { items } = await fetchAnamneses(selectedUuid, {
      query: '',
      page: 1,
      pageSize: selectionPageSize,
    })
    return items
  }

  useEffect(() => {
    if (selectedUuid) {
      void loadAnamneses(selectedUuid, anaQuery, 1)
    } else {
      setAnamneses([])
      setAnaHasNext(false)
      setAnaPage(1)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedUuid])

  const handleCreatePatient = async (input: { firstName: string; lastName: string; phone?: string; address?: string; dateOfBirth?: string; sex?: string }) => {
    const next = await createPatient(input)
    if (next) {
      setSelectedUuid(next.uuid)
      setRecent((prev) => [next.uuid, ...prev.filter((id) => id !== next.uuid)].slice(0, 5))
    }
  }

  const handleAddAnamnesis = (input: { note: string; diagnosis?: string; therapy?: string; otherInfo?: string }) => {
    if (!selectedPatient) return
    void (async () => {
      await createAnamnesis(selectedPatient.uuid, {
        note: input.note,
        diagnosis: input.diagnosis || '',
        therapy: input.therapy || '',
        otherInfo: input.otherInfo || '',
      })
      await loadAnamneses(selectedPatient.uuid, anaQuery, 1)
    })()
  }

  const handleSelectPatient = (uuid: string) => {
    setSelectedUuid(uuid)
    setRecent((prev) => [uuid, ...prev.filter((id) => id !== uuid)].slice(0, 5))
  }

  const handleGeneratePdf = async (anamnesisUuid: string, includes?: string[], onlyCurrent?: boolean) => {
    if (!selectedPatient) return
    const blob = await generateAnamnesisPdf(selectedPatient.uuid, anamnesisUuid, includes, onlyCurrent)
    if (!blob) return
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'anamneza.pdf'
    a.click()
    URL.revokeObjectURL(url)
  }

  const handleBackup = () => {
    navigate('/app/sigurnosna-kopija')
  }

  const handleUpdateIncludes = async (anamnesisUuid: string, includes: string[]) => {
    if (!selectedPatient) return
    await updateAnamnesis(selectedPatient.uuid, anamnesisUuid, { include_visit_uuids: includes })
    await loadAnamneses(selectedPatient.uuid, anaQuery, anaPage)
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
          {loading && <p className="muted-small">Učitavanje pacijenata...</p>}
          {error && <p className="error-text">{error}</p>}
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
              anamneses={anamneses}
              loadSelectionOptions={loadSelectionOptions}
              searchTerm={anaQuery}
              onSearchChange={(term) => {
                setAnaQuery(term)
                if (selectedPatient) void loadAnamneses(selectedPatient.uuid, term, 1)
              }}
              page={anaPage}
              hasNext={anaHasNext}
              onPageChange={(next) => {
                if (!selectedPatient) return
                void loadAnamneses(selectedPatient.uuid, anaQuery, next)
              }}
              disabled={!selectedPatient}
              onAdd={handleAddAnamnesis}
              onDelete={(uuid) => {
                if (!selectedPatient) return
                void (async () => {
                  await deleteAnamnesis(selectedPatient.uuid, uuid)
                  await loadAnamneses(selectedPatient.uuid, anaQuery, 1)
                })()
              }}
              onUpdateIncludes={handleUpdateIncludes}
              onGeneratePdf={handleGeneratePdf}
              onBackup={handleBackup}
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
