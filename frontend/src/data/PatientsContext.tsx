/* eslint-disable react-refresh/only-export-components */
/* eslint-disable react-hooks/exhaustive-deps */
import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'
import { apiRequest } from '../api/client'
import type { AnamnesisDto, ListAnamnesesResponse, ListPatientsResponse, PatientDto } from '../api/dto'
import { useAuth } from '../auth'
import type { Anamnesis, Patient } from '../types'

type PatientsContextValue = {
  patients: Patient[]
  anamneses: Record<string, Anamnesis[]>
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
  fetchAnamneses: (
    patientUuid: string,
    opts?: { query?: string; page?: number; pageSize?: number },
  ) => Promise<{ items: Anamnesis[]; hasNext: boolean }>
  createAnamnesis: (
    patientUuid: string,
    input: { note: string; diagnosis: string; therapy: string; otherInfo: string; includeVisitUuids?: string[] },
  ) => Promise<Anamnesis | null>
  deleteAnamnesis: (patientUuid: string, uuid: string) => Promise<void>
  updateAnamnesis: (
    patientUuid: string,
    uuid: string,
    payload: UpdateAnamnesisPayload,
  ) => Promise<Anamnesis | null>
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

function dtoToAnamnesis(dto: AnamnesisDto): Anamnesis {
  return {
    uuid: dto.uuid,
    patientUuid: dto.patient_uuid,
    note: dto.anamnesis,
    diagnosis: dto.diagnosis,
    therapy: dto.therapy,
    otherInfo: dto.other_info,
    includeVisitUuids: dto.include_visit_uuids ?? undefined,
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
      params.set('page_size', '10')
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
          doctor_uuid: user.doctorUuid,
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
          doctor_uuid: user.doctorUuid,
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

  const fetchAnamneses = async (
    patientUuid: string,
    opts: { query?: string; page?: number; pageSize?: number } = {},
  ): Promise<{ items: Anamnesis[]; hasNext: boolean }> => {
    if (!user?.token) return { items: [], hasNext: false }
    const page = opts.page ?? 1
    const pageSize = opts.pageSize ?? 5
    const params = new URLSearchParams()
    params.set('current_page', String(page))
    params.set('page_size', String(pageSize))
    if (opts.query?.trim()) params.set('query', opts.query.trim())
    const res = await apiRequest<ListAnamnesesResponse>(
      `/patients/${patientUuid}/anamneses?${params.toString()}`,
      { token: user.token },
    )
    const items = res.anamneses.map(dtoToAnamnesis)
    // Simple hasNext heuristic: if we got a full page, we assume there may be more.
    const hasNext = items.length === pageSize
    // Cache by patient
    setAnamneses((prev) => ({ ...prev, [patientUuid]: items }))
    return { items, hasNext }
  }

  const createAnamnesis = async (
    patientUuid: string,
    input: { note: string; diagnosis: string; therapy: string; otherInfo: string; includeVisitUuids?: string[] },
  ) => {
    if (!user?.token) return null
    try {
      const dto = await apiRequest<AnamnesisDto>(`/patients/${patientUuid}/anamneses`, {
        method: 'POST',
        token: user.token,
        body: {
          anamnesis: input.note,
          diagnosis: input.diagnosis,
          therapy: input.therapy,
          other_info: input.otherInfo,
          include_visit_uuids: input.includeVisitUuids ?? [],
        },
      })
      const mapped = dtoToAnamnesis(dto)
      // Optimistically prepend to cache
      setAnamneses((prev) => {
        const existing = prev[patientUuid] ?? []
        return { ...prev, [patientUuid]: [mapped, ...existing] }
      })
      return mapped
    } catch (err) {
      console.error(err)
      setError('Spremanje anamneze nije uspjelo.')
      return null
    }
  }

  const deleteAnamnesis = async (patientUuid: string, uuid: string) => {
    if (!user?.token) return
    try {
      await apiRequest<null>(`/patients/${patientUuid}/anamneses/${uuid}`, {
        method: 'DELETE',
        token: user.token,
      })
      setAnamneses((prev) => ({
        ...prev,
        [patientUuid]: (prev[patientUuid] ?? []).filter((a) => a.uuid !== uuid),
      }))
    } catch (err) {
      console.error(err)
      setError('Brisanje anamneze nije uspjelo.')
    }
  }

  const updateAnamnesis = async (
    patientUuid: string,
    uuid: string,
    payload: UpdateAnamnesisPayload,
  ): Promise<Anamnesis | null> => {
    if (!user?.token) return null
    try {
      const dto = await apiRequest<AnamnesisDto>(`/patients/${patientUuid}/anamneses/${uuid}`, {
        method: 'PATCH',
        token: user.token,
        body: payload,
      })
      const mapped = dtoToAnamnesis(dto)
      setAnamneses((prev) => ({
        ...prev,
        [patientUuid]: (prev[patientUuid] ?? []).map((a) => (a.uuid === uuid ? mapped : a)),
      }))
      return mapped
    } catch (err) {
      console.error(err)
      setError('Ažuriranje anamneze nije uspjelo.')
      return null
    }
  }

  const value = useMemo<PatientsContextValue>(() => {
    return {
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
      fetchAnamneses,
      createAnamnesis,
      deleteAnamnesis,
      updateAnamnesis,
    }
  }, [
    patients,
    anamneses,
    loading,
    error,
    searchTerm,
    refresh,
    createPatient,
    updatePatient,
    deletePatient,
    fetchAnamneses,
    createAnamnesis,
    deleteAnamnesis,
    updateAnamnesis,
  ])

  return <PatientsContext.Provider value={value}>{children}</PatientsContext.Provider>
}

export function usePatients() {
  const ctx = useContext(PatientsContext)
  if (!ctx) throw new Error('usePatients must be used within PatientsProvider')
  return ctx
}
