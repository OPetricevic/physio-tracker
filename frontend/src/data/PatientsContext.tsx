import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'
import { apiRequest } from '../api/client'
import type { ListPatientsResponse, PatientDto } from '../api/dto'
import { useAuth } from '../auth'
import type { Anamnesis, Patient } from '../types'

type PatientsContextValue = {
  patients: Patient[]
  anamneses: Record<string, Anamnesis[]> // still local/offline until backend endpoints exist
  loading: boolean
  error: string | null
  searchTerm: string
  setSearchTerm: (term: string) => void
  refresh: () => Promise<void>
  createPatient: (input: {
    firstName: string
    lastName: string
    phone?: string
    address?: string
    dateOfBirth?: string
    sex?: string
  }) => Promise<Patient | null>
  updatePatient: (
    uuid: string,
    input: { firstName: string; lastName: string; phone?: string; address?: string; dateOfBirth?: string; sex?: string },
  ) => Promise<void>
  deletePatient: (uuid: string) => Promise<void>
  addAnamnesis: (
    patientUuid: string,
    input: { note: string; diagnosis?: string; therapy?: string; otherInfo?: string; visitReason?: string },
  ) => Anamnesis | null
}

const PatientsContext = createContext<PatientsContextValue | undefined>(undefined)

function dtoToPatient(dto: PatientDto): Patient {
  return {
    uuid: dto.uuid,
    firstName: dto.first_name,
    lastName: dto.last_name,
    phone: dto.phone ?? undefined,
    address: dto.address ?? undefined,
    dateOfBirth: dto.date_of_birth ?? undefined,
    sex: dto.sex ?? undefined,
    createdAt: dto.created_at,
    updatedAt: dto.updated_at ?? undefined,
  }
}

export function PatientsProvider({ children }: { children: ReactNode }) {
  const { user } = useAuth()
  const [patients, setPatients] = useState<Patient[]>([])
  const [anamneses, setAnamneses] = useState<Record<string, Anamnesis[]>>({})
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [searchTerm, setSearchTerm] = useState('')

  const refresh = async () => {
    if (!user?.token) {
      setPatients([])
      return
    }
    setLoading(true)
    setError(null)
    try {
      const params = new URLSearchParams()
      params.set('page_size', '200')
      params.set('current_page', '1')
      if (searchTerm.trim()) params.set('query', searchTerm.trim())
      const res = await apiRequest<ListPatientsResponse>(`/patients?${params.toString()}`, {
        token: user.token,
      })
      setPatients(res.patients.map(dtoToPatient))
    } catch (err) {
      console.error(err)
      setError('Ne možemo dohvatiti pacijente.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (user?.token) {
      void refresh()
    } else {
      setPatients([])
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user?.token, searchTerm])

  const createPatient = async (input: {
    firstName: string
    lastName: string
    phone?: string
    address?: string
    dateOfBirth?: string
    sex?: string
  }) => {
    if (!user?.token) return null
    try {
      const dto = await apiRequest<PatientDto>('/patients/create', {
        method: 'POST',
        token: user.token,
        body: {
          first_name: input.firstName,
          last_name: input.lastName,
          phone: input.phone || undefined,
          address: input.address || undefined,
          date_of_birth: input.dateOfBirth || undefined,
          sex: input.sex || undefined,
        },
      })
      const mapped = dtoToPatient(dto)
      setPatients((prev) => [mapped, ...prev])
      return mapped
    } catch (err) {
      console.error(err)
      setError('Spremanje pacijenta nije uspjelo.')
      return null
    }
  }

  const updatePatient = async (
    uuid: string,
    input: { firstName: string; lastName: string; phone?: string; address?: string; dateOfBirth?: string; sex?: string },
  ) => {
    if (!user?.token) return
    try {
      const dto = await apiRequest<PatientDto>(`/patients/${uuid}`, {
        method: 'PATCH',
        token: user.token,
        body: {
          uuid,
          first_name: input.firstName,
          last_name: input.lastName,
          phone: input.phone || undefined,
          address: input.address || undefined,
          date_of_birth: input.dateOfBirth || undefined,
          sex: input.sex || undefined,
        },
      })
      const mapped = dtoToPatient(dto)
      setPatients((prev) => prev.map((p) => (p.uuid === uuid ? mapped : p)))
    } catch (err) {
      console.error(err)
      setError('Ažuriranje pacijenta nije uspjelo.')
    }
  }

  const deletePatient = async (uuid: string) => {
    if (!user?.token) return
    try {
      await apiRequest<null>(`/patients/${uuid}`, { method: 'DELETE', token: user.token })
      setPatients((prev) => prev.filter((p) => p.uuid !== uuid))
      setAnamneses((prev) => {
        const copy = { ...prev }
        delete copy[uuid]
        return copy
      })
    } catch (err) {
      console.error(err)
      setError('Brisanje pacijenta nije uspjelo.')
    }
  }

  const addAnamnesis = (
    patientUuid: string,
    input: { note: string; diagnosis?: string; therapy?: string; otherInfo?: string; visitReason?: string },
  ) => {
    const entry: Anamnesis = {
      uuid: crypto.randomUUID(),
      patientUuid,
      note: input.note,
      diagnosis: input.diagnosis,
      therapy: input.therapy,
      otherInfo: input.otherInfo,
      visitReason: input.visitReason,
      createdAt: new Date().toISOString(),
    }
    setAnamneses((prev) => {
      const existing = prev[patientUuid] ?? []
      return { ...prev, [patientUuid]: [entry, ...existing] }
    })
    return entry
  }

  const value = useMemo<PatientsContextValue>(
    () => ({
      patients,
      anamneses,
      loading,
      error,
      searchTerm,
      setSearchTerm,
      refresh,
      createPatient,
      updatePatient,
      deletePatient,
      addAnamnesis,
    }),
    [patients, anamneses, loading, error, searchTerm],
  )

  return <PatientsContext.Provider value={value}>{children}</PatientsContext.Provider>
}

export function usePatients() {
  const ctx = useContext(PatientsContext)
  if (!ctx) throw new Error('usePatients must be used within PatientsProvider')
  return ctx
}
