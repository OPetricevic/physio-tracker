export interface Patient {
  uuid: string
  firstName: string
  lastName: string
  phone?: string
  address?: string
  dateOfBirth?: string
  sex?: string
  createdAt: string
}

export interface Anamnesis {
  uuid: string
  patientUuid: string
  note: string // Anamneza
  diagnosis?: string
  therapy?: string
  otherInfo?: string
  visitReason?: string
  createdAt: string // Datum dolaska
}
