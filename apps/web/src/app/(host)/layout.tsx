import { redirect } from 'next/navigation'
import Navbar from '@/components/Navbar'
import { getSession } from '@/lib/session'
import { getMe } from '@/lib/api/user'

export default async function HostLayout({ children }: { children: React.ReactNode }) {
  // Middleware handles the redirect in production, but this guard catches the
  // case where middleware isn't running (e.g. local dev with DEV_AUTH_TOKEN).
  // When DEV_AUTH_TOKEN is set, getSession() returns null (no Auth0 cookie),
  // so we only redirect when both are absent.
  const session = await getSession()
  const isDevMode = Boolean(process.env.DEV_AUTH_TOKEN)

  if (!session && !isDevMode) {
    redirect('/api/auth/login')
  }

  const profile = await getMe()
  return (
    <>
      <Navbar profile={profile} />
      {children}
    </>
  )
}
