import { createContext, useContext, useMemo, useState, type ReactNode } from 'react'
import type { Anamnesis, Patient } from '../types'

type PatientsContextValue = {
  patients: Patient[]
  anamneses: Record<string, Anamnesis[]>
  createPatient: (input: { firstName: string; lastName: string; phone?: string }) => Patient
  updatePatient: (uuid: string, input: { firstName: string; lastName: string; phone?: string }) => void
  deletePatient: (uuid: string) => void
  addAnamnesis: (patientUuid: string, note: string) => Anamnesis | null
}

const PatientsContext = createContext<PatientsContextValue | undefined>(undefined)

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

export function PatientsProvider({ children }: { children: ReactNode }) {
  const [patients, setPatients] = useState<Patient[]>(initialPatients)
  const [anamneses, setAnamneses] = useState<Record<string, Anamnesis[]>>(initialAnamneses)

  const createPatient = (input: { firstName: string; lastName: string; phone?: string }) => {
    const next: Patient = {
      uuid: crypto.randomUUID(),
      firstName: input.firstName,
      lastName: input.lastName,
      phone: input.phone,
      createdAt: new Date().toISOString(),
    }
    setPatients((prev) => [next, ...prev])
    return next
  }

  const updatePatient = (uuid: string, input: { firstName: string; lastName: string; phone?: string }) => {
    setPatients((prev) =>
      prev.map((p) => (p.uuid === uuid ? { ...p, firstName: input.firstName, lastName: input.lastName, phone: input.phone } : p)),
    )
  }

  const deletePatient = (uuid: string) => {
    setPatients((prev) => prev.filter((p) => p.uuid !== uuid))
    setAnamneses((prev) => {
      const copy = { ...prev }
      delete copy[uuid]
      return copy
    })
  }

  const addAnamnesis = (patientUuid: string, note: string) => {
    const entry: Anamnesis = {
      uuid: crypto.randomUUID(),
      patientUuid,
      note,
      createdAt: new Date().toISOString(),
    }
    setAnamneses((prev) => {
      const existing = prev[patientUuid] ?? []
      return {
        ...prev,
        [patientUuid]: [entry, ...existing],
      }
    })
    return entry
  }

  const value = useMemo<PatientsContextValue>(
    () => ({
      patients,
      anamneses,
      createPatient,
      updatePatient,
      deletePatient,
      addAnamnesis,
    }),
    [patients, anamneses],
  )

  return <PatientsContext.Provider value={value}>{children}</PatientsContext.Provider>
}

export function usePatients() {
  const ctx = useContext(PatientsContext)
  if (!ctx) throw new Error('usePatients must be used within PatientsProvider')
  return ctx
}
