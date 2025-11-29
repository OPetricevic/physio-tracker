import { BrowserRouter, Navigate, NavLink, Outlet, Route, Routes } from 'react-router-dom'
import type { ReactElement } from 'react'
import { AuthProvider, useAuth } from './auth'
import { PatientsProvider } from './data/PatientsContext'
import { LoginPage } from './pages/Login'
import { RegisterPage } from './pages/Register'
import { WorkspacePage } from './pages/Workspace'
import { SchedulePage } from './pages/Schedule'
import { ManagePatientsPage } from './pages/ManagePatients'

function Protected({ children }: { children: ReactElement }) {
  const { user } = useAuth()
  if (!user) return <Navigate to="/login" replace />
  return children
}

function AppShell() {
  const { logout, user } = useAuth()

  return (
    <div className="app-shell">
      <header className="topbar">
        <div className="brand">
          <span className="logo-mark">PT</span>
          <div>
            <p className="eyebrow">Physio Tracker</p>
            <strong>Ordinacija</strong>
          </div>
        </div>
        <nav className="nav-links">
          <NavLink to="/app/pacijenti" className={({ isActive }) => `nav-link ${isActive ? 'is-active' : ''}`}>
            Pacijenti
          </NavLink>
          <NavLink to="/app/upravljanje" className={({ isActive }) => `nav-link ${isActive ? 'is-active' : ''}`}>
            Upravljanje
          </NavLink>
          <NavLink to="/app/raspored" className={({ isActive }) => `nav-link ${isActive ? 'is-active' : ''}`}>
            Raspored
          </NavLink>
        </nav>
        <div className="nav-actions">
          {user && <span className="pill muted">Prijavljeni: {user.email}</span>}
          <button className="btn ghost small" onClick={logout}>
            Odjava
          </button>
        </div>
      </header>
      <Outlet />
    </div>
  )
}

function App() {
  return (
    <AuthProvider>
      <PatientsProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<Redirector />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route
              path="/app"
              element={
                <Protected>
                  <AppShell />
                </Protected>
              }
            >
              <Route index element={<Navigate to="/app/pacijenti" replace />} />
              <Route path="pacijenti" element={<WorkspacePage />} />
              <Route path="upravljanje" element={<ManagePatientsPage />} />
              <Route path="raspored" element={<SchedulePage />} />
            </Route>
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </PatientsProvider>
    </AuthProvider>
  )
}

function Redirector() {
  const { user } = useAuth()
  if (user) return <Navigate to="/app/pacijenti" replace />
  return <Navigate to="/login" replace />
}

export default App
