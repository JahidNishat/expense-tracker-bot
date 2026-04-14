import { createContext, useState, useCallback, useEffect, type ReactNode } from 'react'
import { setAccessToken, clearAccessToken, getAccessToken } from '../api/client'
import { logout as apiLogout } from '../api/endpoints'

interface AuthContextValue {
  isAuthenticated: boolean
  login: (token: string) => void
  logout: () => void
}

export const AuthContext = createContext<AuthContextValue>({
  isAuthenticated: false,
  login: () => {},
  logout: () => {},
})

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(() => !!getAccessToken())

  const login = useCallback((token: string) => {
    setAccessToken(token)
    setIsAuthenticated(true)
  }, [])

  const logout = useCallback(async () => {
    try { await apiLogout() } catch { /* ignore */ }
    clearAccessToken()
    setIsAuthenticated(false)
  }, [])

  useEffect(() => {
    // Try silent refresh on mount
    if (!isAuthenticated) {
      fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' })
        .then(res => res.ok ? res.json() : null)
        .then(data => {
          if (data?.accessToken) {
            setAccessToken(data.accessToken)
            setIsAuthenticated(true)
          }
        })
        .catch(() => {})
    }
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}
