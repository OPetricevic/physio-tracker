import { useMemo, useState } from 'react'
import { PatientList } from '../components/PatientList'
import { PatientForm } from '../components/PatientForm'
import { AnamnesisPanel } from '../components/AnamnesisPanel'
import type { Anamnesis, Patient } from '../types'
import '../App.css'

const initialPatients: Patient[] = [
  { uuid: '8c3d0c66-7e6c-4a4a-8c5b-7d7d1b5f0a11', firstName: 'Mia', lastName: 'Horvat', phone: '+385 91 123 4567', createdAt: new Date().toISOString() },
  { uuid: 'b5b0c4ad-4a2c-45b0-93f0-6ae9b94e8a22', firstName: 'Luka', lastName: 'Kovač', phone: '+385 98 987 6543', createdAt: new Date().toISOString() },
]

const initialAnamneses: Record<string, Anamnesis[]> = {
  '8c3d0c66-7e6c-4a4a-8c5b-7d7d1b5f0a11': [
    {
      uuid: '2c16b71b-5c9f-4f5c-8a3a-8e3f1d2b5c01',
      patientUuid: '8c3d0c66-7e6c-4a4a-8c5b-7d7d1b5f0a11',
      note: 'Početna evaluacija. Ograničena mobilnost ramena, zadane vježbe istezanja i hlađenje ledom.',
      createdAt: new Date().toISOString(),
    },
  ],
}

export function WorkspacePage() {
  const [patients, setPatients] = useState<Patient[]>(initialPatients)
  const [anamneses, setAnamneses] = useState<Record<string, Anamnesis[]>>(initialAnamneses)
  const [selectedUuid, setSelectedUuid] = useState<string | null>(patients[0]?.uuid ?? null)
  const [searchTerm, setSearchTerm] = useState('')

  const filteredPatients = useMemo(() => {
    const term = searchTerm.trim().toLowerCase()
    if (!term) return patients
    return patients.filter((p) => {
      const name = `${p.firstName} ${p.lastName}`.toLowerCase()
      return name.includes(term) || (p.phone ?? '').toLowerCase().includes(term)
    })
  }, [patients, searchTerm])

  const selectedPatient = patients.find((p) => p.uuid === selectedUuid) ?? null
  const selectedAnamneses = (selectedUuid && anamneses[selectedUuid]) || []

  const handleCreatePatient = (input: { firstName: string; lastName: string; phone?: string }) => {
    const next: Patient = {
      uuid: crypto.randomUUID(),
      firstName: input.firstName,
      lastName: input.lastName,
      phone: input.phone,
      createdAt: new Date().toISOString(),
    }
    setPatients((prev) => [next, ...prev])
    setSelectedUuid(next.uuid)
  }

  const handleAddAnamnesis = (note: string) => {
    if (!selectedPatient) return
    const entry: Anamnesis = {
      uuid: crypto.randomUUID(),
      patientUuid: selectedPatient.uuid,
      note,
      createdAt: new Date().toISOString(),
    }
    setAnamneses((prev) => {
      const existing = prev[selectedPatient.uuid] ?? []
      return {
        ...prev,
        [selectedPatient.uuid]: [entry, ...existing],
      }
    })
  }

  const handleGeneratePdf = (anamnesisUuid: string) => {
    alert(`Generirao bi se PDF za anamnezu ${anamnesisUuid}`)
  }

  const handleBackup = () => {
    alert(`Pokrenula bi se sigurnosna kopija za pacijenta ${selectedPatient?.uuid ?? ''}`)
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
          <PatientForm onCreate={handleCreatePatient} />
          <PatientList
            patients={filteredPatients}
            selectedUuid={selectedUuid}
            onSelect={setSelectedUuid}
            searchTerm={searchTerm}
            onSearchChange={setSearchTerm}
          />
        </div>
        <div className="column wide">
          <AnamnesisPanel
            patientName={
              selectedPatient ? `${selectedPatient.firstName} ${selectedPatient.lastName}` : ''
            }
            anamneses={selectedAnamneses}
            disabled={!selectedPatient}
            onAdd={handleAddAnamnesis}
            onGeneratePdf={handleGeneratePdf}
            onBackup={handleBackup}
          />
        </div>
      </main>
    </div>
  )
}
