export interface Patient {
  uuid: string
  firstName: string
  lastName: string
  phone?: string
  createdAt: string
}

export interface Anamnesis {
  uuid: string
  patientUuid: string
  note: string
  createdAt: string
}
