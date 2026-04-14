import { apiFetch } from './client'
import type {
  Transaction, Wallet, Contact, BudgetStatus, BudgetAlert,
  ChartData, TxnCategory, Profile,
} from '../types'

const API = '/api/v1'

// Auth
export const requestOTP = (identifier: string) =>
  apiFetch(`${API}/auth/request-otp`, { method: 'POST', body: JSON.stringify({ identifier }) })

export const verifyOTP = (identifier: string, code: string) =>
  apiFetch<{ accessToken: string }>(`${API}/auth/verify-otp`, {
    method: 'POST', body: JSON.stringify({ identifier, code }),
  })

export const initQR = () =>
  apiFetch<{ sessionID: string; deepLink: string }>(`${API}/auth/qr/init`, { method: 'POST' })

export const pollQR = (session: string) =>
  apiFetch<{ status: string; accessToken?: string }>(`${API}/auth/qr/status?session=${session}`)

export const logout = () =>
  apiFetch(`${API}/auth/logout`, { method: 'POST' })

// Transactions
export const listTransactions = (params?: Record<string, string>) => {
  const qs = params ? '?' + new URLSearchParams(params).toString() : ''
  return apiFetch<Transaction[]>(`${API}/transactions${qs}`)
}

export const createTransaction = (txn: Partial<Transaction>) =>
  apiFetch(`${API}/transactions`, { method: 'POST', body: JSON.stringify(txn) })

export const updateTransaction = (id: number, txn: Partial<Transaction>) =>
  apiFetch(`${API}/transactions/${id}`, { method: 'PUT', body: JSON.stringify(txn) })

export const deleteTransaction = (id: number) =>
  apiFetch(`${API}/transactions/${id}`, { method: 'DELETE' })

// Wallets
export const listWallets = () => apiFetch<Wallet[]>(`${API}/wallets`)

// Contacts
export const listContacts = () => apiFetch<Contact[]>(`${API}/contacts`)

// Budgets
export const listBudgets = () => apiFetch<BudgetStatus[]>(`${API}/budgets`)

export const setBudget = (categoryId: string, amount: number, alertAt: number) =>
  apiFetch(`${API}/budgets`, { method: 'POST', body: JSON.stringify({ categoryId, amount, alertAt }) })

export const deleteBudget = (categoryId: string) =>
  apiFetch(`${API}/budgets/${categoryId}`, { method: 'DELETE' })

export const getBudgetAlerts = () => apiFetch<BudgetAlert[]>(`${API}/budgets/alerts`)

// Summary
export const getChartData = (year?: number, month?: number, months?: number) => {
  const params = new URLSearchParams()
  if (year) params.set('year', String(year))
  if (month) params.set('month', String(month))
  if (months) params.set('months', String(months))
  return apiFetch<ChartData>(`${API}/summary/charts?${params}`)
}

// Categories
export const listCategories = () => apiFetch<TxnCategory[]>(`${API}/categories`)

// Profile
export const getProfile = () => apiFetch<Profile>(`${API}/profile`)
