import { NavLink } from 'react-router-dom'
import { LayoutDashboard, ArrowLeftRight, Wallet, Target, Settings, LogOut } from 'lucide-react'
import { useAuth } from '../../hooks/useAuth'

const links = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/transactions', icon: ArrowLeftRight, label: 'Transactions' },
  { to: '/wallets', icon: Wallet, label: 'Wallets' },
  { to: '/budgets', icon: Target, label: 'Budgets' },
  { to: '/settings', icon: Settings, label: 'Settings' },
]

export default function Sidebar() {
  const { logout } = useAuth()
  return (
    <aside className="hidden md:flex flex-col w-56 bg-gray-900 text-gray-300 min-h-screen p-4">
      <h2 className="text-lg font-bold text-white mb-8">Expense Tracker</h2>
      <nav className="flex-1 space-y-1">
        {links.map(({ to, icon: Icon, label }) => (
          <NavLink
            key={to}
            to={to}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2 rounded text-sm ${isActive ? 'bg-gray-800 text-white' : 'hover:bg-gray-800'}`
            }
          >
            <Icon size={18} />
            {label}
          </NavLink>
        ))}
      </nav>
      <button
        className="flex items-center gap-3 px-3 py-2 rounded text-sm hover:bg-gray-800 text-red-400"
        onClick={logout}
      >
        <LogOut size={18} /> Logout
      </button>
    </aside>
  )
}
