import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'
import { apiRequest } from './api/client'
import type { AuthLoginResponse } from './api/dto'

const BYPASS_AUTH = import.meta.env.VITE_BYPASS_AUTH === 'true'

type AuthUser = {
  email: string
  doctorUuid: string
  token: string
  expiresAt: string
}

type AuthContextValue = {
  user: AuthUser | null
  login: (identifier: string, password: string) => Promise<void>
  register: (payload: {
    email: string
    username: string
    firstName: string
    lastName: string
    password: string
  }) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

const STORAGE_KEY = 'physio-tracker:user'

type LoginResponse = {
  token: string
  expires_at: string
  doctor_uuid: string
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null)

  useEffect(() => {
    const stored = window.localStorage.getItem(STORAGE_KEY)
    if (stored) {
      try {
        setUser(JSON.parse(stored))
      } catch {
        window.localStorage.removeItem(STORAGE_KEY)
      }
    } else if (BYPASS_AUTH) {
      setAndStore({
        email: 'dev@bypass.local',
        doctorUuid: 'dev-bypass',
        token: 'dev-token',
        expiresAt: '2099-01-01T00:00:00Z',
      })
    }
  }, [])

  const setAndStore = (next: AuthUser | null) => {
    setUser(next)
    if (next) {
      window.localStorage.setItem(STORAGE_KEY, JSON.stringify(next))
    } else {
      window.localStorage.removeItem(STORAGE_KEY)
    }
  }

  const login = async (identifier: string, password: string) => {
    if (BYPASS_AUTH) {
      setAndStore({
        email: identifier,
        doctorUuid: 'dev-bypass',
        token: 'dev-token',
        expiresAt: '2099-01-01T00:00:00Z',
      })
      return
    }
    const res = await apiRequest<LoginResponse>('/auth/login', {
      method: 'POST',
      body: { identifier, password },
    })
    setAndStore({
      email: identifier,
      doctorUuid: res.doctor_uuid,
      token: res.token,
      expiresAt: res.expires_at,
    })
  }

  const register = async (payload: {
    email: string
    username: string
    firstName: string
    lastName: string
    password: string
  }) => {
    if (BYPASS_AUTH) {
      setAndStore({
        email: payload.email,
        doctorUuid: 'dev-bypass',
        token: 'dev-token',
        expiresAt: '2099-01-01T00:00:00Z',
      })
      return
    }
    const res = await apiRequest<LoginResponse>('/auth/register', {
      method: 'POST',
      body: {
        email: payload.email,
        username: payload.username,
        first_name: payload.firstName,
        last_name: payload.lastName,
        password: payload.password,
      },
    })
    setAndStore({
      email: payload.email,
      doctorUuid: res.doctor_uuid,
      token: res.token,
      expiresAt: res.expires_at,
    })
  }

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      login,
      register,
      logout: () => setAndStore(null),
    }),
    [user],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
