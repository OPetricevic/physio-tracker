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

export type AuthLoginResponse = {
  token: string
  expires_at: string
  doctor_uuid: string
  token_uuid?: string
}
