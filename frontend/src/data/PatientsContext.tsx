import { createContext, useContext, useMemo, useState, type ReactNode } from 'react'
import type { Anamnesis, Patient } from '../types'

type PatientsContextValue = {
  patients: Patient[]
  anamneses: Record<string, Anamnesis[]>
  createPatient: (input: {
    firstName: string
    lastName: string
    phone?: string
    address?: string
    dateOfBirth?: string
    sex?: string
  }) => Patient
  updatePatient: (
    uuid: string,
    input: { firstName: string; lastName: string; phone?: string; address?: string; dateOfBirth?: string; sex?: string },
  ) => void
  deletePatient: (uuid: string) => void
  addAnamnesis: (patientUuid: string, note: string) => Anamnesis | null
}

const PatientsContext = createContext<PatientsContextValue | undefined>(undefined)

const initialPatients: Patient[] = [
  {
    uuid: '8c3d0c66-7e6c-4a4a-8c5b-7d7d1b5f0a11',
    firstName: 'Mia',
    lastName: 'Horvat',
    phone: '+385 91 123 4567',
    address: 'Ulica 1, Zagreb',
    dateOfBirth: '1990-05-12',
    sex: 'Ž',
    createdAt: new Date().toISOString(),
  },
  {
    uuid: 'b5b0c4ad-4a2c-45b0-93f0-6ae9b94e8a22',
    firstName: 'Luka',
    lastName: 'Kovač',
    phone: '+385 98 987 6543',
    address: 'Ulica 2, Zagreb',
    dateOfBirth: '1988-11-03',
    sex: 'M',
    createdAt: new Date().toISOString(),
  },
]

const initialAnamneses: Record<string, Anamnesis[]> = {
  '8c3d0c66-7e6c-4a4a-8c5b-7d7d1b5f0a11': [
    {
      uuid: '2c16b71b-5c9f-4f5c-8a3a-8e3f1d2b5c01',
      patientUuid: '8c3d0c66-7e6c-4a4a-8c5b-7d7d1b5f0a11',
      note: 'Početna evaluacija. Ograničena mobilnost ramena, zadane vježbe istezanja i hlađenje ledom.',
      createdAt: new Date().toISOString(),
    },
    {
      uuid: '4a5bcf1e-1b6d-4d2c-a8f9-6f6edc2a9f11',
      patientUuid: '8c3d0c66-7e6c-4a4a-8c5b-7d7d1b5f0a11',
      note: 'Kontrola nakon 7 dana: poboljšana pokretljivost, dodane izometrijske vježbe i lagani otpor.',
      createdAt: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
    },
    {
      uuid: '7c2e1d9f-0f2b-4f0a-9c3d-3e5f2a1b7c22',
      patientUuid: '8c3d0c66-7e6c-4a4a-8c5b-7d7d1b5f0a11',
      note: 'Treći susret: manualna terapija, smanjen bol, zadržati program istezanja i dodati elastičnu traku.',
      createdAt: new Date(Date.now() - 14 * 24 * 60 * 60 * 1000).toISOString(),
    },
    {
      uuid: '9d1a6b3c-2f4e-4c8a-9b1d-5f6e7c8a9b33',
      patientUuid: '8c3d0c66-7e6c-4a4a-8c5b-7d7d1b5f0a11',
      note: 'Četvrti susret: stabilizacijske vježbe za rame, preporuka za hladni oblog nakon vježbanja.',
      createdAt: new Date(Date.now() - 21 * 24 * 60 * 60 * 1000).toISOString(),
    },
    {
      uuid: '1e3f5a7c-9b2d-4a5f-8c7e-0d1b2c3a4d44',
      patientUuid: '8c3d0c66-7e6c-4a4a-8c5b-7d7d1b5f0a11',
      note: 'Peti susret: skoro potpuni opseg pokreta, prelazak na održavanje 2x tjedno, upute za kućni program.',
      createdAt: new Date(Date.now() - 28 * 24 * 60 * 60 * 1000).toISOString(),
    },
  ],
}

export function PatientsProvider({ children }: { children: ReactNode }) {
  const [patients, setPatients] = useState<Patient[]>(initialPatients)
  const [anamneses, setAnamneses] = useState<Record<string, Anamnesis[]>>(initialAnamneses)

  const createPatient = (input: { firstName: string; lastName: string; phone?: string; address?: string; dateOfBirth?: string; sex?: string }) => {
    const next: Patient = {
      uuid: crypto.randomUUID(),
      firstName: input.firstName,
      lastName: input.lastName,
      phone: input.phone,
      address: input.address,
      dateOfBirth: input.dateOfBirth,
      sex: input.sex,
      createdAt: new Date().toISOString(),
    }
    setPatients((prev) => [next, ...prev])
    return next
  }

  const updatePatient = (uuid: string, input: { firstName: string; lastName: string; phone?: string; address?: string; dateOfBirth?: string; sex?: string }) => {
    setPatients((prev) =>
      prev.map((p) =>
        p.uuid === uuid
          ? {
              ...p,
              firstName: input.firstName,
              lastName: input.lastName,
              phone: input.phone,
              address: input.address,
              dateOfBirth: input.dateOfBirth,
              sex: input.sex,
            }
          : p,
      ),
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
