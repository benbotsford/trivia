// Root page — immediately redirects to /banks.
// redirect() in Next.js throws a special error that the framework catches and
// converts into a 307 response, so nothing after it runs.
import { redirect } from 'next/navigation'

export default function Home() {
  redirect('/banks')
}
