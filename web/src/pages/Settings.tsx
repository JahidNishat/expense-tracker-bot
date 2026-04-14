import { useQuery } from '@tanstack/react-query'
import { getProfile } from '../api/endpoints'

export default function Settings() {
  const { data: profile, isLoading } = useQuery({ queryKey: ['profile'], queryFn: getProfile })

  if (isLoading) return <p className="text-gray-500">Loading...</p>
  if (!profile) return <p className="text-gray-400">Could not load profile</p>

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Settings</h1>

      <div className="bg-white rounded-lg shadow p-6 max-w-md">
        <h2 className="text-sm font-semibold mb-4 text-gray-500">Profile</h2>
        <dl className="space-y-3 text-sm">
          <Row label="Username" value={profile.username} />
          <Row label="First Name" value={profile.firstName} />
          <Row label="Last Name" value={profile.lastName} />
          <Row label="Mobile" value={profile.mobileNumber || '-'} />
          <Row label="Timezone" value={profile.timezone || 'UTC'} />
          <Row label="Telegram ID" value={String(profile.telegramId)} />
        </dl>
      </div>
    </div>
  )
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between">
      <dt className="text-gray-500">{label}</dt>
      <dd className="font-medium">{value}</dd>
    </div>
  )
}
