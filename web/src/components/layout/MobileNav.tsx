import { NavLink } from 'react-router-dom'
import { ICONS } from '../ui/Icons'

const items = [
  { to: '/', icon: ICONS.dashboard, label: 'Home' },
  { to: '/transactions', icon: ICONS.transactions, label: 'Txns' },
  { to: '/wallets', icon: ICONS.wallet, label: 'Wallets' },
  { to: '/budgets', icon: ICONS.budget, label: 'Budgets' },
  { to: '/settings', icon: ICONS.settings, label: 'Settings' },
]

export default function MobileNav() {
  return (
    <nav
      style={{
        position: 'fixed',
        bottom: 0,
        left: 0,
        right: 0,
        background: 'var(--color-surface)',
        borderTop: '1px solid var(--color-border)',
        zIndex: 200,
        padding: '6px 8px env(safe-area-inset-bottom, 8px)',
      }}
      className="flex! md:hidden!"
    >
      <div style={{ display: 'flex', justifyContent: 'space-around', width: '100%' }}>
        {items.map(item => (
          <NavLink
            key={item.to}
            to={item.to}
            style={({ isActive }) => ({
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              gap: 2,
              padding: '8px 12px',
              borderRadius: 'var(--radius-sm)',
              textDecoration: 'none',
              color: isActive ? 'var(--color-primary)' : 'var(--color-text-tertiary)',
              fontSize: 10,
              fontWeight: isActive ? 700 : 500,
              transition: 'color var(--transition-fast)',
            })}
          >
            {({ isActive }) => (
              <>
                <span style={{ opacity: isActive ? 1 : 0.6 }}>{item.icon(22)}</span>
                <span>{item.label}</span>
              </>
            )}
          </NavLink>
        ))}
      </div>
    </nav>
  )
}
