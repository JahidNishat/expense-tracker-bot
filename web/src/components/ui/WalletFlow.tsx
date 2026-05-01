interface WalletFlowProps {
  srcId?: string
  dstId?: string
  contactName?: string
}

function capitalize(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1)
}

export default function WalletFlow({ srcId, dstId, contactName }: WalletFlowProps) {
  const from = srcId || ''
  const to = dstId || contactName || ''

  const pillBase: React.CSSProperties = {
    fontSize: 11,
    fontWeight: 700,
    background: 'var(--color-bg)',
    padding: '4px 10px',
    borderRadius: 'var(--radius-sm)',
    color: 'var(--color-text-secondary)',
    whiteSpace: 'nowrap',
  }

  const arrowBase: React.CSSProperties = {
    fontSize: 13,
    fontWeight: 800,
    flexShrink: 0,
  }

  if (from && to) {
    return (
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 6 }}>
        <span style={{ ...pillBase, borderLeft: '3px solid var(--color-danger)', color: 'var(--color-danger)' }}>{capitalize(from)}</span>
        <span style={{ ...arrowBase, color: 'var(--color-primary)' }}>→</span>
        <span style={{ ...pillBase, borderLeft: '3px solid var(--color-success)', color: 'var(--color-success)' }}>{capitalize(to)}</span>
      </div>
    )
  }

  if (to) {
    return (
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 6 }}>
        <span style={{ ...arrowBase, color: 'var(--color-success)' }}>→</span>
        <span style={{ ...pillBase, borderLeft: '3px solid var(--color-success)' }}>{capitalize(to)}</span>
      </div>
    )
  }

  if (from) {
    return (
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 6 }}>
        <span style={{ ...arrowBase, color: 'var(--color-danger)' }}>←</span>
        <span style={{ ...pillBase, borderLeft: '3px solid var(--color-danger)' }}>{capitalize(from)}</span>
      </div>
    )
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center' }}>
      <span style={{ ...pillBase, color: 'var(--color-text-tertiary)' }}>—</span>
    </div>
  )
}
