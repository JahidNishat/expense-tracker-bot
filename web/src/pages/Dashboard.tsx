import { useQuery } from '@tanstack/react-query'
import { getChartData } from '../api/endpoints'
import ExpenseDonut from '../components/charts/ExpenseDonut'
import IncomeVsExpense from '../components/charts/IncomeVsExpense'
import BudgetGauge from '../components/charts/BudgetGauge'
import { useTransactions } from '../hooks/useTransactions'

export default function Dashboard() {
  const { data: charts, isLoading } = useQuery({
    queryKey: ['chartData'],
    queryFn: () => getChartData(),
  })
  const { data: txns } = useTransactions()

  if (isLoading) return <p className="text-gray-500">Loading...</p>
  if (!charts) return null

  const { overview, categories, comparison } = charts
  const recentTxns = (txns ?? []).slice(0, 10)

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Dashboard</h1>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Card label="Total Balance" value={fmt(overview.totalBalance)} color="text-blue-600" />
        <Card label="Month Income" value={fmt(overview.monthIncome)} color="text-green-600" />
        <Card label="Month Expense" value={fmt(overview.monthExpense)} color="text-red-600" />
        <div className="bg-white rounded-lg shadow p-4">
          <p className="text-xs text-gray-500 mb-1">Budget Usage</p>
          <BudgetGauge percent={overview.budgetUsage} />
          <p className="text-sm font-semibold mt-1">{overview.budgetUsage.toFixed(0)}%</p>
        </div>
      </div>

      <div className="grid md:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow p-4">
          <h2 className="text-sm font-semibold mb-2">Expense by Category</h2>
          <ExpenseDonut data={categories ?? []} />
        </div>
        <div className="bg-white rounded-lg shadow p-4">
          <h2 className="text-sm font-semibold mb-2">Income vs Expense</h2>
          <IncomeVsExpense data={comparison ?? []} />
        </div>
      </div>

      <div className="bg-white rounded-lg shadow p-4">
        <h2 className="text-sm font-semibold mb-3">Recent Transactions</h2>
        {recentTxns.length === 0 ? (
          <p className="text-gray-400 text-sm">No recent transactions</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="text-left text-gray-500 border-b">
                  <th className="pb-2">Type</th>
                  <th className="pb-2">Amount</th>
                  <th className="pb-2">Category</th>
                  <th className="pb-2">Date</th>
                </tr>
              </thead>
              <tbody>
                {recentTxns.map(t => (
                  <tr key={t.id} className="border-b last:border-0">
                    <td className="py-2">{t.type}</td>
                    <td className={`py-2 ${t.type === 'Income' ? 'text-green-600' : 'text-red-600'}`}>
                      {fmt(t.amount)}
                    </td>
                    <td className="py-2">{t.subcategoryId}</td>
                    <td className="py-2 text-gray-500">{new Date(t.timestamp * 1000).toLocaleDateString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}

function Card({ label, value, color }: { label: string; value: string; color: string }) {
  return (
    <div className="bg-white rounded-lg shadow p-4">
      <p className="text-xs text-gray-500">{label}</p>
      <p className={`text-xl font-bold ${color}`}>{value}</p>
    </div>
  )
}

function fmt(n: number) {
  return n.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}
