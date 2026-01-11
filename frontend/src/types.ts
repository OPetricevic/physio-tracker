export interface Patient {
  uuid: string
  firstName: string
  lastName: string
  phone?: string
  address?: string
  dateOfBirth?: string
  sex?: string
  createdAt: string
  updatedAt?: string
}

export interface Anamnesis {
  uuid: string
  patientUuid: string
  note: string // Anamneza
  status: string
  diagnosis: string
  therapy: string
  otherInfo: string
  includeVisitUuids?: string[]
  createdAt: string // Datum dolaska
  updatedAt?: string
}
