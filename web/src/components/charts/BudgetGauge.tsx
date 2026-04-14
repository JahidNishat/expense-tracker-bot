export default function BudgetGauge({ percent }: { percent: number }) {
  const color = percent > 90 ? 'bg-red-500' : percent > 70 ? 'bg-yellow-500' : 'bg-green-500'
  return (
    <div className="w-full bg-gray-200 rounded-full h-3">
      <div
        className={`h-3 rounded-full ${color} transition-all`}
        style={{ width: `${Math.min(percent, 100)}%` }}
      />
    </div>
  )
}
