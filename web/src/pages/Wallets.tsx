import { useQuery } from '@tanstack/react-query'
import { listWallets, listContacts } from '../api/endpoints'

export default function Wallets() {
  const { data: wallets, isLoading: wLoading } = useQuery({ queryKey: ['wallets'], queryFn: listWallets })
  const { data: contacts, isLoading: cLoading } = useQuery({ queryKey: ['contacts'], queryFn: listContacts })

  if (wLoading || cLoading) return <p className="text-gray-500">Loading...</p>

  return (
    <div className="space-y-8">
      <section>
        <h1 className="text-2xl font-bold mb-4">Wallets</h1>
        {(!wallets || wallets.length === 0) ? (
          <p className="text-gray-400 text-sm">No wallets</p>
        ) : (
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {wallets.map(w => (
              <div key={w.id} className="bg-white rounded-lg shadow p-4">
                <div className="flex items-center justify-between mb-2">
                  <h3 className="font-semibold">{w.name}</h3>
                  <span className="text-xs bg-gray-100 px-2 py-0.5 rounded">{w.type}</span>
                </div>
                <p className={`text-xl font-bold ${w.balance >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                  {w.balance.toLocaleString(undefined, { minimumFractionDigits: 2 })}
                </p>
                <p className="text-xs text-gray-400 mt-1">{w.shortName}</p>
              </div>
            ))}
          </div>
        )}
      </section>

      <section>
        <h2 className="text-xl font-bold mb-4">Contacts</h2>
        {(!contacts || contacts.length === 0) ? (
          <p className="text-gray-400 text-sm">No contacts</p>
        ) : (
          <div className="bg-white rounded-lg shadow overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="text-left text-gray-500 border-b">
                  <th className="p-3">Nickname</th>
                  <th className="p-3">Full Name</th>
                  <th className="p-3">Email</th>
                  <th className="p-3">Net Balance</th>
                  <th className="p-3">Last Txn</th>
                </tr>
              </thead>
              <tbody>
                {contacts.map(c => (
                  <tr key={c.id} className="border-b last:border-0 hover:bg-gray-50">
                    <td className="p-3 font-medium">{c.nickName}</td>
                    <td className="p-3">{c.fullName}</td>
                    <td className="p-3 text-gray-500">{c.email || '-'}</td>
                    <td className={`p-3 font-medium ${c.netBalance >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                      {c.netBalance.toLocaleString(undefined, { minimumFractionDigits: 2 })}
                    </td>
                    <td className="p-3 text-gray-500">
                      {c.lastTxnTimestamp ? new Date(c.lastTxnTimestamp * 1000).toLocaleDateString() : '-'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </div>
  )
}
