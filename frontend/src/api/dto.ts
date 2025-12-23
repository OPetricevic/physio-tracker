// DTOs that mirror backend JSON shapes.
export type PatientDto = {
  uuid: string
  doctor_uuid: string
  first_name: string
  last_name: string
  phone?: string | null
  address?: string | null
  date_of_birth?: string | null
  sex?: string | null
  created_at: string
  updated_at?: string | null
}

export type ListPatientsResponse = {
  patients: PatientDto[]
}

export type AnamnesisDto = {
  uuid: string
  patient_uuid: string
  anamnesis: string
  diagnosis: string
  therapy: string
  other_info: string
  include_visit_uuids?: string[] | null
  created_at: string
  updated_at?: string | null
}

export type ListAnamnesesResponse = {
  anamneses: AnamnesisDto[]
}

export type UpdateAnamnesisPayload = {
  include_visit_uuids?: string[]
}

export type AuthLoginResponse = {
  token: string
  expires_at: string
  doctor_uuid: string
  token_uuid?: string
}
