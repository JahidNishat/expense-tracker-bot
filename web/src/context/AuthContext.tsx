import { createContext, useState, useCallback, useEffect, type ReactNode } from 'react'
import {
  setAccessToken, clearTokens, getAccessToken,
  setRefreshToken, getRefreshToken, getApiBase,
} from '../api/client'
import { logout as apiLogout } from '../api/endpoints'

interface AuthContextValue {
  isAuthenticated: boolean
  isLoading: boolean
  isAdmin: boolean
  login: (accessToken: string, refreshToken?: string) => void
  logout: () => void
}

function decodeIsAdmin(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return !!payload.is_admin
  } catch {
    return false
  }
}

export const AuthContext = createContext<AuthContextValue>({
  isAuthenticated: false,
  isLoading: true,
  isAdmin: false,
  login: () => {},
  logout: () => {},
})

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(() => !!getAccessToken())
  const [isLoading, setIsLoading] = useState(() => !getAccessToken() && !!getRefreshToken())
  const [isAdmin, setIsAdmin] = useState(() => {
    const t = getAccessToken()
    return t ? decodeIsAdmin(t) : false
  })

  const login = useCallback((accessToken: string, refreshToken?: string) => {
    setAccessToken(accessToken)
    if (refreshToken) setRefreshToken(refreshToken)
    setIsAuthenticated(true)
    setIsAdmin(decodeIsAdmin(accessToken))
  }, [])

  const logout = useCallback(async () => {
    try { await apiLogout() } catch { /* ignore */ }
    clearTokens()
    setIsAuthenticated(false)
    setIsAdmin(false)
  }, [])

  useEffect(() => {
    if (isAuthenticated) {
      setIsLoading(false)
      return
    }
    const rt = getRefreshToken()
    if (!rt) {
      setIsLoading(false)
      return
    }

    fetch(`${getApiBase()}/api/v1/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refreshToken: rt }),
      credentials: 'include',
    })
      .then(res => res.ok ? res.json() : null)
      .then(data => {
        if (data?.accessToken) {
          setAccessToken(data.accessToken)
          if (data.refreshToken) setRefreshToken(data.refreshToken)
          setIsAuthenticated(true)
          setIsAdmin(decodeIsAdmin(data.accessToken))
        } else {
          clearTokens()
        }
      })
      .catch(() => { clearTokens() })
      .finally(() => { setIsLoading(false) })
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <AuthContext.Provider value={{ isAuthenticated, isLoading, isAdmin, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}
