interface LogoProps {
  size?: number
  collapsed?: boolean
}

export default function Logo({ size = 32, collapsed = false }: LogoProps) {
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: collapsed ? 0 : 12 }}>
      <img
        src="/logo-short.svg"
        alt="Expense Tracker"
        style={{ height: size, width: size, borderRadius: size * 0.22 }}
      />
      {!collapsed && (
        <span style={{
          fontSize: size * 0.5,
          fontWeight: 700,
          letterSpacing: '-0.03em',
          color: 'var(--color-text-primary)',
          lineHeight: 1,
        }}>
          Expense<span style={{ color: 'var(--color-primary)' }}> Tracker</span>
        </span>
      )}
    </div>
  )
}
