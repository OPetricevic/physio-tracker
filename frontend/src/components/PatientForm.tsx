import { useState, type FormEvent } from 'react'

type NewPatientInput = {
  firstName: string
  lastName: string
  phone?: string
  address?: string
  dateOfBirth?: string
  sex?: string
}

type Props = {
  onCreate: (input: NewPatientInput) => void
}

export function PatientForm({ onCreate }: Props) {
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [phone, setPhone] = useState('')
  const [address, setAddress] = useState('')
  const [dateOfBirth, setDateOfBirth] = useState('')
  const [sex, setSex] = useState('')

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!firstName.trim() || !lastName.trim()) return
    onCreate({
      firstName: firstName.trim(),
      lastName: lastName.trim(),
      phone: phone.trim() || undefined,
      address: address.trim() || undefined,
      dateOfBirth: dateOfBirth || undefined,
      sex: sex || undefined,
    })
    setFirstName('')
    setLastName('')
    setPhone('')
    setAddress('')
    setDateOfBirth('')
    setSex('')
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
            onInvalid={(e) => e.currentTarget.setCustomValidity('Molimo unesite ime')}
            onInput={(e) => e.currentTarget.setCustomValidity('')}
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
            onInvalid={(e) => e.currentTarget.setCustomValidity('Molimo unesite prezime')}
            onInput={(e) => e.currentTarget.setCustomValidity('')}
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
        <div className="field">
          <label htmlFor="address">Adresa</label>
          <input
            id="address"
            name="address"
            placeholder="Ulica i grad"
            value={address}
            onChange={(e) => setAddress(e.target.value)}
          />
        </div>
        <div className="field">
          <label htmlFor="dateOfBirth">Datum rođenja</label>
          <input
            id="dateOfBirth"
            name="dateOfBirth"
            type="date"
            value={dateOfBirth}
            onChange={(e) => setDateOfBirth(e.target.value)}
          />
        </div>
        <div className="field">
          <label htmlFor="sex">Spol</label>
          <select id="sex" name="sex" value={sex} onChange={(e) => setSex(e.target.value)}>
            <option value="">Odaberite</option>
            <option value="M">M</option>
            <option value="Ž">Ž</option>
          </select>
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
