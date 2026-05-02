import { useEffect } from 'react'
import { Outlet } from 'react-router-dom'
import Sidebar from './Sidebar'
import MobileNav from './MobileNav'

export default function AppLayout() {
  // Initialize dark mode from localStorage
  useEffect(() => {
    const saved = localStorage.getItem('theme')
    if (saved === 'dark') {
      document.documentElement.classList.add('dark')
    }
  }, [])

  return (
    <div style={{ display: 'flex', minHeight: '100vh', background: 'var(--color-bg)' }}>
      <Sidebar />
      <main className="px-4 pb-20 md:px-8 md:pb-10" style={{ flex: 1, minWidth: 0 }}>
        <Outlet />
      </main>
      <MobileNav />
    </div>
  )
}
