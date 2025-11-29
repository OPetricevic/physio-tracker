import { useState, type FormEvent } from 'react'

type NewPatientInput = {
  firstName: string
  lastName: string
  phone?: string
}

type Props = {
  onCreate: (input: NewPatientInput) => void
}

export function PatientForm({ onCreate }: Props) {
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [phone, setPhone] = useState('')

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!firstName.trim() || !lastName.trim()) return
    onCreate({ firstName: firstName.trim(), lastName: lastName.trim(), phone: phone.trim() || undefined })
    setFirstName('')
    setLastName('')
    setPhone('')
  }

  return (
    <form className="panel" onSubmit={handleSubmit}>
      <div className="panel-header">
        <div>
          <p className="eyebrow">Dodaj pacijenta</p>
          <h2>Novi zapis</h2>
        </div>
      </div>
      <div className="field-grid">
        <div className="field">
          <label htmlFor="firstName">Ime</label>
          <input
            id="firstName"
            name="firstName"
            autoComplete="given-name"
            value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
            required
          />
        </div>
        <div className="field">
          <label htmlFor="lastName">Prezime</label>
          <input
            id="lastName"
            name="lastName"
            autoComplete="family-name"
            value={lastName}
            onChange={(e) => setLastName(e.target.value)}
            required
          />
        </div>
        <div className="field">
          <label htmlFor="phone">Telefon</label>
          <input
            id="phone"
            name="phone"
            autoComplete="tel"
            placeholder="+385…"
            value={phone}
            onChange={(e) => setPhone(e.target.value)}
          />
        </div>
      </div>
      <div className="actions">
        <button type="submit" className="btn primary">
          Spremi pacijenta
        </button>
      </div>
    </form>
  )
}
